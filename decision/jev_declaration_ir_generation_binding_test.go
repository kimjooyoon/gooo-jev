package decision

import "testing"

func validDeclarationIRGenerationBinding() ExecutionEnvelopeDeclarationIRGenerationBinding {
	return BindExecutionEnvelopeDeclarationIRGeneration(
		"gooo://gooo-jev/declaration/example",
		"gooo://gooo-jev/contract/example",
		"declaration-digest",
		"ir-digest",
		"generation-digest",
	)
}

func TestBindExecutionEnvelopeDeclarationIRGenerationCreatesValidatedBinding(t *testing.T) {
	binding := validDeclarationIRGenerationBinding()
	if binding.Status != "bound" || binding.BindingDigest == "" || binding.MissingStage != "" {
		t.Fatalf("unexpected binding: %+v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestEvaluateExecutionEnvelopeProvenanceChainBindingReachesReady(t *testing.T) {
	output := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	if output.Status != "ready" || output.MissingStage != "" || output.EvidenceDigest == "" {
		t.Fatalf("unexpected ready chain: %+v", output)
	}
	if !output.NonExecuting || !output.NonAuthorizing || output.DeclarationID == "" || output.IRDigest == "" {
		t.Fatalf("chain lost identity or safety boundary: %+v", output)
	}
}

func TestEvaluateExecutionEnvelopeProvenanceChainBindingPreservesMissingAndTamperedStages(t *testing.T) {
	output := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "reverse_observation" {
		t.Fatalf("unexpected missing reverse stage: %+v", output)
	}

	binding := validDeclarationIRGenerationBinding()
	binding.IRDigest = "tampered-ir-digest"
	output = EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  binding,
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "declaration-ir-generation" {
		t.Fatalf("unexpected tampered binding: %+v", output)
	}
}

func TestBindExecutionEnvelopeDeclarationIRGenerationFailsClosed(t *testing.T) {
	binding := BindExecutionEnvelopeDeclarationIRGeneration(
		"gooo://gooo-jev/declaration/example", "", "declaration-digest", "ir-digest", "generation-digest",
	)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "contract-id" || !binding.NonExecuting {
		t.Fatalf("unexpected incomplete binding: %+v", binding)
	}

	output := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}

func TestExecutionEnvelopeProvenanceChainBindingValidateReplaysEvidence(t *testing.T) {
	ready := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	if err := ready.Validate(); err != nil {
		t.Fatalf("ready chain should validate: %v", err)
	}

	tamperedEvidence := ready
	tamperedEvidence.EvidenceDigest = "tampered-evidence"
	if err := tamperedEvidence.Validate(); err == nil {
		t.Fatal("tampered evidence must fail validation")
	}

	unknown := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown chain should preserve a valid prefix: %v", err)
	}

	tamperedPrefix := unknown
	tamperedPrefix.BindingDigest = "tampered-binding"
	if err := tamperedPrefix.Validate(); err == nil {
		t.Fatal("tampered prefix must fail validation")
	}
}

func TestValidateExecutionEnvelopeProvenanceChainBindingFailsClosed(t *testing.T) {
	ready := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	tampered := ready
	tampered.EvidenceDigest = "tampered-evidence"
	unknown := ValidateExecutionEnvelopeProvenanceChainBinding(tampered)
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "provenance-validation" ||
		unknown.EvidenceDigest != "" || !unknown.NonExecuting || !unknown.NonAuthorizing {
		t.Fatalf("validator retained unsafe evidence: %+v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("validator output should be self-validating: %v", err)
	}

	unsafe := ready
	unsafe.NonAuthorizing = false
	unknown = ValidateExecutionEnvelopeProvenanceChainBinding(unsafe)
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "provenance-validation" || unknown.NonAuthorizing {
		t.Fatalf("validator erased authorization violation: %+v", unknown)
	}
}

func TestMeasureExecutionEnvelopeProvenanceChainMetricReportsExactCoverage(t *testing.T) {
	ready := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	metric := MeasureExecutionEnvelopeProvenanceChainMetric(ready)
	if metric.Status != "complete" || metric.ObservedStageCount != 5 || metric.ExpectedStageCount != 5 ||
		metric.MissingStage != "" || metric.CompletenessDigest == "" || !metric.NonExecuting || !metric.NonAuthorizing {
		t.Fatalf("unexpected complete metric: %+v", metric)
	}

	partial := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		MetricDigest:             "metric-digest",
		NonAuthorizing:           true,
	})
	metric = MeasureExecutionEnvelopeProvenanceChainMetric(partial)
	if metric.Status != "UNKNOWN" || metric.ObservedStageCount != 4 || metric.ExpectedStageCount != 5 ||
		metric.MissingStage != "reverse_observation" || metric.CompletenessDigest == "" {
		t.Fatalf("unexpected partial metric: %+v", metric)
	}

	invalid := ready
	invalid.EvidenceDigest = "tampered-evidence"
	metric = MeasureExecutionEnvelopeProvenanceChainMetric(invalid)
	if metric.Status != "UNKNOWN" || metric.ObservedStageCount != 0 || metric.MissingStage != "provenance-validation" || metric.CompletenessDigest != "" {
		t.Fatalf("unexpected invalid metric: %+v", metric)
	}
}

func TestExecutionEnvelopeProvenanceChainMetricValidateReplaysDigest(t *testing.T) {
	ready := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  validDeclarationIRGenerationBinding(),
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:              "metric-digest",
		NonAuthorizing:            true,
	})
	metric := MeasureExecutionEnvelopeProvenanceChainMetric(ready)
	if err := metric.Validate(); err != nil {
		t.Fatalf("complete metric should validate: %v", err)
	}

	partial := MeasureExecutionEnvelopeProvenanceChainMetric(EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		MetricDigest:             "metric-digest",
		NonAuthorizing:           true,
	}))
	if err := partial.Validate(); err != nil {
		t.Fatalf("partial metric should validate its UNKNOWN evidence: %v", err)
	}

	tampered := metric
	tampered.ObservedStageCount = 4
	if err := tampered.Validate(); err == nil {
		t.Fatal("tampered stage count must fail metric validation")
	}
	tampered = metric
	tampered.EvidenceDigest = "tampered-evidence"
	if err := tampered.Validate(); err == nil {
		t.Fatal("tampered evidence digest must fail metric validation")
	}
}
