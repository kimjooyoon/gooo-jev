package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceTransitionObservation records exact
// provenance continuity and drift between two validated lifecycle snapshots.
type RevisionSelfImprovementProvenanceTransitionObservation struct {
	Status                               string
	MissingStage                         string
	PreviousLifecycleObservationDigest   string
	CurrentLifecycleObservationDigest    string
	PreviousSourceDigest                 string
	CurrentSourceDigest                  string
	PreviousProposedSourceDigest         string
	CurrentProposedSourceDigest          string
	PreviousInputIRDigest                string
	CurrentInputIRDigest                 string
	PreviousProposedIRDigest             string
	CurrentProposedIRDigest              string
	PreviousGeneratedSourceDigest        string
	CurrentGeneratedSourceDigest         string
	PreviousGeneratedIRDigest            string
	CurrentGeneratedIRDigest             string
	PreviousStructureDigest              string
	CurrentStructureDigest               string
	PreviousMetricsBindingDigest         string
	CurrentMetricsBindingDigest          string
	PreviousGenerationAssessmentDigest   string
	CurrentGenerationAssessmentDigest    string
	PreviousChangedByteCount             int
	CurrentChangedByteCount              int
	PreviousChangedLineCount             int
	CurrentChangedLineCount              int
	PreviousIRChanged                    bool
	CurrentIRChanged                     bool
	PreviousExactSourceMatch             bool
	CurrentExactSourceMatch              bool
	PreviousStructureMatch               bool
	CurrentStructureMatch                bool
	SourceChanged                        bool
	ProposedSourceChanged                bool
	InputIRChanged                       bool
	ProposedIRChanged                    bool
	GeneratedSourceChanged               bool
	GeneratedIRChanged                   bool
	StructureChanged                     bool
	MetricsChanged                       bool
	GenerationChanged                    bool
	TransitionSignal                     string
	ObservationDigest                    string
	NonExecuting                         bool
	NonAuthorizing                       bool
}

