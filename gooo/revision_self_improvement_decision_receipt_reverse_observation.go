package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementDecisionReceiptReverseObservation struct {
	Status                    string
	MissingStage              string
	DecisionReceiptDigest     string
	ReverseObservationDigest  string
	DecisionSignal            string
	DecisionReason            string
	ReverseSignal             string
	ProvenanceReverseSignal   string
	DecisionFeedbackSignal    string
	SignalsAligned            bool
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(
	receipt RevisionSelfImprovementDecisionReceiptObservation,
	reverse RevisionSelfImprovementProvenanceReverseObservation,
) (RevisionSelfImprovementDecisionReceiptReverseObservation, error) {
	result := RevisionSelfImprovementDecisionReceiptReverseObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "revision-self-improvement-decision-receipt-reverse",
		DecisionReceiptDigest:    receipt.ObservationDigest,
		ReverseObservationDigest: reverse.ObservationDigest,
		DecisionSignal:           receipt.DecisionSignal,
		DecisionReason:           receipt.DecisionReason,
		ReverseSignal:            reverse.ReverseSignal,
		ProvenanceReverseSignal:  reverse.ProvenanceReverseSignal,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementDecisionReceiptReverseObservation(result)
	}
	setDigest()

	if err := validateRevisionSelfImprovementDecisionReceiptObservation(receipt); err != nil {
		result.MissingStage = "revision-self-improvement-decision-receipt-reverse-receipt"
		setDigest()
		return result, fmt.Errorf("decision receipt is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-decision-receipt-reverse-observation"
		setDigest()
		return result, fmt.Errorf("provenance reverse observation is not valid: %w", err)
	}

	result.SignalsAligned = receipt.DecisionAligned && reverse.SignalsAligned &&
		((receipt.DecisionSignal == revisionSelfImprovementDecisionReceiptAllowNextIteration &&
			reverse.ProvenanceReverseSignal == "provenance-reverse-aligned") ||
			(receipt.DecisionSignal == revisionSelfImprovementDecisionReceiptRequireReplan &&
				reverse.ProvenanceReverseSignal == "provenance-reverse-mismatch"))
	if result.SignalsAligned {
		result.DecisionFeedbackSignal = "decision-receipt-reverse-aligned"
	} else {
		result.DecisionFeedbackSignal = "decision-receipt-reverse-mismatch"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-decision-receipt-reverse"
		result.DecisionFeedbackSignal = "decision-receipt-reverse-mismatch"
		setDigest()
		return result, fmt.Errorf("decision receipt reverse observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementDecisionReceiptReverseObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("decision receipt reverse observation status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound decision receipt reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown decision receipt reverse observation has no missing stage")
	}
	for name, digest := range map[string]string{
		"decision receipt":    o.DecisionReceiptDigest,
		"reverse observation": o.ReverseObservationDigest,
		"observation":         o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("decision receipt reverse observation %s digest is invalid", name)
		}
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("decision receipt reverse decision signal is invalid")
	}
	if o.DecisionReason == "" || o.ReverseSignal == "" || o.ProvenanceReverseSignal == "" {
		return fmt.Errorf("decision receipt reverse evidence is incomplete")
	}
	if o.DecisionFeedbackSignal != "decision-receipt-reverse-aligned" &&
		o.DecisionFeedbackSignal != "decision-receipt-reverse-mismatch" {
		return fmt.Errorf("decision receipt reverse relationship signal is invalid")
	}
	if o.SignalsAligned != (o.DecisionFeedbackSignal == "decision-receipt-reverse-aligned") {
		return fmt.Errorf("decision receipt reverse alignment is inconsistent")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("decision receipt reverse observation must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementDecisionReceiptReverseObservation(o) {
		return fmt.Errorf("decision receipt reverse observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementDecisionReceiptReverseObservation(
	observation RevisionSelfImprovementDecisionReceiptReverseObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.DecisionReceiptDigest,
		observation.ReverseObservationDigest,
		observation.DecisionSignal,
		observation.DecisionReason,
		observation.ReverseSignal,
		observation.ProvenanceReverseSignal,
		observation.DecisionFeedbackSignal,
		strconv.FormatBool(observation.SignalsAligned),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}