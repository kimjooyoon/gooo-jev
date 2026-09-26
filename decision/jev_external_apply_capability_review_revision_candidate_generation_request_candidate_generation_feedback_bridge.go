package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeBound = "generation-request-candidate-generation-feedback-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeInput struct {
	GenerationRequest        JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge
	Direction                JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge
	CandidateDigest          string
	GenerationSource         string
	GenerationEvidenceDigest string
	NonAuthorizing           bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge struct {
	Status                    string
	MissingStage              string
	GenerationRequestStatus   string
	GenerationRequestDigest   string
	GenerationRequestSource   string
	GenerationRequestEvidenceDigest string
	GeneratorIdentity         string
	ProposalDecision          string
	Feedback                  JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback
	GenerationRequestBridgeDigest string
	FeedbackBridgeDigest      string
	BridgeDigest              string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge) Validate() error {
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("JEV generation request candidate feedback bridge must be non-executing and non-authorizing")
	}
	if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		if b.MissingStage == "" || b.BridgeDigest != "" {
			return fmt.Errorf("unknown JEV generation request candidate feedback bridge is inconsistent")
		}
		return nil
	}
	if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeBound ||
		b.BridgeDigest == "" || b.GenerationRequestBridgeDigest == "" || b.FeedbackBridgeDigest == "" {
		return fmt.Errorf("incomplete JEV generation request candidate feedback bridge")
	}
	switch b.ProposalDecision {
	case "revise":
		if b.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBound ||
			b.Feedback.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate ||
			b.Feedback.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated {
			return fmt.Errorf("revise generation request is inconsistent with candidate generation feedback")
		}
	case "retain":
		if b.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestHeld ||
			b.Feedback.Direction != jevExternalApplyCapabilityReviewImprovementDirectionHold ||
			b.Feedback.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld {
			return fmt.Errorf("retain generation request is inconsistent with candidate generation feedback")
		}
	case "rollback":
		if b.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestRejected ||
			b.Feedback.Direction != jevExternalApplyCapabilityReviewImprovementDirectionReject ||
			b.Feedback.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected {
			return fmt.Errorf("rollback generation request is inconsistent with candidate generation feedback")
		}
	default:
		return fmt.Errorf("invalid generation request proposal decision")
	}
	if err := b.Feedback.Validate(); err != nil {
		return fmt.Errorf("invalid candidate generation feedback: %w", err)
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge(
		b.Status, b.GenerationRequestBridgeDigest, b.FeedbackBridgeDigest,
		b.ProposalDecision,
	)
	if b.BridgeDigest != expected {
		return fmt.Errorf("JEV generation request candidate feedback bridge digest mismatch")
	}
	return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge {
	unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge {
		if stage == "" {
			stage = "generation-request-candidate-generation-feedback"
		}
		return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge{
			Status:         jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.GenerationRequest.NonExecuting || !input.GenerationRequest.NonAuthorizing ||
		!input.Direction.NonExecuting || !input.Direction.NonAuthorizing {
		return unknown("capability-boundary")
	}
	if err := input.GenerationRequest.Validate(); err != nil {
		return unknown("generation-request")
	}
	if input.GenerationRequest.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		return unknown(input.GenerationRequest.MissingStage)
	}
	expectedDirection := map[string]string{
		"revise":   jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
		"retain":   jevExternalApplyCapabilityReviewImprovementDirectionHold,
		"rollback": jevExternalApplyCapabilityReviewImprovementDirectionReject,
	}[input.GenerationRequest.ProposalDecision]
	if expectedDirection == "" {
		return unknown("proposal-decision")
	}
	if input.Direction.Direction != expectedDirection {
		return unknown("direction-consistency")
	}
	feedback := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(
		JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackInput{
			Direction:                 input.Direction,
			CandidateDigest:           input.CandidateDigest,
			GenerationSource:          input.GenerationSource,
			GenerationEvidenceDigest: input.GenerationEvidenceDigest,
			NonAuthorizing:            true,
		},
	)
	if err := feedback.Validate(); err != nil {
		return unknown("candidate-generation-feedback")
	}
	output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge{
		Status:                       jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeBound,
		GenerationRequestStatus:      input.GenerationRequest.GenerationRequestStatus,
		GenerationRequestDigest:      input.GenerationRequest.GenerationRequestDigest,
		GenerationRequestSource:      input.GenerationRequest.GenerationRequestSource,
		GenerationRequestEvidenceDigest: input.GenerationRequest.GenerationRequestEvidenceDigest,
		GeneratorIdentity:            input.GenerationRequest.GeneratorIdentity,
		ProposalDecision:             input.GenerationRequest.ProposalDecision,
		Feedback:                     feedback,
		GenerationRequestBridgeDigest: input.GenerationRequest.BridgeDigest,
		FeedbackBridgeDigest:         feedback.BridgeDigest,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge(
		output.Status, output.GenerationRequestBridgeDigest, output.FeedbackBridgeDigest,
		output.ProposalDecision,
	)
	if err := output.Validate(); err != nil {
		return unknown("generation-request-candidate-generation-feedback")
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge(
	status, generationRequestBridgeDigest, feedbackBridgeDigest, proposalDecision string,
) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		status, generationRequestBridgeDigest, feedbackBridgeDigest, proposalDecision,
	}, "|")))
	return hex.EncodeToString(sum[:])
}

