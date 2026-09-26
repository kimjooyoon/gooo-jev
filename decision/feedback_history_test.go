package decision

import (
	"testing"
	"time"
)

func historySummary(t *testing.T, mean float64, at time.Time) FeedbackSummary {
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

func TestBuildFeedbackHistoryReplaysTrendSequence(t *testing.T) {
	summaries := []FeedbackSummary{
		historySummary(t, 0.1, time.Unix(210, 0).UTC()),
		historySummary(t, 0.4, time.Unix(220, 0).UTC()),
		historySummary(t, 0.2, time.Unix(230, 0).UTC()),
	}
	history, err := BuildFeedbackHistory(summaries, time.Unix(240, 0).UTC())
	if err != nil {
		t.Fatalf("BuildFeedbackHistory() error = %v", err)
	}
	if err := history.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if len(history.Trends) != 2 ||
		history.Trends[0].Direction != FeedbackTrendRising ||
		history.Trends[1].Direction != FeedbackTrendFalling {
		t.Fatalf("history = %#v", history)
	}
	history.Summaries[1].MetricMean = 999
	if err := history.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered feedback history")
	}
}
