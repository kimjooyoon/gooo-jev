package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput describes
// one bounded JEV decision without coupling the language to a model provider.
type RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput struct {
	QuestionID      string
	QuestionKind    string
	DecisionDigest  string
	EvidenceDigest  string
	ConfidenceMilli int64
}

// RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation records a
// typed choice, score, or noul signal as a provenance-bearing observation.
type RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation struct {
	Status                     string
	MissingStage               string
	QuestionID                 string
	QuestionKind               string
	DecisionDigest             string
	EvidenceDigest             string
	ConfidenceMilli            int64
	DecisionSignal             string
	DecisionObservationDigest  string
	NonExecuting               bool
	NonAuthorizing             bool
}

func ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
	input RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput,
) (RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation, error) {
	result := RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-typed-decision-signal",
		QuestionID:                input.QuestionID,
		QuestionKind:              input.QuestionKind,
		DecisionDigest:            input.DecisionDigest,
		EvidenceDigest:            input.EvidenceDigest,
		ConfidenceMilli:           input.ConfidenceMilli,
		DecisionSignal:            "jev-typed-decision-unknown",
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.DecisionObservationDigest = digestRevisionSelfImprovementCycleJEVTypedDecisionSignal(result)
	}
	setDigest()

	if err := input.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-typed-decision-signal-input"
		setDigest()
		return result, fmt.Errorf("jev typed decision signal is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.DecisionSignal = "jev-" + input.QuestionKind + "-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-typed-decision-signal"
		result.DecisionSignal = "jev-typed-decision-unknown"
		setDigest()
		return result, fmt.Errorf("jev typed decision observation is not valid: %w", err)
	}
	return result, nil
}

func (i RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput) Validate() error {
	if strings.TrimSpace(i.QuestionID) == "" || strings.Contains(i.QuestionID, "|") {
		return fmt.Errorf("jev typed decision question id is invalid")
	}
	if i.QuestionKind != "choice" && i.QuestionKind != "score" && i.QuestionKind != "noul" {
		return fmt.Errorf("jev typed decision question kind is invalid")
	}
	if !validDigest(i.DecisionDigest) || !validDigest(i.EvidenceDigest) {
		return fmt.Errorf("jev typed decision digest is invalid")
	}
	if i.QuestionKind == "noul" {
		if i.ConfidenceMilli != -1 {
			return fmt.Errorf("jev noul must not carry a separate confidence value")
		}
	} else if i.ConfidenceMilli < 0 || i.ConfidenceMilli > 1000 {
		return fmt.Errorf("jev typed decision confidence must be between 0 and 1000 milli")
	}
	return nil
}

func (o RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev typed decision observation status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev typed decision observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev typed decision observation has no missing stage")
	}
	if err := (RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
		QuestionID:      o.QuestionID,
		QuestionKind:    o.QuestionKind,
		DecisionDigest:  o.DecisionDigest,
		EvidenceDigest:  o.EvidenceDigest,
		ConfidenceMilli: o.ConfidenceMilli,
	}).Validate(); err != nil {
		return fmt.Errorf("jev typed decision observation input is invalid: %w", err)
	}
	if o.DecisionSignal != "jev-choice-observed" &&
		o.DecisionSignal != "jev-score-observed" &&
		o.DecisionSignal != "jev-noul-observed" &&
		o.DecisionSignal != "jev-typed-decision-unknown" {
		return fmt.Errorf("jev typed decision signal is invalid")
	}
	if o.Status == "BOUND" && o.DecisionSignal != "jev-"+o.QuestionKind+"-observed" {
		return fmt.Errorf("bound jev typed decision signal does not match its question kind")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev typed decision observation must remain non-executing and non-authorizing")
	}
	if o.DecisionObservationDigest != digestRevisionSelfImprovementCycleJEVTypedDecisionSignal(o) {
		return fmt.Errorf("jev typed decision observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVTypedDecisionSignal(
	observation RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.QuestionID,
		observation.QuestionKind,
		observation.DecisionDigest,
		observation.EvidenceDigest,
		strconv.FormatInt(observation.ConfidenceMilli, 10),
		observation.DecisionSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
