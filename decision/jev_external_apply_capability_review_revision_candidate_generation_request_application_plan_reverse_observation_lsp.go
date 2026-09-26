package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPProjectionBound = "generation-request-application-plan-reverse-observation-lsp-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPDiagnostic struct {
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
    FeedbackBridgeDigest string
    CandidateGateApplicationPlanBridgeDigest string
    ApplicationPlanReverseObservation JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic
    ApplicationPlanReverseObservationBridgeDigest string
    BridgeDigest string
    ProjectionDigest string
    Publishable bool
    NonExecuting bool
    NonAuthorizing bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV generation request reverse observation LSP diagnostic")
    }
    if !d.NonExecuting || !d.NonAuthorizing {
        return fmt.Errorf("JEV generation request reverse observation LSP diagnostic must be non-executing and non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if d.Publishable || d.MissingStage == "" || d.BridgeDigest != "" || d.ProjectionDigest != "" {
            return fmt.Errorf("UNKNOWN JEV generation request reverse observation diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeBound ||
        !d.Publishable || d.MissingStage != "" || d.GenerationRequestStatus == "" ||
        d.GenerationRequestDigest == "" || d.GenerationRequestSource == "" ||
        d.GenerationRequestEvidenceDigest == "" || d.GeneratorIdentity == "" ||
        d.ProposalDecision == "" || d.FeedbackBridgeDigest == "" ||
        d.CandidateGateApplicationPlanBridgeDigest == "" ||
        d.ApplicationPlanReverseObservationBridgeDigest == "" || d.BridgeDigest == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("bound JEV generation request reverse observation diagnostic is incomplete")
    }
    if err := d.ApplicationPlanReverseObservation.Validate(); err != nil {
        return fmt.Errorf("invalid reverse observation diagnostic: %w", err)
    }
    if !d.ApplicationPlanReverseObservation.Publishable ||
        d.ApplicationPlanReverseObservation.BridgeDigest != d.ApplicationPlanReverseObservationBridgeDigest {
        return fmt.Errorf("generation request reverse observation nested digest is inconsistent")
    }
    switch d.ProposalDecision {
    case "revise":
        if d.ApplicationPlanReverseObservation.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady {
            return fmt.Errorf("revise generation request reverse observation diagnostic is inconsistent")
        }
    case "retain":
        if d.ApplicationPlanReverseObservation.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld {
            return fmt.Errorf("retain generation request reverse observation diagnostic is inconsistent")
        }
    case "rollback":
        if d.ApplicationPlanReverseObservation.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected {
            return fmt.Errorf("rollback generation request reverse observation diagnostic is inconsistent")
        }
    default:
        return fmt.Errorf("invalid generation request proposal decision in reverse observation diagnostic")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPProjection(
        d.Status, d.FeedbackBridgeDigest, d.CandidateGateApplicationPlanBridgeDigest,
        d.ApplicationPlanReverseObservationBridgeDigest, d.BridgeDigest, d.ProposalDecision,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV generation request reverse observation projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPDiagnostic {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPDiagnostic {
        if stage == "" {
            stage = "generation-request-application-plan-reverse-observation-lsp"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPDiagnostic{
            Severity: "warning",
            Code: "jev.provenance.unknown",
            Message: "Generation request reverse observation is UNKNOWN; evidence must be resolved",
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
        return unknown("generation-request-application-plan-reverse-observation")
    }
    reverseObservation := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge{
            Status: input.ApplicationPlanReverseObservationBridgeStatus(),
            ApplicationPlanStatus: input.ApplicationPlanStatus,
            PlanDigest: input.PlanDigest,
            PlanSource: input.PlanSource,
            PlanEvidenceDigest: input.PlanEvidenceDigest,
            ApplicationPlanBridgeDigest: input.CandidateGateApplicationPlanBridgeDigest,
            ReverseObservationStatus: input.ReverseObservationStatus,
            ObservationDigest: input.ObservationDigest,
            ObservationSource: input.ObservationSource,
            ObservationEvidenceDigest: input.ObservationEvidenceDigest,
            ObservationMetricStatus: input.ObservationMetricStatus,
            ObservationMetricDigest: input.ObservationMetricDigest,
            ObservationMetricSource: input.ObservationMetricSource,
            ObservationMetricEvidenceDigest: input.ObservationMetricEvidenceDigest,
            BridgeDigest: input.ApplicationPlanReverseObservationBridgeDigest,
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPDiagnostic{
        Status: input.Status,
        GenerationRequestStatus: input.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GeneratorIdentity,
        ProposalDecision: input.ProposalDecision,
        FeedbackBridgeDigest: input.FeedbackBridgeDigest,
        CandidateGateApplicationPlanBridgeDigest: input.CandidateGateApplicationPlanBridgeDigest,
        ApplicationPlanReverseObservation: reverseObservation,
        ApplicationPlanReverseObservationBridgeDigest: input.ApplicationPlanReverseObservationBridgeDigest,
        BridgeDigest: input.BridgeDigest,
        Publishable: false,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if err := reverseObservation.Validate(); err != nil {
        return unknown("application-plan-reverse-observation-lsp")
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.generation-request-reverse-observation-bound"
    output.Message = "Generation request reverse observation is bound without execution or authorization"
    if input.ProposalDecision == "retain" {
        output.Code = "jev.external-apply.generation-request-reverse-observation-held"
        output.Message = "Generation request reverse observation is held without execution or authorization"
    }
    if input.ProposalDecision == "rollback" {
        output.Severity = "warning"
        output.Code = "jev.external-apply.generation-request-reverse-observation-rejected"
        output.Message = "Generation request reverse observation is rejected without execution or authorization"
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPProjection(
        output.Status, output.FeedbackBridgeDigest, output.CandidateGateApplicationPlanBridgeDigest,
        output.ApplicationPlanReverseObservationBridgeDigest, output.BridgeDigest, output.ProposalDecision,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-application-plan-reverse-observation-lsp")
    }
    return output
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge) ApplicationPlanReverseObservationBridgeStatus() string {
    return jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeBound
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPProjection(status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, applicationPlanReverseObservationBridgeDigest, bridgeDigest, proposalDecision string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, applicationPlanReverseObservationBridgeDigest, bridgeDigest, proposalDecision}, "|")))
    return hex.EncodeToString(sum[:])
}
