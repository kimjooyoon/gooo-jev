package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPProjectionBound = "candidate-gate-revision-proposal-lsp-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPDiagnostic struct {
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
    GenerationRequestObservationMetricCandidateGateBridgeDigest string
    ApplicationPlanStatus string
    ObservationMetricStatus string
    ObservationMetricDigest string
    ObservationMetricSource string
    ObservationMetricEvidenceDigest string
    ObservationStatus string
    ObservationMetricDirectionDigest string
    CandidateGateStatus string
    CandidateDecision string
    CandidateDigest string
    RevisionSource string
    GateDigest string
    FeedbackRevisionProposalStatus string
    FeedbackStatus string
    FeedbackBridgeDigest string
    RevisionProposalStatus string
    RevisionProposalDecision string
    RevisionProposalTarget string
    RevisionProposalDigest string
    RevisionProposalSource string
    RevisionProposalEvidenceDigest string
    FeedbackRevisionProposalBridgeDigest string
    BridgeDigest string
    ProjectionDigest string
    Publishable bool
    NonExecuting bool
    NonAuthorizing bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate gate revision proposal LSP diagnostic")
    }
    if !d.NonExecuting || !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate gate revision proposal LSP diagnostic must be non-executing and non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if d.Publishable || d.MissingStage == "" || d.BridgeDigest != "" || d.ProjectionDigest != "" {
            return fmt.Errorf("UNKNOWN JEV candidate gate revision proposal diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeBound ||
        !d.Publishable || d.MissingStage != "" || d.GenerationRequestStatus == "" ||
        d.GenerationRequestDigest == "" || d.GenerationRequestSource == "" ||
        d.GenerationRequestEvidenceDigest == "" || d.GeneratorIdentity == "" ||
        d.ProposalDecision == "" || d.GenerationRequestObservationMetricCandidateGateBridgeDigest == "" ||
        d.ApplicationPlanStatus == "" || d.ObservationMetricStatus == "" ||
        d.ObservationStatus == "" || d.ObservationMetricDirectionDigest == "" ||
        d.CandidateGateStatus == "" || d.CandidateDecision == "" || d.GateDigest == "" ||
        d.FeedbackRevisionProposalStatus == "" || d.FeedbackStatus == "" ||
        d.FeedbackBridgeDigest == "" || d.RevisionProposalStatus == "" ||
        d.FeedbackRevisionProposalBridgeDigest == "" || d.BridgeDigest == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("bound JEV candidate gate revision proposal diagnostic is incomplete")
    }
    if d.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated ||
        d.FeedbackRevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound {
        return fmt.Errorf("invalid candidate gate revision proposal upstream status")
    }
    switch d.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if d.CandidateDigest == "" || d.RevisionSource == "" ||
            d.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound ||
            d.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound ||
            d.RevisionProposalDecision != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise ||
            d.RevisionProposalTarget == "" || d.RevisionProposalDigest == "" ||
            d.RevisionProposalSource == "" || d.RevisionProposalEvidenceDigest == "" {
            return fmt.Errorf("ready candidate gate revision proposal diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if d.CandidateDigest != "" ||
            d.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld ||
            d.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalHeld ||
            d.RevisionProposalDecision != "" || d.RevisionProposalTarget != "" ||
            d.RevisionProposalDigest != "" || d.RevisionProposalSource != "" ||
            d.RevisionProposalEvidenceDigest != "" {
            return fmt.Errorf("held candidate gate revision proposal diagnostic is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if d.CandidateDigest != "" ||
            d.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected ||
            d.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRejected ||
            d.RevisionProposalDecision != "" || d.RevisionProposalTarget != "" ||
            d.RevisionProposalDigest != "" || d.RevisionProposalSource != "" ||
            d.RevisionProposalEvidenceDigest != "" {
            return fmt.Errorf("rejected candidate gate revision proposal diagnostic is inconsistent")
        }
    default:
        return fmt.Errorf("invalid candidate gate decision in LSP diagnostic")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPProjection(
        d.Status, d.GenerationRequestObservationMetricCandidateGateBridgeDigest,
        d.FeedbackRevisionProposalBridgeDigest, d.BridgeDigest,
        d.ProposalDecision, d.CandidateDecision,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV candidate gate revision proposal projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPDiagnostic {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPDiagnostic {
        if stage == "" {
            stage = "candidate-gate-revision-proposal-lsp"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPDiagnostic{
            Severity: "warning",
            Code: "jev.provenance.unknown",
            Message: "Candidate gate revision proposal is UNKNOWN; evidence must be resolved",
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
        return unknown("candidate-gate-revision-proposal")
    }
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPDiagnostic{
        Status: input.Status,
        GenerationRequestStatus: input.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GeneratorIdentity,
        ProposalDecision: input.ProposalDecision,
        GenerationRequestObservationMetricCandidateGateBridgeDigest: input.GenerationRequestObservationMetricCandidateGateBridgeDigest,
        ApplicationPlanStatus: input.ApplicationPlanStatus,
        ObservationMetricStatus: input.ObservationMetricStatus,
        ObservationMetricDigest: input.ObservationMetricDigest,
        ObservationMetricSource: input.ObservationMetricSource,
        ObservationMetricEvidenceDigest: input.ObservationMetricEvidenceDigest,
        ObservationStatus: input.ObservationStatus,
        ObservationMetricDirectionDigest: input.ObservationMetricDirectionDigest,
        CandidateGateStatus: input.CandidateGateStatus,
        CandidateDecision: input.CandidateDecision,
        CandidateDigest: input.CandidateDigest,
        RevisionSource: input.RevisionSource,
        GateDigest: input.GateDigest,
        FeedbackRevisionProposalStatus: input.FeedbackRevisionProposalStatus,
        FeedbackStatus: input.FeedbackStatus,
        FeedbackBridgeDigest: input.FeedbackBridgeDigest,
        RevisionProposalStatus: input.RevisionProposalStatus,
        RevisionProposalDecision: input.RevisionProposalDecision,
        RevisionProposalTarget: input.RevisionProposalTarget,
        RevisionProposalDigest: input.RevisionProposalDigest,
        RevisionProposalSource: input.RevisionProposalSource,
        RevisionProposalEvidenceDigest: input.RevisionProposalEvidenceDigest,
        FeedbackRevisionProposalBridgeDigest: input.FeedbackRevisionProposalBridgeDigest,
        BridgeDigest: input.BridgeDigest,
        Publishable: false,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-gate-revision-proposal-bound"
    output.Message = "Candidate gate revision proposal provenance is bound without execution or authorization"
    if input.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateHold {
        output.Code = "jev.external-apply.candidate-gate-revision-proposal-held"
        output.Message = "Candidate gate revision proposal is held without execution or authorization"
    }
    if input.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateRejected {
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-gate-revision-proposal-rejected"
        output.Message = "Candidate gate revision proposal is rejected without execution or authorization"
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPProjection(
        output.Status, output.GenerationRequestObservationMetricCandidateGateBridgeDigest,
        output.FeedbackRevisionProposalBridgeDigest, output.BridgeDigest,
        output.ProposalDecision, output.CandidateDecision,
    )
    if err := output.Validate(); err != nil {
        return unknown("candidate-gate-revision-proposal-lsp")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalLSPProjection(status, observationMetricCandidateGateBridgeDigest, feedbackRevisionProposalBridgeDigest, bridgeDigest, proposalDecision, candidateDecision string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, observationMetricCandidateGateBridgeDigest, feedbackRevisionProposalBridgeDigest, bridgeDigest, proposalDecision, candidateDecision}, "|")))
    return hex.EncodeToString(sum[:])
}
