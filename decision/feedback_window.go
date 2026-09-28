package decision

import (
	"errors"
	"math"
	"strings"
)

// FeedbackWindow separates observed quality from missing feedback. Unknown
// observations contribute to coverage counts but never contaminate the known
// metric mean.
type FeedbackWindow struct {
	MetricName         string
	ObservationCount   int
	KnownCount         int
	UnknownCount       int
	ConfirmedCount     int
	RefutedCount       int
	KnownSum           float64
	KnownMean          float64
	KnownMin           float64
	KnownMax           float64
	ObservedCoverage   float64
	NonAuthorizing     bool
	EvidenceDigest     string
}

// MeasureFeedbackWindow builds a bounded observation window without treating
// missing feedback as a quality score.
func MeasureFeedbackWindow(observations []FeedbackObservation, metricName string) (FeedbackWindow, error) {
	if strings.TrimSpace(metricName) == "" {
		return FeedbackWindow{}, errors.New("feedback window metric name is required")
	}
	if len(observations) == 0 {
		return FeedbackWindow{}, errors.New("feedback window requires observations")
	}
	window := FeedbackWindow{
		MetricName:       metricName,
		ObservationCount: len(observations),
		KnownMin:         math.Inf(1),
		KnownMax:         math.Inf(-1),
		NonAuthorizing:   true,
	}
	for _, observation := range observations {
		if err := observation.Validate(); err != nil {
			return FeedbackWindow{}, err
		}
		if observation.MetricName != metricName {
			return FeedbackWindow{}, errors.New("feedback window metric names do not match")
		}
		switch observation.Kind {
		case FeedbackConfirmed:
			window.ConfirmedCount++
		case FeedbackRefuted:
			window.RefutedCount++
		case FeedbackUnknown:
			window.UnknownCount++
			continue
		default:
			return FeedbackWindow{}, errors.New("unsupported feedback kind")
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
	window.EvidenceDigest, _ = Digest(window)
	return window, nil
}

func (window FeedbackWindow) Validate() error {
	if strings.TrimSpace(window.MetricName) == "" || window.ObservationCount <= 0 {
		return errors.New("feedback window is incomplete")
	}
	if window.KnownCount+window.UnknownCount != window.ObservationCount {
		return errors.New("feedback window counts do not match")
	}
	if window.ConfirmedCount+window.RefutedCount != window.KnownCount {
		return errors.New("feedback window known counts do not match")
	}
	if math.IsNaN(window.KnownSum) || math.IsInf(window.KnownSum, 0) ||
		math.IsNaN(window.KnownMean) || math.IsInf(window.KnownMean, 0) ||
		math.IsNaN(window.KnownMin) || math.IsInf(window.KnownMin, 0) ||
		math.IsNaN(window.KnownMax) || math.IsInf(window.KnownMax, 0) {
		return errors.New("feedback window metrics must be finite")
	}
	if window.ObservedCoverage < 0 || window.ObservedCoverage > 1 ||
		window.ObservedCoverage != float64(window.KnownCount)/float64(window.ObservationCount) {
		return errors.New("feedback window observed coverage is invalid")
	}
	if window.KnownCount == 0 {
		if window.KnownSum != 0 || window.KnownMean != 0 || window.KnownMin != 0 || window.KnownMax != 0 {
			return errors.New("empty known feedback window has non-zero metrics")
		}
	} else if window.KnownMin > window.KnownMax ||
		window.KnownMean != window.KnownSum/float64(window.KnownCount) {
		return errors.New("feedback window known metrics are inconsistent")
	}
	if !window.NonAuthorizing {
		return errors.New("feedback window must remain non-authorizing")
	}
	if strings.TrimSpace(window.EvidenceDigest) == "" {
		return errors.New("feedback window evidence digest is required")
	}
	copy := window
	copy.EvidenceDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != window.EvidenceDigest {
		return errors.New("feedback window digest does not match its evidence")
	}
	return nil
}
