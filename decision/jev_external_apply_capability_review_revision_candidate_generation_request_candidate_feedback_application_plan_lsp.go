package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPProjectionBound = "generation-request-candidate-feedback-application-plan-lsp-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPDiagnostic struct {
    Severity string
    Code string
    Message string
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
    ApplicationPlan JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPDiagnostic
    CandidateGateApplicationPlanBridgeDigest string
    BridgeDigest string
    ProjectionDigest string
    Publishable bool
    NonExecuting bool
    NonAuthorizing bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV generation request application plan LSP diagnostic")
    }
    if !d.NonExecuting || !d.NonAuthorizing {
        return fmt.Errorf("JEV generation request application plan LSP diagnostic must be non-executing and non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if d.Publishable || d.MissingStage == "" || d.BridgeDigest != "" || d.ProjectionDigest != "" {
            return fmt.Errorf("UNKNOWN JEV generation request application plan diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeBound ||
        !d.Publishable || d.MissingStage != "" || d.GenerationRequestStatus == "" ||
        d.GenerationRequestDigest == "" || d.GenerationRequestSource == "" ||
        d.GenerationRequestEvidenceDigest == "" || d.GeneratorIdentity == "" ||
        d.FeedbackStatus == "" || d.GeneratedCandidateStatus == "" ||
        d.FeedbackBridgeDigest == "" || d.CandidateGateStatus == "" ||
        d.CandidateDecision == "" || d.CandidateGateDigest == "" ||
        d.CandidateGateApplicationPlanBridgeDigest == "" || d.BridgeDigest == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("bound JEV generation request application plan diagnostic is incomplete")
    }
    if d.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound ||
        d.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("JEV generation request application plan diagnostic has invalid upstream status")
    }
    if err := d.ApplicationPlan.Validate(); err != nil {
        return fmt.Errorf("invalid application plan diagnostic: %w", err)
    }
    if !d.ApplicationPlan.Publishable || d.ApplicationPlan.BridgeDigest != d.CandidateGateApplicationPlanBridgeDigest {
        return fmt.Errorf("generation request application plan nested digest is inconsistent")
    }
    switch d.ProposalDecision {
    case "revise":
        if d.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated ||
            d.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
            d.ApplicationPlan.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady ||
            d.CandidateDigest == "" || d.GenerationSource == "" || d.GenerationEvidenceDigest == "" {
            return fmt.Errorf("revise generation request application plan diagnostic is inconsistent")
        }
    case "retain":
        if d.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackHeld ||
            d.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateHold ||
            d.ApplicationPlan.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld ||
            d.CandidateDigest != "" || d.GenerationSource != "" || d.GenerationEvidenceDigest != "" {
            return fmt.Errorf("retain generation request application plan diagnostic is inconsistent")
        }
    case "rollback":
        if d.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRejected ||
            d.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateRejected ||
            d.ApplicationPlan.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected ||
            d.CandidateDigest != "" || d.GenerationSource != "" || d.GenerationEvidenceDigest != "" {
            return fmt.Errorf("rollback generation request application plan diagnostic is inconsistent")
        }
    default:
        return fmt.Errorf("invalid generation request proposal decision in application plan diagnostic")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPProjection(
        d.Status, d.FeedbackBridgeDigest, d.CandidateGateApplicationPlanBridgeDigest,
        d.BridgeDigest, d.ProposalDecision,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV generation request application plan projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPDiagnostic {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPDiagnostic {
        if stage == "" {
            stage = "generation-request-candidate-feedback-application-plan-lsp"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPDiagnostic{
            Severity: "error",
            Code: "jev.provenance.unknown",
            Message: "Generation request application plan is UNKNOWN; evidence must be resolved",
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: stage,
            NonExecuting: true,
            NonAuthorizing: true,
        }
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        return unknown(input.MissingStage)
    }
    if err := input.Validate(); err != nil {
        return unknown("generation-request-candidate-feedback-application-plan")
    }
    applicationPlan := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge{
            Status: input.ApplicationPlanStatus,
            CandidateGateStatus: input.CandidateGateStatus,
            CandidateDecision: input.CandidateDecision,
            CandidateDigest: input.CandidateDigest,
            CandidateSource: input.GenerationSource,
            CandidateGateDigest: input.CandidateGateDigest,
            ApplicationPlanStatus: input.ApplicationPlanStatus,
            PlanDigest: input.PlanDigest,
            PlanSource: input.PlanSource,
            PlanEvidenceDigest: input.PlanEvidenceDigest,
            BridgeDigest: input.CandidateGateApplicationPlanBridgeDigest,
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPDiagnostic{
        Status: input.Status,
        GenerationRequestStatus: input.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GeneratorIdentity,
        ProposalDecision: input.ProposalDecision,
        FeedbackStatus: input.FeedbackStatus,
        GeneratedCandidateStatus: input.GeneratedCandidateStatus,
        CandidateDigest: input.CandidateDigest,
        GenerationSource: input.GenerationSource,
        GenerationEvidenceDigest: input.GenerationEvidenceDigest,
        FeedbackBridgeDigest: input.FeedbackBridgeDigest,
        CandidateGateStatus: input.CandidateGateStatus,
        CandidateDecision: input.CandidateDecision,
        CandidateGateDigest: input.CandidateGateDigest,
        ApplicationPlan: applicationPlan,
        CandidateGateApplicationPlanBridgeDigest: input.CandidateGateApplicationPlanBridgeDigest,
        BridgeDigest: input.BridgeDigest,
        Publishable: false,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if err := applicationPlan.Validate(); err != nil {
        return unknown("candidate-gate-application-plan-lsp")
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.generation-request-application-plan-bound"
    output.Message = "Generation request and application plan provenance are bound without execution or authorization"
    if input.ProposalDecision == "retain" {
        output.Code = "jev.external-apply.generation-request-application-plan-held"
        output.Message = "Generation request application plan is held without execution or authorization"
    }
    if input.ProposalDecision == "rollback" {
        output.Severity = "warning"
        output.Code = "jev.external-apply.generation-request-application-plan-rejected"
        output.Message = "Generation request application plan is rejected without execution or authorization"
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPProjection(
        output.Status, output.FeedbackBridgeDigest, output.CandidateGateApplicationPlanBridgeDigest,
        output.BridgeDigest, output.ProposalDecision,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-candidate-feedback-application-plan-lsp")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPProjection(status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, bridgeDigest, proposalDecision string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, bridgeDigest, proposalDecision}, "|")))
    return hex.EncodeToString(sum[:])
}
