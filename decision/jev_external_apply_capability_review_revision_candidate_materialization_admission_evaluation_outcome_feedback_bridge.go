package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeBound = "candidate-evaluation-outcome-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackBound = "candidate-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackHeld = "candidate-feedback-held"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackRejected = "candidate-feedback-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeInput struct {
	Outcome               JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge
	FeedbackDirection     string
	FeedbackTarget        string
	FeedbackDigest        string
	FeedbackSource        string
	FeedbackEvidenceDigest string
	NonAuthorizing       bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge
	FeedbackDirection     string
	FeedbackTarget        string
	FeedbackDigest        string
	FeedbackSource        string
	FeedbackEvidenceDigest string
	FeedbackStatus        string
	BridgeDigest          string
	NonExecuting          bool
	NonAuthorizing        bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge) Validate() error {
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("JEV candidate evaluation feedback bridge must be non-executing and non-authorizing")
	}
	if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		if b.MissingStage == "" || b.BridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate evaluation feedback bridge is inconsistent")
		}
		return nil
	}
	if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeBound ||
		b.BridgeDigest == "" {
		return fmt.Errorf("incomplete JEV candidate evaluation feedback bridge")
	}
	if err := b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge.Validate(); err != nil {
		return fmt.Errorf("invalid evaluation outcome bridge: %w", err)
	}
	switch b.AdmissionDecision {
	case "admit":
		if b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackBound ||
			b.FeedbackDirection == "" || b.FeedbackTarget == "" || b.FeedbackDigest == "" ||
			b.FeedbackSource == "" || b.FeedbackEvidenceDigest == "" {
			return fmt.Errorf("admitted JEV candidate feedback lost evidence")
		}
		switch b.FeedbackDirection {
		case "improve", "retain", "rollback":
		default:
			return fmt.Errorf("invalid JEV candidate feedback direction")
		}
	case "hold":
		if b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackHeld ||
			b.FeedbackDirection != "" || b.FeedbackTarget != "" || b.FeedbackDigest != "" ||
			b.FeedbackSource != "" || b.FeedbackEvidenceDigest != "" {
			return fmt.Errorf("held JEV candidate feedback contains evidence")
		}
	case "reject":
		if b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackRejected ||
			b.FeedbackDirection != "" || b.FeedbackTarget != "" || b.FeedbackDigest != "" ||
			b.FeedbackSource != "" || b.FeedbackEvidenceDigest != "" {
			return fmt.Errorf("rejected JEV candidate feedback contains evidence")
		}
	default:
		return fmt.Errorf("invalid JEV candidate evaluation outcome admission decision")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge(
		b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge.BridgeDigest,
		b.FeedbackDirection, b.FeedbackTarget, b.FeedbackDigest, b.FeedbackSource,
		b.FeedbackEvidenceDigest, b.FeedbackStatus, b.NonExecuting, b.NonAuthorizing,
	)
	if b.BridgeDigest != expected {
		return fmt.Errorf("JEV candidate evaluation feedback bridge digest mismatch")
	}
	return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge {
	unknown := func(stage string, prior JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge {
		if stage == "" {
			stage = "evaluation-outcome-feedback-bridge"
		}
		prior.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		prior.MissingStage = stage
		prior.BridgeDigest = ""
		prior.NonExecuting = true
		prior.NonAuthorizing = true
		return JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge: prior,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	prior := input.Outcome
	if !input.NonAuthorizing || !prior.NonExecuting || !prior.NonAuthorizing {
		return unknown("capability-boundary", prior)
	}
	if err := prior.Validate(); err != nil {
		return unknown("evaluation-outcome-bridge", prior)
	}
	if prior.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		return unknown(prior.MissingStage, prior)
	}
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge: prior,
		FeedbackDirection:      input.FeedbackDirection,
		FeedbackTarget:         input.FeedbackTarget,
		FeedbackDigest:         input.FeedbackDigest,
		FeedbackSource:         input.FeedbackSource,
		FeedbackEvidenceDigest: input.FeedbackEvidenceDigest,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	switch prior.AdmissionDecision {
	case "admit":
		output.FeedbackStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackBound
		stages := []struct {
			value string
			stage string
		}{
			{output.FeedbackDirection, "feedback-direction"},
			{output.FeedbackTarget, "feedback-target"},
			{output.FeedbackDigest, "feedback-digest"},
			{output.FeedbackSource, "feedback-source"},
			{output.FeedbackEvidenceDigest, "feedback-evidence"},
		}
		for _, item := range stages {
			if item.value == "" {
				return unknown(item.stage, prior)
			}
		}
		if output.FeedbackDirection != "improve" && output.FeedbackDirection != "retain" && output.FeedbackDirection != "rollback" {
			return unknown("feedback-direction", prior)
		}
	case "hold":
		output.FeedbackStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackHeld
		if input.FeedbackDirection != "" || input.FeedbackTarget != "" || input.FeedbackDigest != "" ||
			input.FeedbackSource != "" || input.FeedbackEvidenceDigest != "" {
			return unknown("held-feedback-evidence", prior)
		}
	case "reject":
		output.FeedbackStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackRejected
		if input.FeedbackDirection != "" || input.FeedbackTarget != "" || input.FeedbackDigest != "" ||
			input.FeedbackSource != "" || input.FeedbackEvidenceDigest != "" {
			return unknown("rejected-feedback-evidence", prior)
		}
	default:
		return unknown("admission-decision", prior)
	}
	output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge.BridgeDigest,
		output.FeedbackDirection, output.FeedbackTarget, output.FeedbackDigest,
		output.FeedbackSource, output.FeedbackEvidenceDigest, output.FeedbackStatus,
		output.NonExecuting, output.NonAuthorizing,
	)
	if err := output.Validate(); err != nil {
		return unknown("evaluation-outcome-feedback-bridge", prior)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge(
	outcomeBridgeDigest, feedbackDirection, feedbackTarget, feedbackDigest, feedbackSource,
	feedbackEvidenceDigest, feedbackStatus string, nonExecuting, nonAuthorizing bool,
) string {
	values := []string{
		outcomeBridgeDigest, feedbackDirection, feedbackTarget, feedbackDigest, feedbackSource,
		feedbackEvidenceDigest, feedbackStatus, strconv.FormatBool(nonExecuting),
		strconv.FormatBool(nonAuthorizing),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

