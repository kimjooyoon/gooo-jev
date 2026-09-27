package gooo

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

type UsageTrendDirection string

const (
	UsageTrendRising  UsageTrendDirection = "rising"
	UsageTrendFalling UsageTrendDirection = "falling"
	UsageTrendFlat    UsageTrendDirection = "flat"
)

type UsageObservationWindowTrend struct {
	MetricName             string              `json:"metric_name"`
	PreviousWindowDigest   string              `json:"previous_window_digest"`
	CurrentWindowDigest    string              `json:"current_window_digest"`
	PreviousKnownMean      float64             `json:"previous_known_mean"`
	CurrentKnownMean       float64             `json:"current_known_mean"`
	KnownMeanDelta         float64             `json:"known_mean_delta"`
	PreviousCoverage       float64             `json:"previous_coverage"`
	CurrentCoverage        float64             `json:"current_coverage"`
	CoverageDelta          float64             `json:"coverage_delta"`
	PreviousUnknownCount   int                 `json:"previous_unknown_count"`
	CurrentUnknownCount    int                 `json:"current_unknown_count"`
	Direction              UsageTrendDirection `json:"direction"`
	ComparedAt             time.Time           `json:"compared_at"`
	NonExecuting           bool                `json:"non_executing"`
	NonAuthorizing         bool                `json:"non_authorizing"`
	TrendDigest            string              `json:"trend_digest"`
}

// CompareUsageObservationWindows compares known quality separately from coverage.
func CompareUsageObservationWindows(previous, current UsageObservationWindow, comparedAt time.Time) (UsageObservationWindowTrend, error) {
	if err := previous.Validate(); err != nil {
		return UsageObservationWindowTrend{}, err
	}
	if err := current.Validate(); err != nil {
		return UsageObservationWindowTrend{}, err
	}
	if previous.MetricName != current.MetricName {
		return UsageObservationWindowTrend{}, errors.New("usage window trend metric names do not match")
	}
	if comparedAt.IsZero() {
		return UsageObservationWindowTrend{}, errors.New("usage window trend comparison time is required")
	}
	trend := UsageObservationWindowTrend{
		MetricName:           previous.MetricName,
		PreviousWindowDigest: previous.WindowDigest,
		CurrentWindowDigest:  current.WindowDigest,
		PreviousKnownMean:    previous.KnownMean,
		CurrentKnownMean:     current.KnownMean,
		KnownMeanDelta:       current.KnownMean - previous.KnownMean,
		PreviousCoverage:     previous.ObservedCoverage,
		CurrentCoverage:      current.ObservedCoverage,
		CoverageDelta:        current.ObservedCoverage - previous.ObservedCoverage,
		PreviousUnknownCount: previous.UnknownCount,
		CurrentUnknownCount:  current.UnknownCount,
		Direction:            UsageTrendFlat,
		ComparedAt:           comparedAt.UTC(),
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	switch {
	case trend.KnownMeanDelta > 0:
		trend.Direction = UsageTrendRising
	case trend.KnownMeanDelta < 0:
		trend.Direction = UsageTrendFalling
	}
	trend.TrendDigest = digestUsageObservationWindowTrend(trend)
	return trend, nil
}

func (trend UsageObservationWindowTrend) Validate() error {
	if strings.TrimSpace(trend.MetricName) == "" ||
		!validDigest(trend.PreviousWindowDigest) || !validDigest(trend.CurrentWindowDigest) ||
		!validDigest(trend.TrendDigest) || trend.ComparedAt.IsZero() {
		return errors.New("usage window trend is incomplete")
	}
	if math.IsNaN(trend.PreviousKnownMean) || math.IsInf(trend.PreviousKnownMean, 0) ||
		math.IsNaN(trend.CurrentKnownMean) || math.IsInf(trend.CurrentKnownMean, 0) ||
		math.IsNaN(trend.KnownMeanDelta) || math.IsInf(trend.KnownMeanDelta, 0) ||
		math.IsNaN(trend.PreviousCoverage) || math.IsInf(trend.PreviousCoverage, 0) ||
		math.IsNaN(trend.CurrentCoverage) || math.IsInf(trend.CurrentCoverage, 0) ||
		math.IsNaN(trend.CoverageDelta) || math.IsInf(trend.CoverageDelta, 0) {
		return errors.New("usage window trend metrics must be finite")
	}
	if trend.PreviousCoverage < 0 || trend.PreviousCoverage > 1 ||
		trend.CurrentCoverage < 0 || trend.CurrentCoverage > 1 ||
		trend.KnownMeanDelta != trend.CurrentKnownMean-trend.PreviousKnownMean ||
		trend.CoverageDelta != trend.CurrentCoverage-trend.PreviousCoverage ||
		trend.PreviousUnknownCount < 0 || trend.CurrentUnknownCount < 0 {
		return errors.New("usage window trend metrics are inconsistent")
	}
	switch trend.Direction {
	case UsageTrendRising:
		if trend.KnownMeanDelta <= 0 {
			return errors.New("rising usage trend has no positive delta")
		}
	case UsageTrendFalling:
		if trend.KnownMeanDelta >= 0 {
			return errors.New("falling usage trend has no negative delta")
		}
	case UsageTrendFlat:
		if trend.KnownMeanDelta != 0 {
			return errors.New("flat usage trend has a non-zero delta")
		}
	default:
		return errors.New("usage window trend direction is invalid")
	}
	if !trend.NonExecuting || !trend.NonAuthorizing {
		return errors.New("usage window trend crossed a capability boundary")
	}
	if digestUsageObservationWindowTrend(trend) != trend.TrendDigest {
		return errors.New("usage window trend digest does not match its evidence")
	}
	return nil
}

func digestUsageObservationWindowTrend(trend UsageObservationWindowTrend) string {
	return digestString(trend.MetricName + "|" + trend.PreviousWindowDigest + "|" + trend.CurrentWindowDigest + "|" +
		strconv.FormatFloat(trend.PreviousKnownMean, 'g', -1, 64) + "|" +
		strconv.FormatFloat(trend.CurrentKnownMean, 'g', -1, 64) + "|" +
		strconv.FormatFloat(trend.KnownMeanDelta, 'g', -1, 64) + "|" +
		strconv.FormatFloat(trend.PreviousCoverage, 'g', -1, 64) + "|" +
		strconv.FormatFloat(trend.CurrentCoverage, 'g', -1, 64) + "|" +
		strconv.FormatFloat(trend.CoverageDelta, 'g', -1, 64) + "|" +
		strconv.Itoa(trend.PreviousUnknownCount) + "|" + strconv.Itoa(trend.CurrentUnknownCount) + "|" +
		string(trend.Direction) + "|" + trend.ComparedAt.UTC().Format(time.RFC3339Nano) + "|" +
		strconv.FormatBool(trend.NonExecuting) + "|" + strconv.FormatBool(trend.NonAuthorizing))
}

