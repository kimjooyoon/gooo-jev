package decision

import "testing"

func TestBindExecutionEnvelopeFullProvenanceFromSource(t *testing.T) {
	input := ExecutionEnvelopeFullProvenanceSourceInput{
		DeclarationID:          "decl-1",
		ContractID:             "contract-1",
		DeclarationSource:      "entity Example id gooo://example",
		IRSource:               "ir: Example declaration",
		GenerationSource:       "generation: Example artifact",
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "evidence-1",
		ExpectedEvidenceDigest: "evidence-1",
		ReverseObservationSource: "reverse_observation: ready evidence-1",
		MetricSource:           "metric: all five provenance stages observed",
		NonAuthorizing:         true,
	}
	complete := BindExecutionEnvelopeFullProvenanceFromSource(input)
	if complete.Status != "complete" || complete.MissingStage != "" {
		t.Fatalf("complete = %#v, want complete", complete)
	}
	if complete.DeclarationDigest == "" || complete.IRDigest == "" || complete.GenerationDigest == "" ||
		complete.BindingDigest == "" || complete.ReverseObservationDigest == "" ||
		complete.MetricDigest == "" || complete.EvidenceDigest == "" || complete.CompletenessDigest == "" {
		t.Fatal("complete source path must retain every evidence digest")
	}
	changedMetric := input
	changedMetric.MetricSource += "\n"
	changed := BindExecutionEnvelopeFullProvenanceFromSource(changedMetric)
	if changed.Status != "complete" || changed.MetricDigest == complete.MetricDigest || changed.EvidenceDigest == complete.EvidenceDigest {
		t.Fatal("changed metric source must change the complete evidence chain")
	}
	missingGeneration := input
	missingGeneration.GenerationSource = ""
	generationUnknown := BindExecutionEnvelopeFullProvenanceFromSource(missingGeneration)
	if generationUnknown.Status != "UNKNOWN" || generationUnknown.MissingStage != "generation-source" {
		t.Fatalf("missing generation = %#v, want UNKNOWN at generation-source", generationUnknown)
	}
	missingReverse := input
	missingReverse.ReverseObservationSource = ""
	reverseUnknown := BindExecutionEnvelopeFullProvenanceFromSource(missingReverse)
	if reverseUnknown.Status != "UNKNOWN" || reverseUnknown.MissingStage != "reverse-observation-source" {
		t.Fatalf("missing reverse source = %#v, want UNKNOWN at reverse-observation-source", reverseUnknown)
	}
	counterexample := input
	counterexample.ExpectedEvidenceDigest = "different-evidence"
	counterexampleResult := BindExecutionEnvelopeFullProvenanceFromSource(counterexample)
	if counterexampleResult.Status != "UNKNOWN" || counterexampleResult.MissingStage != "reverse-observation" {
		t.Fatalf("counterexample = %#v, want UNKNOWN at reverse-observation", counterexampleResult)
	}
	missingMetric := input
	missingMetric.MetricSource = ""
	metricUnknown := BindExecutionEnvelopeFullProvenanceFromSource(missingMetric)
	if metricUnknown.Status != "UNKNOWN" || metricUnknown.MissingStage != "metric-source" {
		t.Fatalf("missing metric = %#v, want UNKNOWN at metric-source", metricUnknown)
	}
}
