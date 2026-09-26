package decision

import "testing"

func TestCompareExecutionEnvelopeFeedbackMetricPreservesDirection(t *testing.T) {
	previous := ExecutionEnvelopeFeedbackMetric{FeedbackDigest: "previous", RejectedCount: 1, NonAuthorizing: true}
	current := ExecutionEnvelopeFeedbackMetric{FeedbackDigest: "current", AnalysisReadyCount: 1, NonAuthorizing: true}
	comparison := CompareExecutionEnvelopeFeedbackMetric(previous, current)
	if comparison.Delta != "increased" || !comparison.NonAuthorizing {
		t.Fatalf("analysis-ready increase was not preserved: %#v", comparison)
	}

	current.AnalysisReadyCount = 0
	current.RejectedCount = 1
	comparison = CompareExecutionEnvelopeFeedbackMetric(previous, current)
	if comparison.Delta != "unchanged" || !comparison.NonAuthorizing {
		t.Fatalf("unchanged feedback state was not preserved: %#v", comparison)
	}

	previous.AnalysisReadyCount = 1
	previous.RejectedCount = 0
	current.AnalysisReadyCount = 0
	current.RejectedCount = 0
	current.HoldCount = 1
	comparison = CompareExecutionEnvelopeFeedbackMetric(previous, current)
	if comparison.Delta != "declined" || !comparison.NonAuthorizing {
		t.Fatalf("analysis-ready decline was not preserved: %#v", comparison)
	}
}

func TestCompareExecutionEnvelopeFeedbackMetricRejectsInvalidObservation(t *testing.T) {
	previous := ExecutionEnvelopeFeedbackMetric{FeedbackDigest: "previous", AnalysisReadyCount: 1, HoldCount: 1, NonAuthorizing: true}
	current := ExecutionEnvelopeFeedbackMetric{FeedbackDigest: "current", AnalysisReadyCount: 1, NonAuthorizing: true}
	comparison := CompareExecutionEnvelopeFeedbackMetric(previous, current)
	if comparison.Delta != "inconclusive" || !comparison.NonAuthorizing {
		t.Fatalf("invalid one-hot metric escaped inconclusive: %#v", comparison)
	}
}

func TestCompareExecutionEnvelopeFeedbackMetricRejectsAuthorizationClaim(t *testing.T) {
	previous := ExecutionEnvelopeFeedbackMetric{FeedbackDigest: "previous", AnalysisReadyCount: 1, NonAuthorizing: false}
	current := ExecutionEnvelopeFeedbackMetric{FeedbackDigest: "current", AnalysisReadyCount: 1, NonAuthorizing: true}
	comparison := CompareExecutionEnvelopeFeedbackMetric(previous, current)
	if comparison.Delta != "UNKNOWN" || comparison.NonAuthorizing {
		t.Fatalf("authorization claim escaped metric comparison: %#v", comparison)
	}
}
