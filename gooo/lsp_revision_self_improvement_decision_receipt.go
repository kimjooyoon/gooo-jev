package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementDecisionReceiptObservation struct {
	Status                 string
	MissingStage           string
	DecisionReceiptDigest  string
	QuestionDigest         string
	DecisionDigest         string
	DecisionSignal         string
	DecisionReason         string
	ReceiptSignal          string
	DecisionAligned        bool
	ProjectionDigest       string
	ReadOnly               bool
	NonExecuting           bool
	NonAuthorizing         bool
}

func ObserveLSPRevisionSelfImprovementDecisionReceipt(
	receipt RevisionSelfImprovementDecisionReceiptObservation,
) (LSPRevisionSelfImprovementDecisionReceiptObservation, error) {
	result := LSPRevisionSelfImprovementDecisionReceiptObservation{
		Status:                "UNKNOWN",
		MissingStage:          "lsp-revision-self-improvement-decision-receipt",
		DecisionReceiptDigest: receipt.ObservationDigest,
		QuestionDigest:        receipt.QuestionDigest,
		DecisionDigest:        receipt.DecisionDigest,
		DecisionSignal:        receipt.DecisionSignal,
		DecisionReason:        receipt.DecisionReason,
		ReceiptSignal:         receipt.ReceiptSignal,
		DecisionAligned:       receipt.DecisionAligned,
		ReadOnly:              true,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementDecisionReceipt(result)
	}
	setDigest()

	if err := validateRevisionSelfImprovementDecisionReceiptObservation(receipt); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-decision-receipt-input"
		setDigest()
		return result, fmt.Errorf("decision receipt is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-decision-receipt"
		setDigest()
		return result, fmt.Errorf("lsp decision receipt projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementDecisionReceiptObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp decision receipt status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp decision receipt has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp decision receipt has no missing stage")
	}
	for name, digest := range map[string]string{
		"decision receipt": o.DecisionReceiptDigest,
		"question":         o.QuestionDigest,
		"decision":         o.DecisionDigest,
		"projection":       o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp decision receipt %s digest is invalid", name)
		}
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("lsp decision receipt signal is invalid")
	}
	if o.ReceiptSignal != revisionSelfImprovementDecisionReceiptAcceptedSignal {
		return fmt.Errorf("lsp decision receipt acceptance signal is invalid")
	}
	if o.DecisionReason == "" || !o.DecisionAligned {
		return fmt.Errorf("lsp decision receipt alignment is invalid")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp decision receipt must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementDecisionReceipt(o) {
		return fmt.Errorf("lsp decision receipt projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementDecisionReceipt(
	observation LSPRevisionSelfImprovementDecisionReceiptObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.DecisionReceiptDigest,
		observation.QuestionDigest,
		observation.DecisionDigest,
		observation.DecisionSignal,
		observation.DecisionReason,
		observation.ReceiptSignal,
		strconv.FormatBool(observation.DecisionAligned),
		strconv.FormatBool(observation.ReadOnly),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}