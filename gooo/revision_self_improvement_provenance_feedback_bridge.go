package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceFeedbackBridgeObservation explicitly links
// provenance decision evidence to the existing window feedback vocabulary.
type RevisionSelfImprovementProvenanceFeedbackBridgeObservation struct {
	Status             string
	MissingStage       string
	HistoryDigest      string
	DecisionDigest     string
	FeedbackDigest     string
	WindowDigest       string
	HistorySignal      string
	ComparisonSignal   string
	DecisionSignal     string
	FeedbackSignal     string
	DecisionReason     string
	MetricsChangeCount int
	GenerationChangeCount int
	RequiresObservation bool
	RequiresInspection  bool
	RequiresMeasurement bool
	BridgeSignal        string
	BridgeDigest        string
	NonExecuting        bool
	NonAuthorizing      bool
}

// ObserveRevisionSelfImprovementProvenanceFeedbackBridge verifies the bounded
// mapping from provenance history and decision signals to legacy feedback.
func ObserveRevisionSelfImprovementProvenanceFeedbackBridge(
	history RevisionSelfImprovementProvenanceHistoryObservation,
	decision RevisionSelfImprovementProvenanceDecisionObservation,
	feedback RevisionSelfImprovementFeedback,
) (RevisionSelfImprovementProvenanceFeedbackBridgeObservation, error) {
	result := RevisionSelfImprovementProvenanceFeedbackBridgeObservation{
		Status:                "UNKNOWN",
		MissingStage:          "revision-self-improvement-provenance-feedback-bridge",
		HistoryDigest:         history.HistoryDigest,
		DecisionDigest:        decision.DecisionDigest,
		FeedbackDigest:        feedback.FeedbackDigest,
		WindowDigest:          feedback.WindowDigest,
		HistorySignal:         history.HistorySignal,
		ComparisonSignal:      feedback.ComparisonSignal,
		DecisionSignal:        decision.DecisionSignal,
		FeedbackSignal:        feedback.FeedbackSignal,
		DecisionReason:        decision.DecisionReason,
		MetricsChangeCount:    history.MetricsChangeCount,
		GenerationChangeCount: history.GenerationChangeCount,
		RequiresObservation:   decision.RequiresObservation,
		RequiresInspection:    decision.RequiresInspection,
		RequiresMeasurement:   decision.RequiresMeasurement,
		BridgeSignal:          "provenance-feedback-bridge-unknown",
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setBridgeDigest := func() {
		result.BridgeDigest = digestRevisionSelfImprovementProvenanceFeedbackBridge(result)
	}
	setBridgeDigest()

	if err := history.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-history"
		setBridgeDigest()
		return result, fmt.Errorf("provenance history is not valid: %w", err)
	}
	if err := decision.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-decision"
		setBridgeDigest()
		return result, fmt.Errorf("provenance decision is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-feedback"
		setBridgeDigest()
		return result, fmt.Errorf("legacy feedback is not valid: %w", err)
	}
	if decision.HistoryDigest != history.HistoryDigest {
		result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-history-link"
		setBridgeDigest()
		return result, fmt.Errorf("provenance decision is not linked to history")
	}

	switch history.HistorySignal {
	case "stable":
		if decision.DecisionSignal != "observe" ||
			feedback.ComparisonSignal != "stable" ||
			feedback.FeedbackSignal != "observe" ||
			decision.RequiresInspection || decision.RequiresMeasurement ||
			feedback.RequiresInspection || feedback.RequiresMeasurement {
			result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-mapping"
			setBridgeDigest()
			return result, fmt.Errorf("stable provenance does not map to stable observe feedback")
		}
	case "transitioned", "mixed":
		if decision.DecisionSignal != "inspect" ||
			feedback.ComparisonSignal != "mixed" ||
			feedback.FeedbackSignal != "inspect" ||
			!decision.RequiresInspection || !feedback.RequiresInspection ||
			decision.RequiresMeasurement != (history.MetricsChangeCount > 0 || history.GenerationChangeCount > 0) ||
			feedback.RequiresMeasurement != decision.RequiresMeasurement {
			result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-mapping"
			setBridgeDigest()
			return result, fmt.Errorf("changed provenance does not map to mixed inspect feedback")
		}
	default:
		result.MissingStage = "revision-self-improvement-provenance-feedback-bridge-signal"
		setBridgeDigest()
		return result, fmt.Errorf("provenance history signal is not recognized")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.BridgeSignal = "provenance-feedback-bridge-bound"
	setBridgeDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-feedback-bridge"
		result.BridgeSignal = "provenance-feedback-bridge-unknown"
		setBridgeDigest()
		return result, fmt.Errorf("provenance feedback bridge is not valid: %w", err)
	}
	return result, nil
}

