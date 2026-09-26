package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeBound = "candidate-gate-revision-proposal-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeInput struct {
    ObservationMetricCandidateGate JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge
    FeedbackRevisionProposal JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge struct {
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
    NonExecuting bool
    NonAuthorizing bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge) Validate() error {
    if !b.NonExecuting || !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate gate revision proposal bridge must be non-executing and non-authorizing")
    }
    if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if b.MissingStage == "" || b.BridgeDigest != "" {
            return fmt.Errorf("unknown JEV candidate gate revision proposal bridge is inconsistent")
        }
        return nil
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeBound ||
        b.BridgeDigest == "" || b.GenerationRequestObservationMetricCandidateGateBridgeDigest == "" ||
        b.FeedbackRevisionProposalBridgeDigest == "" || b.GenerationRequestStatus == "" ||
        b.GenerationRequestDigest == "" || b.GenerationRequestSource == "" ||
        b.GenerationRequestEvidenceDigest == "" || b.GeneratorIdentity == "" ||
        b.ProposalDecision == "" || b.ApplicationPlanStatus == "" ||
        b.ObservationMetricStatus == "" || b.ObservationStatus == "" ||
        b.ObservationMetricDirectionDigest == "" || b.CandidateGateStatus == "" ||
        b.CandidateDecision == "" || b.GateDigest == "" || b.FeedbackRevisionProposalStatus == "" ||
        b.FeedbackStatus == "" || b.FeedbackBridgeDigest == "" || b.RevisionProposalStatus == "" {
        return fmt.Errorf("incomplete JEV candidate gate revision proposal bridge")
    }
    if b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated ||
        b.FeedbackRevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound {
        return fmt.Errorf("invalid candidate gate revision proposal upstream status")
    }
    if b.FeedbackRevisionProposalBridgeDigest == b.GateDigest {
        return fmt.Errorf("candidate gate and proposal bridge digests must remain distinct")
    }
    switch b.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if b.CandidateDigest == "" || b.RevisionSource == "" ||
            b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound ||
            b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound ||
            b.RevisionProposalDecision != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise ||
            b.RevisionProposalTarget == "" || b.RevisionProposalDigest == "" ||
            b.RevisionProposalSource == "" || b.RevisionProposalEvidenceDigest == "" {
            return fmt.Errorf("ready candidate gate revision proposal bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if b.CandidateDigest != "" ||
            b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld ||
            b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalHeld ||
            b.RevisionProposalDecision != "" || b.RevisionProposalTarget != "" ||
            b.RevisionProposalDigest != "" || b.RevisionProposalSource != "" ||
            b.RevisionProposalEvidenceDigest != "" {
            return fmt.Errorf("held candidate gate revision proposal bridge is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if b.CandidateDigest != "" ||
            b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected ||
            b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRejected ||
            b.RevisionProposalDecision != "" || b.RevisionProposalTarget != "" ||
            b.RevisionProposalDigest != "" || b.RevisionProposalSource != "" ||
            b.RevisionProposalEvidenceDigest != "" {
            return fmt.Errorf("rejected candidate gate revision proposal bridge is inconsistent")
        }
    default:
        return fmt.Errorf("invalid candidate gate decision")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge(
        b.Status, b.GenerationRequestObservationMetricCandidateGateBridgeDigest,
        b.FeedbackRevisionProposalBridgeDigest, b.ProposalDecision,
        b.CandidateDecision, b.RevisionProposalStatus, b.RevisionProposalDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate gate revision proposal bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge {
        if stage == "" {
            stage = "candidate-gate-revision-proposal"
        }
        return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: stage,
            NonExecuting: true,
            NonAuthorizing: true,
        }
    }
    if !input.NonAuthorizing || !input.ObservationMetricCandidateGate.NonExecuting ||
        !input.ObservationMetricCandidateGate.NonAuthorizing || !input.FeedbackRevisionProposal.NonExecuting ||
        !input.FeedbackRevisionProposal.NonAuthorizing {
        return unknown("capability-boundary")
    }
    if err := input.ObservationMetricCandidateGate.Validate(); err != nil {
        return unknown("observation-metric-candidate-gate")
    }
    if err := input.FeedbackRevisionProposal.Validate(); err != nil {
        return unknown("feedback-revision-proposal")
    }
    if input.FeedbackRevisionProposal.CandidateGateDigest != input.ObservationMetricCandidateGate.GateDigest {
        return unknown("candidate-gate-digest-consistency")
    }
    if input.ObservationMetricCandidateGate.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateReady &&
        (input.FeedbackRevisionProposal.CandidateDigest != input.ObservationMetricCandidateGate.CandidateDigest ||
            input.FeedbackRevisionProposal.CandidateSource != input.ObservationMetricCandidateGate.RevisionSource) {
        return unknown("candidate-provenance-consistency")
    }
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge{
        Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeBound,
        GenerationRequestStatus: input.ObservationMetricCandidateGate.GenerationRequestStatus,
        GenerationRequestDigest: input.ObservationMetricCandidateGate.GenerationRequestDigest,
        GenerationRequestSource: input.ObservationMetricCandidateGate.GenerationRequestSource,
        GenerationRequestEvidenceDigest: input.ObservationMetricCandidateGate.GenerationRequestEvidenceDigest,
        GeneratorIdentity: input.ObservationMetricCandidateGate.GeneratorIdentity,
        ProposalDecision: input.ObservationMetricCandidateGate.ProposalDecision,
        GenerationRequestObservationMetricCandidateGateBridgeDigest: input.ObservationMetricCandidateGate.BridgeDigest,
        ApplicationPlanStatus: input.ObservationMetricCandidateGate.ApplicationPlanStatus,
        ObservationMetricStatus: input.ObservationMetricCandidateGate.ObservationMetricStatus,
        ObservationMetricDigest: input.ObservationMetricCandidateGate.ObservationMetricDigest,
        ObservationMetricSource: input.ObservationMetricCandidateGate.ObservationMetricSource,
        ObservationMetricEvidenceDigest: input.ObservationMetricCandidateGate.ObservationMetricEvidenceDigest,
        ObservationStatus: input.ObservationMetricCandidateGate.ObservationStatus,
        ObservationMetricDirectionDigest: input.ObservationMetricCandidateGate.ObservationMetricDirectionDigest,
        CandidateGateStatus: input.ObservationMetricCandidateGate.CandidateGateStatus,
        CandidateDecision: input.ObservationMetricCandidateGate.CandidateDecision,
        CandidateDigest: input.ObservationMetricCandidateGate.CandidateDigest,
        RevisionSource: input.ObservationMetricCandidateGate.RevisionSource,
        GateDigest: input.ObservationMetricCandidateGate.GateDigest,
        FeedbackRevisionProposalStatus: input.FeedbackRevisionProposal.Status,
        FeedbackStatus: input.FeedbackRevisionProposal.FeedbackStatus,
        FeedbackBridgeDigest: input.FeedbackRevisionProposal.FeedbackBridgeDigest,
        RevisionProposalStatus: input.FeedbackRevisionProposal.RevisionProposalStatus,
        RevisionProposalDecision: input.FeedbackRevisionProposal.ProposalDecision,
        RevisionProposalTarget: input.FeedbackRevisionProposal.ProposalTarget,
        RevisionProposalDigest: input.FeedbackRevisionProposal.ProposalDigest,
        RevisionProposalSource: input.FeedbackRevisionProposal.ProposalSource,
        RevisionProposalEvidenceDigest: input.FeedbackRevisionProposal.ProposalEvidenceDigest,
        FeedbackRevisionProposalBridgeDigest: input.FeedbackRevisionProposal.BridgeDigest,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge(
        output.Status, output.GenerationRequestObservationMetricCandidateGateBridgeDigest,
        output.FeedbackRevisionProposalBridgeDigest, output.ProposalDecision,
        output.CandidateDecision, output.RevisionProposalStatus, output.RevisionProposalDigest,
    )
    if err := output.Validate(); err != nil {
        return unknown("candidate-gate-revision-proposal-evidence")
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge(status, observationMetricCandidateGateBridgeDigest, feedbackRevisionProposalBridgeDigest, proposalDecision, candidateDecision, revisionProposalStatus, revisionProposalDigest string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{status, observationMetricCandidateGateBridgeDigest, feedbackRevisionProposalBridgeDigest, proposalDecision, candidateDecision, revisionProposalStatus, revisionProposalDigest}, "|")))
    return hex.EncodeToString(sum[:])
}