// ObserveRevisionSelfImprovementProvenanceTransition compares two lifecycle
// snapshots without executing, authorizing, or inferring a change.
func ObserveRevisionSelfImprovementProvenanceTransition(
	previous RevisionSelfImprovementLifecycleObservation,
	current RevisionSelfImprovementLifecycleObservation,
) (RevisionSelfImprovementProvenanceTransitionObservation, error) {
	result := RevisionSelfImprovementProvenanceTransitionObservation{
		Status:                             "UNKNOWN",
		MissingStage:                       "revision-self-improvement-provenance-transition",
		PreviousLifecycleObservationDigest: previous.ObservationDigest,
		CurrentLifecycleObservationDigest:  current.ObservationDigest,
		PreviousSourceDigest:               previous.SourceDigest,
		CurrentSourceDigest:                current.SourceDigest,
		PreviousProposedSourceDigest:       previous.ProposedSourceDigest,
		CurrentProposedSourceDigest:         current.ProposedSourceDigest,
		PreviousInputIRDigest:              previous.InputIRDigest,
		CurrentInputIRDigest:               current.InputIRDigest,
		PreviousProposedIRDigest:           previous.ProposedIRDigest,
		CurrentProposedIRDigest:             current.ProposedIRDigest,
		PreviousGeneratedSourceDigest:      previous.GeneratedSourceDigest,
		CurrentGeneratedSourceDigest:       current.GeneratedSourceDigest,
		PreviousGeneratedIRDigest:          previous.GeneratedIRDigest,
		CurrentGeneratedIRDigest:            current.GeneratedIRDigest,
		PreviousStructureDigest:            previous.StructureDigest,
		CurrentStructureDigest:              current.StructureDigest,
		PreviousMetricsBindingDigest:       previous.MetricsBindingDigest,
		CurrentMetricsBindingDigest:         current.MetricsBindingDigest,
		PreviousGenerationAssessmentDigest: previous.GenerationAssessmentDigest,
		CurrentGenerationAssessmentDigest:  current.GenerationAssessmentDigest,
		PreviousChangedByteCount:            previous.ChangedByteCount,
		CurrentChangedByteCount:             current.ChangedByteCount,
		PreviousChangedLineCount:            previous.ChangedLineCount,
		CurrentChangedLineCount:              current.ChangedLineCount,
		PreviousIRChanged:                  previous.IRChanged,
		CurrentIRChanged:                    current.IRChanged,
		PreviousExactSourceMatch:            previous.ExactSourceMatch,
		CurrentExactSourceMatch:             current.ExactSourceMatch,
		PreviousStructureMatch:              previous.StructureMatch,
		CurrentStructureMatch:               current.StructureMatch,
		SourceChanged:                       previous.SourceDigest != current.SourceDigest,
		ProposedSourceChanged:               previous.ProposedSourceDigest != current.ProposedSourceDigest,
		InputIRChanged:                      previous.InputIRDigest != current.InputIRDigest,
		ProposedIRChanged:                   previous.ProposedIRDigest != current.ProposedIRDigest,
		GeneratedSourceChanged:              previous.GeneratedSourceDigest != current.GeneratedSourceDigest,
		GeneratedIRChanged:                  previous.GeneratedIRDigest != current.GeneratedIRDigest,
		StructureChanged:                    previous.StructureDigest != current.StructureDigest,
		MetricsChanged: previous.MetricsBindingDigest != current.MetricsBindingDigest ||
			previous.ChangedByteCount != current.ChangedByteCount ||
			previous.ChangedLineCount != current.ChangedLineCount ||
			previous.IRChanged != current.IRChanged ||
			previous.ExactSourceMatch != current.ExactSourceMatch ||
			previous.StructureMatch != current.StructureMatch,
		GenerationChanged: previous.GenerationAssessmentDigest != current.GenerationAssessmentDigest ||
			previous.GeneratedSourceDigest != current.GeneratedSourceDigest ||
			previous.GeneratedIRDigest != current.GeneratedIRDigest ||
			previous.StructureDigest != current.StructureDigest,
		TransitionSignal: "provenance-transition-unknown",
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementProvenanceTransition(result)
	}
	setObservationDigest()

	if err := previous.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-transition-previous"
		setObservationDigest()
		return result, fmt.Errorf("previous self-improvement lifecycle is not valid: %w", err)
	}
	if err := current.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-transition-current"
		setObservationDigest()
		return result, fmt.Errorf("current self-improvement lifecycle is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	if result.SourceChanged || result.ProposedSourceChanged || result.InputIRChanged ||
		result.ProposedIRChanged || result.GeneratedSourceChanged ||
		result.GeneratedIRChanged || result.StructureChanged || result.MetricsChanged ||
		result.GenerationChanged {
		result.TransitionSignal = "provenance-transition-bound"
	} else {
		result.TransitionSignal = "provenance-stable"
	}
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-transition"
		result.TransitionSignal = "provenance-transition-unknown"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement provenance transition is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenanceTransitionObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("self-improvement provenance transition status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound self-improvement provenance transition has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown self-improvement provenance transition has no missing stage")
	}
	for name, digest := range map[string]string{
		"previous lifecycle":           o.PreviousLifecycleObservationDigest,
		"current lifecycle":            o.CurrentLifecycleObservationDigest,
		"previous source":              o.PreviousSourceDigest,
		"current source":               o.CurrentSourceDigest,
		"previous proposed source":     o.PreviousProposedSourceDigest,
		"current proposed source":      o.CurrentProposedSourceDigest,
		"previous input IR":            o.PreviousInputIRDigest,
		"current input IR":             o.CurrentInputIRDigest,
		"previous proposed IR":         o.PreviousProposedIRDigest,
		"current proposed IR":          o.CurrentProposedIRDigest,
		"previous generated source":    o.PreviousGeneratedSourceDigest,
		"current generated source":     o.CurrentGeneratedSourceDigest,
		"previous generated IR":        o.PreviousGeneratedIRDigest,
		"current generated IR":         o.CurrentGeneratedIRDigest,
		"previous structure":            o.PreviousStructureDigest,
		"current structure":             o.CurrentStructureDigest,
		"previous metrics":              o.PreviousMetricsBindingDigest,
		"current metrics":               o.CurrentMetricsBindingDigest,
		"previous generation":           o.PreviousGenerationAssessmentDigest,
		"current generation":            o.CurrentGenerationAssessmentDigest,
		"observation":                   o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("self-improvement provenance transition %s digest is invalid", name)
		}
	}
	if o.TransitionSignal != "provenance-transition-bound" &&
		o.TransitionSignal != "provenance-stable" &&
		o.TransitionSignal != "provenance-transition-unknown" {
		return fmt.Errorf("self-improvement provenance transition signal is invalid")
	}
	if o.SourceChanged != (o.PreviousSourceDigest != o.CurrentSourceDigest) ||
		o.ProposedSourceChanged != (o.PreviousProposedSourceDigest != o.CurrentProposedSourceDigest) ||
		o.InputIRChanged != (o.PreviousInputIRDigest != o.CurrentInputIRDigest) ||
		o.ProposedIRChanged != (o.PreviousProposedIRDigest != o.CurrentProposedIRDigest) ||
		o.GeneratedSourceChanged != (o.PreviousGeneratedSourceDigest != o.CurrentGeneratedSourceDigest) ||
		o.GeneratedIRChanged != (o.PreviousGeneratedIRDigest != o.CurrentGeneratedIRDigest) ||
		o.StructureChanged != (o.PreviousStructureDigest != o.CurrentStructureDigest) {
		return fmt.Errorf("self-improvement provenance transition digest change flags are inconsistent")
	}
	metricsChanged := o.PreviousMetricsBindingDigest != o.CurrentMetricsBindingDigest ||
		o.PreviousChangedByteCount != o.CurrentChangedByteCount ||
		o.PreviousChangedLineCount != o.CurrentChangedLineCount ||
		o.PreviousIRChanged != o.CurrentIRChanged ||
		o.PreviousExactSourceMatch != o.CurrentExactSourceMatch ||
		o.PreviousStructureMatch != o.CurrentStructureMatch
	if o.MetricsChanged != metricsChanged {
		return fmt.Errorf("self-improvement provenance transition metrics flag is inconsistent")
	}
	generationChanged := o.PreviousGenerationAssessmentDigest != o.CurrentGenerationAssessmentDigest ||
		o.PreviousGeneratedSourceDigest != o.CurrentGeneratedSourceDigest ||
		o.PreviousGeneratedIRDigest != o.CurrentGeneratedIRDigest ||
		o.PreviousStructureDigest != o.CurrentStructureDigest
	if o.GenerationChanged != generationChanged {
		return fmt.Errorf("self-improvement provenance transition generation flag is inconsistent")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("self-improvement provenance transition must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceTransition(o) != o.ObservationDigest {
		return fmt.Errorf("self-improvement provenance transition digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceTransition(observation RevisionSelfImprovementProvenanceTransitionObservation) string {
	fields := []string{
		observation.Status,
		observation.MissingStage,
		observation.PreviousLifecycleObservationDigest,
		observation.CurrentLifecycleObservationDigest,
		observation.PreviousSourceDigest,
		observation.CurrentSourceDigest,
		observation.PreviousProposedSourceDigest,
		observation.CurrentProposedSourceDigest,
		observation.PreviousInputIRDigest,
		observation.CurrentInputIRDigest,
		observation.PreviousProposedIRDigest,
		observation.CurrentProposedIRDigest,
		observation.PreviousGeneratedSourceDigest,
		observation.CurrentGeneratedSourceDigest,
		observation.PreviousGeneratedIRDigest,
		observation.CurrentGeneratedIRDigest,
		observation.PreviousStructureDigest,
		observation.CurrentStructureDigest,
		observation.PreviousMetricsBindingDigest,
		observation.CurrentMetricsBindingDigest,
		observation.PreviousGenerationAssessmentDigest,
		observation.CurrentGenerationAssessmentDigest,
		strconv.Itoa(observation.PreviousChangedByteCount),
		strconv.Itoa(observation.CurrentChangedByteCount),
		strconv.Itoa(observation.PreviousChangedLineCount),
		strconv.Itoa(observation.CurrentChangedLineCount),
		strconv.FormatBool(observation.PreviousIRChanged),
		strconv.FormatBool(observation.CurrentIRChanged),
		strconv.FormatBool(observation.PreviousExactSourceMatch),
		strconv.FormatBool(observation.CurrentExactSourceMatch),
		strconv.FormatBool(observation.PreviousStructureMatch),
		strconv.FormatBool(observation.CurrentStructureMatch),
		strconv.FormatBool(observation.SourceChanged),
		strconv.FormatBool(observation.ProposedSourceChanged),
		strconv.FormatBool(observation.InputIRChanged),
		strconv.FormatBool(observation.ProposedIRChanged),
		strconv.FormatBool(observation.GeneratedSourceChanged),
		strconv.FormatBool(observation.GeneratedIRChanged),
		strconv.FormatBool(observation.StructureChanged),
		strconv.FormatBool(observation.MetricsChanged),
		strconv.FormatBool(observation.GenerationChanged),
		observation.TransitionSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}
	return digestString(strings.Join(fields, "|"))
}
