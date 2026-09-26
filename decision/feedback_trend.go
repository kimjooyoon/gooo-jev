package decision

import (
	"errors"
	"math"
	"strings"
	"time"
)

type FeedbackTrendDirection string

const (
	FeedbackTrendRising FeedbackTrendDirection = "rising"
	FeedbackTrendFalling FeedbackTrendDirection = "falling"
	FeedbackTrendFlat FeedbackTrendDirection = "flat"
)

type FeedbackTrend struct {
	MetricName             string
	PreviousSummaryDigest  string
	CurrentSummaryDigest   string
	PreviousMean           float64
	CurrentMean            float64
	Delta                  float64
	Direction              FeedbackTrendDirection
	ComparedAt             time.Time
	NonAuthorizing         bool
	TrendDigest            string
}

func CompareFeedbackSummaries(previous, current FeedbackSummary, comparedAt time.Time) (FeedbackTrend, error) {
	if err := previous.Validate(); err != nil {
		return FeedbackTrend{}, err
	}
	if err := current.Validate(); err != nil {
		return FeedbackTrend{}, err
	}
	if previous.MetricName != current.MetricName {
		return FeedbackTrend{}, errors.New("feedback trend metric names do not match")
	}
	if comparedAt.IsZero() {
		return FeedbackTrend{}, errors.New("feedback trend comparison time is required")
	}
	delta := current.MetricMean - previous.MetricMean
	direction := FeedbackTrendFlat
	switch {
	case delta > 0:
		direction = FeedbackTrendRising
	case delta < 0:
		direction = FeedbackTrendFalling
	}
	trend := FeedbackTrend{
		MetricName:            previous.MetricName,
		PreviousSummaryDigest: previous.SummaryDigest,
		CurrentSummaryDigest:  current.SummaryDigest,
		PreviousMean:          previous.MetricMean,
		CurrentMean:           current.MetricMean,
		Delta:                 delta,
		Direction:             direction,
		ComparedAt:            comparedAt.UTC(),
		NonAuthorizing:        true,
	}
	trend.TrendDigest, _ = Digest(trend)
	return trend, nil
}

func (trend FeedbackTrend) Validate() error {
	if strings.TrimSpace(trend.MetricName) == "" ||
		strings.TrimSpace(trend.PreviousSummaryDigest) == "" ||
		strings.TrimSpace(trend.CurrentSummaryDigest) == "" ||
		strings.TrimSpace(trend.TrendDigest) == "" {
		return errors.New("feedback trend is incomplete")
	}
	if math.IsNaN(trend.PreviousMean) || math.IsInf(trend.PreviousMean, 0) ||
		math.IsNaN(trend.CurrentMean) || math.IsInf(trend.CurrentMean, 0) ||
		math.IsNaN(trend.Delta) || math.IsInf(trend.Delta, 0) {
		return errors.New("feedback trend metrics must be finite")
	}
	if trend.ComparedAt.IsZero() {
		return errors.New("feedback trend comparison time is required")
	}
	if !trend.NonAuthorizing {
		return errors.New("feedback trend must remain non-authorizing")
	}
	if trend.Delta != trend.CurrentMean-trend.PreviousMean {
		return errors.New("feedback trend delta does not match its means")
	}
	expectedDirection := FeedbackTrendFlat
	switch {
	case trend.Delta > 0:
		expectedDirection = FeedbackTrendRising
	case trend.Delta < 0:
		expectedDirection = FeedbackTrendFalling
	}
	if trend.Direction != expectedDirection {
		return errors.New("feedback trend direction does not match its delta")
	}
	copy := trend
	copy.TrendDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != trend.TrendDigest {
		return errors.New("feedback trend digest does not match its evidence")
	}
	return nil
}
