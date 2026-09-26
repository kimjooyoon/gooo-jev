package decision

import "testing"

func verifiedActionDerivedCandidateForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding {
	t.Helper()
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationInput{
		Guard: validActionDerivedCandidateGuard(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "reproduced",
			EvidenceDigest: "action-derived-candidate-reverse-evidence",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func reviewActionDerivedCandidateForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding {
	t.Helper()
	return VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationInput{
		Guard: validActionDerivedCandidateGuard(t),
		ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
			Status:         "counterexample",
			EvidenceDigest: "action-derived-candidate-counterexample",
			FirstMismatch:  "generated-output",
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPClear(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verifiedActionDerivedCandidateForLSP(t),
		EvidencePrefixDigest: "candidate-prefix-clear",
		NonAuthorizing:      true,
	})
	if output.Status != "clear" || output.Code != "jev-candidate-verified" ||
		output.Publishable || output.MissingStageIndex != -1 || output.CandidateEvidenceDigest == "" {
		t.Fatalf("output = %#v, want clear candidate LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPCounterexample(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        reviewActionDerivedCandidateForLSP(t),
		MissingStageIndex:   11,
		EvidencePrefixDigest: "candidate-prefix-counterexample",
		NonAuthorizing:      true,
	})
	if output.Status != "publishable" || output.Code != "jev-reverse-counterexample" ||
		!output.Publishable || output.Severity != "error" || output.MissingStageIndex != 11 {
		t.Fatalf("output = %#v, want counterexample candidate LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPRequiresEvidencePrefix(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:   reviewActionDerivedCandidateForLSP(t),
		MissingStageIndex: 11,
		NonAuthorizing:  true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "lsp-diagnostic-evidence" {
		t.Fatalf("output = %#v, want lsp-diagnostic-evidence UNKNOWN", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPRejectsTamperedVerification(t *testing.T) {
	verification := verifiedActionDerivedCandidateForLSP(t)
	verification.EvidenceDigest = "tampered"
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verification,
		EvidencePrefixDigest: "candidate-prefix-tampered",
		NonAuthorizing:      true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-verification" {
		t.Fatalf("output = %#v, want candidate verification UNKNOWN", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPRejectsAuthorizationBoundary(t *testing.T) {
	verification := verifiedActionDerivedCandidateForLSP(t)
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verification,
		EvidencePrefixDigest: "candidate-prefix-authorization",
		NonAuthorizing:      false,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "authorization-boundary" {
		t.Fatalf("output = %#v, want authorization-boundary UNKNOWN", output)
	}
}