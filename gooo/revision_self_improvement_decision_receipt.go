package gooo

import (
	"errors"
	"strconv"
	"strings"
)

const (
	revisionSelfImprovementDecisionReceiptBoundStatus = "BOUND"
	revisionSelfImprovementDecisionReceiptUnknownStatus = "UNKNOWN"

	revisionSelfImprovementDecisionReceiptAllowNextIteration = "allow-next-iteration"
	revisionSelfImprovementDecisionReceiptRequireReplan = "require-replan"
	revisionSelfImprovementDecisionReceiptAcceptedSignal = "provenance-decision-accepted"

	revisionSelfImprovementDecisionReceiptMissingInput = "revision-self-improvement-provenance-decision-receipt-input"
	revisionSelfImprovementDecisionReceiptMissingQuestion = "revision-self-improvement-provenance-decision-receipt-question"
	revisionSelfImprovementDecisionReceiptMissingReason = "revision-self-improvement-provenance-decision-receipt-reason"
	revisionSelfImprovementDecisionReceiptMissingSignal = "revision-self-improvement-provenance-decision-receipt-signal"
	revisionSelfImprovementDecisionReceiptMissingAlignment = "revision-self-improvement-provenance-decision-receipt-alignment"
	revisionSelfImprovementDecisionReceiptMissingReceipt = "revision-self-improvement-provenance-decision-receipt"

)

type RevisionSelfImprovementDecisionReceiptObservation struct {
	Status                 string
	MissingStage           string
	InputObservationDigest string
	QuestionDigest         string
	DecisionDigest         string
	DecisionSignal         string
	DecisionReason         string
	ReceiptSignal          string
	DecisionAligned        bool
	ObservationDigest      string
	NonExecuting           bool
	NonAuthorizing         bool
}

func ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
	input RevisionSelfImprovementProvenanceNextIterationObservation,
	questionDigest string,
	decisionSignal string,
	decisionReason string,
) (RevisionSelfImprovementDecisionReceiptObservation, error) {
	observation := RevisionSelfImprovementDecisionReceiptObservation{
		Status:                 revisionSelfImprovementDecisionReceiptUnknownStatus,
		MissingStage:           revisionSelfImprovementDecisionReceiptMissingInput,
		InputObservationDigest: input.ObservationDigest,
		QuestionDigest:         questionDigest,
		DecisionSignal:         decisionSignal,
		DecisionReason:         decisionReason,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	if err := input.Validate(); err != nil {
		observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
		return observation, errors.Join(errors.New(revisionSelfImprovementDecisionReceiptMissingInput), err)
	}
	if questionDigest == "" {
		observation.MissingStage = revisionSelfImprovementDecisionReceiptMissingQuestion
		observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
		return observation, errors.New(revisionSelfImprovementDecisionReceiptMissingQuestion)
	}
	if decisionReason == "" {
		observation.MissingStage = revisionSelfImprovementDecisionReceiptMissingReason
		observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
		return observation, errors.New(revisionSelfImprovementDecisionReceiptMissingReason)
	}

	expectedSignal := revisionSelfImprovementDecisionReceiptExpectedSignal(input)
	if expectedSignal == "" || !revisionSelfImprovementDecisionReceiptValidSignal(decisionSignal) {
		observation.MissingStage = revisionSelfImprovementDecisionReceiptMissingSignal
		observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
		return observation, errors.New(revisionSelfImprovementDecisionReceiptMissingSignal)
	}
	if decisionSignal != expectedSignal {
		observation.MissingStage = revisionSelfImprovementDecisionReceiptMissingAlignment
		observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
		return observation, errors.New(revisionSelfImprovementDecisionReceiptMissingAlignment)
	}

	observation.Status = revisionSelfImprovementDecisionReceiptBoundStatus
	observation.MissingStage = ""
	observation.DecisionAligned = true
	observation.ReceiptSignal = revisionSelfImprovementDecisionReceiptAcceptedSignal
	observation.DecisionDigest = digestString(strings.Join([]string{
		observation.InputObservationDigest,
		observation.QuestionDigest,
		observation.DecisionSignal,
		observation.DecisionReason,
	}, "|"))
	observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
	if err := validateRevisionSelfImprovementDecisionReceiptObservation(observation); err != nil {
		observation.Status = revisionSelfImprovementDecisionReceiptUnknownStatus
		observation.MissingStage = revisionSelfImprovementDecisionReceiptMissingReceipt
		observation.DecisionAligned = false
		observation.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(observation)
		return observation, errors.Join(errors.New(revisionSelfImprovementDecisionReceiptMissingReceipt), err)
	}
	return observation, nil
}

func revisionSelfImprovementDecisionReceiptExpectedSignal(
	input RevisionSelfImprovementProvenanceNextIterationObservation,
) string {
	switch input.TransitionSignal {
	case "provenance-next-iteration-aligned":
		return revisionSelfImprovementDecisionReceiptAllowNextIteration
	case "provenance-next-iteration-mismatch":
		return revisionSelfImprovementDecisionReceiptRequireReplan
	default:
		return ""
	}
}

func revisionSelfImprovementDecisionReceiptValidSignal(signal string) bool {
	return signal == revisionSelfImprovementDecisionReceiptAllowNextIteration ||
		signal == revisionSelfImprovementDecisionReceiptRequireReplan
}

func validateRevisionSelfImprovementDecisionReceiptObservation(
	observation RevisionSelfImprovementDecisionReceiptObservation,
) error {
	if observation.Status != revisionSelfImprovementDecisionReceiptBoundStatus &&
		observation.Status != revisionSelfImprovementDecisionReceiptUnknownStatus {
		return errors.New("invalid decision receipt status")
	}
	if observation.Status == revisionSelfImprovementDecisionReceiptBoundStatus && observation.MissingStage != "" {
		return errors.New("bound decision receipt has missing stage")
	}
	if observation.Status == revisionSelfImprovementDecisionReceiptUnknownStatus && observation.MissingStage == "" {
		return errors.New("unknown decision receipt has no missing stage")
	}
	if observation.Status == revisionSelfImprovementDecisionReceiptBoundStatus &&
		(!observation.DecisionAligned ||
			observation.ReceiptSignal != revisionSelfImprovementDecisionReceiptAcceptedSignal) {
		return errors.New("bound decision receipt is not aligned")
	}
	if observation.InputObservationDigest == "" ||
		observation.QuestionDigest == "" ||
		observation.DecisionSignal == "" ||
		observation.DecisionReason == "" ||
		observation.ObservationDigest == "" {
		return errors.New("decision receipt is missing required evidence")
	}
	if observation.Status == revisionSelfImprovementDecisionReceiptBoundStatus &&
		observation.DecisionDigest == "" {
		return errors.New("bound decision receipt has no decision digest")
	}
	if !observation.NonExecuting || !observation.NonAuthorizing {
		return errors.New("decision receipt must remain non-executing and non-authorizing")
	}
	if observation.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		observation.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return errors.New("invalid decision receipt signal")
	}
	if observation.Status == revisionSelfImprovementDecisionReceiptBoundStatus &&
		observation.ObservationDigest != digestRevisionSelfImprovementDecisionReceipt(observation) {
		return errors.New("decision receipt digest mismatch")
	}
	return nil
}

func digestRevisionSelfImprovementDecisionReceipt(
	observation RevisionSelfImprovementDecisionReceiptObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.InputObservationDigest,
		observation.QuestionDigest,
		observation.DecisionDigest,
		observation.DecisionSignal,
		observation.DecisionReason,
		observation.ReceiptSignal,
		strconv.FormatBool(observation.DecisionAligned),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}