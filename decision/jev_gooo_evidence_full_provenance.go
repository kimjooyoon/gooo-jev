package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceInput connects extended .gooo
// evidence declarations to the existing full provenance source path.
type ExecutionEnvelopeGoooEvidenceFullProvenanceInput struct {
	DeclarationID      string
	ContractID         string
	SourceText         string
	EvidenceGeneration GoooEvidenceDeclarationIRGeneration
	ObservedStatus     string
	ExpectedStatus     string
	ObservedMissingStage string
	ExpectedMissingStage string
	ObservedEvidenceDigest string
	ExpectedEvidenceDigest string
	ReverseObservationSource string
	MetricSource       string
	NonAuthorizing     bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceBinding preserves legacy full
// provenance and the extended evidence declaration digest separately.
type ExecutionEnvelopeGoooEvidenceFullProvenanceBinding struct {
	Status                    string
	MissingStage              string
	DeclarationID             string
	ContractID                string
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	BaseBindingDigest         string
	EvidenceBindingDigest     string
	EvidenceDeclarationDigest string
	ReverseObservationDigest  string
	MetricDigest              string
	ProvenanceEvidenceDigest  string
	CompletenessDigest        string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceBinding) Validate() error {
	if b.Status != "complete" ||
		b.DeclarationID == "" ||
		b.ContractID == "" ||
		b.DeclarationDigest == "" ||
		b.IRDigest == "" ||
		b.GenerationDigest == "" ||
		b.BaseBindingDigest == "" ||
		b.EvidenceBindingDigest == "" ||
		b.EvidenceDeclarationDigest == "" ||
		b.ReverseObservationDigest == "" ||
		b.MetricDigest == "" ||
		b.ProvenanceEvidenceDigest == "" ||
		b.CompletenessDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo evidence full provenance binding")
	}
	if !b.NonExecuting {
		return fmt.Errorf("Gooo evidence full provenance binding must be non-executing")
	}
	if !b.NonAuthorizing {
		return fmt.Errorf("Gooo evidence full provenance binding must be non-authorizing")
	}
	bindingDigest, err := Digest(struct {
		BaseBindingDigest         string
		EvidenceBindingDigest     string
		EvidenceDeclarationDigest string
		ProvenanceEvidenceDigest  string
	}{
		BaseBindingDigest:         b.BaseBindingDigest,
		EvidenceBindingDigest:     b.EvidenceBindingDigest,
		EvidenceDeclarationDigest: b.EvidenceDeclarationDigest,
		ProvenanceEvidenceDigest:  b.ProvenanceEvidenceDigest,
	})
	if err != nil || b.EvidenceBindingDigest != bindingDigest {
		return fmt.Errorf("Gooo evidence full provenance base binding digest mismatch")
	}
	expectedEvidence, err := Digest(struct {
		DeclarationDigest         string
		IRDigest                  string
		GenerationDigest          string
		BaseBindingDigest         string
		EvidenceBindingDigest     string
		EvidenceDeclarationDigest string
		ReverseObservationDigest  string
		MetricDigest              string
		ProvenanceEvidenceDigest  string
		CompletenessDigest        string
	}{
		DeclarationDigest:         b.DeclarationDigest,
		IRDigest:                  b.IRDigest,
		GenerationDigest:          b.GenerationDigest,
		BaseBindingDigest:         b.BaseBindingDigest,
		EvidenceDeclarationDigest: b.EvidenceDeclarationDigest,
		ReverseObservationDigest:  b.ReverseObservationDigest,
		MetricDigest:              b.MetricDigest,
		ProvenanceEvidenceDigest:  b.ProvenanceEvidenceDigest,
		CompletenessDigest:        b.CompletenessDigest,
	})
	if err != nil || b.EvidenceDigest != expectedEvidence {
		return fmt.Errorf("Gooo evidence full provenance evidence digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenance replays both the extended
// evidence syntax and the legacy source chain before emitting complete status.
func BindExecutionEnvelopeGoooEvidenceFullProvenance(input ExecutionEnvelopeGoooEvidenceFullProvenanceInput) ExecutionEnvelopeGoooEvidenceFullProvenanceBinding {
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceBinding{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.EvidenceGeneration.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.EvidenceGeneration.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if strings.TrimSpace(input.SourceText) == "" {
		output.MissingStage = "declaration-source"
		return output
	}
	derived := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText:     input.SourceText,
		NonAuthorizing: true,
	})
	if derived.Status != "ready" ||
		input.EvidenceGeneration.Status != "ready" ||
		derived.SourceDigest != input.EvidenceGeneration.SourceDigest ||
		derived.BaseIRDigest != input.EvidenceGeneration.BaseIRDigest ||
		derived.EvidenceDigest != input.EvidenceGeneration.EvidenceDigest ||
		derived.GeneratedSource != input.EvidenceGeneration.GeneratedSource ||
		derived.GenerationDigest != input.EvidenceGeneration.GenerationDigest {
		output.MissingStage = "evidence-ir-generation-replay"
		return output
	}
	if len(derived.Evidence) == 0 {
		output.MissingStage = "evidence-declaration"
		return output
	}
	baseSource, err := GenerateGoooDeclaration(derived.BaseIR)
	if err != nil {
		output.MissingStage = "legacy-generation"
		return output
	}
	baseGeneration := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     baseSource,
		NonAuthorizing: true,
	})
	if baseGeneration.Status != "ready" {
		output.MissingStage = "legacy-ir-generation"
		return output
	}
	full := BindExecutionEnvelopeFullProvenanceFromGooo(ExecutionEnvelopeGoooFullProvenanceInput{
		DeclarationID:            input.DeclarationID,
		ContractID:               input.ContractID,
		SourceText:               baseSource,
		IRGeneration:             baseGeneration,
		ObservedStatus:            input.ObservedStatus,
		ExpectedStatus:            input.ExpectedStatus,
		ObservedMissingStage:      input.ObservedMissingStage,
		ExpectedMissingStage:      input.ExpectedMissingStage,
		ObservedEvidenceDigest:    input.ObservedEvidenceDigest,
		ExpectedEvidenceDigest:    input.ExpectedEvidenceDigest,
		ReverseObservationSource: input.ReverseObservationSource,
		MetricSource:              input.MetricSource,
		NonAuthorizing:            true,
	})
	if full.Status != "complete" {
		output.MissingStage = full.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "full-provenance"
		}
		return output
	}
	output.Status = "complete"
	output.DeclarationID = full.DeclarationID
	output.ContractID = full.ContractID
	output.DeclarationDigest = full.DeclarationDigest
	output.IRDigest = full.IRDigest
	output.GenerationDigest = full.GenerationDigest
	output.BaseBindingDigest = full.BindingDigest
	output.EvidenceDeclarationDigest = derived.EvidenceDigest
	output.ReverseObservationDigest = full.ReverseObservationDigest
	output.MetricDigest = full.MetricDigest
	output.ProvenanceEvidenceDigest = full.EvidenceDigest
	output.CompletenessDigest = full.CompletenessDigest
	output.EvidenceBindingDigest, err = Digest(struct {
		BaseBindingDigest         string
		EvidenceBindingDigest     string
		EvidenceDeclarationDigest string
		ProvenanceEvidenceDigest  string
	}{
		BaseBindingDigest:         output.BaseBindingDigest,
		EvidenceBindingDigest:     output.EvidenceBindingDigest,
		EvidenceDeclarationDigest: output.EvidenceDeclarationDigest,
		ProvenanceEvidenceDigest:  output.ProvenanceEvidenceDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "evidence-binding"
		return output
	}
	_ = bindingDigest
	output.EvidenceDigest, err = Digest(struct {
		DeclarationDigest         string
		IRDigest                  string
		GenerationDigest          string
		BaseBindingDigest         string
		EvidenceBindingDigest     string
		EvidenceDeclarationDigest string
		ReverseObservationDigest  string
		MetricDigest              string
		ProvenanceEvidenceDigest  string
		CompletenessDigest        string
	}{
		DeclarationDigest:         output.DeclarationDigest,
		IRDigest:                  output.IRDigest,
		GenerationDigest:          output.GenerationDigest,
		BaseBindingDigest:         output.BaseBindingDigest,
		EvidenceDeclarationDigest: output.EvidenceDeclarationDigest,
		ReverseObservationDigest:  output.ReverseObservationDigest,
		MetricDigest:              output.MetricDigest,
		ProvenanceEvidenceDigest:  output.ProvenanceEvidenceDigest,
		CompletenessDigest:        output.CompletenessDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "evidence-binding"
		output.EvidenceDigest = ""
		return output
	}
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "evidence-binding"
		output.EvidenceDigest = ""
	}
	return output
}