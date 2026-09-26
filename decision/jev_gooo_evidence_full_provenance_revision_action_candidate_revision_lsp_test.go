package decision

import "testing"

func reviewActionDerivedCandidateRevisionForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationBinding {
	t.Helper()
	guard := verifiedActionDerivedCandidateRevisionForVerification(t)
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: guard,
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "candidate-revision-lsp-counterexample",
			FirstMismatch:  "generated-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPClear(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput{
		Verification:        VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
			Guard: verifiedActionDerivedCandidateRevisionForVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "candidate-revision-lsp-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "candidate-revision-lsp-clear",
		NonAuthorizing:      true,
	})
	if output.Status != "clear" || output.Code != "jev-candidate-verified" ||
		output.Publishable || output.MissingStageIndex != -1 ||
		output.SourceCandidateID == "" {
		t.Fatalf("output = %#v, want clear candidate revision LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPCounterexample(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput{
		Verification:        reviewActionDerivedCandidateRevisionForLSP(t),
		MissingStageIndex:   17,
		EvidencePrefixDigest: "candidate-revision-lsp-counterexample",
		NonAuthorizing:      true,
	})
	if output.Status != "publishable" || output.Code != "jev-reverse-counterexample" ||
		!output.Publishable || output.Severity != "error" || output.MissingStageIndex != 17 {
		t.Fatalf("output = %#v, want counterexample candidate revision LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPRequiresEvidencePrefix(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput{
		Verification:        reviewActionDerivedCandidateRevisionForLSP(t),
		MissingStageIndex:   17,
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lsp-diagnostic-evidence" {
		t.Fatalf("output = %#v, want lsp-diagnostic-evidence UNKNOWN", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPRejectsTampering(t *testing.T) {
	verification := VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
		Guard: verifiedActionDerivedCandidateRevisionForVerification(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "candidate-revision-lsp-tampered",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	verification.EvidenceDigest = "tampered"
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput{
		Verification:        verification,
		EvidencePrefixDigest: "candidate-revision-lsp-tampered",
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "action-candidate-revision-verification" {
		t.Fatalf("output = %#v, want verification UNKNOWN", output)
	}
}