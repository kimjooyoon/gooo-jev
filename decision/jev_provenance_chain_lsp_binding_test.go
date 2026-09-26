package decision

import "testing"

func validProvenanceChainForLSP() ExecutionEnvelopeProvenanceChainBinding {
	return EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
}

func TestProjectExecutionEnvelopeProvenanceChainLSPProjectsReadyAndMissingStages(t *testing.T) {
	ready := validProvenanceChainForLSP()
	output := ProjectExecutionEnvelopeProvenanceChainLSP(ExecutionEnvelopeProvenanceChainLSPBindingInput{
		Chain:                ready,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "clear" || output.Publishable || output.Code != "provenance-complete" || output.ChainEvidenceDigest != ready.EvidenceDigest || output.ChainBindingDigest != ready.BindingDigest {
		t.Fatalf("unexpected ready projection: %+v", output)
	}

	missing := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	output = ProjectExecutionEnvelopeProvenanceChainLSP(ExecutionEnvelopeProvenanceChainLSPBindingInput{
		Chain:                missing,
		MissingStageIndex:    99,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "publishable" || !output.Publishable || output.Severity != "error" || output.MissingStage != "reverse_observation" || output.MissingStageIndex != 3 {
		t.Fatalf("unexpected missing-stage projection: %+v", output)
	}
}

func TestProjectExecutionEnvelopeProvenanceChainLSPKeepsIncompleteEvidenceUnknown(t *testing.T) {
	missing := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		NonAuthorizing:            true,
	})
	output := ProjectExecutionEnvelopeProvenanceChainLSP(ExecutionEnvelopeProvenanceChainLSPBindingInput{
		Chain:          missing,
		MissingStageIndex: -1,
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.Publishable || output.Code != "lsp-diagnostic-evidence" {
		t.Fatalf("unexpected incomplete projection: %+v", output)
	}

	missing.MissingStage = ""
	output = ProjectExecutionEnvelopeProvenanceChainLSP(ExecutionEnvelopeProvenanceChainLSPBindingInput{
		Chain:                missing,
		MissingStageIndex:    1,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.Publishable || output.Code != "chain-integrity" {
		t.Fatalf("unexpected stage-less projection: %+v", output)
	}
}

func TestProjectExecutionEnvelopeProvenanceChainLSPFailsClosedOnAuthorization(t *testing.T) {
	output := ProjectExecutionEnvelopeProvenanceChainLSP(ExecutionEnvelopeProvenanceChainLSPBindingInput{
		Chain: validProvenanceChainForLSP(),
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.Code != "authorization-boundary" || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization projection: %+v", output)
	}
}

func TestProjectExecutionEnvelopeProvenanceChainLSPRejectsTamperedChain(t *testing.T) {
	chain := validProvenanceChainForLSP()
	chain.BindingDigest = "tampered-binding"
	output := ProjectExecutionEnvelopeProvenanceChainLSP(ExecutionEnvelopeProvenanceChainLSPBindingInput{
		Chain:                chain,
		MissingStageIndex:    -1,
		EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.Publishable || output.Code != "chain-integrity" || output.ChainBindingDigest != "" {
		t.Fatalf("tampered chain was projected: %+v", output)
	}
}

func TestProvenanceChainMissingStageIndexDerivesCanonicalOrder(t *testing.T) {
	cases := map[string]int{
		"declaration":         0,
		"ir":                  1,
		"generation":          2,
		"reverse_observation": 3,
		"metric":              4,
	}
	for stage, want := range cases {
		got, ok := ProvenanceChainMissingStageIndex(stage)
		if !ok || got != want {
			t.Fatalf("stage %q index = (%d, %v), want (%d, true)", stage, got, ok, want)
		}
	}
	if got, ok := ProvenanceChainMissingStageIndex("provenance-validation"); ok || got != -1 {
		t.Fatalf("unsupported stage index = (%d, %v), want (-1, false)", got, ok)
	}
}
