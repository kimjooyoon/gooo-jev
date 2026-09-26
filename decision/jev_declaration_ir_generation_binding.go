package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeDeclarationSourceDigestInput identifies the exact .gooo
// declaration text whose provenance should be recorded without execution.
type ExecutionEnvelopeDeclarationSourceDigestInput struct {
	DeclarationID  string
	ContractID     string
	SourceText     string
	NonAuthorizing bool
}

// ExecutionEnvelopeDeclarationSourceDigest is the content-addressed origin
// of a declaration before IR and generation are attached.
type ExecutionEnvelopeDeclarationSourceDigest struct {
	Status            string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

// ComputeExecutionEnvelopeDeclarationSourceDigest derives declaration origin
// from the exact source text instead of accepting an opaque digest.
func ComputeExecutionEnvelopeDeclarationSourceDigest(input ExecutionEnvelopeDeclarationSourceDigestInput) ExecutionEnvelopeDeclarationSourceDigest {
	output := ExecutionEnvelopeDeclarationSourceDigest{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	stages := []struct {
		name  string
		value string
	}{
		{name: "declaration-id", value: input.DeclarationID},
		{name: "contract-id", value: input.ContractID},
		{name: "declaration-source", value: input.SourceText},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	digest, err := Digest(struct {
		DeclarationID string
		ContractID    string
		SourceText    string
	}{
		DeclarationID: input.DeclarationID,
		ContractID:    input.ContractID,
		SourceText:    input.SourceText,
	})
	if err != nil {
		output.MissingStage = "declaration-source-digest"
		return output
	}
	output.Status = "derived"
	output.DeclarationID = input.DeclarationID
	output.ContractID = input.ContractID
	output.DeclarationDigest = digest
	return output
}

// ExecutionEnvelopeDeclarationIRGenerationBinding binds a .gooo declaration
// identity to the IR and generation artifacts derived from it.
type ExecutionEnvelopeDeclarationIRGenerationBinding struct {
	Status            string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	IRDigest          string
	GenerationDigest  string
	BindingDigest     string
	NonExecuting      bool
	NonAuthorizing    bool
	MissingStage      string
}

// BindExecutionEnvelopeDeclarationIRGeneration creates a content-addressed
// declaration-to-IR-to-generation boundary without executing the plan.
// ExecutionEnvelopeDeclarationIRGenerationSourceInput connects exact
// declaration text directly to the downstream IR and generation artifacts.
type ExecutionEnvelopeDeclarationIRGenerationSourceInput struct {
	DeclarationID    string
	ContractID       string
	SourceText       string
	IRDigest         string
	GenerationDigest string
	NonAuthorizing   bool
}

// BindExecutionEnvelopeDeclarationIRGenerationFromSource computes the
// declaration origin first, then reuses the existing content-addressed bind.
func BindExecutionEnvelopeDeclarationIRGenerationFromSource(input ExecutionEnvelopeDeclarationIRGenerationSourceInput) ExecutionEnvelopeDeclarationIRGenerationBinding {
	source := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  input.DeclarationID,
		ContractID:     input.ContractID,
		SourceText:     input.SourceText,
		NonAuthorizing: input.NonAuthorizing,
	})
	if source.Status != "derived" {
		return ExecutionEnvelopeDeclarationIRGenerationBinding{
			Status:         "UNKNOWN",
			MissingStage:   source.MissingStage,
			NonExecuting:   source.NonExecuting,
			NonAuthorizing: source.NonAuthorizing,
		}
	}
	return BindExecutionEnvelopeDeclarationIRGeneration(
		source.DeclarationID,
		source.ContractID,
		source.DeclarationDigest,
		input.IRDigest,
		input.GenerationDigest,
	)
}

func BindExecutionEnvelopeDeclarationIRGeneration(
	declarationID string,
	contractID string,
	declarationDigest string,
	irDigest string,
	generationDigest string,
) ExecutionEnvelopeDeclarationIRGenerationBinding {
	output := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	stages := []struct {
		name  string
		value string
	}{
		{name: "declaration-id", value: declarationID},
		{name: "contract-id", value: contractID},
		{name: "declaration", value: declarationDigest},
		{name: "ir", value: irDigest},
		{name: "generation", value: generationDigest},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	bindingDigest, err := Digest(struct {
		DeclarationID     string
		ContractID        string
		DeclarationDigest string
		IRDigest          string
		GenerationDigest  string
	}{
		DeclarationID:     declarationID,
		ContractID:        contractID,
		DeclarationDigest: declarationDigest,
		IRDigest:          irDigest,
		GenerationDigest:  generationDigest,
	})
	if err != nil {
		output.MissingStage = "binding-evidence"
		return output
	}
	output.Status = "bound"
	output.DeclarationID = declarationID
	output.ContractID = contractID
	output.DeclarationDigest = declarationDigest
	output.IRDigest = irDigest
	output.GenerationDigest = generationDigest
	output.BindingDigest = bindingDigest
	return output
}

func (binding ExecutionEnvelopeDeclarationIRGenerationBinding) Validate() error {
	if binding.Status != "bound" || !binding.NonExecuting || !binding.NonAuthorizing ||
		strings.TrimSpace(binding.DeclarationID) == "" ||
		strings.TrimSpace(binding.ContractID) == "" ||
		strings.TrimSpace(binding.DeclarationDigest) == "" ||
		strings.TrimSpace(binding.IRDigest) == "" ||
		strings.TrimSpace(binding.GenerationDigest) == "" ||
		strings.TrimSpace(binding.BindingDigest) == "" ||
		strings.TrimSpace(binding.MissingStage) != "" {
		return fmt.Errorf("declaration IR generation binding is incomplete")
	}
	expected, err := Digest(struct {
		DeclarationID     string
		ContractID        string
		DeclarationDigest string
		IRDigest          string
		GenerationDigest  string
	}{
		DeclarationID:     binding.DeclarationID,
		ContractID:        binding.ContractID,
		DeclarationDigest: binding.DeclarationDigest,
		IRDigest:          binding.IRDigest,
		GenerationDigest:  binding.GenerationDigest,
	})
	if err != nil {
		return err
	}
	if expected != binding.BindingDigest {
		return fmt.Errorf("declaration IR generation binding digest mismatch")
	}
	return nil
}

// ExecutionEnvelopeProvenanceChainBindingInput connects the declaration/IR/
// generation binding to reverse observation and metric evidence.
type ExecutionEnvelopeProvenanceChainBindingInput struct {
	DeclarationIRGeneration ExecutionEnvelopeDeclarationIRGenerationBinding
	ReverseObservationDigest string
	MetricDigest             string
	NonAuthorizing           bool
}

// ExecutionEnvelopeProvenanceChainBinding records the full evidence chain
// without making an authorization or execution claim.
type ExecutionEnvelopeProvenanceChainBinding struct {
	Status                  string
	MissingStage            string
	DeclarationID           string
	ContractID              string
	DeclarationDigest       string
	IRDigest                string
	GenerationDigest        string
	BindingDigest           string
	ReverseObservationDigest string
	MetricDigest            string
	EvidenceDigest          string
	NonExecuting             bool
	NonAuthorizing           bool
}

// EvaluateExecutionEnvelopeProvenanceChainBinding reuses the existing
// five-stage provenance gate after validating the declaration-to-generation
// binding.
func EvaluateExecutionEnvelopeProvenanceChainBinding(input ExecutionEnvelopeProvenanceChainBindingInput) ExecutionEnvelopeProvenanceChainBinding {
	output := ExecutionEnvelopeProvenanceChainBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if err := input.DeclarationIRGeneration.Validate(); err != nil {
		output.MissingStage = "declaration-ir-generation"
		return output
	}
	binding := input.DeclarationIRGeneration
	gate := EvaluateExecutionEnvelopeProvenanceGate(ExecutionEnvelopeProvenanceGateInput{
		DeclarationDigest:        binding.DeclarationDigest,
		IRDigest:                 binding.IRDigest,
		GenerationDigest:         binding.GenerationDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		MetricDigest:             input.MetricDigest,
		NonAuthorizing:           true,
	})
	output.Status = gate.Status
	output.MissingStage = gate.MissingStage
	output.DeclarationID = binding.DeclarationID
	output.ContractID = binding.ContractID
	output.DeclarationDigest = binding.DeclarationDigest
	output.IRDigest = binding.IRDigest
	output.GenerationDigest = binding.GenerationDigest
	output.BindingDigest = binding.BindingDigest
	output.ReverseObservationDigest = input.ReverseObservationDigest
	output.MetricDigest = input.MetricDigest
	output.EvidenceDigest = gate.EvidenceDigest
	return output
}

// Validate replays the full provenance gate so a stored result cannot be
// treated as ready after any declaration, binding, observation, metric, or
// evidence field has been changed.
func (binding ExecutionEnvelopeProvenanceChainBinding) Validate() error {
	if !binding.NonExecuting || !binding.NonAuthorizing {
		return fmt.Errorf("provenance chain binding crosses an execution or authorization boundary")
	}
	if strings.TrimSpace(binding.EvidenceDigest) != "" && binding.Status != "ready" {
		return fmt.Errorf("unknown provenance chain cannot carry evidence digest")
	}
	if binding.Status == "UNKNOWN" {
		if strings.TrimSpace(binding.MissingStage) == "" {
			return fmt.Errorf("unknown provenance chain is missing its first unresolved stage")
		}
		if binding.MissingStage == "authorization-boundary" || binding.MissingStage == "declaration-ir-generation" || binding.MissingStage == "provenance-validation" {
			if strings.TrimSpace(binding.DeclarationID) != "" || strings.TrimSpace(binding.ContractID) != "" ||
				strings.TrimSpace(binding.DeclarationDigest) != "" || strings.TrimSpace(binding.IRDigest) != "" ||
				strings.TrimSpace(binding.GenerationDigest) != "" || strings.TrimSpace(binding.BindingDigest) != "" ||
				strings.TrimSpace(binding.ReverseObservationDigest) != "" || strings.TrimSpace(binding.MetricDigest) != "" {
				return fmt.Errorf("pre-binding unknown provenance chain contains evidence fields")
			}
			return nil
		}
	} else if binding.Status != "ready" {
		return fmt.Errorf("unsupported provenance chain status %q", binding.Status)
	}

	declaration := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status:            "bound",
		DeclarationID:     binding.DeclarationID,
		ContractID:        binding.ContractID,
		DeclarationDigest: binding.DeclarationDigest,
		IRDigest:          binding.IRDigest,
		GenerationDigest:  binding.GenerationDigest,
		BindingDigest:     binding.BindingDigest,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if err := declaration.Validate(); err != nil {
		return fmt.Errorf("declaration-to-generation prefix is invalid: %w", err)
	}

	gate := EvaluateExecutionEnvelopeProvenanceGate(ExecutionEnvelopeProvenanceGateInput{
		DeclarationDigest:        binding.DeclarationDigest,
		IRDigest:                 binding.IRDigest,
		GenerationDigest:         binding.GenerationDigest,
		ReverseObservationDigest: binding.ReverseObservationDigest,
		MetricDigest:             binding.MetricDigest,
		NonAuthorizing:           true,
	})
	if gate.Status != binding.Status || gate.MissingStage != binding.MissingStage || gate.EvidenceDigest != binding.EvidenceDigest {
		return fmt.Errorf("provenance gate replay mismatch")
	}
	return nil
}

// ValidateExecutionEnvelopeProvenanceChainBinding converts unverifiable or
// unsafe stored evidence into an explicit UNKNOWN result without retaining
// stale digests. A valid result is returned unchanged.
func ValidateExecutionEnvelopeProvenanceChainBinding(binding ExecutionEnvelopeProvenanceChainBinding) ExecutionEnvelopeProvenanceChainBinding {
	if err := binding.Validate(); err == nil {
		return binding
	}
	return ExecutionEnvelopeProvenanceChainBinding{
		Status:         "UNKNOWN",
		MissingStage:   "provenance-validation",
		NonExecuting:   binding.NonExecuting,
		NonAuthorizing: binding.NonAuthorizing,
	}
}

// ExecutionEnvelopeProvenanceChainMetric reports exact evidence coverage for
// the five-stage declaration-to-metric chain without inferring readiness.
type ExecutionEnvelopeProvenanceChainMetric struct {
	Status              string
	ObservedStageCount  int
	ExpectedStageCount  int
	MissingStage        string
	CompletenessDigest  string
	EvidenceDigest      string
	NonExecuting        bool
	NonAuthorizing      bool
}

// MeasureExecutionEnvelopeProvenanceChainMetric counts only present stages;
// partial or unverifiable chains remain UNKNOWN.
func MeasureExecutionEnvelopeProvenanceChainMetric(binding ExecutionEnvelopeProvenanceChainBinding) ExecutionEnvelopeProvenanceChainMetric {
	output := ExecutionEnvelopeProvenanceChainMetric{
		Status:             "UNKNOWN",
		ExpectedStageCount: 5,
		NonExecuting:       binding.NonExecuting,
		NonAuthorizing:     binding.NonAuthorizing,
	}
	if !binding.NonExecuting || !binding.NonAuthorizing {
		output.MissingStage = "authorization-boundary"
		return output
	}
	if binding.Status == "UNKNOWN" && (binding.MissingStage == "declaration-ir-generation" || binding.MissingStage == "provenance-validation") {
		output.MissingStage = binding.MissingStage
		return output
	}
	if err := binding.Validate(); err != nil {
		output.MissingStage = "provenance-validation"
		return output
	}
	output.EvidenceDigest = binding.EvidenceDigest

	stages := []struct {
		name  string
		value string
	}{
		{name: "declaration", value: binding.DeclarationDigest},
		{name: "ir", value: binding.IRDigest},
		{name: "generation", value: binding.GenerationDigest},
		{name: "reverse_observation", value: binding.ReverseObservationDigest},
		{name: "metric", value: binding.MetricDigest},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) != "" {
			output.ObservedStageCount++
		} else if output.MissingStage == "" {
			output.MissingStage = stage.name
		}
	}
	if output.ObservedStageCount == output.ExpectedStageCount {
		output.Status = "complete"
		output.MissingStage = ""
	}
	digest, err := Digest(struct {
		Status             string
		ObservedStageCount int
		ExpectedStageCount int
		MissingStage       string
		EvidenceDigest     string
	}{
		Status:             output.Status,
		ObservedStageCount: output.ObservedStageCount,
		ExpectedStageCount: output.ExpectedStageCount,
		MissingStage:       output.MissingStage,
		EvidenceDigest:     binding.EvidenceDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "metric-evidence"
		return output
	}
	output.CompletenessDigest = digest
	return output
}

// Validate ensures the reported coverage count and digest can be replayed
// without access to mutable execution state.
func (metric ExecutionEnvelopeProvenanceChainMetric) Validate() error {
	if !metric.NonExecuting || !metric.NonAuthorizing ||
		metric.ExpectedStageCount != 5 || metric.ObservedStageCount < 0 ||
		metric.ObservedStageCount > metric.ExpectedStageCount ||
		strings.TrimSpace(metric.CompletenessDigest) == "" {
		return fmt.Errorf("provenance chain metric is incomplete")
	}
	if metric.Status == "complete" {
		if metric.ObservedStageCount != metric.ExpectedStageCount || strings.TrimSpace(metric.MissingStage) != "" {
			return fmt.Errorf("complete provenance chain metric has partial coverage")
		}
	} else if metric.Status == "UNKNOWN" {
		if metric.ObservedStageCount == metric.ExpectedStageCount || strings.TrimSpace(metric.MissingStage) == "" {
			return fmt.Errorf("unknown provenance chain metric has no unresolved stage")
		}
	} else {
		return fmt.Errorf("unsupported provenance chain metric status %q", metric.Status)
	}
	expected, err := Digest(struct {
		Status             string
		ObservedStageCount int
		ExpectedStageCount int
		MissingStage       string
		EvidenceDigest     string
	}{
		Status:             metric.Status,
		ObservedStageCount: metric.ObservedStageCount,
		ExpectedStageCount: metric.ExpectedStageCount,
		MissingStage:       metric.MissingStage,
		EvidenceDigest:     metric.EvidenceDigest,
	})
	if err != nil {
		return err
	}
	if expected != metric.CompletenessDigest {
		return fmt.Errorf("provenance chain metric digest mismatch")
	}
	return nil
}

// ExecutionEnvelopeProvenanceChainReverseObservationInput replaces an opaque
// reverse-observation string with the validated binding that produced it.
type ExecutionEnvelopeProvenanceChainReverseObservationInput struct {
	DeclarationIRGeneration ExecutionEnvelopeDeclarationIRGenerationBinding
	ReverseObservation      ExecutionEnvelopeReverseObservationBinding
	MetricDigest            string
	NonAuthorizing          bool
}

// EvaluateExecutionEnvelopeProvenanceChainBindingFromReverseObservation
// admits only reproduced reverse observations into the ready chain.
func EvaluateExecutionEnvelopeProvenanceChainBindingFromReverseObservation(input ExecutionEnvelopeProvenanceChainReverseObservationInput) ExecutionEnvelopeProvenanceChainBinding {
	output := ExecutionEnvelopeProvenanceChainBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if err := input.ReverseObservation.Validate(); err != nil {
		output.MissingStage = "reverse-observation-binding"
		return output
	}
	if input.ReverseObservation.Status != "reproduced" {
		output.MissingStage = "reverse-observation"
		return output
	}
	return EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  input.DeclarationIRGeneration,
		ReverseObservationDigest: input.ReverseObservation.ObservationDigest,
		MetricDigest:              input.MetricDigest,
		NonAuthorizing:            true,
	})
}

