package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound = "feedback-revision-proposal-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound = "revision-proposal-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalHeld = "revision-proposal-held"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRejected = "revision-proposal-rejected"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise = "revise"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRetain = "retain"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRollback = "rollback"

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeInput struct {
    FeedbackBridge          JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge
    ProposalDecision        string
    ProposalTarget          string
    ProposalDigest          string
    ProposalSource          string
    ProposalEvidenceDigest  string
    NonAuthorizing          bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge struct {
    Status                         string
    MissingStage                   string
    CandidateDigest                string
    CandidateSource                string
    CandidateGateDigest            string
    FeedbackStatus                string
    FeedbackDirection              string
    FeedbackDigest                 string
    FeedbackSource                 string
    FeedbackEvidenceDigest         string
    FeedbackBridgeDigest           string
    RevisionProposalStatus         string
    ProposalDecision               string
    ProposalTarget                 string
    ProposalDigest                 string
    ProposalSource                 string
    ProposalEvidenceDigest         string
    BridgeDigest                   string
    NonExecuting                   bool
    NonAuthorizing                 bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge) Validate() error {
    if b.Status == "" || b.FeedbackStatus == "" || b.FeedbackBridgeDigest == "" ||
        b.RevisionProposalStatus == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV feedback revision proposal bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound {
        return fmt.Errorf("invalid JEV feedback revision proposal bridge status")
    }
    switch b.FeedbackStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound:
        if b.CandidateDigest == "" || b.CandidateSource == "" || b.CandidateGateDigest == "" ||
            !validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackDirection(b.FeedbackDirection) ||
            b.FeedbackDigest == "" || b.FeedbackSource == "" || b.FeedbackEvidenceDigest == "" ||
            b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound ||
            !validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalDecision(b.ProposalDecision, b.FeedbackDirection) ||
            b.ProposalTarget == "" || b.ProposalDigest == "" || b.ProposalSource == "" ||
            b.ProposalEvidenceDigest == "" {
            return fmt.Errorf("bound feedback revision proposal bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld:
        if b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalHeld ||
            b.ProposalDecision != "" || b.ProposalTarget != "" || b.ProposalDigest != "" ||
            b.ProposalSource != "" || b.ProposalEvidenceDigest != "" {
            return fmt.Errorf("held feedback revision proposal bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected:
        if b.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRejected ||
            b.ProposalDecision != "" || b.ProposalTarget != "" || b.ProposalDigest != "" ||
            b.ProposalSource != "" || b.ProposalEvidenceDigest != "" {
            return fmt.Errorf("rejected feedback revision proposal bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid feedback status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV feedback revision proposal bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV feedback revision proposal bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(
        b.Status,
        b.CandidateDigest,
        b.CandidateSource,
        b.CandidateGateDigest,
        b.FeedbackStatus,
        b.FeedbackDirection,
        b.FeedbackDigest,
        b.FeedbackSource,
        b.FeedbackEvidenceDigest,
        b.FeedbackBridgeDigest,
        b.RevisionProposalStatus,
        b.ProposalDecision,
        b.ProposalTarget,
        b.ProposalDigest,
        b.ProposalSource,
        b.ProposalEvidenceDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV feedback revision proposal bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.FeedbackBridge.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.FeedbackBridge.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.FeedbackBridge.Validate(); err != nil {
        output.MissingStage = "candidate-feedback"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeBound
    output.CandidateDigest = input.FeedbackBridge.CandidateDigest
    output.CandidateSource = input.FeedbackBridge.CandidateSource
    output.CandidateGateDigest = input.FeedbackBridge.CandidateGateDigest
    output.FeedbackStatus = input.FeedbackBridge.FeedbackStatus
    output.FeedbackDirection = input.FeedbackBridge.FeedbackDirection
    output.FeedbackDigest = input.FeedbackBridge.FeedbackDigest
    output.FeedbackSource = input.FeedbackBridge.FeedbackSource
    output.FeedbackEvidenceDigest = input.FeedbackBridge.FeedbackEvidenceDigest
    output.FeedbackBridgeDigest = input.FeedbackBridge.BridgeDigest

    switch input.FeedbackBridge.FeedbackStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound:
        if input.ProposalDecision == "" {
            output.MissingStage = "proposal-decision"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        if !validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalDecision(input.ProposalDecision, input.FeedbackBridge.FeedbackDirection) {
            output.MissingStage = "proposal-decision-consistency"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        if input.ProposalTarget == "" {
            output.MissingStage = "proposal-target"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        if input.ProposalDigest == "" {
            output.MissingStage = "proposal-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        if input.ProposalSource == "" {
            output.MissingStage = "proposal-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        if input.ProposalEvidenceDigest == "" {
            output.MissingStage = "proposal-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        output.RevisionProposalStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound
        output.ProposalDecision = input.ProposalDecision
        output.ProposalTarget = input.ProposalTarget
        output.ProposalDigest = input.ProposalDigest
        output.ProposalSource = input.ProposalSource
        output.ProposalEvidenceDigest = input.ProposalEvidenceDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld:
        if input.ProposalDecision != "" || input.ProposalTarget != "" || input.ProposalDigest != "" ||
            input.ProposalSource != "" || input.ProposalEvidenceDigest != "" {
            output.MissingStage = "held-proposal-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        output.RevisionProposalStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected:
        if input.ProposalDecision != "" || input.ProposalTarget != "" || input.ProposalDigest != "" ||
            input.ProposalSource != "" || input.ProposalEvidenceDigest != "" {
            output.MissingStage = "rejected-proposal-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
            return output
        }
        output.RevisionProposalStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRejected
    default:
        output.MissingStage = "feedback-status"
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(
        output.Status,
        output.CandidateDigest,
        output.CandidateSource,
        output.CandidateGateDigest,
        output.FeedbackStatus,
        output.FeedbackDirection,
        output.FeedbackDigest,
        output.FeedbackSource,
        output.FeedbackEvidenceDigest,
        output.FeedbackBridgeDigest,
        output.RevisionProposalStatus,
        output.ProposalDecision,
        output.ProposalTarget,
        output.ProposalDigest,
        output.ProposalSource,
        output.ProposalEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown
        output.MissingStage = "feedback-revision-proposal-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalDecision(decision, direction string) bool {
    switch direction {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove:
        return decision == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRetain:
        return decision == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRetain
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRollback:
        return decision == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRollback
    default:
        return false
    }
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(status, candidateDigest, candidateSource, candidateGateDigest, feedbackStatus, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackBridgeDigest, revisionProposalStatus, proposalDecision, proposalTarget, proposalDigest, proposalSource, proposalEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{
        status, candidateDigest, candidateSource, candidateGateDigest, feedbackStatus,
        feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest,
        feedbackBridgeDigest, revisionProposalStatus, proposalDecision, proposalTarget,
        proposalDigest, proposalSource, proposalEvidenceDigest,
    }, "|")))
    return hex.EncodeToString(sum[:])
}
