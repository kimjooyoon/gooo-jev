package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceDecisionObservation records a bounded next
// observation signal derived from ordered provenance history.
type RevisionSelfImprovementProvenanceDecisionObservation struct {
	Status                 string
	MissingStage           string
	HistoryDigest          string
	ObservationCount       int
	StableCount            int
	TransitionedCount      int
	MetricsChangeCount     int
	GenerationChangeCount  int
	HistorySignal          string
	DecisionSignal         string
	DecisionReason         string
	RequiresObservation    bool
	RequiresInspection     bool
	RequiresMeasurement    bool
	DecisionDigest         string
	NonExecuting           bool
	NonAuthorizing         bool
}

// ObserveRevisionSelfImprovementProvenanceDecision converts ordered
// provenance evidence into a bounded next-observation signal.
func ObserveRevisionSelfImprovementProvenanceDecision(
	history RevisionSelfImprovementProvenanceHistoryObservation,
) (RevisionSelfImprovementProvenanceDecisionObservation, error) {
	result := RevisionSelfImprovementProvenanceDecisionObservation{
		Status:                "UNKNOWN",
		MissingStage:          "revision-self-improvement-provenance-decision",
		HistoryDigest:         history.HistoryDigest,
		ObservationCount:      history.ObservationCount,
		StableCount:            history.StableCount,
		TransitionedCount:     history.TransitionedCount,
		MetricsChangeCount:    history.MetricsChangeCount,
		GenerationChangeCount: history.GenerationChangeCount,
		HistorySignal:         history.HistorySignal,
		RequiresObservation:   true,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setDecisionDigest := func() {
		result.DecisionDigest = digestRevisionSelfImprovementProvenanceDecision(result)
	}
	setDecisionDigest()

	if err := history.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-decision-history"
		setDecisionDigest()
		return result, fmt.Errorf("self-improvement provenance history is not valid: %w", err)
	}

	switch history.HistorySignal {
	case "stable":
		result.DecisionSignal = "observe"
		result.DecisionReason = "provenance-stable-observation"
	case "transitioned":
		result.DecisionSignal = "inspect"
		result.DecisionReason = "provenance-transition-requires-inspection"
		result.RequiresInspection = true
		result.RequiresMeasurement = history.MetricsChangeCount > 0 || history.GenerationChangeCount > 0
	case "mixed":
		result.DecisionSignal = "inspect"
		result.DecisionReason = "mixed-provenance-requires-inspection"
		result.RequiresInspection = true
		result.RequiresMeasurement = history.MetricsChangeCount > 0 || history.GenerationChangeCount > 0
	default:
		result.MissingStage = "revision-self-improvement-provenance-decision-signal"
		setDecisionDigest()
		return result, fmt.Errorf("provenance history signal is not recognized")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDecisionDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-decision"
		setDecisionDigest()
		return result, fmt.Errorf("self-improvement provenance decision is not valid: %w", err)
	}
	return result, nil
}

func (d RevisionSelfImprovementProvenanceDecisionObservation) Validate() error {
	if d.Status == "" {
		return fmt.Errorf("self-improvement provenance decision status is empty")
	}
	if d.Status == "BOUND" && d.MissingStage != "" {
		return fmt.Errorf("bound self-improvement provenance decision has a missing stage")
	}
	if d.Status == "UNKNOWN" && d.MissingStage == "" {
		return fmt.Errorf("unknown self-improvement provenance decision has no missing stage")
	}
	if !validDigest(d.HistoryDigest) || !validDigest(d.DecisionDigest) {
		return fmt.Errorf("self-improvement provenance decision digest is invalid")
	}
	if d.ObservationCount < 1 || d.StableCount < 0 || d.TransitionedCount < 0 ||
		d.StableCount+d.TransitionedCount != d.ObservationCount {
		return fmt.Errorf("self-improvement provenance decision counts are invalid")
	}
	if d.MetricsChangeCount < 0 || d.MetricsChangeCount > d.ObservationCount ||
		d.GenerationChangeCount < 0 || d.GenerationChangeCount > d.ObservationCount {
		return fmt.Errorf("self-improvement provenance decision change counts are invalid")
	}
	if d.HistorySignal != "stable" && d.HistorySignal != "transitioned" &&
		d.HistorySignal != "mixed" {
		return fmt.Errorf("self-improvement provenance decision history signal is invalid")
	}
	if d.DecisionSignal != "observe" && d.DecisionSignal != "inspect" {
		return fmt.Errorf("self-improvement provenance decision signal is invalid")
	}
	if d.Status == "BOUND" && d.DecisionReason == "" {
		return fmt.Errorf("bound self-improvement provenance decision has no reason")
	}
	if !d.RequiresObservation {
		return fmt.Errorf("self-improvement provenance decision must require observation")
	}
	switch d.HistorySignal {
	case "stable":
		if d.DecisionSignal != "observe" || d.RequiresInspection || d.RequiresMeasurement {
			return fmt.Errorf("stable provenance decision does not match its signal")
		}
	case "transitioned", "mixed":
		if d.DecisionSignal != "inspect" || !d.RequiresInspection {
			return fmt.Errorf("transitioned provenance decision does not match its signal")
		}
		if d.RequiresMeasurement != (d.MetricsChangeCount > 0 || d.GenerationChangeCount > 0) {
			return fmt.Errorf("provenance decision measurement flag is not linked")
		}
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("self-improvement provenance decision must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceDecision(d) != d.DecisionDigest {
		return fmt.Errorf("self-improvement provenance decision digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceDecision(
	decision RevisionSelfImprovementProvenanceDecisionObservation,
) string {
	fields := []string{
		decision.Status,
		decision.MissingStage,
		decision.HistoryDigest,
		strconv.Itoa(decision.ObservationCount),
		strconv.Itoa(decision.StableCount),
		strconv.Itoa(decision.TransitionedCount),
		strconv.Itoa(decision.MetricsChangeCount),
		strconv.Itoa(decision.GenerationChangeCount),
		decision.HistorySignal,
		decision.DecisionSignal,
		decision.DecisionReason,
		strconv.FormatBool(decision.RequiresObservation),
		strconv.FormatBool(decision.RequiresInspection),
		strconv.FormatBool(decision.RequiresMeasurement),
		strconv.FormatBool(decision.NonExecuting),
		strconv.FormatBool(decision.NonAuthorizing),
	}
	return digestString(strings.Join(fields, "|"))
}
