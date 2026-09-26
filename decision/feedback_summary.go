package decision

import (
	"errors"
	"math"
	"strings"
	"time"
)

type FeedbackSummary struct {
	MetricName      string
	Count           int
	ConfirmedCount  int
	RefutedCount    int
	UnknownCount    int
	MetricSum       float64
	MetricMean      float64
	MetricMin       float64
	MetricMax       float64
	LastRecordedAt  time.Time
	NonAuthorizing  bool
	SummaryDigest   string
}

func SummarizeFeedback(observations []FeedbackObservation, metricName string) (FeedbackSummary, error) {
	if strings.TrimSpace(metricName) == "" {
		return FeedbackSummary{}, errors.New("feedback summary metric name is required")
	}
	if len(observations) == 0 {
		return FeedbackSummary{}, errors.New("feedback summary requires observations")
	}
	summary := FeedbackSummary{
		MetricName:     metricName,
		MetricMin:      math.Inf(1),
		MetricMax:      math.Inf(-1),
		NonAuthorizing: true,
	}
	for index, observation := range observations {
		if err := observation.Validate(); err != nil {
			return FeedbackSummary{}, err
		}
		if observation.MetricName != metricName {
			return FeedbackSummary{}, errors.New("feedback summary metric names do not match")
		}
		summary.Count++
		summary.MetricSum += observation.MetricValue
		if observation.MetricValue < summary.MetricMin {
			summary.MetricMin = observation.MetricValue
		}
		if observation.MetricValue > summary.MetricMax {
			summary.MetricMax = observation.MetricValue
		}
		if observation.RecordedAt.After(summary.LastRecordedAt) {
			summary.LastRecordedAt = observation.RecordedAt
		}
		switch observation.Kind {
		case FeedbackConfirmed:
			summary.ConfirmedCount++
		case FeedbackRefuted:
			summary.RefutedCount++
		case FeedbackUnknown:
			summary.UnknownCount++
		default:
			return FeedbackSummary{}, errors.New("unsupported feedback kind")
		}
		if index == 0 && summary.LastRecordedAt.IsZero() {
			summary.LastRecordedAt = observation.RecordedAt
		}
	}
	summary.MetricMean = summary.MetricSum / float64(summary.Count)
	summary.SummaryDigest, _ = Digest(summary)
	return summary, nil
}

func (summary FeedbackSummary) Validate() error {
	if strings.TrimSpace(summary.MetricName) == "" || summary.Count <= 0 {
		return errors.New("feedback summary is incomplete")
	}
	if summary.ConfirmedCount+summary.RefutedCount+summary.UnknownCount != summary.Count {
		return errors.New("feedback summary kind counts do not match")
	}
	if summary.LastRecordedAt.IsZero() {
		return errors.New("feedback summary last recorded time is required")
	}
	if !summary.NonAuthorizing {
		return errors.New("feedback summary must remain non-authorizing")
	}
	if math.IsNaN(summary.MetricSum) || math.IsInf(summary.MetricSum, 0) ||
		math.IsNaN(summary.MetricMean) || math.IsInf(summary.MetricMean, 0) ||
		math.IsNaN(summary.MetricMin) || math.IsInf(summary.MetricMin, 0) ||
		math.IsNaN(summary.MetricMax) || math.IsInf(summary.MetricMax, 0) {
		return errors.New("feedback summary metrics must be finite")
	}
	copy := summary
	copy.SummaryDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != summary.SummaryDigest {
		return errors.New("feedback summary digest does not match its evidence")
	}
	return nil
}
