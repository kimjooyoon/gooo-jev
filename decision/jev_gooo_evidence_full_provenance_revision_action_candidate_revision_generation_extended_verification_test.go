package decision

import (
	"testing"
	"time"
)

func admittedGeneratedExtendedCandidateRevisionForVerification(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding {
	t.Helper()
	return GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(validGeneratedExtendedCandidateRevisionGuardInput(t))
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
		Guard: admittedGeneratedExtendedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "generated-extended-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "verified" || output.VerificationStatus != "verified" ||
		output.MissingStage != "" || output.SourceCandidateID == "" ||
		output.VerificationEvidenceDigest == "" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want verified extended generated candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedKeepsCounterexampleReviewable(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
		Guard: admittedGeneratedExtendedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "generated-extended-candidate-counterexample",
			FirstMismatch:  "generated-extended-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "review" || output.MissingStage != "reverse-observation" ||
		output.FirstMismatch != "generated-extended-output" {
		t.Fatalf("output = %#v, want reverse-observation review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRejectsStaleGuard(t *testing.T) {
	input := validGeneratedExtendedCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
		Guard: GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(input),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "generated-extended-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale UNKNOWN", output)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRequiresReverseEvidence(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
		Guard: admittedGeneratedExtendedCandidateRevisionForVerification(t),
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