// ExecutionEnvelopeArtifactDigestInput identifies an exact IR or generation
// artifact without executing or authorizing it.
type ExecutionEnvelopeArtifactDigestInput struct {
	ArtifactKind   string
	ArtifactText   string
	NonAuthorizing bool
}

// ExecutionEnvelopeArtifactDigest records content-addressed IR/generation
// provenance before it is attached to a declaration binding.
type ExecutionEnvelopeArtifactDigest struct {
	Status         string
	ArtifactKind   string
	ArtifactDigest string
	MissingStage   string
	NonExecuting   bool
	NonAuthorizing bool
}

// ComputeExecutionEnvelopeArtifactDigest derives a digest from the artifact
// kind and exact text, preventing IR and generation evidence from aliasing.
func ComputeExecutionEnvelopeArtifactDigest(input ExecutionEnvelopeArtifactDigestInput) ExecutionEnvelopeArtifactDigest {
	output := ExecutionEnvelopeArtifactDigest{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.ArtifactKind != "ir" && input.ArtifactKind != "generation" {
		output.MissingStage = "artifact-kind"
		return output
	}
	if strings.TrimSpace(input.ArtifactText) == "" {
		output.ArtifactKind = input.ArtifactKind
		output.MissingStage = "artifact-source"
		return output
	}
	digest, err := Digest(struct {
		ArtifactKind string
		ArtifactText string
	}{
		ArtifactKind: input.ArtifactKind,
		ArtifactText: input.ArtifactText,
	})
	if err != nil {
		output.ArtifactKind = input.ArtifactKind
		output.MissingStage = "artifact-digest"
		return output
	}
	output.Status = "derived"
	output.ArtifactKind = input.ArtifactKind
	output.ArtifactDigest = digest
	return output
}

// ExecutionEnvelopeDeclarationIRGenerationArtifactSourceInput connects the
// exact declaration, IR, and generation texts in one non-executing boundary.
type ExecutionEnvelopeDeclarationIRGenerationArtifactSourceInput struct {
	DeclarationID       string
	ContractID          string
	DeclarationSource   string
	IRSource            string
	GenerationSource    string
	NonAuthorizing      bool
}

// BindExecutionEnvelopeDeclarationIRGenerationFromSources computes every
// source digest before reusing the existing declaration-to-generation bind.
func BindExecutionEnvelopeDeclarationIRGenerationFromSources(input ExecutionEnvelopeDeclarationIRGenerationArtifactSourceInput) ExecutionEnvelopeDeclarationIRGenerationBinding {
	source := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  input.DeclarationID,
		ContractID:     input.ContractID,
		SourceText:     input.DeclarationSource,
		NonAuthorizing: input.NonAuthorizing,
	})
	if source.Status != "derived" {
		return ExecutionEnvelopeDeclarationIRGenerationBinding{
			Status: "UNKNOWN", MissingStage: source.MissingStage,
			NonExecuting: source.NonExecuting, NonAuthorizing: source.NonAuthorizing,
		}
	}
	ir := ComputeExecutionEnvelopeArtifactDigest(ExecutionEnvelopeArtifactDigestInput{
		ArtifactKind: "ir", ArtifactText: input.IRSource, NonAuthorizing: true,
	})
	if ir.Status != "derived" {
		return ExecutionEnvelopeDeclarationIRGenerationBinding{
			Status: "UNKNOWN", MissingStage: "ir-source",
			NonExecuting: true, NonAuthorizing: true,
		}
	}
	generation := ComputeExecutionEnvelopeArtifactDigest(ExecutionEnvelopeArtifactDigestInput{
		ArtifactKind: "generation", ArtifactText: input.GenerationSource, NonAuthorizing: true,
	})
	if generation.Status != "derived" {
		return ExecutionEnvelopeDeclarationIRGenerationBinding{
			Status: "UNKNOWN", MissingStage: "generation-source",
			NonExecuting: true, NonAuthorizing: true,
		}
	}
	return BindExecutionEnvelopeDeclarationIRGeneration(
		source.DeclarationID, source.ContractID, source.DeclarationDigest,
		ir.ArtifactDigest, generation.ArtifactDigest,
	)
}
