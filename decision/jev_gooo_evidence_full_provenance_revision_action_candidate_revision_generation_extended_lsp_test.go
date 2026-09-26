package decision

import "testing"

func reviewGeneratedExtendedCandidateRevisionForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding {
	t.Helper()
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
		Guard: admittedGeneratedExtendedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "generated-extended-candidate-lsp-counterexample",
			FirstMismatch:  "generated-extended-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPClear(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
			Guard: admittedGeneratedExtendedCandidateRevisionForVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "generated-extended-candidate-lsp-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "generated-extended-candidate-lsp-clear",
		NonAuthorizing:      true,
	})
	if output.Status != "clear" || output.Code != "jev-candidate-verified" ||
		output.Publishable || output.MissingStageIndex != -1 ||
		output.SourceCandidateID == "" || output.GuardEvidenceDigest == "" {
		t.Fatalf("output = %#v, want clear extended candidate LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPCounterexample(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput{
		Verification:        reviewGeneratedExtendedCandidateRevisionForLSP(t),
		MissingStageIndex:   21,
		EvidencePrefixDigest: "generated-extended-candidate-lsp-counterexample",
		NonAuthorizing:      true,
	})
	if output.Status != "publishable" || output.Code != "jev-reverse-counterexample" ||
		!output.Publishable || output.Severity != "error" || output.MissingStageIndex != 21 ||
		output.ReverseEvidenceDigest == "" {
		t.Fatalf("output = %#v, want counterexample extended candidate LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPRequiresEvidencePrefix(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput{
		Verification:        reviewGeneratedExtendedCandidateRevisionForLSP(t),
		MissingStageIndex:   21,
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lsp-diagnostic-evidence" {
		t.Fatalf("output = %#v, want lsp-diagnostic-evidence UNKNOWN", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPRejectsTampering(t *testing.T) {
	verification := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput{
		Guard: admittedGeneratedExtendedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "generated-extended-candidate-lsp-tampered",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	verification.ReverseEvidenceDigest = "tampered"
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput{
		Verification:         verification,
		EvidencePrefixDigest: "generated-extended-candidate-lsp-tampered",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "action-candidate-generation-extended-verification" {
		t.Fatalf("output = %#v, want verification validation UNKNOWN", output)
	}
}
