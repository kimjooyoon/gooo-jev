package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeBound = "candidate-evaluation-feedback-revision-proposal-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBound = "candidate-revision-proposal-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalHeld = "candidate-revision-proposal-held"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalRejected = "candidate-revision-proposal-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeInput struct {
	FeedbackBridge         JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge
	ProposalDecision       string
	ProposalTarget         string
	ProposalDigest         string
	ProposalSource         string
	ProposalEvidenceDigest string
	NonAuthorizing         bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge
	RevisionProposalStatus string
	ProposalDecision       string
	ProposalTarget         string
	ProposalDigest         string
	ProposalSource         string
	ProposalEvidenceDigest string
	BridgeDigest           string
	NonExecuting           bool
	NonAuthorizing         bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge) Validate() error {
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("JEV candidate revision proposal bridge must be non-executing and non-authorizing")
	}
	if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		if b.MissingStage == "" || b.BridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate revision proposal bridge is inconsistent")
		}
		return nil
	}
	if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeBound ||
		b.BridgeDigest == "" {
		return fmt.Errorf("incomplete JEV candidate revision proposal bridge")
	}
	if b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge.BridgeDigest == "" {
		return fmt.Errorf("incomplete candidate feedback bridge")
	}
	switch b.AdmissionDecision {
	case "admit":
		if b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackBound ||
			b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBound ||
			b.ProposalTarget == "" || b.ProposalDigest == "" || b.ProposalSource == "" ||
			b.ProposalEvidenceDigest == "" || !validJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalDecision(b.ProposalDecision, b.FeedbackDirection) {
			return fmt.Errorf("admitted candidate revision proposal lost evidence or consistency")
		}
	case "hold":
		if b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalHeld ||
			b.ProposalDecision != "" || b.ProposalTarget != "" || b.ProposalDigest != "" ||
			b.ProposalSource != "" || b.ProposalEvidenceDigest != "" {
			return fmt.Errorf("held candidate revision proposal contains evidence")
		}
	case "reject":
		if b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalRejected ||
			b.ProposalDecision != "" || b.ProposalTarget != "" || b.ProposalDigest != "" ||
			b.ProposalSource != "" || b.ProposalEvidenceDigest != "" {
			return fmt.Errorf("rejected candidate revision proposal contains evidence")
		}
	default:
		return fmt.Errorf("invalid candidate revision proposal admission decision")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge(
		b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge.BridgeDigest,
		b.RevisionProposalStatus, b.ProposalDecision, b.ProposalTarget, b.ProposalDigest,
		b.ProposalSource, b.ProposalEvidenceDigest,
	)
	if b.BridgeDigest != expected {
		return fmt.Errorf("JEV candidate revision proposal bridge digest mismatch")
	}
	return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge {
	unknown := func(stage string, prior JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge {
		if stage == "" {
			stage = "evaluation-outcome-feedback-revision-proposal"
		}
		prior.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		prior.MissingStage = stage
		prior.BridgeDigest = ""
		prior.NonExecuting = true
		prior.NonAuthorizing = true
		return JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge: prior,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	prior := input.FeedbackBridge
	if !input.NonAuthorizing || !prior.NonExecuting || !prior.NonAuthorizing {
		return unknown("capability-boundary", prior)
	}
	if err := prior.Validate(); err != nil {
		return unknown("evaluation-outcome-feedback", prior)
	}
	if prior.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		return unknown(prior.MissingStage, prior)
	}
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge: prior,
		ProposalDecision:       input.ProposalDecision,
		ProposalTarget:         input.ProposalTarget,
		ProposalDigest:         input.ProposalDigest,
		ProposalSource:         input.ProposalSource,
		ProposalEvidenceDigest: input.ProposalEvidenceDigest,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	switch prior.AdmissionDecision {
	case "admit":
		output.RevisionProposalStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBound
		stages := []struct {
			value string
			stage string
		}{
			{output.ProposalDecision, "proposal-decision"},
			{output.ProposalTarget, "proposal-target"},
			{output.ProposalDigest, "proposal-digest"},
			{output.ProposalSource, "proposal-source"},
			{output.ProposalEvidenceDigest, "proposal-evidence"},
		}
		for _, item := range stages {
			if item.value == "" {
				return unknown(item.stage, prior)
			}
		}
		if !validJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalDecision(output.ProposalDecision, prior.FeedbackDirection) {
			return unknown("proposal-decision-consistency", prior)
		}
	case "hold":
		output.RevisionProposalStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalHeld
		if input.ProposalDecision != "" || input.ProposalTarget != "" || input.ProposalDigest != "" ||
			input.ProposalSource != "" || input.ProposalEvidenceDigest != "" {
			return unknown("held-proposal-evidence", prior)
		}
	case "reject":
		output.RevisionProposalStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalRejected
		if input.ProposalDecision != "" || input.ProposalTarget != "" || input.ProposalDigest != "" ||
			input.ProposalSource != "" || input.ProposalEvidenceDigest != "" {
			return unknown("rejected-proposal-evidence", prior)
		}
	default:
		return unknown("admission-decision", prior)
	}
	output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeBound
	output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge.BridgeDigest,
		output.RevisionProposalStatus, output.ProposalDecision, output.ProposalTarget,
		output.ProposalDigest, output.ProposalSource, output.ProposalEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("evaluation-outcome-feedback-revision-proposal", prior)
	}
	return output
}

func validJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalDecision(decision, direction string) bool {
	switch direction {
	case "improve":
		return decision == "revise"
	case "retain":
		return decision == "retain"
	case "rollback":
		return decision == "rollback"
	default:
		return false
	}
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge(
	feedbackBridgeDigest, revisionProposalStatus, proposalDecision, proposalTarget, proposalDigest,
	proposalSource, proposalEvidenceDigest string,
) string {
	values := []string{
		feedbackBridgeDigest, revisionProposalStatus, proposalDecision, proposalTarget,
		proposalDigest, proposalSource, proposalEvidenceDigest,
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

