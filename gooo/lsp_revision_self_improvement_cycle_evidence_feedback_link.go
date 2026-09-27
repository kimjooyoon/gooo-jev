package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedback struct {
	Status                   string
	MissingStage             string
	BridgeDigest             string
	CoverageMetricDigest     string
	FeedbackDigest           string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	CoverageSignal           string
	FeedbackSignal           string
	FeedbackReason           string
	RequiresObservation      bool
	RequiresReview           bool
	RequiresReplan           bool
	RequiresMeasurement      bool
	FeedbackAligned          bool
	LinkSignal               string
	ProjectionDigest         string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(
	link RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedback, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedback{
		Status:                   "UNKNOWN",
		MissingStage:             "lsp-revision-self-improvement-cycle-evidence-coverage-feedback",
		BridgeDigest:             link.BridgeDigest,
		CoverageMetricDigest:     link.CoverageMetricDigest,
		FeedbackDigest:           link.FeedbackDigest,
		SourceDigest:             link.SourceDigest,
		CandidateSourceDigest:    link.CandidateSourceDigest,
		GeneratedIRDigest:        link.GeneratedIRDigest,
		ReverseObservationDigest: link.ReverseObservationDigest,
		CoverageSignal:           link.CoverageSignal,
		FeedbackSignal:           link.FeedbackSignal,
		FeedbackReason:           link.FeedbackReason,
		RequiresObservation:      link.RequiresObservation,
		RequiresReview:           link.RequiresReview,
		RequiresReplan:           link.RequiresReplan,
		RequiresMeasurement:       link.RequiresMeasurement,
		FeedbackAligned:          link.FeedbackAligned,
		LinkSignal:               link.LinkSignal,
		ReadOnly:                 true,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(result)
	}
	setDigest()

	if err := link.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage feedback link is not valid: %w", err)
	}
	if link.Status != "BOUND" || link.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage feedback link is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback"
		setDigest()
		return result, fmt.Errorf("lsp cycle evidence coverage feedback projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedback) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle evidence coverage feedback status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle evidence coverage feedback has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle evidence coverage feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"bridge":              o.BridgeDigest,
		"coverage metric":    o.CoverageMetricDigest,
		"feedback":           o.FeedbackDigest,
		"source":              o.SourceDigest,
		"candidate source":   o.CandidateSourceDigest,
		"generated ir":       o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"projection":         o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle evidence coverage feedback %s digest is invalid", name)
		}
	}
	if o.CoverageSignal != "evidence-complete" && o.CoverageSignal != "evidence-incomplete" {
		return fmt.Errorf("lsp cycle evidence coverage feedback coverage signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("lsp cycle evidence coverage feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("lsp cycle evidence coverage feedback reason is empty")
	}
	if o.Status == "BOUND" {
		if o.CoverageSignal != "evidence-complete" ||
			!o.FeedbackAligned ||
			o.LinkSignal != "cycle-evidence-feedback-linked" {
			return fmt.Errorf("bound lsp cycle evidence coverage feedback is not linked")
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle evidence coverage feedback must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(o) {
		return fmt.Errorf("lsp cycle evidence coverage feedback projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedback(
	link LSPRevisionSelfImprovementCycleEvidenceCoverageFeedback,
) string {
	return digestString(strings.Join([]string{
		link.Status,
		link.MissingStage,
		link.BridgeDigest,
		link.CoverageMetricDigest,
		link.FeedbackDigest,
		link.SourceDigest,
		link.CandidateSourceDigest,
		link.GeneratedIRDigest,
		link.ReverseObservationDigest,
		link.CoverageSignal,
		link.FeedbackSignal,
		link.FeedbackReason,
		strconv.FormatBool(link.RequiresObservation),
		strconv.FormatBool(link.RequiresReview),
		strconv.FormatBool(link.RequiresReplan),
		strconv.FormatBool(link.RequiresMeasurement),
		strconv.FormatBool(link.FeedbackAligned),
		link.LinkSignal,
		strconv.FormatBool(link.ReadOnly),
		strconv.FormatBool(link.NonExecuting),
		strconv.FormatBool(link.NonAuthorizing),
	}, "|"))
}
