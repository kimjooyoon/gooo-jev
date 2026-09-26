package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded = "origin-mutation-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationProposed = "origin-mutation-proposed"
const jevExternalApplyCapabilityReviewReplayOriginMutationUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput struct {
    Consistency          JEVExternalApplyCapabilityReviewReplayOriginConsistency
    ProposedOriginDigest string
    Reason               string
    NonAuthorizing       bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationProposal struct {
    Status                    string
    MissingStage              string
    ConsistencyStatus         string
    DeclaredOriginDigest      string
    ReverseOriginDigest       string
    ProposedOriginDigest      string
    Reason                    string
    ConsistencyDigest         string
    MutationProposalDigest    string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (p JEVExternalApplyCapabilityReviewReplayOriginMutationProposal) Validate() error {
    if p.Status == "" || p.ConsistencyStatus == "" || p.DeclaredOriginDigest == "" || p.ReverseOriginDigest == "" || p.ConsistencyDigest == "" || p.MutationProposalDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin mutation proposal")
    }
    if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation proposal status")
    }
    if p.ConsistencyStatus != jevExternalApplyCapabilityReviewReplayOriginConsistent &&
        p.ConsistencyStatus != jevExternalApplyCapabilityReviewReplayOriginMismatch {
        return fmt.Errorf("invalid JEV external apply capability review replay origin mutation consistency status")
    }
    if p.ConsistencyStatus == jevExternalApplyCapabilityReviewReplayOriginConsistent {
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded {
            return fmt.Errorf("consistent origin requires no mutation proposal")
        }
        if p.ProposedOriginDigest != "" || p.Reason != "" {
            return fmt.Errorf("mutation-not-needed result must not contain a proposal")
        }
    }
    if p.ConsistencyStatus == jevExternalApplyCapabilityReviewReplayOriginMismatch {
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
            return fmt.Errorf("mismatched origin requires an explicit mutation proposal")
        }
        if p.ProposedOriginDigest == "" || p.Reason == "" {
            return fmt.Errorf("origin mutation proposal requires proposed digest and reason")
        }
    }
    if !p.NonExecuting {
        return fmt.Errorf("JEV external apply capability review replay origin mutation proposal must be non-executing")
    }
    if !p.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review replay origin mutation proposal must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        p.Status,
        p.ConsistencyStatus,
        p.DeclaredOriginDigest,
        p.ReverseOriginDigest,
        p.ProposedOriginDigest,
        p.Reason,
        p.ConsistencyDigest,
    )
    if p.MutationProposalDigest != expected {
        return fmt.Errorf("JEV external apply capability review replay origin mutation proposal digest mismatch")
    }
    return nil
}

func ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(input JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput) JEVExternalApplyCapabilityReviewReplayOriginMutationProposal {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationProposal{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Consistency.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Consistency.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Consistency.Validate(); err != nil {
        output.MissingStage = "replay-origin-consistency"
        return output
    }
    output.ConsistencyStatus = input.Consistency.Status
    output.DeclaredOriginDigest = input.Consistency.DeclaredOriginDigest
    output.ReverseOriginDigest = input.Consistency.ReverseOriginDigest
    output.ConsistencyDigest = input.Consistency.ConsistencyDigest
    switch input.Consistency.Status {
    case jevExternalApplyCapabilityReviewReplayOriginConsistent:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMismatch:
        if input.ProposedOriginDigest == "" {
            output.MissingStage = "mutation-proposed-origin"
            return output
        }
        if input.Reason == "" {
            output.MissingStage = "mutation-reason"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationProposed
        output.ProposedOriginDigest = input.ProposedOriginDigest
        output.Reason = input.Reason
    default:
        output.MissingStage = "replay-origin-consistency-status"
        return output
    }
    output.MutationProposalDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        output.Status,
        output.ConsistencyStatus,
        output.DeclaredOriginDigest,
        output.ReverseOriginDigest,
        output.ProposedOriginDigest,
        output.Reason,
        output.ConsistencyDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationUnknown
        output.MissingStage = "origin-mutation-proposal-evidence"
        output.MutationProposalDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(status, consistencyStatus, declaredOriginDigest, reverseOriginDigest, proposedOriginDigest, reason, consistencyDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", status, consistencyStatus, declaredOriginDigest, reverseOriginDigest, proposedOriginDigest, reason, consistencyDigest)))
    return hex.EncodeToString(sum[:])
}
