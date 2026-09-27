package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleObservation struct {
	Status                         string
	MissingStage                   string
	NextIterationObservationDigest string
	DecisionReceiptDigest          string
	DecisionReceiptReverseDigest   string
	SourceDigest                   string
	NextSourceDigest               string
	IRDigest                       string
	NextIRDigest                   string
	HistoryDigest                  string
	FeedbackDigest                 string
	ReverseObservationDigest       string
	TransitionSignal               string
	DecisionSignal                 string
	ReverseSignal                  string
	CycleSignal                    string
	SignalsAligned                 bool
	ObservationCount               int
	ObservationDigest              string
	NonExecuting                   bool
	NonAuthorizing                 bool
}

func ObserveRevisionSelfImprovementCycle(
	nextIteration RevisionSelfImprovementProvenanceNextIterationObservation,
	receipt RevisionSelfImprovementDecisionReceiptObservation,
	reverse RevisionSelfImprovementDecisionReceiptReverseObservation,
) (RevisionSelfImprovementCycleObservation, error) {
	result := RevisionSelfImprovementCycleObservation{
		Status:                          "UNKNOWN",
		MissingStage:                    "revision-self-improvement-cycle",
		NextIterationObservationDigest: nextIteration.ObservationDigest,
		DecisionReceiptDigest:           receipt.ObservationDigest,
		DecisionReceiptReverseDigest:    reverse.ObservationDigest,
		SourceDigest:                    nextIteration.SourceDigest,
		NextSourceDigest:                nextIteration.NextSourceDigest,
		IRDigest:                        nextIteration.IRDigest,
		NextIRDigest:                    nextIteration.NextIRDigest,
		HistoryDigest:                   nextIteration.HistoryDigest,
		FeedbackDigest:                  nextIteration.FeedbackDigest,
		ReverseObservationDigest:        reverse.ReverseObservationDigest,
		TransitionSignal:                nextIteration.TransitionSignal,
		DecisionSignal:                  receipt.DecisionSignal,
		ReverseSignal:                   reverse.ReverseSignal,
		ObservationCount:                nextIteration.ObservationCount,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycle(result)
	}
	setDigest()

	if err := nextIteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-next-iteration"
		setDigest()
		return result, fmt.Errorf("next iteration observation is not valid: %w", err)
	}
	if err := validateRevisionSelfImprovementDecisionReceiptObservation(receipt); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-decision-receipt"
		setDigest()
		return result, fmt.Errorf("decision receipt is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-reverse"
		setDigest()
		return result, fmt.Errorf("decision receipt reverse observation is not valid: %w", err)
	}
	if receipt.InputObservationDigest != nextIteration.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-decision-receipt-link"
		setDigest()
		return result, fmt.Errorf("decision receipt is not linked to the next iteration observation")
	}
	if receipt.DecisionSignal != revisionSelfImprovementDecisionReceiptExpectedSignal(nextIteration) {
		result.MissingStage = "revision-self-improvement-cycle-decision-link"
		setDigest()
		return result, fmt.Errorf("decision receipt signal is not linked to the next iteration transition")
	}
	if reverse.DecisionReceiptDigest != receipt.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-reverse-link"
		setDigest()
		return result, fmt.Errorf("reverse observation is not linked to the decision receipt")
	}

	result.SignalsAligned = nextIteration.SignalsAligned &&
		receipt.DecisionAligned &&
		reverse.SignalsAligned
	if result.SignalsAligned {
		result.CycleSignal = "self-improvement-cycle-closed"
	} else {
		result.CycleSignal = "self-improvement-cycle-replan-required"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle"
		result.CycleSignal = "self-improvement-cycle-replan-required"
		setDigest()
		return result, fmt.Errorf("self-improvement cycle observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("self-improvement cycle status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound self-improvement cycle has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown self-improvement cycle has no missing stage")
	}
	for name, digest := range map[string]string{
		"next iteration":           o.NextIterationObservationDigest,
		"decision receipt":         o.DecisionReceiptDigest,
		"decision receipt reverse": o.DecisionReceiptReverseDigest,
		"source":                   o.SourceDigest,
		"next source":              o.NextSourceDigest,
		"ir":                       o.IRDigest,
		"next ir":                  o.NextIRDigest,
		"history":                  o.HistoryDigest,
		"feedback":                 o.FeedbackDigest,
		"reverse observation":      o.ReverseObservationDigest,
		"observation":              o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("self-improvement cycle %s digest is invalid", name)
		}
	}
	if o.TransitionSignal != "provenance-next-iteration-aligned" &&
		o.TransitionSignal != "provenance-next-iteration-mismatch" {
		return fmt.Errorf("self-improvement cycle transition signal is invalid")
	}
	if o.DecisionSignal != "allow-next-iteration" &&
		o.DecisionSignal != "require-replan" {
		return fmt.Errorf("self-improvement cycle decision signal is invalid")
	}
	if o.ReverseSignal == "" {
		return fmt.Errorf("self-improvement cycle reverse signal is empty")
	}
	if o.CycleSignal != "self-improvement-cycle-closed" &&
		o.CycleSignal != "self-improvement-cycle-replan-required" {
		return fmt.Errorf("self-improvement cycle signal is invalid")
	}
	if o.SignalsAligned != (o.CycleSignal == "self-improvement-cycle-closed") {
		return fmt.Errorf("self-improvement cycle alignment is inconsistent")
	}
	if o.ObservationCount < 1 {
		return fmt.Errorf("self-improvement cycle observation count must be positive")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("self-improvement cycle must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycle(o) {
		return fmt.Errorf("self-improvement cycle digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycle(o RevisionSelfImprovementCycleObservation) string {
	return digestString(strings.Join([]string{
		o.Status,
		o.MissingStage,
		o.NextIterationObservationDigest,
		o.DecisionReceiptDigest,
		o.DecisionReceiptReverseDigest,
		o.SourceDigest,
		o.NextSourceDigest,
		o.IRDigest,
		o.NextIRDigest,
		o.HistoryDigest,
		o.FeedbackDigest,
		o.ReverseObservationDigest,
		o.TransitionSignal,
		o.DecisionSignal,
		o.ReverseSignal,
		o.CycleSignal,
		strconv.FormatBool(o.SignalsAligned),
		strconv.Itoa(o.ObservationCount),
		strconv.FormatBool(o.NonExecuting),
		strconv.FormatBool(o.NonAuthorizing),
	}, "|"))
}