func (b RevisionSelfImprovementProvenanceFeedbackBridgeObservation) Validate() error {
	if b.Status == "" {
		return fmt.Errorf("provenance feedback bridge status is empty")
	}
	if b.Status == "BOUND" && b.MissingStage != "" {
		return fmt.Errorf("bound provenance feedback bridge has a missing stage")
	}
	if b.Status == "UNKNOWN" && b.MissingStage == "" {
		return fmt.Errorf("unknown provenance feedback bridge has no missing stage")
	}
	for name, digest := range map[string]string{
		"history":   b.HistoryDigest,
		"decision":  b.DecisionDigest,
		"feedback":  b.FeedbackDigest,
		"window":    b.WindowDigest,
		"bridge":    b.BridgeDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("provenance feedback bridge %s digest is invalid", name)
		}
	}
	if b.HistorySignal != "stable" && b.HistorySignal != "transitioned" &&
		b.HistorySignal != "mixed" {
		return fmt.Errorf("provenance feedback bridge history signal is invalid")
	}
	if b.ComparisonSignal != "stable" && b.ComparisonSignal != "mixed" {
		return fmt.Errorf("provenance feedback bridge comparison signal is invalid")
	}
	if b.DecisionSignal != "observe" && b.DecisionSignal != "inspect" {
		return fmt.Errorf("provenance feedback bridge decision signal is invalid")
	}
	if b.FeedbackSignal != "observe" && b.FeedbackSignal != "inspect" {
		return fmt.Errorf("provenance feedback bridge feedback signal is invalid")
	}
	if b.BridgeSignal != "provenance-feedback-bridge-bound" &&
		b.BridgeSignal != "provenance-feedback-bridge-unknown" {
		return fmt.Errorf("provenance feedback bridge signal is invalid")
	}
	if b.MetricsChangeCount < 0 || b.GenerationChangeCount < 0 {
		return fmt.Errorf("provenance feedback bridge change counts are invalid")
	}
	if !b.RequiresObservation {
		return fmt.Errorf("provenance feedback bridge must require observation")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("provenance feedback bridge must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceFeedbackBridge(b) != b.BridgeDigest {
		return fmt.Errorf("provenance feedback bridge digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceFeedbackBridge(
	bridge RevisionSelfImprovementProvenanceFeedbackBridgeObservation,
) string {
	fields := []string{
		bridge.Status,
		bridge.MissingStage,
		bridge.HistoryDigest,
		bridge.DecisionDigest,
		bridge.FeedbackDigest,
		bridge.WindowDigest,
		bridge.HistorySignal,
		bridge.ComparisonSignal,
		bridge.DecisionSignal,
		bridge.FeedbackSignal,
		bridge.DecisionReason,
		strconv.Itoa(bridge.MetricsChangeCount),
		strconv.Itoa(bridge.GenerationChangeCount),
		strconv.FormatBool(bridge.RequiresObservation),
		strconv.FormatBool(bridge.RequiresInspection),
		strconv.FormatBool(bridge.RequiresMeasurement),
		bridge.BridgeSignal,
		strconv.FormatBool(bridge.NonExecuting),
		strconv.FormatBool(bridge.NonAuthorizing),
	}
	return digestString(strings.Join(fields, "|"))
}
