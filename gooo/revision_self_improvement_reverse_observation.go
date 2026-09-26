package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementReverseObservation records a non-executing
// re-observation of generated source, IR, and structure against lifecycle
// provenance.
type RevisionSelfImprovementReverseObservation struct {
	Status                        string
	MissingStage                  string
	LifecycleObservationDigest    string
	MetricsBindingDigest          string
	GenerationAssessmentDigest    string
	SourceDigest                  string
	ProposedSourceDigest          string
	InputIRDigest                 string
	ProposedIRDigest              string
	GeneratedSourceDigest         string
	GeneratedIRDigest              string
	StructureDigest               string
	ObservedGeneratedSourceDigest string
	ObservedGeneratedIRDigest     string
	ObservedStructureDigest       string
	ExactSourceMatch              bool
	ExactIRMatch                  bool
	ExactStructureMatch           bool
	ReverseSignal                 string
	ObservationDigest             string
	NonExecuting                  bool
	NonAuthorizing                bool
}

// ObserveRevisionSelfImprovementReverseGeneration re-observes a generation
// receipt against lifecycle provenance without executing or authorizing it.
func ObserveRevisionSelfImprovementReverseGeneration(
	lifecycle RevisionSelfImprovementLifecycleObservation,
	generation GenerationReceipt,
) (RevisionSelfImprovementReverseObservation, error) {
	result := RevisionSelfImprovementReverseObservation{
		Status:                        "UNKNOWN",
		MissingStage:                  "revision-self-improvement-reverse-observation",
		LifecycleObservationDigest:    lifecycle.ObservationDigest,
		MetricsBindingDigest:          lifecycle.MetricsBindingDigest,
		GenerationAssessmentDigest:    lifecycle.GenerationAssessmentDigest,
		SourceDigest:                  lifecycle.SourceDigest,
		ProposedSourceDigest:           lifecycle.ProposedSourceDigest,
		InputIRDigest:                 lifecycle.InputIRDigest,
		ProposedIRDigest:              lifecycle.ProposedIRDigest,
		GeneratedSourceDigest:         lifecycle.GeneratedSourceDigest,
		GeneratedIRDigest:             lifecycle.GeneratedIRDigest,
		StructureDigest:               lifecycle.StructureDigest,
		ObservedGeneratedSourceDigest: generation.GeneratedSourceDigest,
		ObservedGeneratedIRDigest:     generation.GeneratedIRDigest,
		ObservedStructureDigest:       generation.StructureDigest,
		ExactSourceMatch: generation.GeneratedSourceDigest == lifecycle.GeneratedSourceDigest &&
			digestString(generation.GeneratedSource) == generation.GeneratedSourceDigest,
		ExactIRMatch:        generation.GeneratedIRDigest == lifecycle.GeneratedIRDigest,
		ExactStructureMatch: generation.StructureDigest == lifecycle.StructureDigest && generation.StructureMatch,
		ReverseSignal:       "reverse-observation-unknown",
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementReverseObservation(result)
	}
	setObservationDigest()

	if err := lifecycle.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-reverse-lifecycle"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement lifecycle is not valid: %w", err)
	}
	if generation.Status != "BOUND" || generation.MissingStage != "" ||
		!generation.NonExecuting || !generation.NonAuthorizing {
		result.MissingStage = "revision-self-improvement-reverse-generation"
		setObservationDigest()
		return result, fmt.Errorf("generation receipt is not a bounded non-executing observation")
	}
	if generation.SourceDigest != lifecycle.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-reverse-source-link"
		setObservationDigest()
		return result, fmt.Errorf("generation source is not linked to the proposed source")
	}
	if digestString(generation.GeneratedSource) != generation.GeneratedSourceDigest {
		result.MissingStage = "revision-self-improvement-reverse-generated-source-integrity"
		setObservationDigest()
		return result, fmt.Errorf("generated source digest does not match observed generated source")
	}
	if generation.GeneratedSourceDigest != lifecycle.GeneratedSourceDigest {
		result.MissingStage = "revision-self-improvement-reverse-generated-source"
		setObservationDigest()
		return result, fmt.Errorf("generated source does not match lifecycle provenance")
	}
	if generation.GeneratedIRDigest != lifecycle.GeneratedIRDigest {
		result.MissingStage = "revision-self-improvement-reverse-generated-ir"
		setObservationDigest()
		return result, fmt.Errorf("generated IR does not match lifecycle provenance")
	}
	if generation.StructureDigest != lifecycle.StructureDigest || !generation.StructureMatch {
		result.MissingStage = "revision-self-improvement-reverse-structure"
		setObservationDigest()
		return result, fmt.Errorf("generated structure does not match lifecycle provenance")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ReverseSignal = "reverse-observed"
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-reverse-observation"
		result.ReverseSignal = "reverse-observation-unknown"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement reverse observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementReverseObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("self-improvement reverse observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound self-improvement reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown self-improvement reverse observation has no missing stage")
	}
	for name, digest := range map[string]string{
		"lifecycle observation":     o.LifecycleObservationDigest,
		"metrics binding":           o.MetricsBindingDigest,
		"generation assessment":     o.GenerationAssessmentDigest,
		"source":                    o.SourceDigest,
		"proposed source":           o.ProposedSourceDigest,
		"input IR":                  o.InputIRDigest,
		"proposed IR":               o.ProposedIRDigest,
		"generated source":          o.GeneratedSourceDigest,
		"generated IR":              o.GeneratedIRDigest,
		"structure":                 o.StructureDigest,
		"observed generated source": o.ObservedGeneratedSourceDigest,
		"observed generated IR":     o.ObservedGeneratedIRDigest,
		"observed structure":        o.ObservedStructureDigest,
		"observation":               o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("self-improvement reverse observation %s digest is invalid", name)
		}
	}
	if o.ReverseSignal != "reverse-observed" && o.ReverseSignal != "reverse-observation-unknown" &&
		o.ReverseSignal != "reverse-source-mismatch" && o.ReverseSignal != "reverse-ir-mismatch" &&
		o.ReverseSignal != "reverse-structure-mismatch" {
		return fmt.Errorf("self-improvement reverse observation signal is invalid")
	}
	if o.ExactSourceMatch != (o.ObservedGeneratedSourceDigest == o.GeneratedSourceDigest) {
		return fmt.Errorf("self-improvement reverse observation source match is inconsistent")
	}
	if o.ExactIRMatch != (o.ObservedGeneratedIRDigest == o.GeneratedIRDigest) {
		return fmt.Errorf("self-improvement reverse observation IR match is inconsistent")
	}
	if o.ExactStructureMatch != (o.ObservedStructureDigest == o.StructureDigest) {
		return fmt.Errorf("self-improvement reverse observation structure match is inconsistent")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("self-improvement reverse observation must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementReverseObservation(o) != o.ObservationDigest {
		return fmt.Errorf("self-improvement reverse observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementReverseObservation(observation RevisionSelfImprovementReverseObservation) string {
	fields := []string{
		observation.Status,
		observation.MissingStage,
		observation.LifecycleObservationDigest,
		observation.MetricsBindingDigest,
		observation.GenerationAssessmentDigest,
		observation.SourceDigest,
		observation.ProposedSourceDigest,
		observation.InputIRDigest,
		observation.ProposedIRDigest,
		observation.GeneratedSourceDigest,
		observation.GeneratedIRDigest,
		observation.StructureDigest,
		observation.ObservedGeneratedSourceDigest,
		observation.ObservedGeneratedIRDigest,
		observation.ObservedStructureDigest,
		strconv.FormatBool(observation.ExactSourceMatch),
		strconv.FormatBool(observation.ExactIRMatch),
		strconv.FormatBool(observation.ExactStructureMatch),
		observation.ReverseSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}
	return digestString(strings.Join(fields, "|"))
}
