package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleMetricFeedback struct {
	Status                   string
	MissingStage             string
	CycleMetricDigest        string
	CycleDigest              string
	MetricName               string
	Direction                string
	BaselineValue            int64
	CandidateValue           int64
	MetricDelta              int64
	MetricSignal             string
	FeedbackSignal           string
	FeedbackReason           string
	RequiresObservation      bool
	RequiresReview           bool
	RequiresReplan           bool
	RequiresMeasurement      bool
	ProjectionDigest         string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveLSPRevisionSelfImprovementCycleMetricFeedback(
	feedback RevisionSelfImprovementCycleMetricFeedback,
) (LSPRevisionSelfImprovementCycleMetricFeedback, error) {
	result := LSPRevisionSelfImprovementCycleMetricFeedback{
		Status:                   "UNKNOWN",
		MissingStage:             "lsp-revision-self-improvement-cycle-metric-feedback",
		CycleMetricDigest:        feedback.CycleMetricDigest,
		CycleDigest:              feedback.CycleDigest,
		MetricName:               feedback.MetricName,
		Direction:                feedback.Direction,
		BaselineValue:            feedback.BaselineValue,
		CandidateValue:           feedback.CandidateValue,
		MetricDelta:              feedback.MetricDelta,
		MetricSignal:              feedback.MetricSignal,
		FeedbackSignal:           feedback.FeedbackSignal,
		FeedbackReason:           feedback.FeedbackReason,
		RequiresObservation:      feedback.RequiresObservation,
		RequiresReview:            feedback.RequiresReview,
		RequiresReplan:            feedback.RequiresReplan,
		RequiresMeasurement:      feedback.RequiresMeasurement,
		ReadOnly:                 true,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleMetricFeedback(result)
	}
	setDigest()

	if err := feedback.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-metric-feedback-input"
		setDigest()
		return result, fmt.Errorf("cycle metric feedback is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-metric-feedback"
		setDigest()
		return result, fmt.Errorf("lsp cycle metric feedback projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleMetricFeedback) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle metric feedback status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle metric feedback has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle metric feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle metric": o.CycleMetricDigest,
		"cycle":        o.CycleDigest,
		"projection":   o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle metric feedback %s digest is invalid", name)
		}
	}
	if o.MetricName == "" || o.FeedbackReason == "" {
		return fmt.Errorf("lsp cycle metric feedback evidence is incomplete")
	}
	if o.Direction != "higher-is-better" && o.Direction != "lower-is-better" {
		return fmt.Errorf("lsp cycle metric feedback direction is invalid")
	}
	if o.MetricDelta != o.CandidateValue-o.BaselineValue {
		return fmt.Errorf("lsp cycle metric feedback delta is invalid")
	}
	if o.MetricSignal != "metric-improved" &&
		o.MetricSignal != "metric-regressed" &&
		o.MetricSignal != "metric-stable" {
		return fmt.Errorf("lsp cycle metric feedback metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("lsp cycle metric feedback signal is invalid")
	}
	if !o.RequiresObservation || !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle metric feedback must remain read-only and non-executing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleMetricFeedback(o) {
		return fmt.Errorf("lsp cycle metric feedback projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleMetricFeedback(
	feedback LSPRevisionSelfImprovementCycleMetricFeedback,
) string {
	return digestString(strings.Join([]string{
		feedback.Status,
		feedback.MissingStage,
		feedback.CycleMetricDigest,
		feedback.CycleDigest,
		feedback.MetricName,
		feedback.Direction,
		strconv.FormatInt(feedback.BaselineValue, 10),
		strconv.FormatInt(feedback.CandidateValue, 10),
		strconv.FormatInt(feedback.MetricDelta, 10),
		feedback.MetricSignal,
		feedback.FeedbackSignal,
		feedback.FeedbackReason,
		strconv.FormatBool(feedback.RequiresObservation),
		strconv.FormatBool(feedback.RequiresReview),
		strconv.FormatBool(feedback.RequiresReplan),
		strconv.FormatBool(feedback.RequiresMeasurement),
		strconv.FormatBool(feedback.ReadOnly),
		strconv.FormatBool(feedback.NonExecuting),
		strconv.FormatBool(feedback.NonAuthorizing),
	}, "|"))
}