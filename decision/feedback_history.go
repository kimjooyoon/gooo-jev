package decision

import (
	"errors"
	"strings"
	"time"
)

type FeedbackHistory struct {
	MetricName     string
	Summaries      []FeedbackSummary
	Trends         []FeedbackTrend
	NonAuthorizing bool
	HistoryDigest  string
}

func BuildFeedbackHistory(summaries []FeedbackSummary, comparedAt time.Time) (FeedbackHistory, error) {
	if len(summaries) == 0 {
		return FeedbackHistory{}, errors.New("feedback history requires summaries")
	}
	if comparedAt.IsZero() {
		return FeedbackHistory{}, errors.New("feedback history comparison time is required")
	}
	metricName := summaries[0].MetricName
	if strings.TrimSpace(metricName) == "" {
		return FeedbackHistory{}, errors.New("feedback history metric name is required")
	}
	for _, summary := range summaries {
		if err := summary.Validate(); err != nil {
			return FeedbackHistory{}, err
		}
		if summary.MetricName != metricName {
			return FeedbackHistory{}, errors.New("feedback history metric names do not match")
		}
	}
	trends := make([]FeedbackTrend, 0, len(summaries)-1)
	for index := 1; index < len(summaries); index++ {
		trend, err := CompareFeedbackSummaries(summaries[index-1], summaries[index], comparedAt)
		if err != nil {
			return FeedbackHistory{}, err
		}
		trends = append(trends, trend)
	}
	history := FeedbackHistory{
		MetricName:     metricName,
		Summaries:      append([]FeedbackSummary(nil), summaries...),
		Trends:         trends,
		NonAuthorizing: true,
	}
	history.HistoryDigest, _ = Digest(history)
	return history, nil
}

func (history FeedbackHistory) Validate() error {
	if strings.TrimSpace(history.MetricName) == "" ||
		len(history.Summaries) == 0 ||
		len(history.Trends) != len(history.Summaries)-1 ||
		!history.NonAuthorizing ||
		strings.TrimSpace(history.HistoryDigest) == "" {
		return errors.New("feedback history is incomplete")
	}
	for index, summary := range history.Summaries {
		if err := summary.Validate(); err != nil {
			return err
		}
		if summary.MetricName != history.MetricName {
			return errors.New("feedback history contains another metric")
		}
		if index == 0 {
			continue
		}
		trend := history.Trends[index-1]
		if trend.PreviousSummaryDigest != history.Summaries[index-1].SummaryDigest ||
			trend.CurrentSummaryDigest != summary.SummaryDigest {
			return errors.New("feedback history trend does not bind adjacent summaries")
		}
		expected, err := CompareFeedbackSummaries(history.Summaries[index-1], summary, trend.ComparedAt)
		if err != nil {
			return err
		}
		if expected.TrendDigest != trend.TrendDigest {
			return errors.New("feedback history trend digest does not match adjacent summaries")
		}
	}
	copy := history
	copy.HistoryDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != history.HistoryDigest {
		return errors.New("feedback history digest does not match its evidence")
	}
	return nil
}
