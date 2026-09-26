package decision

import "testing"

func reproducedJEVActionReverseObservation() ExecutionEnvelopeReverseObservationOutput {
	return ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedMissingStage:   "",
		ExpectedMissingStage:   "",
		ObservedEvidenceDigest: "reverse-evidence",
		ExpectedEvidenceDigest: "reverse-evidence",
		NonAuthorizing:         true,
	})
}

func TestVerifyJEVActionCandidateRequiresReproducedReverseObservation(t *testing.T) {
	output := VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard:              GuardJEVActionCandidate(validJEVActionCandidateInput()),
		ReverseObservation: reproducedJEVActionReverseObservation(),
		NonAuthorizing:     true,
	})
	if output.Status != "verified" || output.MissingStage != "" {
		t.Fatalf("unexpected verified output: %+v", output)
	}
	if output.EvidenceDigest == "" || !output.NonExecuting || !output.NonAuthorizing {
		t.Fatalf("verification lost evidence or safety boundary: %+v", output)
	}
}

func TestVerifyJEVActionCandidatePreservesCounterexamplesAndGuardReview(t *testing.T) {
	reverse := ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "observed",
		ExpectedEvidenceDigest: "expected",
		NonAuthorizing:         true,
	})
	output := VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard:              GuardJEVActionCandidate(validJEVActionCandidateInput()),
		ReverseObservation: reverse,
		NonAuthorizing:     true,
	})
	if output.Status != "review" || output.MissingStage != "reverse-observation" || output.FirstMismatch != "evidence_digest" {
		t.Fatalf("unexpected counterexample output: %+v", output)
	}

	input := validJEVActionCandidateInput()
	input.Now = input.ObservationAt.Add(2 * 60 * 1000000000)
	output = VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard:              GuardJEVActionCandidate(input),
		ReverseObservation: reproducedJEVActionReverseObservation(),
		NonAuthorizing:     true,
	})
	if output.Status != "review" || output.MissingStage != "candidate-guard" {
		t.Fatalf("unexpected guard review output: %+v", output)
	}
}

func TestVerifyJEVActionCandidateFailsClosedOnUnknownOrAuthorization(t *testing.T) {
	input := validJEVActionCandidateInput()
	output := VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard: JEVActionCandidateGuard{
			Status: "admitted", CandidateID: "candidate-1", EvidenceDigest: "guard-evidence",
			NonExecuting: true, NonAuthorizing: true,
		},
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status: "UNKNOWN", NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "reverse-evidence" {
		t.Fatalf("unexpected unknown output: %+v", output)
	}

	input.NonAuthorizing = false
	output = VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard:              GuardJEVActionCandidate(input),
		ReverseObservation: reproducedJEVActionReverseObservation(),
		NonAuthorizing:     false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
