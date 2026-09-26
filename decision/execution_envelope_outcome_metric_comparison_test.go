package decision

import "testing"

func TestCompareExecutionEnvelopeOutcomeMetricPreservesDirection(t *testing.T) {
	previous := ExecutionEnvelopeOutcomeMetric{OutcomeDigest: "previous", RejectedCount: 1, NonAuthorizing: true}
	current := ExecutionEnvelopeOutcomeMetric{OutcomeDigest: "current", AcceptedCount: 1, NonAuthorizing: true}
	comparison := CompareExecutionEnvelopeOutcomeMetric(previous, current)
	if comparison.Delta != "increased" || !comparison.NonAuthorizing {
		t.Fatalf("accepted increase was not preserved: %#v", comparison)
	}

	current.AcceptedCount = 0
	current.RejectedCount = 1
	comparison = CompareExecutionEnvelopeOutcomeMetric(previous, current)
	if comparison.Delta != "unchanged" || !comparison.NonAuthorizing {
		t.Fatalf("unchanged outcome state was not preserved: %#v", comparison)
	}

	previous.AcceptedCount = 1
	previous.RejectedCount = 0
	current.AcceptedCount = 0
	current.RejectedCount = 0
	current.HoldCount = 1
	comparison = CompareExecutionEnvelopeOutcomeMetric(previous, current)
	if comparison.Delta != "declined" || !comparison.NonAuthorizing {
		t.Fatalf("accepted decline was not preserved: %#v", comparison)
	}
}

func TestCompareExecutionEnvelopeOutcomeMetricRejectsInvalidObservation(t *testing.T) {
	previous := ExecutionEnvelopeOutcomeMetric{OutcomeDigest: "previous", AcceptedCount: 1, HoldCount: 1, NonAuthorizing: true}
	current := ExecutionEnvelopeOutcomeMetric{OutcomeDigest: "current", AcceptedCount: 1, NonAuthorizing: true}
	comparison := CompareExecutionEnvelopeOutcomeMetric(previous, current)
	if comparison.Delta != "inconclusive" || !comparison.NonAuthorizing {
		t.Fatalf("invalid one-hot outcome escaped inconclusive: %#v", comparison)
	}
}

func TestCompareExecutionEnvelopeOutcomeMetricRejectsAuthorizationClaim(t *testing.T) {
	previous := ExecutionEnvelopeOutcomeMetric{OutcomeDigest: "previous", AcceptedCount: 1, NonAuthorizing: false}
	current := ExecutionEnvelopeOutcomeMetric{OutcomeDigest: "current", AcceptedCount: 1, NonAuthorizing: true}
	comparison := CompareExecutionEnvelopeOutcomeMetric(previous, current)
	if comparison.Delta != "UNKNOWN" || comparison.NonAuthorizing {
		t.Fatalf("authorization claim escaped outcome comparison: %#v", comparison)
	}
}
