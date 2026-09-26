package decision

import "testing"

func verifiedRevisionActionForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding {
	t.Helper()
	guard := validRevisionActionVerificationGuard(t)
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationInput{
		Guard: guard,
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "reverse-observation-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func reviewRevisionActionForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding {
	t.Helper()
	guard := validRevisionActionVerificationGuard(t)
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationInput{
		Guard: guard,
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "counterexample-evidence",
			FirstMismatch:  "generated-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPClear(t *testing.T) {
	verification := verifiedRevisionActionForLSP(t)
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verification,
		EvidencePrefixDigest: "prefix-clear",
		NonAuthorizing:      true,
	})
	if output.Status != "clear" || output.Code != "jev-candidate-verified" ||
		output.Publishable || output.MissingStageIndex != -1 {
		t.Fatalf("output = %#v, want clear LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPCounterexample(t *testing.T) {
	verification := reviewRevisionActionForLSP(t)
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verification,
		MissingStageIndex:   7,
		EvidencePrefixDigest: "prefix-counterexample",
		NonAuthorizing:      true,
	})
	if output.Status != "publishable" || output.Code != "jev-reverse-counterexample" ||
		!output.Publishable || output.Severity != "error" || output.MissingStageIndex != 7 {
		t.Fatalf("output = %#v, want counterexample LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPRequiresEvidencePrefix(t *testing.T) {
	verification := reviewRevisionActionForLSP(t)
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:   verification,
		MissingStageIndex: 7,
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lsp-diagnostic-evidence" {
		t.Fatalf("output = %#v, want lsp-diagnostic-evidence UNKNOWN", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPFailsClosedForTamperedVerification(t *testing.T) {
	verification := verifiedRevisionActionForLSP(t)
	verification.EvidenceDigest = "tampered"
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verification,
		EvidencePrefixDigest: "prefix-tampered",
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-verification" {
		t.Fatalf("output = %#v, want revision-action-verification UNKNOWN", output)
	}
}