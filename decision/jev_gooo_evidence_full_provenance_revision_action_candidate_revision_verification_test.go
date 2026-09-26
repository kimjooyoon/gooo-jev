package decision

import (
	"testing"
	"time"
)

func verifiedActionDerivedCandidateRevisionForVerification(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardBinding {
	t.Helper()
	return GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(validActionDerivedCandidateRevisionGuardInput(t))
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: verifiedActionDerivedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "candidate-revision-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "verified" || output.VerificationStatus != "verified" ||
		output.MissingStage != "" || output.SourceCandidateID == "" {
		t.Fatalf("output = %#v, want verified candidate revision", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionKeepsCounterexampleReviewable(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: verifiedActionDerivedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "candidate-revision-counterexample",
			FirstMismatch:  "generated-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "review" || output.MissingStage != "reverse-observation" ||
		output.FirstMismatch != "generated-output" {
		t.Fatalf("output = %#v, want reverse-observation review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRejectsStaleGuard(t *testing.T) {
	input := validActionDerivedCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	guard := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(input)
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: guard,
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "candidate-revision-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale UNKNOWN", output)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRequiresReverseEvidence(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: verifiedActionDerivedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "reverse-evidence" {
		t.Fatalf("output = %#v, want reverse-evidence UNKNOWN", output)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRejectsTamperedGuard(t *testing.T) {
	guard := verifiedActionDerivedCandidateRevisionForVerification(t)
	guard.GuardEvidenceDigest = "tampered"
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: guard,
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "candidate-revision-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "action-candidate-revision-guard-validation" {
		t.Fatalf("output = %#v, want guard validation UNKNOWN", output)
	}
}