package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleFeedbackDecisionLink struct {
	Status            string
	MissingStage      string
	BridgeDigest      string
	DecisionDigest    string
	QuestionDigest    string
	FeedbackSignal    string
	DecisionSignal    string
	DecisionReason    string
	ReceiptSignal     string
	DecisionAligned   bool
	LinkSignal        string
	ProjectionDigest  string
	ReadOnly          bool
	NonExecuting      bool
	NonAuthorizing    bool
}

func ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionLink(
	link RevisionSelfImprovementCycleFeedbackDecisionLinkObservation,
) (LSPRevisionSelfImprovementCycleFeedbackDecisionLink, error) {
	result := LSPRevisionSelfImprovementCycleFeedbackDecisionLink{
		Status:           "UNKNOWN",
		MissingStage:     "lsp-revision-self-improvement-cycle-feedback-decision-link",
		BridgeDigest:     link.BridgeDigest,
		DecisionDigest:   link.DecisionDigest,
		QuestionDigest:   link.QuestionDigest,
		FeedbackSignal:   link.FeedbackSignal,
		DecisionSignal:   link.DecisionSignal,
		DecisionReason:   link.DecisionReason,
		ReceiptSignal:    link.ReceiptSignal,
		DecisionAligned:  link.DecisionAligned,
		LinkSignal:       link.LinkSignal,
		ReadOnly:         true,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleFeedbackDecisionLink(result)
	}
	setDigest()

	if err := link.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-feedback-decision-link-input"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision link is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-feedback-decision-link"
		setDigest()
		return result, fmt.Errorf("lsp cycle feedback decision link projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleFeedbackDecisionLink) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle feedback decision link status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle feedback decision link has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle feedback decision link has no missing stage")
	}
	for name, digest := range map[string]string{
		"bridge":      o.BridgeDigest,
		"decision":    o.DecisionDigest,
		"question":    o.QuestionDigest,
		"projection":  o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle feedback decision link %s digest is invalid", name)
		}
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "replan" && o.FeedbackSignal != "review" {
		return fmt.Errorf("lsp cycle feedback decision link feedback signal is invalid")
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("lsp cycle feedback decision link decision signal is invalid")
	}
	if o.DecisionReason == "" || o.ReceiptSignal != revisionSelfImprovementDecisionReceiptAcceptedSignal {
		return fmt.Errorf("lsp cycle feedback decision link receipt evidence is incomplete")
	}
	if o.LinkSignal != "cycle-feedback-decision-linked" || !o.DecisionAligned {
		return fmt.Errorf("lsp cycle feedback decision link is not aligned")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle feedback decision link must remain read-only and non-executing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleFeedbackDecisionLink(o) {
		return fmt.Errorf("lsp cycle feedback decision link projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleFeedbackDecisionLink(
	link LSPRevisionSelfImprovementCycleFeedbackDecisionLink,
) string {
	return digestString(strings.Join([]string{
		link.Status,
		link.MissingStage,
		link.BridgeDigest,
		link.DecisionDigest,
		link.QuestionDigest,
		link.FeedbackSignal,
		link.DecisionSignal,
		link.DecisionReason,
		link.ReceiptSignal,
		strconv.FormatBool(link.DecisionAligned),
		link.LinkSignal,
		strconv.FormatBool(link.ReadOnly),
		strconv.FormatBool(link.NonExecuting),
		strconv.FormatBool(link.NonAuthorizing),
	}, "|"))
}
