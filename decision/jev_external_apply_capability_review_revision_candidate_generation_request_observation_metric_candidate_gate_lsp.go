package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPProjectionBound = "generation-request-observation-metric-candidate-gate-lsp-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPDiagnostic struct {
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
    GenerationRequestReverseObservationBridgeDigest string
    ApplicationPlanStatus string
    ObservationMetricStatus string
    ObservationMetricDigest string
    ObservationMetricSource string
    ObservationMetricEvidenceDigest string
    ObservationStatus string
    ObservationMetricDirectionDigest string
    ObservationMetricCandidateGateBridgeDigest string
    CandidateGateStatus string
    CandidateDecision string
    CandidateDigest string
    RevisionSource string
    GateDigest string
    ApplicationPlanReverseObservation JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic
    ApplicationPlanReverseObservationBridgeDigest string
    BridgeDigest string
    ProjectionDigest string
    Publishable bool
    NonExecuting bool
    NonAuthorizing bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV generation request candidate gate LSP diagnostic")
    }
    if !d.NonExecuting || !d.NonAuthorizing {
        return fmt.Errorf("JEV generation request candidate gate LSP diagnostic must be non-executing and non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if d.Publishable || d.MissingStage == "" || d.BridgeDigest != "" || d.ProjectionDigest != "" {
            return fmt.Errorf("UNKNOWN JEV generation request candidate gate diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeBound ||
        !d.Publishable || d.MissingStage != "" || d.GenerationRequestStatus == "" ||
        d.GenerationRequestDigest == "" || d.GenerationRequestSource == "" ||
        d.GenerationRequestEvidenceDigest == "" || d.GeneratorIdentity == "" ||
        d.ProposalDecision == "" || d.GenerationRequestReverseObservationBridgeDigest == "" ||
        d.ApplicationPlanStatus == "" || d.ObservationMetricStatus == "" ||
        d.ObservationStatus == "" || d.ObservationMetricDirectionDigest == "" ||
        d.ObservationMetricCandidateGateBridgeDigest == "" || d.CandidateGateStatus == "" ||
        d.CandidateDecision == "" || d.GateDigest == "" || d.ApplicationPlanReverseObservationBridgeDigest == "" ||
        d.BridgeDigest == "" || d.ProjectionDigest == "" {
        return fmt.Errorf("bound JEV generation request candidate gate diagnostic is incomplete")
    }
    if err := d.ApplicationPlanReverseObservation.Validate(); err != nil {
        return fmt.Errorf("invalid generation request reverse observation diagnostic: %w", err)
    }
    if !d.ApplicationPlanReverseObservation.Publishable ||
        d.ApplicationPlanReverseObservation.BridgeDigest != d.ApplicationPlanReverseObservationBridgeDigest {
        return fmt.Errorf("generation request candidate gate nested digest is inconsistent")
    }
    if d.ApplicationPlanReverseObservation.ApplicationPlanStatus != d.ApplicationPlanStatus ||
        d.ApplicationPlanReverseObservation.ObservationMetricStatus != d.ObservationMetricStatus {
        return fmt.Errorf("generation request candidate gate status continuity is inconsistent")
    }
    switch d.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if d.CandidateDigest == "" || d.RevisionSource == "" {
            return fmt.Errorf("ready generation request candidate gate diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold, jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if d.CandidateDigest != "" {
            return fmt.Errorf("held or rejected generation request candidate gate diagnostic exposes candidate digest")
        }
    default:
        return fmt.Errorf("invalid candidate gate decision in generation request diagnostic")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPProjection(
        d.Status, d.GenerationRequestReverseObservationBridgeDigest,
        d.ObservationMetricCandidateGateBridgeDigest, d.ApplicationPlanReverseObservationBridgeDigest,
        d.BridgeDigest, d.ProposalDecision, d.CandidateDecision,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV generation request candidate gate projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPDiagnostic {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPDiagnostic {
        if stage == "" {
            stage = "generation-request-observation-metric-candidate-gate-lsp"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPDiagnostic{
            Severity: "warning",
            Code: "jev.provenance.unknown",
            Message: "Generation request candidate gate is UNKNOWN; evidence must be resolved",
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
        return unknown("generation-request-observation-metric-candidate-gate")
    }
    reverseObservation := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge{
            Status: input.ApplicationPlanReverseObservationBridgeStatus(),
            ApplicationPlanStatus: input.ApplicationPlanStatus,
            PlanDigest: input.PlanDigestFromGenerationRequest(),
            PlanSource: input.PlanSourceFromGenerationRequest(),
            PlanEvidenceDigest: input.PlanEvidenceFromGenerationRequest(),
            ApplicationPlanBridgeDigest: input.GenerationRequestApplicationPlanBridgeDigest(),
            ReverseObservationStatus: input.ReverseObservationStatus,
            ObservationDigest: input.ObservationDigestFromGenerationRequest(),
            ObservationSource: input.ObservationSourceFromGenerationRequest(),
            ObservationEvidenceDigest: input.ObservationEvidenceFromGenerationRequest(),
            ObservationMetricStatus: input.ObservationMetricStatus,
            ObservationMetricDigest: input.ObservationMetricDigest,
            ObservationMetricSource: input.ObservationMetricSource,
            ObservationMetricEvidenceDigest: input.ObservationMetricEvidenceDigest,
            BridgeDigest: input.ApplicationPlanReverseObservationBridgeDigest,
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPDiagnostic{
        Status: input.Status,
        GenerationRequestStatus: input.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GeneratorIdentity,
        ProposalDecision: input.ProposalDecision,
        GenerationRequestReverseObservationBridgeDigest: input.GenerationRequestReverseObservationBridgeDigest,
        ApplicationPlanStatus: input.ApplicationPlanStatus,
        ObservationMetricStatus: input.ObservationMetricStatus,
        ObservationMetricDigest: input.ObservationMetricDigest,
        ObservationMetricSource: input.ObservationMetricSource,
        ObservationMetricEvidenceDigest: input.ObservationMetricEvidenceDigest,
        ObservationStatus: input.ObservationStatus,
        ObservationMetricDirectionDigest: input.ObservationMetricDirectionDigest,
        ObservationMetricCandidateGateBridgeDigest: input.ObservationMetricCandidateGateBridgeDigest,
        CandidateGateStatus: input.CandidateGateStatus,
        CandidateDecision: input.CandidateDecision,
        CandidateDigest: input.CandidateDigest,
        RevisionSource: input.RevisionSource,
        GateDigest: input.GateDigest,
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
    output.Code = "jev.external-apply.generation-request-candidate-gate-bound"
    output.Message = "Generation request observation metric and candidate gate provenance are bound without execution or authorization"
    if input.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateHold {
        output.Code = "jev.external-apply.generation-request-candidate-gate-held"
        output.Message = "Generation request candidate gate is held without execution or authorization"
    }
    if input.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateRejected {
        output.Severity = "warning"
        output.Code = "jev.external-apply.generation-request-candidate-gate-rejected"
        output.Message = "Generation request candidate gate is rejected without execution or authorization"
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPProjection(
        output.Status, output.GenerationRequestReverseObservationBridgeDigest,
        output.ObservationMetricCandidateGateBridgeDigest, output.ApplicationPlanReverseObservationBridgeDigest,
        output.BridgeDigest, output.ProposalDecision, output.CandidateDecision,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-observation-metric-candidate-gate-lsp")
    }
    return output
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) ApplicationPlanReverseObservationBridgeStatus() string {
    return jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeBound
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) PlanDigestFromGenerationRequest() string { return "" }
func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) PlanSourceFromGenerationRequest() string { return "" }
func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) PlanEvidenceFromGenerationRequest() string { return "" }
func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) GenerationRequestApplicationPlanBridgeDigest() string { return "" }
func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) ObservationDigestFromGenerationRequest() string { return "" }
func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) ObservationSourceFromGenerationRequest() string { return "" }
func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) ObservationEvidenceFromGenerationRequest() string { return "" }

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPProjection(status, generationRequestReverseObservationBridgeDigest, observationMetricCandidateGateBridgeDigest, applicationPlanReverseObservationBridgeDigest, bridgeDigest, proposalDecision, candidateDecision string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, generationRequestReverseObservationBridgeDigest, observationMetricCandidateGateBridgeDigest, applicationPlanReverseObservationBridgeDigest, bridgeDigest, proposalDecision, candidateDecision}, "|")))
    return hex.EncodeToString(sum[:])
}
