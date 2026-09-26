package decision

import "testing"

func TestObserveExecutionEnvelopeOutcomeMetricPreservesOneHotOutcomes(t *testing.T) {
	cases := []struct {
		name     string
		status   string
		accepted uint64
		rejected uint64
		hold     uint64
	}{
		{name: "accepted", status: "accepted", accepted: 1},
		{name: "rejected", status: "rejected", rejected: 1},
		{name: "hold", status: "hold", hold: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			metric := ObserveExecutionEnvelopeOutcomeMetric(ExecutionEnvelopeOutcomeMetricInput{
				Status:         tc.status,
				OutcomeDigest:  "outcome",
				NonAuthorizing: true,
			})
			if metric.Status != tc.status || metric.AcceptedCount != tc.accepted || metric.RejectedCount != tc.rejected || metric.HoldCount != tc.hold || !metric.NonAuthorizing {
				t.Fatalf("outcome metric was not preserved: %#v", metric)
			}
		})
	}
}

func TestObserveExecutionEnvelopeOutcomeMetricHoldsMissingEvidence(t *testing.T) {
	metric := ObserveExecutionEnvelopeOutcomeMetric(ExecutionEnvelopeOutcomeMetricInput{
		Status:         "accepted",
		MissingStage:   "review-outcome",
		NonAuthorizing: true,
	})
	if metric.Status != "UNKNOWN" || metric.HoldCount != 1 || metric.MissingStage != "review-outcome" || !metric.NonAuthorizing {
		t.Fatalf("missing outcome evidence escaped hold: %#v", metric)
	}
}

func TestObserveExecutionEnvelopeOutcomeMetricRejectsAuthorizationClaim(t *testing.T) {
	metric := ObserveExecutionEnvelopeOutcomeMetric(ExecutionEnvelopeOutcomeMetricInput{
		Status:         "accepted",
		OutcomeDigest:  "outcome",
		NonAuthorizing: false,
	})
	if metric.Status != "UNKNOWN" || metric.MissingStage != "authorization-boundary" || metric.NonAuthorizing {
		t.Fatalf("authorization claim escaped outcome metric: %#v", metric)
	}
}
