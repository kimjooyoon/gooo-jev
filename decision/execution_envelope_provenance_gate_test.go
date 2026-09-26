package decision

import "testing"

func TestEvaluateExecutionEnvelopeProvenanceGatePreservesFirstMissingStage(t *testing.T) {
	result := EvaluateExecutionEnvelopeProvenanceGate(ExecutionEnvelopeProvenanceGateInput{
		NonAuthorizing: true,
		IRDigest:       "ir",
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "declaration" || !result.NonAuthorizing {
		t.Fatalf("partial provenance escaped fail-closed state: %#v", result)
	}
}

func TestEvaluateExecutionEnvelopeProvenanceGateBindsAllStages(t *testing.T) {
	result := EvaluateExecutionEnvelopeProvenanceGate(ExecutionEnvelopeProvenanceGateInput{
		DeclarationDigest:        "declaration",
		IRDigest:                 "ir",
		GenerationDigest:         "generation",
		ReverseObservationDigest: "reverse",
		MetricDigest:             "metric",
		NonAuthorizing:           true,
	})
	if result.Status != "ready" || result.MissingStage != "" || result.EvidenceDigest == "" || !result.NonAuthorizing {
		t.Fatalf("complete provenance was not bound: %#v", result)
	}
}

func TestEvaluateExecutionEnvelopeProvenanceGateRejectsAuthorizationClaim(t *testing.T) {
	result := EvaluateExecutionEnvelopeProvenanceGate(ExecutionEnvelopeProvenanceGateInput{
		DeclarationDigest: "declaration",
		NonAuthorizing:    false,
	})
	if result.Status != "UNKNOWN" || result.MissingStage != "authorization-boundary" || result.NonAuthorizing {
		t.Fatalf("authorization claim was not rejected: %#v", result)
	}
}
