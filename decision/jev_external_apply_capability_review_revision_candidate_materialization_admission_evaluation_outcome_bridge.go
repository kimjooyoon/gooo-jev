package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeBound = "candidate-evaluation-outcome-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeHeld = "candidate-evaluation-outcome-held"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeRejected = "candidate-evaluation-outcome-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeInput struct {
	Evaluation                 JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge
	Outcome                    string
	BaselineMetricDigest       string
	ObservedMetricDigest       string
	OutcomeSource              string
	OutcomeEvidenceDigest      string
	NonAuthorizing             bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge
	Outcome               string
	BaselineMetricDigest  string
	ObservedMetricDigest  string
	OutcomeSource         string
	OutcomeEvidenceDigest string
	OutcomeStatus         string
	BridgeDigest          string
	NonExecuting          bool
	NonAuthorizing        bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge) Validate() error {
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("JEV candidate evaluation outcome bridge must be non-executing and non-authorizing")
	}
	if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		if b.MissingStage == "" || b.BridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate evaluation outcome bridge is inconsistent")
		}
		return nil
	}
	if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound ||
		b.BridgeDigest == "" {
		return fmt.Errorf("incomplete JEV candidate evaluation outcome bridge")
	}
	if err := b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge.Validate(); err != nil {
		return fmt.Errorf("invalid evaluation bridge: %w", err)
	}
	switch b.AdmissionDecision {
	case "admit":
		if b.OutcomeStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeBound ||
			b.Outcome == "" || b.BaselineMetricDigest == "" || b.ObservedMetricDigest == "" ||
			b.OutcomeSource == "" || b.OutcomeEvidenceDigest == "" {
			return fmt.Errorf("admitted JEV candidate outcome lost metric evidence")
		}
		switch b.Outcome {
		case "improved", "regressed", "inconclusive":
		default:
			return fmt.Errorf("invalid JEV candidate evaluation outcome")
		}
	case "hold":
		if b.OutcomeStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeHeld ||
			b.Outcome != "" || b.BaselineMetricDigest != "" || b.ObservedMetricDigest != "" ||
			b.OutcomeSource != "" || b.OutcomeEvidenceDigest != "" {
			return fmt.Errorf("held JEV candidate outcome contains metric evidence")
		}
	case "reject":
		if b.OutcomeStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeRejected ||
			b.Outcome != "" || b.BaselineMetricDigest != "" || b.ObservedMetricDigest != "" ||
			b.OutcomeSource != "" || b.OutcomeEvidenceDigest != "" {
			return fmt.Errorf("rejected JEV candidate outcome contains metric evidence")
		}
	default:
		return fmt.Errorf("invalid JEV candidate evaluation admission decision")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(
		b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge.BridgeDigest,
		b.Outcome, b.BaselineMetricDigest, b.ObservedMetricDigest, b.OutcomeSource,
		b.OutcomeEvidenceDigest, b.OutcomeStatus, b.NonExecuting, b.NonAuthorizing,
	)
	if b.BridgeDigest != expected {
		return fmt.Errorf("JEV candidate evaluation outcome bridge digest mismatch")
	}
	return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge {
	unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge {
		if stage == "" {
			stage = "evaluation-outcome-bridge"
		}
		return JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
				Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
				MissingStage:  stage,
				NonExecuting:  true,
				NonAuthorizing: true,
			},
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Evaluation.NonExecuting || !input.Evaluation.NonAuthorizing {
		return unknown("capability-boundary")
	}
	if err := input.Evaluation.Validate(); err != nil {
		return unknown("evaluation-bridge")
	}
	if input.Evaluation.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		return unknown(input.Evaluation.MissingStage)
	}
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: input.Evaluation,
		Outcome:               input.Outcome,
		BaselineMetricDigest:  input.BaselineMetricDigest,
		ObservedMetricDigest:  input.ObservedMetricDigest,
		OutcomeSource:         input.OutcomeSource,
		OutcomeEvidenceDigest: input.OutcomeEvidenceDigest,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	switch input.Evaluation.AdmissionDecision {
	case "admit":
		output.OutcomeStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeBound
		stages := []struct {
			value string
			stage string
		}{
			{output.Outcome, "evaluation-outcome"},
			{output.BaselineMetricDigest, "baseline-metric-digest"},
			{output.ObservedMetricDigest, "observed-metric-digest"},
			{output.OutcomeSource, "outcome-source"},
			{output.OutcomeEvidenceDigest, "outcome-evidence"},
		}
		for _, item := range stages {
			if item.value == "" {
				return unknown(item.stage)
			}
		}
		if output.Outcome != "improved" && output.Outcome != "regressed" && output.Outcome != "inconclusive" {
			return unknown("evaluation-outcome")
		}
	case "hold":
		output.OutcomeStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeHeld
		if input.Outcome != "" || input.BaselineMetricDigest != "" || input.ObservedMetricDigest != "" ||
			input.OutcomeSource != "" || input.OutcomeEvidenceDigest != "" {
			return unknown("held-outcome-evidence")
		}
	case "reject":
		output.OutcomeStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeRejected
		if input.Outcome != "" || input.BaselineMetricDigest != "" || input.ObservedMetricDigest != "" ||
			input.OutcomeSource != "" || input.OutcomeEvidenceDigest != "" {
			return unknown("rejected-outcome-evidence")
		}
	default:
		return unknown("admission-decision")
	}
	output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge.BridgeDigest,
		output.Outcome, output.BaselineMetricDigest, output.ObservedMetricDigest,
		output.OutcomeSource, output.OutcomeEvidenceDigest, output.OutcomeStatus,
		output.NonExecuting, output.NonAuthorizing,
	)
	if err := output.Validate(); err != nil {
		return unknown("evaluation-outcome-bridge")
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(
	evaluationBridgeDigest, outcome, baselineMetricDigest, observedMetricDigest, outcomeSource,
	outcomeEvidenceDigest, outcomeStatus string, nonExecuting, nonAuthorizing bool,
) string {
	values := []string{
		evaluationBridgeDigest, outcome, baselineMetricDigest, observedMetricDigest,
		outcomeSource, outcomeEvidenceDigest, outcomeStatus,
		fmt.Sprint(nonExecuting), fmt.Sprint(nonAuthorizing),
	}
	return digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeValues(values)
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeValues(values []string) string {
	return digestStringsJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(values)
}


func digestStringsJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(values []string) string {
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}
