package decision

import "testing"

func admittedGeneratedCandidateRevisionForVerification(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationGuardBinding {
	t.Helper()
	return GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(validGeneratedCandidateRevisionGuardInput(t))
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput{
		Guard: admittedGeneratedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "generated-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "verified" || output.VerificationStatus != "verified" ||
		output.MissingStage != "" || output.SourceCandidateID == "" {
		t.Fatalf("output = %#v, want verified generated candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationKeepsCounterexampleReviewable(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput{
		Guard: admittedGeneratedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "generated-candidate-counterexample",
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

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationRejectsStaleGuard(t *testing.T) {
	input := validGeneratedCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 *  time.Minute)
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput{
		Guard: GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(input),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "generated-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale UNKNOWN", output)
	}
}

func TestVerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationRequiresReverseEvidence(t *testing.T) {
	output := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput{
		Guard: admittedGeneratedCandidateRevisionForVerification(t),
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