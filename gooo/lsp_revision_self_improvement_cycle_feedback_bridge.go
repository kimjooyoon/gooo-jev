package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleFeedbackBridge struct {
	Status                   string
	MissingStage             string
	CycleDigest              string
	MetricDigest             string
	FeedbackDigest           string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	MetricSignal             string
	FeedbackSignal           string
	FeedbackReason           string
	BridgeSignal             string
	RequiresObservation      bool
	RequiresReview           bool
	RequiresReplan           bool
	RequiresMeasurement      bool
	SignalsAligned           bool
	ProjectionDigest         string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveLSPRevisionSelfImprovementCycleFeedbackBridge(
	bridge RevisionSelfImprovementCycleFeedbackBridgeObservation,
) (LSPRevisionSelfImprovementCycleFeedbackBridge, error) {
	result := LSPRevisionSelfImprovementCycleFeedbackBridge{
		Status:                   "UNKNOWN",
		MissingStage:             "lsp-revision-self-improvement-cycle-feedback-bridge",
		CycleDigest:              bridge.CycleDigest,
		MetricDigest:             bridge.MetricDigest,
		FeedbackDigest:           bridge.FeedbackDigest,
		SourceDigest:             bridge.SourceDigest,
		CandidateSourceDigest:    bridge.CandidateSourceDigest,
		GeneratedIRDigest:        bridge.GeneratedIRDigest,
		ReverseObservationDigest: bridge.ReverseObservationDigest,
		MetricSignal:             bridge.MetricSignal,
		FeedbackSignal:           bridge.FeedbackSignal,
		FeedbackReason:           bridge.FeedbackReason,
		BridgeSignal:             bridge.BridgeSignal,
		RequiresObservation:      bridge.RequiresObservation,
		RequiresReview:           bridge.RequiresReview,
		RequiresReplan:            bridge.RequiresReplan,
		RequiresMeasurement:      bridge.RequiresMeasurement,
		SignalsAligned:            bridge.SignalsAligned,
		ReadOnly:                 true,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleFeedbackBridge(result)
	}
	setDigest()

	if err := bridge.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-feedback-bridge-input"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not valid: %w", err)
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-feedback-bridge"
		setDigest()
		return result, fmt.Errorf("lsp cycle feedback bridge projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleFeedbackBridge) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle feedback bridge status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle feedback bridge has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle feedback bridge has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle":              o.CycleDigest,
		"metric":             o.MetricDigest,
		"feedback":           o.FeedbackDigest,
		"source":             o.SourceDigest,
		"candidate source":   o.CandidateSourceDigest,
		"generated ir":       o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"projection":         o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle feedback bridge %s digest is invalid", name)
		}
	}
	if o.MetricSignal != "metric-improved" &&
		o.MetricSignal != "metric-regressed" &&
		o.MetricSignal != "metric-stable" {
		return fmt.Errorf("lsp cycle feedback bridge metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("lsp cycle feedback bridge feedback signal is invalid")
	}
	if o.FeedbackReason == "" || o.BridgeSignal != "cycle-feedback-linked" {
		return fmt.Errorf("lsp cycle feedback bridge evidence is incomplete")
	}
	if !o.SignalsAligned || !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle feedback bridge must remain aligned, read-only, and non-executing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleFeedbackBridge(o) {
		return fmt.Errorf("lsp cycle feedback bridge projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleFeedbackBridge(
	bridge LSPRevisionSelfImprovementCycleFeedbackBridge,
) string {
	return digestString(strings.Join([]string{
		bridge.Status,
		bridge.MissingStage,
		bridge.CycleDigest,
		bridge.MetricDigest,
		bridge.FeedbackDigest,
		bridge.SourceDigest,
		bridge.CandidateSourceDigest,
		bridge.GeneratedIRDigest,
		bridge.ReverseObservationDigest,
		bridge.MetricSignal,
		bridge.FeedbackSignal,
		bridge.FeedbackReason,
		bridge.BridgeSignal,
		strconv.FormatBool(bridge.RequiresObservation),
		strconv.FormatBool(bridge.RequiresReview),
		strconv.FormatBool(bridge.RequiresReplan),
		strconv.FormatBool(bridge.RequiresMeasurement),
		strconv.FormatBool(bridge.SignalsAligned),
		strconv.FormatBool(bridge.ReadOnly),
		strconv.FormatBool(bridge.NonExecuting),
		strconv.FormatBool(bridge.NonAuthorizing),
	}, "|"))
}
