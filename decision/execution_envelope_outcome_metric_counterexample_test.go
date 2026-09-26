package decision

import "testing"

func TestClassifyExecutionEnvelopeOutcomeMetricCounterexamplePreservesDisposition(t *testing.T) {
	cases := []struct {
		name           string
		delta          string
		status         string
		reviewRequired bool
	}{
		{name: "declined", delta: "declined", status: "review-required", reviewRequired: true},
		{name: "increased", delta: "increased", status: "observation-only"},
		{name: "unchanged", delta: "unchanged", status: "observation-only"},
		{name: "inconclusive", delta: "inconclusive", status: "hold"},
		{name: "unknown", delta: "UNKNOWN", status: "UNKNOWN"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ClassifyExecutionEnvelopeOutcomeMetricCounterexample(ExecutionEnvelopeOutcomeMetricCounterexampleInput{
				PreviousOutcomeDigest: "previous",
				CurrentOutcomeDigest:  "current",
				Delta:                 tc.delta,
				NonAuthorizing:        true,
			})
			if result.Status != tc.status || result.ReviewRequired != tc.reviewRequired || result.ComparisonDigest == "" || !result.NonAuthorizing {
				t.Fatalf("outcome disposition was not preserved: %#v", result)
			}
		})
	}
}

func TestClassifyExecutionEnvelopeOutcomeMetricCounterexampleRejectsAuthorizationClaim(t *testing.T) {
	result := ClassifyExecutionEnvelopeOutcomeMetricCounterexample(ExecutionEnvelopeOutcomeMetricCounterexampleInput{
		PreviousOutcomeDigest: "previous",
		CurrentOutcomeDigest:  "current",
		Delta:                 "declined",
		NonAuthorizing:        false,
	})
	if result.Status != "UNKNOWN" || result.ComparisonDigest != "" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped outcome counterexample: %#v", result)
	}
}

func TestClassifyExecutionEnvelopeOutcomeMetricCounterexampleHoldsMissingDigest(t *testing.T) {
	result := ClassifyExecutionEnvelopeOutcomeMetricCounterexample(ExecutionEnvelopeOutcomeMetricCounterexampleInput{
		Delta:          "declined",
		NonAuthorizing: true,
	})
	if result.Status != "UNKNOWN" || result.ComparisonDigest != "" || !result.NonAuthorizing {
		t.Fatalf("missing outcome comparison escaped UNKNOWN: %#v", result)
	}
}
