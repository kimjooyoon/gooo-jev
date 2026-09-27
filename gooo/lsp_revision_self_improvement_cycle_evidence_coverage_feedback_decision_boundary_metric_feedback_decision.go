package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation struct {
	Status              string
	MissingStage        string
	FeedbackDigest      string
	MetricSignal        string
	FeedbackSignal      string
	FeedbackReason      string
	DecisionSignal      string
	DecisionReason      string
	RequiresObservation bool
	RequiresReplan      bool
	DecisionAligned     bool
	DecisionDigest      string
	ProjectionDigest    string
	ReadOnly            bool
	NonExecuting        bool
	NonAuthorizing      bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(
	decision RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation{
		Status:              "UNKNOWN",
		MissingStage:        "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision",
		FeedbackDigest:      decision.FeedbackDigest,
		MetricSignal:        decision.MetricSignal,
		FeedbackSignal:      decision.FeedbackSignal,
		FeedbackReason:      decision.FeedbackReason,
		DecisionSignal:      decision.DecisionSignal,
		DecisionReason:      decision.DecisionReason,
		RequiresObservation: decision.RequiresObservation,
		RequiresReplan:      decision.RequiresReplan,
		DecisionAligned:     decision.DecisionAligned,
		DecisionDigest:      decision.DecisionDigest,
		ReadOnly:            true,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecision(result)
	}
	setDigest()

	if err := decision.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision is not valid: %w", err)
	}
	if decision.Status != "BOUND" || decision.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision"
		setDigest()
		return result, fmt.Errorf("lsp cycle evidence feedback decision projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle evidence feedback decision status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle evidence feedback decision has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle evidence feedback decision has no missing stage")
	}
	for name, digest := range map[string]string{
		"feedback":    o.FeedbackDigest,
		"decision":    o.DecisionDigest,
		"projection":  o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle evidence feedback decision %s digest is invalid", name)
		}
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("lsp cycle evidence feedback decision metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "replan" {
		return fmt.Errorf("lsp cycle evidence feedback decision feedback signal is invalid")
	}
	if o.FeedbackReason == "" || !o.RequiresObservation {
		return fmt.Errorf("lsp cycle evidence feedback decision feedback evidence is incomplete")
	}
	if o.Status == "BOUND" {
		switch o.FeedbackSignal {
		case "observe":
			if o.DecisionSignal != "observe-next-cycle" ||
				o.RequiresReplan ||
				!o.DecisionAligned {
				return fmt.Errorf("observe feedback does not match its lsp decision observation")
			}
		case "replan":
			if o.DecisionSignal != "replan-next-cycle" ||
				!o.RequiresReplan ||
				!o.DecisionAligned {
				return fmt.Errorf("replan feedback does not match its lsp decision observation")
			}
		}
		if o.DecisionDigest != digestString(strings.Join([]string{
			o.FeedbackDigest,
			o.DecisionSignal,
			o.DecisionReason,
		}, "|")) {
			return fmt.Errorf("lsp cycle evidence feedback decision digest does not match its fields")
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle evidence feedback decision must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecision(o) {
		return fmt.Errorf("lsp cycle evidence feedback decision projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecision(
	decision LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation,
) string {
	return digestString(strings.Join([]string{
		decision.Status,
		decision.MissingStage,
		decision.FeedbackDigest,
		decision.MetricSignal,
		decision.FeedbackSignal,
		decision.FeedbackReason,
		decision.DecisionSignal,
		decision.DecisionReason,
		strconv.FormatBool(decision.RequiresObservation),
		strconv.FormatBool(decision.RequiresReplan),
		strconv.FormatBool(decision.DecisionAligned),
		strconv.FormatBool(decision.ReadOnly),
		strconv.FormatBool(decision.NonExecuting),
		strconv.FormatBool(decision.NonAuthorizing),
	}, "|"))
}
