package decision

import (
	"testing"
	"time"
)

func trendSummary(t *testing.T, mean float64, at time.Time) FeedbackSummary {
	t.Helper()
	summary := FeedbackSummary{
		MetricName:     "utility.delta",
		Count:          1,
		ConfirmedCount: 1,
		MetricSum:      mean,
		MetricMean:     mean,
		MetricMin:      mean,
		MetricMax:      mean,
		LastRecordedAt: at,
		NonAuthorizing: true,
	}
	var err error
	summary.SummaryDigest, err = Digest(summary)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	return summary
}

func TestCompareFeedbackSummariesPreservesDirection(t *testing.T) {
	previous := trendSummary(t, 0.4, time.Unix(180, 0).UTC())
	current := trendSummary(t, 0.1, time.Unix(190, 0).UTC())
	trend, err := CompareFeedbackSummaries(previous, current, time.Unix(200, 0).UTC())
	if err != nil {
		t.Fatalf("CompareFeedbackSummaries() error = %v", err)
	}
	if err := trend.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if trend.Direction != FeedbackTrendFalling || trend.Delta >= 0 {
		t.Fatalf("trend = %#v", trend)
	}
	trend.Direction = FeedbackTrendRising
	if err := trend.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered trend direction")
	}
}
