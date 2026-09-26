package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeBound = "generation-request-observation-metric-candidate-gate-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeInput struct {
    GenerationRequestReverseObservation JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge
    ObservationMetricCandidateGate JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge struct {
    Status string
    MissingStage string
    GenerationRequestStatus string
    GenerationRequestDigest string
    GenerationRequestSource string
    GenerationRequestEvidenceDigest string
    GeneratorIdentity string
    ProposalDecision string
    GenerationRequestReverseObservationBridgeDigest string
    CandidateGateApplicationPlanBridgeDigest string
    ApplicationPlanStatus string
    PlanDigest string
    PlanSource string
    PlanEvidenceDigest string
    ReverseObservationStatus string
    ObservationDigest string
    ObservationSource string
    ObservationEvidenceDigest string
    ObservationMetricStatus string
    ObservationMetricDigest string
    ObservationMetricSource string
    ObservationMetricEvidenceDigest string
    ApplicationPlanReverseObservationBridgeDigest string
    ObservationStatus string
    ObservationMetricDirectionDigest string
    ObservationMetricCandidateGateBridgeDigest string
    CandidateGateStatus string
    CandidateDecision string
    CandidateDigest string
    RevisionSource string
    GateDigest string
    BridgeDigest string
    NonExecuting bool
    NonAuthorizing bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge) Validate() error {
    if !b.NonExecuting || !b.NonAuthorizing {
        return fmt.Errorf("JEV generation request observation metric candidate gate bridge must be non-executing and non-authorizing")
    }
    if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if b.MissingStage == "" || b.BridgeDigest != "" {
            return fmt.Errorf("unknown JEV generation request observation metric candidate gate bridge is inconsistent")
        }
        return nil
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeBound ||
        b.BridgeDigest == "" || b.GenerationRequestReverseObservationBridgeDigest == "" ||
        b.ObservationMetricCandidateGateBridgeDigest == "" || b.GenerationRequestStatus == "" ||
        b.GenerationRequestDigest == "" || b.GenerationRequestSource == "" ||
        b.GenerationRequestEvidenceDigest == "" || b.GeneratorIdentity == "" ||
        b.ProposalDecision == "" || b.CandidateGateApplicationPlanBridgeDigest == "" ||
        b.ApplicationPlanStatus == "" || b.ApplicationPlanReverseObservationBridgeDigest == "" ||
        b.ReverseObservationStatus == "" || b.ObservationMetricStatus == "" ||
        b.ObservationStatus == "" || b.ObservationMetricDirectionDigest == "" ||
        b.CandidateGateStatus == "" || b.CandidateDecision == "" || b.GateDigest == "" {
        return fmt.Errorf("incomplete JEV generation request observation metric candidate gate bridge")
    }
    if b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("invalid JEV generation request candidate gate status")
    }
    switch b.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if b.CandidateDigest == "" || b.RevisionSource == "" {
            return fmt.Errorf("ready candidate gate requires candidate provenance")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold, jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if b.CandidateDigest != "" {
            return fmt.Errorf("held or rejected candidate gate must not expose candidate digest")
        }
    default:
        return fmt.Errorf("invalid candidate gate decision")
    }
    switch b.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound ||
            b.ObservationMetricDigest == "" || b.ObservationMetricSource == "" ||
            b.ObservationMetricEvidenceDigest == "" {
            return fmt.Errorf("ready generation request observation metric candidate gate bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationHeld ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricHeld ||
            b.ObservationMetricDigest != "" || b.ObservationMetricSource != "" ||
            b.ObservationMetricEvidenceDigest != "" {
            return fmt.Errorf("held generation request observation metric candidate gate bridge is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationRejected ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricRejected ||
            b.ObservationMetricDigest != "" || b.ObservationMetricSource != "" ||
            b.ObservationMetricEvidenceDigest != "" {
            return fmt.Errorf("rejected generation request observation metric candidate gate bridge is inconsistent")
        }
    default:
        return fmt.Errorf("invalid application plan status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge(
        b.Status, b.GenerationRequestReverseObservationBridgeDigest,
        b.ObservationMetricCandidateGateBridgeDigest, b.ProposalDecision,
        b.ApplicationPlanStatus, b.CandidateDecision,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV generation request observation metric candidate gate bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge {
        if stage == "" {
            stage = "generation-request-observation-metric-candidate-gate"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: stage,
            NonExecuting: true,
            NonAuthorizing: true,
        }
    }
    if !input.NonAuthorizing || !input.GenerationRequestReverseObservation.NonExecuting ||
        !input.GenerationRequestReverseObservation.NonAuthorizing || !input.ObservationMetricCandidateGate.NonExecuting ||
        !input.ObservationMetricCandidateGate.NonAuthorizing {
        return unknown("capability-boundary")
    }
    if err := input.GenerationRequestReverseObservation.Validate(); err != nil {
        return unknown("generation-request-reverse-observation")
    }
    if err := input.ObservationMetricCandidateGate.Validate(); err != nil {
        return unknown("observation-metric-candidate-gate")
    }
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge{
        Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeBound,
        GenerationRequestStatus: input.GenerationRequestReverseObservation.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestReverseObservation.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestReverseObservation.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestReverseObservation.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GenerationRequestReverseObservation.GeneratorIdentity,
        ProposalDecision: input.GenerationRequestReverseObservation.ProposalDecision,
        GenerationRequestReverseObservationBridgeDigest: input.GenerationRequestReverseObservation.BridgeDigest,
        CandidateGateApplicationPlanBridgeDigest: input.GenerationRequestReverseObservation.CandidateGateApplicationPlanBridgeDigest,
        ApplicationPlanStatus: input.GenerationRequestReverseObservation.ApplicationPlanStatus,
        PlanDigest: input.GenerationRequestReverseObservation.PlanDigest,
        PlanSource: input.GenerationRequestReverseObservation.PlanSource,
        PlanEvidenceDigest: input.GenerationRequestReverseObservation.PlanEvidenceDigest,
        ReverseObservationStatus: input.GenerationRequestReverseObservation.ReverseObservationStatus,
        ObservationDigest: input.GenerationRequestReverseObservation.ObservationDigest,
        ObservationSource: input.GenerationRequestReverseObservation.ObservationSource,
        ObservationEvidenceDigest: input.GenerationRequestReverseObservation.ObservationEvidenceDigest,
        ObservationMetricStatus: input.GenerationRequestReverseObservation.ObservationMetricStatus,
        ObservationMetricDigest: input.GenerationRequestReverseObservation.ObservationMetricDigest,
        ObservationMetricSource: input.GenerationRequestReverseObservation.ObservationMetricSource,
        ObservationMetricEvidenceDigest: input.GenerationRequestReverseObservation.ObservationMetricEvidenceDigest,
        ApplicationPlanReverseObservationBridgeDigest: input.GenerationRequestReverseObservation.BridgeDigest,
        ObservationStatus: input.ObservationMetricCandidateGate.ObservationStatus,
        ObservationMetricDirectionDigest: input.ObservationMetricCandidateGate.ObservationMetricDirectionDigest,
        ObservationMetricCandidateGateBridgeDigest: input.ObservationMetricCandidateGate.BridgeDigest,
        CandidateGateStatus: input.ObservationMetricCandidateGate.CandidateGateStatus,
        CandidateDecision: input.ObservationMetricCandidateGate.CandidateDecision,
        CandidateDigest: input.ObservationMetricCandidateGate.CandidateDigest,
        RevisionSource: input.ObservationMetricCandidateGate.RevisionSource,
        GateDigest: input.ObservationMetricCandidateGate.GateDigest,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if err := output.Validate(); err != nil {
        return unknown("generation-request-observation-metric-candidate-gate-consistency")
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge(
        output.Status, output.GenerationRequestReverseObservationBridgeDigest,
        output.ObservationMetricCandidateGateBridgeDigest, output.ProposalDecision,
        output.ApplicationPlanStatus, output.CandidateDecision,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-observation-metric-candidate-gate-evidence")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge(status, generationRequestReverseObservationBridgeDigest, observationMetricCandidateGateBridgeDigest, proposalDecision, applicationPlanStatus, candidateDecision string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, generationRequestReverseObservationBridgeDigest, observationMetricCandidateGateBridgeDigest, proposalDecision, applicationPlanStatus, candidateDecision}, "|")))
    return hex.EncodeToString(sum[:])
}
