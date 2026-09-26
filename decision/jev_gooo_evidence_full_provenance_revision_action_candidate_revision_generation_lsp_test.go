package decision

import "testing"

func reviewGeneratedCandidateRevisionForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding {
	t.Helper()
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput{
		Guard: admittedGeneratedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "generated-candidate-lsp-counterexample",
			FirstMismatch:  "generated-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPClear(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput{
			Guard: admittedGeneratedCandidateRevisionForVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "generated-candidate-lsp-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "generated-candidate-lsp-clear",
		NonAuthorizing:      true,
	})
	if output.Status != "clear" || output.Code != "jev-candidate-verified" ||
		output.Publishable || output.MissingStageIndex != -1 ||
		output.SourceCandidateID == "" {
		t.Fatalf("output = %#v, want clear generated candidate LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPCounterexample(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPInput{
		Verification:        reviewGeneratedCandidateRevisionForLSP(t),
		MissingStageIndex:   19,
		EvidencePrefixDigest: "generated-candidate-lsp-counterexample",
		NonAuthorizing:      true,
	})
	if output.Status != "publishable" || output.Code != "jev-reverse-counterexample" ||
		!output.Publishable || output.Severity != "error" || output.MissingStageIndex != 19 {
		t.Fatalf("output = %#v, want counterexample generated candidate LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPRequiresEvidencePrefix(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPInput{
		Verification:        reviewGeneratedCandidateRevisionForLSP(t),
		MissingStageIndex:   19,
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lsp-diagnostic-evidence" {
		t.Fatalf("output = %#v, want lsp-diagnostic-evidence UNKNOWN", output)
	}
}