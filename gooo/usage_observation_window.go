package gooo

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type UsageObservationWindow struct {
	MetricName       string  `json:"metric_name"`
	ObservationCount int     `json:"observation_count"`
	KnownCount       int     `json:"known_count"`
	UnknownCount     int     `json:"unknown_count"`
	ConfirmedCount   int     `json:"confirmed_count"`
	RefutedCount     int     `json:"refuted_count"`
	KnownSum         float64 `json:"known_sum"`
	KnownMean        float64 `json:"known_mean"`
	KnownMin         float64 `json:"known_min"`
	KnownMax         float64 `json:"known_max"`
	ObservedCoverage float64 `json:"observed_coverage"`
	SourceDigest     string  `json:"source_digest"`
	DiscoveryDigest  string  `json:"discovery_digest"`
	PlanDigest       string  `json:"plan_digest"`
	NonExecuting     bool    `json:"non_executing"`
	NonAuthorizing   bool    `json:"non_authorizing"`
	WindowDigest     string  `json:"window_digest"`
}

// MeasureUsageObservationWindow excludes UNKNOWN values from quality statistics
// while preserving their count in coverage.
func MeasureUsageObservationWindow(observations []UsageObservation, metricName string) (UsageObservationWindow, error) {
	if strings.TrimSpace(metricName) == "" {
		return UsageObservationWindow{}, errors.New("usage observation window metric name is required")
	}
	if len(observations) == 0 {
		return UsageObservationWindow{}, errors.New("usage observation window requires observations")
	}
	window := UsageObservationWindow{
		MetricName:       metricName,
		ObservationCount: len(observations),
		KnownMin:         math.Inf(1),
		KnownMax:         math.Inf(-1),
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	for index, observation := range observations {
		if err := observation.Validate(); err != nil {
			return UsageObservationWindow{}, err
		}
		if observation.MetricName != metricName {
			return UsageObservationWindow{}, errors.New("usage observation window metric names do not match")
		}
		if index == 0 {
			window.SourceDigest = observation.SourceDigest
			window.DiscoveryDigest = observation.DiscoveryDigest
			window.PlanDigest = observation.PlanDigest
		} else if observation.SourceDigest != window.SourceDigest ||
			observation.DiscoveryDigest != window.DiscoveryDigest ||
			observation.PlanDigest != window.PlanDigest {
			return UsageObservationWindow{}, errors.New("usage observation window provenance does not match")
		}
		switch observation.Kind {
		case UsageObservationConfirmed:
			window.ConfirmedCount++
		case UsageObservationRefuted:
			window.RefutedCount++
		case UsageObservationUnknown:
			window.UnknownCount++
			continue
		default:
			return UsageObservationWindow{}, errors.New("unsupported usage observation kind")
		}
		window.KnownCount++
		window.KnownSum += observation.MetricValue
		if observation.MetricValue < window.KnownMin {
			window.KnownMin = observation.MetricValue
		}
		if observation.MetricValue > window.KnownMax {
			window.KnownMax = observation.MetricValue
		}
	}
	if window.KnownCount == 0 {
		window.KnownMin = 0
		window.KnownMax = 0
	} else {
		window.KnownMean = window.KnownSum / float64(window.KnownCount)
	}
	window.ObservedCoverage = float64(window.KnownCount) / float64(window.ObservationCount)
	window.WindowDigest = digestUsageObservationWindow(window)
	return window, nil
}

func (window UsageObservationWindow) Validate() error {
	if strings.TrimSpace(window.MetricName) == "" || window.ObservationCount <= 0 ||
		window.KnownCount < 0 || window.UnknownCount < 0 ||
		window.ConfirmedCount < 0 || window.RefutedCount < 0 ||
		!validDigest(window.SourceDigest) || !validDigest(window.DiscoveryDigest) ||
		!validDigest(window.PlanDigest) || !validDigest(window.WindowDigest) {
		return errors.New("usage observation window is incomplete")
	}
	if window.KnownCount+window.UnknownCount != window.ObservationCount ||
		window.ConfirmedCount+window.RefutedCount != window.KnownCount {
		return errors.New("usage observation window counts do not match")
	}
	if math.IsNaN(window.KnownSum) || math.IsInf(window.KnownSum, 0) ||
		math.IsNaN(window.KnownMean) || math.IsInf(window.KnownMean, 0) ||
		math.IsNaN(window.KnownMin) || math.IsInf(window.KnownMin, 0) ||
		math.IsNaN(window.KnownMax) || math.IsInf(window.KnownMax, 0) ||
		math.IsNaN(window.ObservedCoverage) || math.IsInf(window.ObservedCoverage, 0) {
		return errors.New("usage observation window metrics must be finite")
	}
	if window.ObservedCoverage < 0 || window.ObservedCoverage > 1 ||
		window.ObservedCoverage != float64(window.KnownCount)/float64(window.ObservationCount) {
		return errors.New("usage observation window coverage is invalid")
	}
	if window.KnownCount == 0 {
		if window.KnownSum != 0 || window.KnownMean != 0 || window.KnownMin != 0 || window.KnownMax != 0 {
			return errors.New("empty usage observation window has non-zero metrics")
		}
	} else if window.KnownMin > window.KnownMax ||
		window.KnownMean != window.KnownSum/float64(window.KnownCount) {
		return errors.New("usage observation window known metrics are inconsistent")
	}
	if !window.NonExecuting || !window.NonAuthorizing {
		return errors.New("usage observation window crossed a capability boundary")
	}
	if digestUsageObservationWindow(window) != window.WindowDigest {
		return errors.New("usage observation window digest does not match its evidence")
	}
	return nil
}

func digestUsageObservationWindow(window UsageObservationWindow) string {
	return digestString(window.MetricName + "|" + strconv.Itoa(window.ObservationCount) + "|" +
		strconv.Itoa(window.KnownCount) + "|" + strconv.Itoa(window.UnknownCount) + "|" +
		strconv.Itoa(window.ConfirmedCount) + "|" + strconv.Itoa(window.RefutedCount) + "|" +
		strconv.FormatFloat(window.KnownSum, 'g', -1, 64) + "|" +
		strconv.FormatFloat(window.KnownMean, 'g', -1, 64) + "|" +
		strconv.FormatFloat(window.KnownMin, 'g', -1, 64) + "|" +
		strconv.FormatFloat(window.KnownMax, 'g', -1, 64) + "|" +
		strconv.FormatFloat(window.ObservedCoverage, 'g', -1, 64) + "|" +
		window.SourceDigest + "|" + window.DiscoveryDigest + "|" + window.PlanDigest + "|" +
		strconv.FormatBool(window.NonExecuting) + "|" + strconv.FormatBool(window.NonAuthorizing))
}

