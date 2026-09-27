package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary struct {
	Status                   string
	MissingStage             string
	CycleDigest              string
	BridgeDigest             string
	EvidenceFeedbackDigest   string
	DecisionBoundaryDigest   string
	CoverageSignal           string
	FeedbackSignal           string
	FeedbackReason           string
	DecisionSignal           string
	BoundarySignal           string
	EvidenceFeedbackAligned  bool
	BoundaryAligned          bool
	ProjectionDigest         string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(
	boundary RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary{
		Status:                 "UNKNOWN",
		MissingStage:           "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary",
		CycleDigest:            boundary.CycleDigest,
		BridgeDigest:           boundary.BridgeDigest,
		EvidenceFeedbackDigest: boundary.EvidenceFeedbackDigest,
		DecisionBoundaryDigest: boundary.DecisionBoundaryDigest,
		CoverageSignal:         boundary.CoverageSignal,
		FeedbackSignal:         boundary.FeedbackSignal,
		FeedbackReason:         boundary.FeedbackReason,
		DecisionSignal:         boundary.DecisionSignal,
		BoundarySignal:         boundary.BoundarySignal,
		EvidenceFeedbackAligned: boundary.EvidenceFeedbackAligned,
		BoundaryAligned:        boundary.BoundaryAligned,
		ReadOnly:               true,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(result)
	}
	setDigest()

	if err := boundary.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary is not valid: %w", err)
	}
	if boundary.Status != "BOUND" || boundary.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary"
		setDigest()
		return result, fmt.Errorf("lsp cycle evidence feedback decision boundary projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle evidence feedback decision boundary has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle evidence feedback decision boundary has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle":             o.CycleDigest,
		"bridge":            o.BridgeDigest,
		"evidence feedback": o.EvidenceFeedbackDigest,
		"decision boundary": o.DecisionBoundaryDigest,
		"projection":        o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle evidence feedback decision boundary %s digest is invalid", name)
		}
	}
	if o.CoverageSignal != "evidence-complete" && o.CoverageSignal != "evidence-incomplete" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary coverage signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary feedback reason is empty")
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary decision signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.CoverageSignal != "evidence-complete" ||
			!o.EvidenceFeedbackAligned ||
			!o.BoundaryAligned ||
			o.BoundarySignal != "cycle-evidence-feedback-decision-boundary-linked") {
		return fmt.Errorf("bound lsp cycle evidence feedback decision boundary is not aligned")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(o) {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(
	boundary LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary,
) string {
	return digestString(strings.Join([]string{
		boundary.Status,
		boundary.MissingStage,
		boundary.CycleDigest,
		boundary.BridgeDigest,
		boundary.EvidenceFeedbackDigest,
		boundary.DecisionBoundaryDigest,
		boundary.CoverageSignal,
		boundary.FeedbackSignal,
		boundary.FeedbackReason,
		boundary.DecisionSignal,
		boundary.BoundarySignal,
		strconv.FormatBool(boundary.EvidenceFeedbackAligned),
		strconv.FormatBool(boundary.BoundaryAligned),
		strconv.FormatBool(boundary.ReadOnly),
		strconv.FormatBool(boundary.NonExecuting),
		strconv.FormatBool(boundary.NonAuthorizing),
	}, "|"))
}
