package decision

import "testing"

func makeImprovementReviewEvidenceGateInput(t *testing.T, decision ImprovementReviewDecision) ImprovementReviewEvidenceGateInput {
	t.Helper()
	review := ImprovementReviewReceipt{
		CandidateDigest:      "candidate-digest",
		ReviewerReference:    "reviewer://one",
		ReviewEvidenceDigest: "review-evidence",
		Decision:             decision,
		ExecutionGranted:     false,
	}
	digest, err := Digest(review)
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	review.ReviewDigest = digest
	return ImprovementReviewEvidenceGateInput{
		CandidateAdmissionStatus: "proposed",
		CandidateDigest:          "candidate-digest",
		ReviewCandidateDigest:    review.CandidateDigest,
		ReviewerReference:        review.ReviewerReference,
		ReviewEvidenceDigest:     review.ReviewEvidenceDigest,
		ReviewDecision:           review.Decision,
		ExecutionGranted:         review.ExecutionGranted,
		ReviewDigest:             review.ReviewDigest,
		NonAuthorizing:           true,
	}
}

func TestGateImprovementReviewEvidence(t *testing.T) {
	output := GateImprovementReviewEvidence(makeImprovementReviewEvidenceGateInput(t, ReviewPassed))
	if output.Status != "accepted-for-observation" || !output.ObservationOnly ||
		output.ExecutionGranted || output.CandidateDigest != "candidate-digest" ||
		output.ReviewDigest == "" || !output.NonAuthorizing {
		t.Fatalf("unexpected accepted gate: %#v", output)
	}

	output = GateImprovementReviewEvidence(makeImprovementReviewEvidenceGateInput(t, ReviewFailed))
	if output.Status != "rejected" || output.ObservationOnly || output.ExecutionGranted {
		t.Fatalf("unexpected rejected gate: %#v", output)
	}

	output = GateImprovementReviewEvidence(makeImprovementReviewEvidenceGateInput(t, ReviewRequired))
	if output.Status != "review" || output.MissingStage != "review-required" {
		t.Fatalf("unexpected required gate: %#v", output)
	}

	input := makeImprovementReviewEvidenceGateInput(t, ReviewPassed)
	input.ReviewCandidateDigest = "other-candidate"
	output = GateImprovementReviewEvidence(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-review-binding" {
		t.Fatalf("unexpected binding gate: %#v", output)
	}

	input = makeImprovementReviewEvidenceGateInput(t, ReviewPassed)
	input.ReviewDigest = "tampered"
	output = GateImprovementReviewEvidence(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "review-receipt" {
		t.Fatalf("unexpected tampered gate: %#v", output)
	}

	input = makeImprovementReviewEvidenceGateInput(t, ReviewPassed)
	input.ExecutionGranted = true
	output = GateImprovementReviewEvidence(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing ||
		output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization gate: %#v", output)
	}
}