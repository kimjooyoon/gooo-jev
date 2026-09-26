package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeBound = "generation-request-application-plan-reverse-observation-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeInput struct {
    GenerationRequestApplicationPlan JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge
    ApplicationPlanReverseObservation JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge struct {
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
    BridgeDigest string
    NonExecuting bool
    NonAuthorizing bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge) Validate() error {
    if !b.NonExecuting || !b.NonAuthorizing {
        return fmt.Errorf("JEV generation request application plan reverse observation bridge must be non-executing and non-authorizing")
    }
    if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if b.MissingStage == "" || b.BridgeDigest != "" {
            return fmt.Errorf("unknown JEV generation request application plan reverse observation bridge is inconsistent")
        }
        return nil
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeBound ||
        b.BridgeDigest == "" || b.FeedbackBridgeDigest == "" ||
        b.CandidateGateApplicationPlanBridgeDigest == "" ||
        b.ApplicationPlanReverseObservationBridgeDigest == "" ||
        b.GenerationRequestStatus == "" || b.GenerationRequestDigest == "" ||
        b.GenerationRequestSource == "" || b.GenerationRequestEvidenceDigest == "" ||
        b.GeneratorIdentity == "" || b.ProposalDecision == "" ||
        b.ApplicationPlanStatus == "" {
        return fmt.Errorf("incomplete JEV generation request application plan reverse observation bridge")
    }
    if b.ApplicationPlanStatus == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady {
        if b.PlanDigest == "" || b.PlanSource == "" || b.PlanEvidenceDigest == "" ||
            b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound ||
            b.ObservationDigest == "" || b.ObservationSource == "" || b.ObservationEvidenceDigest == "" ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound ||
            b.ObservationMetricDigest == "" || b.ObservationMetricSource == "" ||
            b.ObservationMetricEvidenceDigest == "" {
            return fmt.Errorf("ready generation request reverse observation bridge is incomplete")
        }
    } else if b.ApplicationPlanStatus == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld {
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationHeld ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricHeld ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" ||
            b.ObservationDigest != "" || b.ObservationSource != "" || b.ObservationEvidenceDigest != "" ||
            b.ObservationMetricDigest != "" || b.ObservationMetricSource != "" || b.ObservationMetricEvidenceDigest != "" {
            return fmt.Errorf("held generation request reverse observation bridge is inconsistent")
        }
    } else if b.ApplicationPlanStatus == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected {
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationRejected ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricRejected ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" ||
            b.ObservationDigest != "" || b.ObservationSource != "" || b.ObservationEvidenceDigest != "" ||
            b.ObservationMetricDigest != "" || b.ObservationMetricSource != "" || b.ObservationMetricEvidenceDigest != "" {
            return fmt.Errorf("rejected generation request reverse observation bridge is inconsistent")
        }
    } else {
        return fmt.Errorf("invalid generation request application plan status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge(
        b.Status, b.FeedbackBridgeDigest, b.CandidateGateApplicationPlanBridgeDigest,
        b.ApplicationPlanReverseObservationBridgeDigest, b.ProposalDecision, b.ApplicationPlanStatus,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV generation request application plan reverse observation bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge {
        if stage == "" {
            stage = "generation-request-application-plan-reverse-observation"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: stage,
            NonExecuting: true,
            NonAuthorizing: true,
        }
    }
    if !input.NonAuthorizing || !input.GenerationRequestApplicationPlan.NonExecuting ||
        !input.GenerationRequestApplicationPlan.NonAuthorizing || !input.ApplicationPlanReverseObservation.NonExecuting ||
        !input.ApplicationPlanReverseObservation.NonAuthorizing {
        return unknown("capability-boundary")
    }
    if err := input.GenerationRequestApplicationPlan.Validate(); err != nil {
        return unknown("generation-request-application-plan")
    }
    if err := input.ApplicationPlanReverseObservation.Validate(); err != nil {
        return unknown("application-plan-reverse-observation")
    }
    if input.GenerationRequestApplicationPlan.BridgeDigest != input.ApplicationPlanReverseObservation.ApplicationPlanBridgeDigest {
        return unknown("application-plan-digest-consistency")
    }
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge{
        Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeBound,
        GenerationRequestStatus: input.GenerationRequestApplicationPlan.GenerationRequestStatus,
        GenerationRequestDigest: input.GenerationRequestApplicationPlan.GenerationRequestDigest,
        GenerationRequestSource: input.GenerationRequestApplicationPlan.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.GenerationRequestApplicationPlan.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.GenerationRequestApplicationPlan.GeneratorIdentity,
        ProposalDecision: input.GenerationRequestApplicationPlan.ProposalDecision,
        FeedbackBridgeDigest: input.GenerationRequestApplicationPlan.FeedbackBridgeDigest,
        CandidateGateApplicationPlanBridgeDigest: input.GenerationRequestApplicationPlan.CandidateGateApplicationPlanBridgeDigest,
        ApplicationPlanStatus: input.GenerationRequestApplicationPlan.ApplicationPlanStatus,
        PlanDigest: input.ApplicationPlanReverseObservation.PlanDigest,
        PlanSource: input.ApplicationPlanReverseObservation.PlanSource,
        PlanEvidenceDigest: input.ApplicationPlanReverseObservation.PlanEvidenceDigest,
        ReverseObservationStatus: input.ApplicationPlanReverseObservation.ReverseObservationStatus,
        ObservationDigest: input.ApplicationPlanReverseObservation.ObservationDigest,
        ObservationSource: input.ApplicationPlanReverseObservation.ObservationSource,
        ObservationEvidenceDigest: input.ApplicationPlanReverseObservation.ObservationEvidenceDigest,
        ObservationMetricStatus: input.ApplicationPlanReverseObservation.ObservationMetricStatus,
        ObservationMetricDigest: input.ApplicationPlanReverseObservation.ObservationMetricDigest,
        ObservationMetricSource: input.ApplicationPlanReverseObservation.ObservationMetricSource,
        ObservationMetricEvidenceDigest: input.ApplicationPlanReverseObservation.ObservationMetricEvidenceDigest,
        ApplicationPlanReverseObservationBridgeDigest: input.ApplicationPlanReverseObservation.BridgeDigest,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if err := output.Validate(); err != nil {
        return unknown("generation-request-application-plan-reverse-observation-consistency")
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge(
        output.Status, output.FeedbackBridgeDigest, output.CandidateGateApplicationPlanBridgeDigest,
        output.ApplicationPlanReverseObservationBridgeDigest, output.ProposalDecision, output.ApplicationPlanStatus,
    )
    if err := output.Validate(); err != nil {
        return unknown("generation-request-application-plan-reverse-observation-evidence")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge(status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, applicationPlanReverseObservationBridgeDigest, proposalDecision, applicationPlanStatus string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, feedbackBridgeDigest, candidateGateApplicationPlanBridgeDigest, applicationPlanReverseObservationBridgeDigest, proposalDecision, applicationPlanStatus}, "|")))
    return hex.EncodeToString(sum[:])
}
