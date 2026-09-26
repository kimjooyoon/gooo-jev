package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeBound = "generation-request-candidate-feedback-application-plan-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeInput struct {
    GenerationRequestFeedback JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge
    CandidateGateApplicationPlan JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge struct {
    Status string
    MissingStage string
    GenerationRequestStatus string
    GenerationRequestDigest string
    GenerationRequestSource string
    GenerationRequestEvidenceDigest string
    GeneratorIdentity string
    ProposalDecision string
    FeedbackStatus string
    GeneratedCandidateStatus string
    CandidateDigest string
    GenerationSource string
    GenerationEvidenceDigest string
    FeedbackBridgeDigest string
    CandidateGateStatus string
    CandidateDecision string
    CandidateGateDigest string
    ApplicationPlanStatus string
    PlanDigest string
    PlanSource string
    PlanEvidenceDigest string
    CandidateGateApplicationPlanBridgeDigest string
    BridgeDigest string
    NonExecuting bool
    NonAuthorizing bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge) Validate() error {
    if !b.NonExecuting || !b.NonAuthorizing {
        return fmt.Errorf("JEV generation request candidate feedback application plan bridge must be non-executing and non-authorizing")
    }
    if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if b.MissingStage == "" || b.BridgeDigest != "" {
            return fmt.Errorf("unknown JEV generation request candidate feedback application plan bridge is inconsistent")
        }
        return nil
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeBound ||
        b.BridgeDigest == "" || b.FeedbackBridgeDigest == "" || b.CandidateGateApplicationPlanBridgeDigest == "" {
        return fmt.Errorf("incomplete JEV generation request candidate feedback application plan bridge")
    }
    if b.GenerationRequestStatus == "" || b.GenerationRequestDigest == "" || b.GenerationRequestSource == "" ||
        b.GenerationRequestEvidenceDigest == "" || b.GeneratorIdentity == "" || b.FeedbackStatus == "" ||
        b.CandidateGateStatus == "" || b.CandidateDecision == "" || b.CandidateGateDigest == "" ||
        b.ApplicationPlanStatus == "" {
        return fmt.Errorf("JEV generation request candidate feedback application plan bridge is missing provenance")
    }
    if b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound ||
        b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("JEV generation request candidate feedback application plan bridge has invalid upstream status")
    }
    switch b.ProposalDecision {
    case "revise":
        if b.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated ||
            b.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
            b.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady ||
            b.CandidateDigest == "" || b.GenerationSource == "" || b.GenerationEvidenceDigest == "" ||
            b.CandidateDigest != b.PlanDigest || b.GenerationSource != b.PlanSource ||
            b.PlanDigest == "" || b.PlanSource == "" || b.PlanEvidenceDigest == "" {
            return fmt.Errorf("revise generation request application plan is inconsistent")
        }
    case "retain":
        if b.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld ||
            b.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateHold ||
            b.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld ||
            b.CandidateDigest != "" || b.GenerationSource != "" || b.GenerationEvidenceDigest != "" ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" {
            return fmt.Errorf("retain generation request application plan is inconsistent")
        }
    case "rollback":
        if b.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected ||
            b.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateRejected ||
            b.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected ||
            b.CandidateDigest != "" || b.GenerationSource != "" || b.GenerationEvidenceDigest != "" ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" {
            return fmt.Errorf("rollback generation request application plan is inconsistent")
        }
    default:
        return fmt.Errorf("invalid generation request proposal decision")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge(
        b.Status, b.FeedbackBridgeDigest, b.CandidateGateApplicationPlanBridgeDigest,
        b.ProposalDecision, b.ApplicationPlanStatus,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV generation request candidate feedback application plan bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge {
        if stage == "" {
            stage = "generation-request-candidate-feedback-application-plan"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: stage,
            NonExecuting: true,
            NonAuthorizing: true,
        }
    }
    if !input.NonAuthorizing || !input.GenerationRequestFeedback.NonExecuting ||
        !input.GenerationRequestFeedback.NonAuthorizing || !input.CandidateGateApplicationPlan.NonExecuting ||
        !input.CandidateGateApplicationPlan.NonAuthorizing {
        return unknown("capability-boundary")
    }
    if err := input.GenerationRequestFeedback.Validate(); err != nil {
        return unknown("generation-request-candidate-feedback")
    }
    if err := input.CandidateGateApplicationPlan.Validate(); err != nil {
        return unknown("candidate-gate-application-plan")
    }
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge{
        Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeBound,
        GenerationRequestStatus: input.GenerationRequestFeedback.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestFeedback.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestFeedback.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestFeedback.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GenerationRequestFeedback.GeneratorIdentity,
        ProposalDecision: input.GenerationRequestFeedback.ProposalDecision,
        FeedbackStatus: input.GenerationRequestFeedback.Feedback.Status,
        GeneratedCandidateStatus: input.GenerationRequestFeedback.Feedback.GeneratedCandidateStatus,
        CandidateDigest: input.GenerationRequestFeedback.Feedback.CandidateDigest,
        GenerationSource: input.GenerationRequestFeedback.Feedback.GenerationSource,
        GenerationEvidenceDigest: input.GenerationRequestFeedback.Feedback.GenerationEvidenceDigest,
        FeedbackBridgeDigest: input.GenerationRequestFeedback.FeedbackBridgeDigest,
        CandidateGateStatus: input.CandidateGateApplicationPlan.CandidateGateStatus,
        CandidateDecision: input.CandidateGateApplicationPlan.CandidateDecision,
        CandidateGateDigest: input.CandidateGateApplicationPlan.CandidateGateDigest,
        ApplicationPlanStatus: input.CandidateGateApplicationPlan.ApplicationPlanStatus,
        PlanDigest: input.CandidateGateApplicationPlan.PlanDigest,
        PlanSource: input.CandidateGateApplicationPlan.PlanSource,
        PlanEvidenceDigest: input.CandidateGateApplicationPlan.PlanEvidenceDigest,
        CandidateGateApplicationPlanBridgeDigest: input.CandidateGateApplicationPlan.BridgeDigest,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge(
        output.Status, output.FeedbackBridgeDigest, output.CandidateGateApplicationPlanBridgeDigest,
        output.ProposalDecision, output.ApplicationPlanStatus,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-candidate-feedback-application-plan-evidence")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge(status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, proposalDecision, applicationPlanStatus string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, proposalDecision, applicationPlanStatus}, "|")))
    return hex.EncodeToString(sum[:])
}
