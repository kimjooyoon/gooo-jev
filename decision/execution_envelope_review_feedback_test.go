package decision

import "testing"

func TestRecordExecutionEnvelopeReviewFeedbackBindsAcceptedReview(t *testing.T) {
	result := RecordExecutionEnvelopeReviewFeedback(ExecutionEnvelopeReviewFeedbackInput{
		DispositionStatus:      "review-required",
		ReviewRequired:         true,
		MismatchStage:          "status",
		EvidenceDigest:         "observed",
		ReviewStatus:           "accepted",
		ReviewerEvidenceDigest: "reviewer",
		NonAuthorizing:         true,
	})
	if result.Status != "reviewed" || result.FeedbackDigest == "" || result.MissingStage != "" || !result.NonAuthorizing {
		t.Fatalf("accepted review was not bound: %#v", result)
	}
}

func TestRecordExecutionEnvelopeReviewFeedbackPreservesRejectedReview(t *testing.T) {
	result := RecordExecutionEnvelopeReviewFeedback(ExecutionEnvelopeReviewFeedbackInput{
		DispositionStatus:      "review-required",
		ReviewRequired:         true,
		MismatchStage:          "status",
		EvidenceDigest:         "observed",
		ReviewStatus:           "rejected",
		ReviewerEvidenceDigest: "reviewer",
		NonAuthorizing:         true,
	})
	if result.Status != "review-rejected" || result.FeedbackDigest == "" || !result.NonAuthorizing {
		t.Fatalf("rejected review was not preserved: %#v", result)
	}
}

func TestRecordExecutionEnvelopeReviewFeedbackHoldsMissingOutcome(t *testing.T) {
	result := RecordExecutionEnvelopeReviewFeedback(ExecutionEnvelopeReviewFeedbackInput{
		DispositionStatus: "review-required",
		ReviewRequired:    true,
		MismatchStage:     "status",
		EvidenceDigest:    "observed",
		NonAuthorizing:    true,
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "review-outcome" || !result.NonAuthorizing {
		t.Fatalf("missing review outcome escaped UNKNOWN: %#v", result)
	}
}

func TestRecordExecutionEnvelopeReviewFeedbackKeepsObservationOnly(t *testing.T) {
	result := RecordExecutionEnvelopeReviewFeedback(ExecutionEnvelopeReviewFeedbackInput{
		DispositionStatus: "observation-only",
		EvidenceDigest:    "observed",
		NonAuthorizing:    true,
	})
	if result.Status != "observation-only" || result.FeedbackDigest != "observed" || result.ReviewRequired || !result.NonAuthorizing {
		t.Fatalf("observation-only feedback changed disposition: %#v", result)
	}
}
