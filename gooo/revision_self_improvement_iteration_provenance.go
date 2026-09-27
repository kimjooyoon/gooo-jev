package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementIterationProvenanceObservation binds the legacy
// iteration boundary to provenance history without executing or authorizing a
// candidate change.
type RevisionSelfImprovementIterationProvenanceObservation struct {
	Status                       string
	MissingStage                 string
	StageCount                   int
	IterationDigest              string
	SourceDigest                 string
	IRDigest                     string
	LegacyHistoryDigest          string
	ProvenanceHistoryDigest      string
	WindowDigest                 string
	FeedbackDigest               string
	BridgeDigest                 string
	CandidateSourceDigest        string
	CandidateProposedSourceDigest string
	CandidateGeneratedIRDigest   string
	LegacyHistorySignal           string
	ProvenanceHistorySignal       string
	DecisionSignal                string
	FeedbackSignal                string
	BridgeSignal                  string
	MetricsChangeCount            int
	GenerationChangeCount         int
	RequiresObservation           bool
	RequiresInspection            bool
	RequiresMeasurement           bool
	ObservationSignal             string
	ObservationDigest             string
	NonExecuting                  bool
	NonAuthorizing                bool
}

// ObserveRevisionSelfImprovementIterationProvenance links legacy iteration
// evidence with the provenance feedback bridge while preserving explicit
// signal mappings and unresolved stages.
func ObserveRevisionSelfImprovementIterationProvenance(
	iteration RevisionSelfImprovementIteration,
	history RevisionSelfImprovementProvenanceHistoryObservation,
	bridge RevisionSelfImprovementProvenanceFeedbackBridgeObservation,
) (RevisionSelfImprovementIterationProvenanceObservation, error) {
	result := RevisionSelfImprovementIterationProvenanceObservation{
		Status:                       "UNKNOWN",
		MissingStage:                 "revision-self-improvement-iteration-provenance",
		StageCount:                   4,
		IterationDigest:              iteration.IterationDigest,
		SourceDigest:                 iteration.SourceDigest,
		IRDigest:                     iteration.IRDigest,
		LegacyHistoryDigest:           iteration.HistoryDigest,
		ProvenanceHistoryDigest:       history.HistoryDigest,
		WindowDigest:                 iteration.WindowDigest,
		FeedbackDigest:               iteration.FeedbackDigest,
		BridgeDigest:                 bridge.BridgeDigest,
		CandidateSourceDigest:         iteration.CandidateSourceDigest,
		CandidateProposedSourceDigest: iteration.CandidateProposedSourceDigest,
		CandidateGeneratedIRDigest:   iteration.CandidateGeneratedIRDigest,
		LegacyHistorySignal:           iteration.HistorySignal,
		ProvenanceHistorySignal:       history.HistorySignal,
		DecisionSignal:                bridge.DecisionSignal,
		FeedbackSignal:               bridge.FeedbackSignal,
		BridgeSignal:                 bridge.BridgeSignal,
		MetricsChangeCount:            bridge.MetricsChangeCount,
		GenerationChangeCount:         bridge.GenerationChangeCount,
		RequiresObservation:           bridge.RequiresObservation,
		RequiresInspection:            bridge.RequiresInspection,
		RequiresMeasurement:           bridge.RequiresMeasurement,
		ObservationSignal:             "iteration-provenance-unknown",
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementIterationProvenance(result)
	}
	setObservationDigest()

	if err := iteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-iteration-provenance-iteration"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement iteration is not valid: %w", err)
	}
	if err := history.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-iteration-provenance-history"
		setObservationDigest()
		return result, fmt.Errorf("provenance history is not valid: %w", err)
	}
	if err := bridge.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-iteration-provenance-bridge"
		setObservationDigest()
		return result, fmt.Errorf("provenance feedback bridge is not valid: %w", err)
	}
	if bridge.HistoryDigest != history.HistoryDigest {
		result.MissingStage = "revision-self-improvement-iteration-provenance-history-link"
		setObservationDigest()
		return result, fmt.Errorf("provenance feedback bridge is not linked to provenance history")
	}
	if bridge.FeedbackDigest != iteration.FeedbackDigest ||
		bridge.WindowDigest != iteration.WindowDigest {
		result.MissingStage = "revision-self-improvement-iteration-provenance-feedback-link"
		setObservationDigest()
		return result, fmt.Errorf("provenance feedback bridge is not linked to legacy iteration feedback")
	}
	if bridge.MetricsChangeCount != history.MetricsChangeCount ||
		bridge.GenerationChangeCount != history.GenerationChangeCount {
		result.MissingStage = "revision-self-improvement-iteration-provenance-metric-link"
		setObservationDigest()
		return result, fmt.Errorf("provenance change counts are not linked to history")
	}

	switch history.HistorySignal {
	case "stable":
		if iteration.HistorySignal != "stable" ||
			iteration.FeedbackSignal != "observe" ||
			bridge.DecisionSignal != "observe" ||
			bridge.FeedbackSignal != "observe" {
			result.MissingStage = "revision-self-improvement-iteration-provenance-signal-mapping"
			setObservationDigest()
			return result, fmt.Errorf("stable provenance does not map to stable legacy iteration")
		}
	case "transitioned", "mixed":
		if iteration.HistorySignal != "mixed" ||
			iteration.FeedbackSignal != "inspect" ||
			bridge.DecisionSignal != "inspect" ||
			bridge.FeedbackSignal != "inspect" {
			result.MissingStage = "revision-self-improvement-iteration-provenance-signal-mapping"
			setObservationDigest()
			return result, fmt.Errorf("changed provenance does not map to mixed inspecting iteration")
		}
	default:
		result.MissingStage = "revision-self-improvement-iteration-provenance-signal"
		setObservationDigest()
		return result, fmt.Errorf("provenance history signal is not recognized")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ObservationSignal = "iteration-provenance-bound"
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-iteration-provenance"
		result.ObservationSignal = "iteration-provenance-unknown"
		setObservationDigest()
		return result, fmt.Errorf("iteration provenance observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementIterationProvenanceObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("iteration provenance observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound iteration provenance observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown iteration provenance observation has no missing stage")
	}
	if o.StageCount != 4 {
		return fmt.Errorf("iteration provenance observation stage count must be 4")
	}
	for name, digest := range map[string]string{
		"iteration":                o.IterationDigest,
		"source":                   o.SourceDigest,
		"ir":                       o.IRDigest,
		"legacy-history":           o.LegacyHistoryDigest,
		"provenance-history":       o.ProvenanceHistoryDigest,
		"window":                   o.WindowDigest,
		"feedback":                 o.FeedbackDigest,
		"bridge":                   o.BridgeDigest,
		"candidate-source":         o.CandidateSourceDigest,
		"candidate-proposed":       o.CandidateProposedSourceDigest,
		"candidate-generated-ir":   o.CandidateGeneratedIRDigest,
		"observation":              o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("iteration provenance observation %s digest is invalid", name)
		}
	}
	if o.LegacyHistorySignal != "stable" && o.LegacyHistorySignal != "narrower" &&
		o.LegacyHistorySignal != "wider" && o.LegacyHistorySignal != "mixed" {
		return fmt.Errorf("iteration provenance observation legacy history signal is invalid")
	}
	if o.ProvenanceHistorySignal != "stable" &&
		o.ProvenanceHistorySignal != "transitioned" &&
		o.ProvenanceHistorySignal != "mixed" {
		return fmt.Errorf("iteration provenance observation provenance history signal is invalid")
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "inspect" {
		return fmt.Errorf("iteration provenance observation decision signal is invalid")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "inspect" {
		return fmt.Errorf("iteration provenance observation feedback signal is invalid")
	}
	if o.BridgeSignal != "provenance-feedback-bridge-bound" &&
		o.BridgeSignal != "provenance-feedback-bridge-unknown" {
		return fmt.Errorf("iteration provenance observation bridge signal is invalid")
	}
	if o.ObservationSignal != "iteration-provenance-bound" &&
		o.ObservationSignal != "iteration-provenance-unknown" {
		return fmt.Errorf("iteration provenance observation signal is invalid")
	}
	if o.MetricsChangeCount < 0 || o.GenerationChangeCount < 0 {
		return fmt.Errorf("iteration provenance observation change counts are invalid")
	}
	if !o.RequiresObservation {
		return fmt.Errorf("iteration provenance observation must require observation")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("iteration provenance observation must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementIterationProvenance(o) != o.ObservationDigest {
		return fmt.Errorf("iteration provenance observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementIterationProvenance(
	observation RevisionSelfImprovementIterationProvenanceObservation,
) string {
	fields := []string{
		observation.Status,
		observation.MissingStage,
		strconv.Itoa(observation.StageCount),
		observation.IterationDigest,
		observation.SourceDigest,
		observation.IRDigest,
		observation.LegacyHistoryDigest,
		observation.ProvenanceHistoryDigest,
		observation.WindowDigest,
		observation.FeedbackDigest,
		observation.BridgeDigest,
		observation.CandidateSourceDigest,
		observation.CandidateProposedSourceDigest,
		observation.CandidateGeneratedIRDigest,
		observation.LegacyHistorySignal,
		observation.ProvenanceHistorySignal,
		observation.DecisionSignal,
		observation.FeedbackSignal,
		observation.BridgeSignal,
		strconv.Itoa(observation.MetricsChangeCount),
		strconv.Itoa(observation.GenerationChangeCount),
		strconv.FormatBool(observation.RequiresObservation),
		strconv.FormatBool(observation.RequiresInspection),
		strconv.FormatBool(observation.RequiresMeasurement),
		observation.ObservationSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}
	return digestString(strings.Join(fields, "|"))
}
