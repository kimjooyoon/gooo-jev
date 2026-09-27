package gooo

import (
	"fmt"
	"strings"
)

// LSPRevisionSelfImprovementCycleJEVTypedDecisionSignalObservation carries a
// typed JEV decision into an editor-facing, read-only provenance projection.
type LSPRevisionSelfImprovementCycleJEVTypedDecisionSignalObservation struct {
	Status                    string
	MissingStage              string
	QuestionID                string
	QuestionKind              string
	DecisionDigest             string
	EvidenceDigest             string
	ConfidenceMilli            int64
	DecisionSignal             string
	DecisionObservationDigest string
	ProjectionSignal          string
	ProjectionDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(
	signal RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation,
) (LSPRevisionSelfImprovementCycleJEVTypedDecisionSignalObservation, error) {
	result := LSPRevisionSelfImprovementCycleJEVTypedDecisionSignalObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "lsp-revision-self-improvement-cycle-jev-typed-decision-signal",
		QuestionID:                signal.QuestionID,
		QuestionKind:               signal.QuestionKind,
		DecisionDigest:             signal.DecisionDigest,
		EvidenceDigest:             signal.EvidenceDigest,
		ConfidenceMilli:            signal.ConfidenceMilli,
		DecisionSignal:             signal.DecisionSignal,
		DecisionObservationDigest: signal.DecisionObservationDigest,
		ProjectionSignal:          "jev-typed-decision-lsp-unknown",
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(result)
	}
	setDigest()

	if err := signal.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-typed-decision-signal-input"
		setDigest()
		return result, fmt.Errorf("jev typed decision signal is not valid for lsp projection: %w", err)
	}
	if signal.Status != "BOUND" || signal.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-typed-decision-signal-input-status"
		setDigest()
		return result, fmt.Errorf("jev typed decision signal is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ProjectionSignal = "jev-typed-decision-lsp-projected"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-typed-decision-signal"
		result.ProjectionSignal = "jev-typed-decision-lsp-unknown"
		setDigest()
		return result, fmt.Errorf("jev typed decision lsp projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleJEVTypedDecisionSignalObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev typed decision lsp status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev typed decision lsp projection has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev typed decision lsp projection has no missing stage")
	}
	if err := (RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
		QuestionID:      o.QuestionID,
		QuestionKind:     o.QuestionKind,
		DecisionDigest:   o.DecisionDigest,
		EvidenceDigest:  o.EvidenceDigest,
		ConfidenceMilli: o.ConfidenceMilli,
	}).Validate(); err != nil {
		return fmt.Errorf("jev typed decision lsp input is invalid: %w", err)
	}
	if o.DecisionSignal != "jev-choice-observed" &&
		o.DecisionSignal != "jev-score-observed" &&
		o.DecisionSignal != "jev-noul-observed" &&
		o.DecisionSignal != "jev-typed-decision-unknown" {
		return fmt.Errorf("jev typed decision lsp signal is invalid")
	}
	if o.Status == "BOUND" && o.DecisionSignal != "jev-"+o.QuestionKind+"-observed" {
		return fmt.Errorf("bound jev typed decision lsp signal does not match its question kind")
	}
	if o.DecisionObservationDigest != digestRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation{
			Status:                    o.Status,
			MissingStage:              o.MissingStage,
			QuestionID:                o.QuestionID,
			QuestionKind:              o.QuestionKind,
			DecisionDigest:            o.DecisionDigest,
			EvidenceDigest:            o.EvidenceDigest,
			ConfidenceMilli:           o.ConfidenceMilli,
			DecisionSignal:             o.DecisionSignal,
			NonExecuting:              o.NonExecuting,
			NonAuthorizing:            o.NonAuthorizing,
		},
	) {
		return fmt.Errorf("jev typed decision lsp source observation digest does not match its fields")
	}
	if o.ProjectionSignal != "jev-typed-decision-lsp-projected" &&
		o.ProjectionSignal != "jev-typed-decision-lsp-unknown" {
		return fmt.Errorf("jev typed decision lsp projection signal is invalid")
	}
	if o.Status == "BOUND" && o.ProjectionSignal != "jev-typed-decision-lsp-projected" {
		return fmt.Errorf("bound jev typed decision lsp projection has an unknown signal")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev typed decision lsp projection must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(o) {
		return fmt.Errorf("jev typed decision lsp projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(
	observation LSPRevisionSelfImprovementCycleJEVTypedDecisionSignalObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.QuestionID,
		observation.QuestionKind,
		observation.DecisionDigest,
		observation.EvidenceDigest,
		fmt.Sprintf("%d", observation.ConfidenceMilli),
		observation.DecisionSignal,
		observation.DecisionObservationDigest,
		observation.ProjectionSignal,
		fmt.Sprintf("%t", observation.ReadOnly),
		fmt.Sprintf("%t", observation.NonExecuting),
		fmt.Sprintf("%t", observation.NonAuthorizing),
	}, "|"))
}
