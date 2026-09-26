package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeBound = "revision-proposal-candidate-materialization-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized = "candidate-materialized"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateHeld = "candidate-held"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateRollback = "candidate-rollback"

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeInput struct {
    RevisionProposal          JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge
    GeneratedCandidateDigest  string
    GeneratedCandidateSource  string
    GenerationInputDigest     string
    GenerationEvidenceDigest  string
    GeneratorIdentity         string
    NonAuthorizing             bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge struct {
    Status                      string
    MissingStage                string
    CandidateDigest             string
    CandidateSource             string
    CandidateGateDigest         string
    FeedbackDirection            string
    FeedbackDigest               string
    FeedbackSource               string
    FeedbackEvidenceDigest       string
    FeedbackBridgeDigest         string
    ProposalDecision             string
    ProposalTarget               string
    ProposalDigest              string
    ProposalSource              string
    ProposalEvidenceDigest      string
    RevisionProposalBridgeDigest string
    MaterializationStatus        string
    GeneratedCandidateDigest     string
    GeneratedCandidateSource     string
    GenerationInputDigest        string
    GenerationEvidenceDigest     string
    GeneratorIdentity            string
    BridgeDigest                 string
    NonExecuting                 bool
    NonAuthorizing               bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge) Validate() error {
    if b.Status == "" || b.FeedbackBridgeDigest == "" ||
        b.RevisionProposalBridgeDigest == "" || b.MaterializationStatus == "" ||
        b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV revision proposal candidate materialization bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeBound {
        return fmt.Errorf("invalid JEV revision proposal candidate materialization bridge status")
    }
    switch b.ProposalDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise:
        if b.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized ||
            b.CandidateDigest == "" || b.CandidateSource == "" || b.CandidateGateDigest == "" ||
            !validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackDirection(b.FeedbackDirection) ||
            b.FeedbackDigest == "" || b.FeedbackSource == "" || b.FeedbackEvidenceDigest == "" ||
            b.ProposalTarget == "" || b.ProposalDigest == "" || b.ProposalSource == "" ||
            b.ProposalEvidenceDigest == "" || b.GeneratedCandidateDigest == "" ||
            b.GeneratedCandidateSource == "" || b.GenerationInputDigest == "" ||
            b.GenerationEvidenceDigest == "" || b.GeneratorIdentity == "" {
            return fmt.Errorf("revise candidate materialization bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRetain:
        if b.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateHeld ||
            b.GeneratedCandidateDigest != "" || b.GeneratedCandidateSource != "" ||
            b.GenerationInputDigest != "" || b.GenerationEvidenceDigest != "" ||
            b.GeneratorIdentity != "" {
            return fmt.Errorf("retain candidate materialization bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRollback:
        if b.MaterializationStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateRollback ||
            b.GeneratedCandidateDigest != "" || b.GeneratedCandidateSource != "" ||
            b.GenerationInputDigest != "" || b.GenerationEvidenceDigest != "" ||
            b.GeneratorIdentity != "" {
            return fmt.Errorf("rollback candidate materialization bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid revision proposal decision")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV revision proposal candidate materialization bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV revision proposal candidate materialization bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(
        b.Status,
        b.CandidateDigest,
        b.CandidateSource,
        b.CandidateGateDigest,
        b.FeedbackDirection,
        b.FeedbackDigest,
        b.FeedbackSource,
        b.FeedbackEvidenceDigest,
        b.FeedbackBridgeDigest,
        b.ProposalDecision,
        b.ProposalTarget,
        b.ProposalDigest,
        b.ProposalSource,
        b.ProposalEvidenceDigest,
        b.RevisionProposalBridgeDigest,
        b.MaterializationStatus,
        b.GeneratedCandidateDigest,
        b.GeneratedCandidateSource,
        b.GenerationInputDigest,
        b.GenerationEvidenceDigest,
        b.GeneratorIdentity,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV revision proposal candidate materialization bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.RevisionProposal.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.RevisionProposal.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.RevisionProposal.Validate(); err != nil {
        output.MissingStage = "revision-proposal"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeBound
    output.CandidateDigest = input.RevisionProposal.CandidateDigest
    output.CandidateSource = input.RevisionProposal.CandidateSource
    output.CandidateGateDigest = input.RevisionProposal.CandidateGateDigest
    output.FeedbackDirection = input.RevisionProposal.FeedbackDirection
    output.FeedbackDigest = input.RevisionProposal.FeedbackDigest
    output.FeedbackSource = input.RevisionProposal.FeedbackSource
    output.FeedbackEvidenceDigest = input.RevisionProposal.FeedbackEvidenceDigest
    output.FeedbackBridgeDigest = input.RevisionProposal.FeedbackBridgeDigest
    output.ProposalDecision = input.RevisionProposal.ProposalDecision
    output.ProposalTarget = input.RevisionProposal.ProposalTarget
    output.ProposalDigest = input.RevisionProposal.ProposalDigest
    output.ProposalSource = input.RevisionProposal.ProposalSource
    output.ProposalEvidenceDigest = input.RevisionProposal.ProposalEvidenceDigest
    output.RevisionProposalBridgeDigest = input.RevisionProposal.BridgeDigest

    switch input.RevisionProposal.ProposalDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise:
        if input.GeneratedCandidateDigest == "" {
            output.MissingStage = "generated-candidate-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        if input.GeneratedCandidateSource == "" {
            output.MissingStage = "generated-candidate-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        if input.GenerationInputDigest == "" {
            output.MissingStage = "generation-input"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        if input.GeneratorIdentity == "" {
            output.MissingStage = "generator-identity"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        if input.GenerationEvidenceDigest == "" {
            output.MissingStage = "generation-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        output.MaterializationStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized
        output.GeneratedCandidateDigest = input.GeneratedCandidateDigest
        output.GeneratedCandidateSource = input.GeneratedCandidateSource
        output.GenerationInputDigest = input.GenerationInputDigest
        output.GenerationEvidenceDigest = input.GenerationEvidenceDigest
        output.GeneratorIdentity = input.GeneratorIdentity
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRetain:
        if input.GeneratedCandidateDigest != "" || input.GeneratedCandidateSource != "" ||
            input.GenerationInputDigest != "" || input.GenerationEvidenceDigest != "" ||
            input.GeneratorIdentity != "" {
            output.MissingStage = "retain-generation-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        output.MaterializationStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRollback:
        if input.GeneratedCandidateDigest != "" || input.GeneratedCandidateSource != "" ||
            input.GenerationInputDigest != "" || input.GenerationEvidenceDigest != "" ||
            input.GeneratorIdentity != "" {
            output.MissingStage = "rollback-generation-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
            return output
        }
        output.MaterializationStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateRollback
    default:
        output.MissingStage = "proposal-decision"
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(
        output.Status,
        output.CandidateDigest,
        output.CandidateSource,
        output.CandidateGateDigest,
        output.FeedbackDirection,
        output.FeedbackDigest,
        output.FeedbackSource,
        output.FeedbackEvidenceDigest,
        output.FeedbackBridgeDigest,
        output.ProposalDecision,
        output.ProposalTarget,
        output.ProposalDigest,
        output.ProposalSource,
        output.ProposalEvidenceDigest,
        output.RevisionProposalBridgeDigest,
        output.MaterializationStatus,
        output.GeneratedCandidateDigest,
        output.GeneratedCandidateSource,
        output.GenerationInputDigest,
        output.GenerationEvidenceDigest,
        output.GeneratorIdentity,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridgeUnknown
        output.MissingStage = "candidate-materialization-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge(status, candidateDigest, candidateSource, candidateGateDigest, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackBridgeDigest, proposalDecision, proposalTarget, proposalDigest, proposalSource, proposalEvidenceDigest, revisionProposalBridgeDigest, materializationStatus, generatedCandidateDigest, generatedCandidateSource, generationInputDigest, generationEvidenceDigest, generatorIdentity string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{
        status, candidateDigest, candidateSource, candidateGateDigest, feedbackDirection,
        feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackBridgeDigest,
        proposalDecision, proposalTarget, proposalDigest, proposalSource, proposalEvidenceDigest,
        revisionProposalBridgeDigest, materializationStatus, generatedCandidateDigest,
        generatedCandidateSource, generationInputDigest, generationEvidenceDigest, generatorIdentity,
    }, "|")))
    return hex.EncodeToString(sum[:])
}
