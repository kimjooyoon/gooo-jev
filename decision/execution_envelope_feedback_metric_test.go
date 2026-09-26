package decision

import "testing"

func TestObserveExecutionEnvelopeFeedbackMetricPreservesReviewStates(t *testing.T) {
	cases := []struct {
		name          string
		status        string
		analysisReady uint64
		rejected      uint64
		hold          uint64
		observation   uint64
		metricStatus  string
	}{
		{name: "reviewed", status: "reviewed", analysisReady: 1, metricStatus: "analysis-ready"},
		{name: "review-rejected", status: "review-rejected", rejected: 1, metricStatus: "rejected"},
		{name: "hold", status: "hold", hold: 1, metricStatus: "hold"},
		{name: "observation-only", status: "observation-only", observation: 1, metricStatus: "observation-only"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			metric := ObserveExecutionEnvelopeFeedbackMetric(ExecutionEnvelopeFeedbackMetricInput{
				Status:         tc.status,
				FeedbackDigest: "feedback",
				NonAuthorizing: true,
			})
			if metric.Status != tc.metricStatus || metric.AnalysisReadyCount != tc.analysisReady || metric.RejectedCount != tc.rejected || metric.HoldCount != tc.hold || metric.ObservationOnlyCount != tc.observation || !metric.NonAuthorizing {
				t.Fatalf("feedback state was not preserved: %#v", metric)
			}
		})
	}
}

func TestObserveExecutionEnvelopeFeedbackMetricHoldsMissingEvidence(t *testing.T) {
	metric := ObserveExecutionEnvelopeFeedbackMetric(ExecutionEnvelopeFeedbackMetricInput{
		Status:       "UNKNOWN",
		MissingStage: "review-outcome",
		NonAuthorizing: true,
	})
	if metric.Status != "UNKNOWN" || metric.HoldCount != 1 || metric.MissingStage != "review-outcome" || !metric.NonAuthorizing {
		t.Fatalf("missing feedback evidence escaped hold: %#v", metric)
	}
}

func TestObserveExecutionEnvelopeFeedbackMetricRejectsAuthorizationClaim(t *testing.T) {
	metric := ObserveExecutionEnvelopeFeedbackMetric(ExecutionEnvelopeFeedbackMetricInput{
		Status:         "reviewed",
		FeedbackDigest: "feedback",
		NonAuthorizing: false,
	})
	if metric.Status != "UNKNOWN" || metric.MissingStage != "authorization-boundary" || metric.NonAuthorizing {
		t.Fatalf("authorization claim escaped metric: %#v", metric)
	}
}
