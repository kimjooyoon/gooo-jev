package decision

import "testing"

func reviewExtendedLineageCandidateRevisionForRevisionLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding {
	t.Helper()
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput{
		Guard: admittedExtendedLineageCandidateRevisionForRevisionVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "extended-lineage-candidate-revision-lsp-counterexample",
			FirstMismatch:  "extended-lineage-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPClear(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput{
			Guard: admittedExtendedLineageCandidateRevisionForRevisionVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "extended-lineage-candidate-revision-lsp-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "extended-lineage-candidate-revision-lsp-clear",
		MissingStageIndex:   -1,
		NonAuthorizing:      true,
	})
	if output.Status != "clear" || output.Code != "jev-candidate-verified" ||
		output.Publishable || output.MissingStageIndex != -1 ||
		output.SourceCandidateID == "" || output.ParentCandidateDigest == "" ||
		output.CandidateDigest == "" {
		t.Fatalf("output = %#v, want clear extended lineage candidate revision LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPCounterexample(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPInput{
		Verification:        reviewExtendedLineageCandidateRevisionForRevisionLSP(t),
		MissingStageIndex:   28,
		EvidencePrefixDigest: "extended-lineage-candidate-revision-lsp-counterexample",
		NonAuthorizing:      true,
	})
	if output.Status != "publishable" || output.Code != "jev-reverse-counterexample" ||
		!output.Publishable || output.Severity != "error" || output.MissingStageIndex != 28 ||
		output.ReverseEvidenceDigest == "" {
		t.Fatalf("output = %#v, want counterexample extended lineage candidate revision LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPRequiresEvidencePrefix(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPInput{
		Verification:        reviewExtendedLineageCandidateRevisionForRevisionLSP(t),
		MissingStageIndex:   28,
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lsp-diagnostic-evidence" {
		t.Fatalf("output = %#v, want lsp-diagnostic-evidence UNKNOWN", output)
	}
}