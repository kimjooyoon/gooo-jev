package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementDecisionReceiptReverseObservation struct {
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
	ProjectionDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservation(
	observation RevisionSelfImprovementDecisionReceiptReverseObservation,
) (LSPRevisionSelfImprovementDecisionReceiptReverseObservation, error) {
	result := LSPRevisionSelfImprovementDecisionReceiptReverseObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "lsp-revision-self-improvement-decision-receipt-reverse",
		DecisionReceiptDigest:    observation.DecisionReceiptDigest,
		ReverseObservationDigest: observation.ReverseObservationDigest,
		DecisionSignal:           observation.DecisionSignal,
		DecisionReason:           observation.DecisionReason,
		ReverseSignal:            observation.ReverseSignal,
		ProvenanceReverseSignal:  observation.ProvenanceReverseSignal,
		DecisionFeedbackSignal:   observation.DecisionFeedbackSignal,
		SignalsAligned:           observation.SignalsAligned,
		ReadOnly:                 true,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementDecisionReceiptReverseObservation(result)
	}
	setDigest()

	if err := observation.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-decision-receipt-reverse-input"
		setDigest()
		return result, fmt.Errorf("decision receipt reverse observation is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-decision-receipt-reverse"
		setDigest()
		return result, fmt.Errorf("lsp decision receipt reverse projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementDecisionReceiptReverseObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp decision receipt reverse status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp decision receipt reverse has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp decision receipt reverse has no missing stage")
	}
	for name, digest := range map[string]string{
		"decision receipt":    o.DecisionReceiptDigest,
		"reverse observation": o.ReverseObservationDigest,
		"projection":         o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp decision receipt reverse %s digest is invalid", name)
		}
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("lsp decision receipt reverse decision signal is invalid")
	}
	if o.DecisionReason == "" || o.ReverseSignal == "" ||
		o.ProvenanceReverseSignal == "" || o.DecisionFeedbackSignal == "" {
		return fmt.Errorf("lsp decision receipt reverse evidence is incomplete")
	}
	if o.DecisionFeedbackSignal != "decision-receipt-reverse-aligned" &&
		o.DecisionFeedbackSignal != "decision-receipt-reverse-mismatch" {
		return fmt.Errorf("lsp decision receipt reverse relationship signal is invalid")
	}
	if o.SignalsAligned != (o.DecisionFeedbackSignal == "decision-receipt-reverse-aligned") {
		return fmt.Errorf("lsp decision receipt reverse alignment is inconsistent")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp decision receipt reverse projection must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementDecisionReceiptReverseObservation(o) {
		return fmt.Errorf("lsp decision receipt reverse projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementDecisionReceiptReverseObservation(
	observation LSPRevisionSelfImprovementDecisionReceiptReverseObservation,
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
		strconv.FormatBool(observation.ReadOnly),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}