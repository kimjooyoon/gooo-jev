package decision

import (
	"testing"
	"time"
)

func admittedExtendedLineageCandidateRevisionForVerification(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding {
	t.Helper()
	return GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(validExtendedLineageCandidateRevisionGuardInput(t))
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput{
		Guard: admittedExtendedLineageCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "extended-lineage-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "verified" || output.VerificationStatus != "verified" ||
		output.MissingStage != "" || output.SourceCandidateID == "" {
		t.Fatalf("output = %#v, want verified extended lineage candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageKeepsCounterexampleReviewable(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput{
		Guard: admittedExtendedLineageCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "extended-lineage-candidate-counterexample",
			FirstMismatch:  "extended-lineage-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "review" || output.MissingStage != "reverse-observation" ||
		output.FirstMismatch != "extended-lineage-output" {
		t.Fatalf("output = %#v, want reverse-observation review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRejectsStaleGuard(t *testing.T) {
	input := validExtendedLineageCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput{
		Guard: GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "extended-lineage-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale UNKNOWN", output)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRequiresReverseEvidence(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput{
		Guard: admittedExtendedLineageCandidateRevisionForVerification(t),
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
