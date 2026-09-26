package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded = "origin-mutation-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved = "origin-mutation-review-approved"
const jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold = "origin-mutation-review-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected = "origin-mutation-review-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationReviewUnknown = "UNKNOWN"

const jevExternalApplyCapabilityReviewReplayOriginMutationReviewApprove = "approve"
const jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold = "hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationReviewReject = "reject"

type JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput struct {
    Proposal             JEVExternalApplyCapabilityReviewReplayOriginMutationProposal
    ReviewDecision       string
    ReviewEvidenceDigest string
    NonAuthorizing       bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate struct {
    Status               string
    MissingStage         string
    ProposalStatus       string
    ReviewDecision       string
    ProposalDigest       string
    ReviewEvidenceDigest string
    DecisionDigest       string
    NonExecuting         bool
    NonAuthorizing       bool
}

func (g JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate) Validate() error {
    if g.Status == "" || g.ProposalStatus == "" || g.ProposalDigest == "" || g.DecisionDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin mutation review gate")
    }
    if g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation review gate status")
    }
    if g.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded &&
        g.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation proposal status")
    }
    if g.ProposalStatus == jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded {
        if g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded ||
            g.ReviewDecision != "" ||
            g.ReviewEvidenceDigest != "" {
            return fmt.Errorf("mutation-not-needed proposal requires no review decision")
        }
    }
    if g.ProposalStatus == jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        if g.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApprove &&
            g.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
            g.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewReject {
            return fmt.Errorf("invalid JEV external apply capability review decision")
        }
        if g.ReviewEvidenceDigest == "" {
            return fmt.Errorf("review decision requires review evidence digest")
        }
        if g.ReviewDecision == jevExternalApplyCapabilityReviewReplayOriginMutationReviewApprove &&
            g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved {
            return fmt.Errorf("approve decision requires approved review status")
        }
        if g.ReviewDecision == jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
            g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold {
            return fmt.Errorf("hold decision requires hold review status")
        }
        if g.ReviewDecision == jevExternalApplyCapabilityReviewReplayOriginMutationReviewReject &&
            g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
            return fmt.Errorf("reject decision requires rejected review status")
        }
    }
    if !g.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay origin mutation review gate must be non-executing")
    }
    if !g.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay origin mutation review gate must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate(
        g.Status,
        g.ProposalStatus,
        g.ReviewDecision,
        g.ProposalDigest,
        g.ReviewEvidenceDigest,
    )
    if g.DecisionDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay origin mutation review gate digest mismatch")
    }
    return nil
}

func GateJEVExternalApplyCapabilityReviewReplayOriginMutation(input JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput) JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationReviewUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Proposal.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Proposal.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Proposal.Validate(); err != nil {
        output.MissingStage = "origin-mutation-proposal"
        return output
    }
    output.ProposalStatus = input.Proposal.Status
    output.ProposalDigest = input.Proposal.MutationProposalDigest
    switch input.Proposal.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationProposed:
        if input.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewApprove &&
            input.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold &&
            input.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationReviewReject {
            output.MissingStage = "review-decision"
            return output
        }
        if input.ReviewEvidenceDigest == "" {
            output.MissingStage = "review-evidence"
            return output
        }
        output.ReviewDecision = input.ReviewDecision
        output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
        switch input.ReviewDecision {
        case jevExternalApplyCapabilityReviewReplayOriginMutationReviewApprove:
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved
        case jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold:
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold
        case jevExternalApplyCapabilityReviewReplayOriginMutationReviewReject:
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected
        }
    default:
        output.MissingStage = "mutation-review-proposal-status"
        return output
    }
    output.DecisionDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate(
        output.Status,
        output.ProposalStatus,
        output.ReviewDecision,
        output.ProposalDigest,
        output.ReviewEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationReviewUnknown
        output.MissingStage = "mutation-review-gate-evidence"
        output.DecisionDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGate(status, proposalStatus, reviewDecision, proposalDigest, reviewEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, proposalStatus, reviewDecision, proposalDigest, reviewEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
