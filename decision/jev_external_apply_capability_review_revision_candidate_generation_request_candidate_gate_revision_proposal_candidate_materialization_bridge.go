package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeBound = "candidate-gate-revision-proposal-materialization-bound"

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeInput struct {
    RevisionProposal JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge
    CandidateMaterialization JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge struct {
    Status string
    MissingStage string
    GenerationRequestStatus string
    GenerationRequestDigest string
    GenerationRequestSource string
    GenerationRequestEvidenceDigest string
    GeneratorIdentity string
    ProposalDecision string
    GenerationRequestCandidateGateRevisionProposalBridgeDigest string
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
    MaterializationStatus string
    GeneratedCandidateDigest string
    GeneratedCandidateSource string
    GenerationInputDigest string
    GenerationEvidenceDigest string
    CandidateGeneratorIdentity string
    CandidateMaterializationBridgeDigest string
    BridgeDigest string
    NonExecuting bool
    NonAuthorizing bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge) Validate() error {
    if !b.NonExecuting || !b.NonAuthorizing { return fmt.Errorf("JEV candidate materialization bridge must be non-executing and non-authorizing") }
    if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
        if b.MissingStage == "" || b.BridgeDigest != "" { return fmt.Errorf("unknown JEV candidate materialization bridge is inconsistent") }
        return nil
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeBound || b.BridgeDigest == "" || b.GenerationRequestCandidateGateRevisionProposalBridgeDigest == "" || b.CandidateMaterializationBridgeDigest == "" || b.GenerationRequestStatus == "" || b.GenerationRequestDigest == "" || b.GenerationRequestSource == "" || b.GenerationRequestEvidenceDigest == "" || b.GeneratorIdentity == "" || b.ProposalDecision == "" || b.ApplicationPlanStatus == "" || b.ObservationMetricStatus == "" || b.ObservationStatus == "" || b.ObservationMetricDirectionDigest == "" || b.CandidateGateStatus == "" || b.CandidateDecision == "" || b.GateDigest == "" || b.FeedbackRevisionProposalStatus == "" || b.FeedbackStatus == "" || b.FeedbackBridgeDigest == "" || b.RevisionProposalStatus == "" || b.MaterializationStatus == "" { return fmt.Errorf("incomplete JEV candidate materialization bridge") }
    if b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated || b.FeedbackRevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound || b.CandidateMaterializationBridgeDigest == b.FeedbackRevisionProposalBridgeDigest { return fmt.Errorf("invalid candidate materialization upstream status") }
    switch b.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if b.CandidateDigest == "" || b.RevisionSource == "" || b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound || b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound || b.RevisionProposalDecision != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise || b.RevisionProposalTarget == "" || b.RevisionProposalDigest == "" || b.RevisionProposalSource == "" || b.RevisionProposalEvidenceDigest == "" || b.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized || b.GeneratedCandidateDigest == "" || b.GeneratedCandidateSource == "" || b.GenerationInputDigest == "" || b.GenerationEvidenceDigest == "" || b.CandidateGeneratorIdentity == "" { return fmt.Errorf("materialized candidate bridge is incomplete") }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if b.CandidateDigest != "" || b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld || b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalHeld || b.RevisionProposalDecision != "" || b.RevisionProposalTarget != "" || b.RevisionProposalDigest != "" || b.RevisionProposalSource != "" || b.RevisionProposalEvidenceDigest != "" || b.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateHeld || b.GeneratedCandidateDigest != "" || b.GeneratedCandidateSource != "" || b.GenerationInputDigest != "" || b.GenerationEvidenceDigest != "" || b.CandidateGeneratorIdentity != "" { return fmt.Errorf("held candidate bridge is inconsistent") }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if b.CandidateDigest != "" || b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected || b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRejected || b.RevisionProposalDecision != "" || b.RevisionProposalTarget != "" || b.RevisionProposalDigest != "" || b.RevisionProposalSource != "" || b.RevisionProposalEvidenceDigest != "" || b.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateRollback || b.GeneratedCandidateDigest != "" || b.GeneratedCandidateSource != "" || b.GenerationInputDigest != "" || b.GenerationEvidenceDigest != "" || b.CandidateGeneratorIdentity != "" { return fmt.Errorf("rollback candidate bridge is inconsistent") }
    default: return fmt.Errorf("invalid candidate decision")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge(b.Status,b.GenerationRequestCandidateGateRevisionProposalBridgeDigest,b.FeedbackRevisionProposalBridgeDigest,b.CandidateMaterializationBridgeDigest,b.ProposalDecision,b.CandidateDecision,b.MaterializationStatus)
    if b.BridgeDigest != expected { return fmt.Errorf("JEV candidate materialization bridge digest mismatch") }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge {
    unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge { if stage == "" { stage = "candidate-gate-revision-proposal-materialization" }; return JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge{Status:jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,MissingStage:stage,NonExecuting:true,NonAuthorizing:true} }
    if !input.NonAuthorizing || !input.RevisionProposal.NonExecuting || !input.RevisionProposal.NonAuthorizing || !input.CandidateMaterialization.NonExecuting || !input.CandidateMaterialization.NonAuthorizing { return unknown("capability-boundary") }
    if err := input.RevisionProposal.Validate(); err != nil { return unknown("candidate-gate-revision-proposal") }
    if err := input.CandidateMaterialization.Validate(); err != nil { return unknown("candidate-materialization") }
    if input.CandidateMaterialization.CandidateGateDigest != input.RevisionProposal.GateDigest { return unknown("candidate-gate-digest-consistency") }
    if input.CandidateMaterialization.FeedbackBridgeDigest != input.RevisionProposal.FeedbackBridgeDigest { return unknown("feedback-digest-consistency") }
    if input.CandidateMaterialization.RevisionProposalBridgeDigest != input.RevisionProposal.FeedbackRevisionProposalBridgeDigest { return unknown("revision-proposal-digest-consistency") }
    if input.CandidateMaterialization.ProposalDecision != input.RevisionProposal.RevisionProposalDecision || input.CandidateMaterialization.ProposalTarget != input.RevisionProposal.RevisionProposalTarget || input.CandidateMaterialization.ProposalDigest != input.RevisionProposal.RevisionProposalDigest || input.CandidateMaterialization.ProposalSource != input.RevisionProposal.RevisionProposalSource || input.CandidateMaterialization.ProposalEvidenceDigest != input.RevisionProposal.RevisionProposalEvidenceDigest { return unknown("proposal-evidence-consistency") }
    if input.RevisionProposal.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateReady && (input.CandidateMaterialization.CandidateDigest != input.RevisionProposal.CandidateDigest || input.CandidateMaterialization.CandidateSource != input.RevisionProposal.RevisionSource) { return unknown("candidate-provenance-consistency") }
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge{Status:jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeBound,GenerationRequestStatus:input.RevisionProposal.GenerationRequestStatus,GenerationRequestDigest:input.RevisionProposal.GenerationRequestDigest,GenerationRequestSource:input.RevisionProposal.GenerationRequestSource,GenerationRequestEvidenceDigest:input.RevisionProposal.GenerationRequestEvidenceDigest,GeneratorIdentity:input.RevisionProposal.GeneratorIdentity,ProposalDecision:input.RevisionProposal.ProposalDecision,GenerationRequestCandidateGateRevisionProposalBridgeDigest:input.RevisionProposal.BridgeDigest,ApplicationPlanStatus:input.RevisionProposal.ApplicationPlanStatus,ObservationMetricStatus:input.RevisionProposal.ObservationMetricStatus,ObservationMetricDigest:input.RevisionProposal.ObservationMetricDigest,ObservationMetricSource:input.RevisionProposal.ObservationMetricSource,ObservationMetricEvidenceDigest:input.RevisionProposal.ObservationMetricEvidenceDigest,ObservationStatus:input.RevisionProposal.ObservationStatus,ObservationMetricDirectionDigest:input.RevisionProposal.ObservationMetricDirectionDigest,CandidateGateStatus:input.RevisionProposal.CandidateGateStatus,CandidateDecision:input.RevisionProposal.CandidateDecision,CandidateDigest:input.RevisionProposal.CandidateDigest,RevisionSource:input.RevisionProposal.RevisionSource,GateDigest:input.RevisionProposal.GateDigest,FeedbackRevisionProposalStatus:input.RevisionProposal.FeedbackRevisionProposalStatus,FeedbackStatus:input.RevisionProposal.FeedbackStatus,FeedbackBridgeDigest:input.RevisionProposal.FeedbackBridgeDigest,RevisionProposalStatus:input.RevisionProposal.RevisionProposalStatus,RevisionProposalDecision:input.RevisionProposal.RevisionProposalDecision,RevisionProposalTarget:input.RevisionProposal.RevisionProposalTarget,RevisionProposalDigest:input.RevisionProposal.RevisionProposalDigest,RevisionProposalSource:input.RevisionProposal.RevisionProposalSource,RevisionProposalEvidenceDigest:input.RevisionProposal.RevisionProposalEvidenceDigest,FeedbackRevisionProposalBridgeDigest:input.RevisionProposal.FeedbackRevisionProposalBridgeDigest,MaterializationStatus:input.CandidateMaterialization.MaterializationStatus,GeneratedCandidateDigest:input.CandidateMaterialization.GeneratedCandidateDigest,GeneratedCandidateSource:input.CandidateMaterialization.GeneratedCandidateSource,GenerationInputDigest:input.CandidateMaterialization.GenerationInputDigest,GenerationEvidenceDigest:input.CandidateMaterialization.GenerationEvidenceDigest,CandidateGeneratorIdentity:input.CandidateMaterialization.GeneratorIdentity,CandidateMaterializationBridgeDigest:input.CandidateMaterialization.BridgeDigest,NonExecuting:true,NonAuthorizing:true}
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge(output.Status,output.GenerationRequestCandidateGateRevisionProposalBridgeDigest,output.FeedbackRevisionProposalBridgeDigest,output.CandidateMaterializationBridgeDigest,output.ProposalDecision,output.CandidateDecision,output.MaterializationStatus)
    if err := output.Validate(); err != nil { return unknown("candidate-materialization-consistency") }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge(status,generationRequestCandidateGateRevisionProposalBridgeDigest,feedbackRevisionProposalBridgeDigest,candidateMaterializationBridgeDigest,proposalDecision,candidateDecision,materializationStatus string) string { sum := sha256.Sum256([]byte(strings.Join([]string{status,generationRequestCandidateGateRevisionProposalBridgeDigest,feedbackRevisionProposalBridgeDigest,candidateMaterializationBridgeDigest,proposalDecision,candidateDecision,materializationStatus},"|"))); return hex.EncodeToString(sum[:]) }
