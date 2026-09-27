package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback struct {
	Status              string
	MissingStage        string
	BoundaryMetricDigest string
	MetricName          string
	MetricSignal        string
	FeedbackSignal      string
	FeedbackReason      string
	RequiresObservation bool
	RequiresReview      bool
	RequiresReplan      bool
	RequiresMeasurement bool
	FeedbackDigest      string
	ProjectionDigest    string
	ReadOnly            bool
	NonExecuting        bool
	NonAuthorizing      bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(
	feedback RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback{
		Status:              "UNKNOWN",
		MissingStage:        "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback",
		BoundaryMetricDigest: feedback.BoundaryMetricDigest,
		MetricName:          feedback.MetricName,
		MetricSignal:        feedback.MetricSignal,
		FeedbackSignal:      feedback.FeedbackSignal,
		FeedbackReason:      feedback.FeedbackReason,
		RequiresObservation: feedback.RequiresObservation,
		RequiresReview:      feedback.RequiresReview,
		RequiresReplan:      feedback.RequiresReplan,
		RequiresMeasurement: feedback.RequiresMeasurement,
		FeedbackDigest:      feedback.FeedbackDigest,
		ReadOnly:            true,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(result)
	}
	setDigest()

	if err := feedback.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric feedback is not valid: %w", err)
	}
	if feedback.Status != "BOUND" || feedback.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric feedback is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback"
		setDigest()
		return result, fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle evidence feedback decision boundary metric feedback has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle evidence feedback decision boundary metric feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"boundary metric": o.BoundaryMetricDigest,
		"feedback":        o.FeedbackDigest,
		"projection":      o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback metric name is invalid")
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "replan" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback reason is empty")
	}
	if !o.RequiresObservation {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback must require observation")
	}
	if o.Status == "BOUND" {
		switch o.MetricSignal {
		case "boundary-complete":
			if o.FeedbackSignal != "observe" ||
				o.RequiresReview ||
				o.RequiresReplan ||
				o.RequiresMeasurement {
				return fmt.Errorf("complete lsp boundary feedback does not match its signal")
			}
		case "boundary-incomplete":
			if o.FeedbackSignal != "replan" ||
				!o.RequiresReplan ||
				o.RequiresReview ||
				o.RequiresMeasurement {
				return fmt.Errorf("incomplete lsp boundary feedback does not match its signal")
			}
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback must remain read-only and non-executing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(o) {
		return fmt.Errorf("lsp cycle evidence feedback decision boundary metric feedback projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(
	feedback LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback,
) string {
	return digestString(strings.Join([]string{
		feedback.Status,
		feedback.MissingStage,
		feedback.BoundaryMetricDigest,
		feedback.MetricName,
		feedback.MetricSignal,
		feedback.FeedbackSignal,
		feedback.FeedbackReason,
		strconv.FormatBool(feedback.RequiresObservation),
		strconv.FormatBool(feedback.RequiresReview),
		strconv.FormatBool(feedback.RequiresReplan),
		strconv.FormatBool(feedback.RequiresMeasurement),
		feedback.FeedbackDigest,
		strconv.FormatBool(feedback.ReadOnly),
		strconv.FormatBool(feedback.NonExecuting),
		strconv.FormatBool(feedback.NonAuthorizing),
	}, "|"))
}
