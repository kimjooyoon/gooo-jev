package decision

import "testing"

func TestCreateExecutionEnvelopeReviewHandoffBindsReviewRequired(t *testing.T) {
	result := CreateExecutionEnvelopeReviewHandoff(ExecutionEnvelopeReviewHandoffInput{
		Status:           "review-required",
		ReviewRequired:   true,
		ComparisonDigest: "comparison",
		NonAuthorizing:   true,
	})
	if result.Status != "handoff-required" || result.HandoffDigest == "" || !result.ReviewRequired || !result.NonAuthorizing {
		t.Fatalf("review-required disposition was not handed off: %#v", result)
	}
}

func TestCreateExecutionEnvelopeReviewHandoffKeepsObservationOnly(t *testing.T) {
	result := CreateExecutionEnvelopeReviewHandoff(ExecutionEnvelopeReviewHandoffInput{
		Status:           "observation-only",
		ComparisonDigest: "comparison",
		NonAuthorizing:   true,
	})
	if result.Status != "observation-only" || result.HandoffDigest != "comparison" || result.ReviewRequired || !result.NonAuthorizing {
		t.Fatalf("observation-only state became handoff-required: %#v", result)
	}
}

func TestCreateExecutionEnvelopeReviewHandoffHoldsIncompleteState(t *testing.T) {
	result := CreateExecutionEnvelopeReviewHandoff(ExecutionEnvelopeReviewHandoffInput{
		Status:         "review-required",
		ReviewRequired: true,
		NonAuthorizing: true,
	})
	if result.Status != "UNKNOWN" || result.HandoffDigest != "" || !result.NonAuthorizing {
		t.Fatalf("incomplete handoff escaped UNKNOWN: %#v", result)
	}
}

func TestCreateExecutionEnvelopeReviewHandoffRejectsAuthorizationClaim(t *testing.T) {
	result := CreateExecutionEnvelopeReviewHandoff(ExecutionEnvelopeReviewHandoffInput{
		Status:           "review-required",
		ReviewRequired:   true,
		ComparisonDigest: "comparison",
		NonAuthorizing:   false,
	})
	if result.Status != "UNKNOWN" || result.HandoffDigest != "" || result.NonAuthorizing {
		t.Fatalf("authorization claim escaped handoff: %#v", result)
	}
}
