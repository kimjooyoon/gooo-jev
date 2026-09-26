package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded = "revision-plan-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady = "revision-plan-ready"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold = "revision-plan-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected = "revision-plan-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanInput struct {
    Proposal        JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal
    PlanSource      string
    NonAuthorizing  bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan struct {
    Status           string
    MissingStage     string
    ProposalStatus   string
    CandidateDigest  string
    ProposalDigest   string
    RevisionSource   string
    PlanSource       string
    IntentDigest     string
    PlanDigest       string
    NonExecuting     bool
    NonAuthorizing   bool
}

func (p JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan) Validate() error {
    if p.Status == "" || p.ProposalStatus == "" || p.ProposalDigest == "" || p.PlanDigest == "" {
        return fmt.Errorf("incomplete JEV revision plan")
    }
    if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected {
        return fmt.Errorf("invalid JEV revision plan status")
    }
    switch p.ProposalStatus {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded:
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded ||
            p.CandidateDigest != "" || p.RevisionSource != "" ||
            p.PlanSource != "" || p.IntentDigest != "" {
            return fmt.Errorf("not-needed proposal requires an empty revision plan")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady:
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady ||
            p.CandidateDigest == "" || p.RevisionSource == "" ||
            p.PlanSource == "" || p.IntentDigest == "" {
            return fmt.Errorf("ready proposal requires complete revision plan evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold:
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold {
            return fmt.Errorf("held proposal requires revision-plan-hold status")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected:
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected {
            return fmt.Errorf("rejected proposal requires revision-plan-rejected status")
        }
    default:
        return fmt.Errorf("invalid JEV revision proposal status")
    }
    if !p.NonExecuting {
        return fmt.Errorf("JEV revision plan must be non-executing")
    }
    if !p.NonAuthorizing {
        return fmt.Errorf("JEV revision plan must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(
        p.Status,
        p.ProposalStatus,
        p.ProposalDigest,
        p.PlanSource,
        p.IntentDigest,
    )
    if p.PlanDigest != expected {
        return fmt.Errorf("JEV revision plan digest mismatch")
    }
    return nil
}

func PlanJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanInput) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown
        ProposalStatus: input.Proposal.Status,
        ProposalDigest: input.Proposal.ProposalDigest,
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
        output.MissingStage = "revision-proposal"
        return output
    }
    output.ProposalStatus = input.Proposal.Status
    output.CandidateDigest = input.Proposal.CandidateDigest
    output.RevisionSource = input.Proposal.RevisionSource
    switch input.Proposal.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady:
        if input.PlanSource == "" {
            output.MissingStage = "revision-plan-source"
            return output
        }
        if input.PlanSource != input.Proposal.RevisionPlanSource {
            output.MissingStage = "revision-plan-source-consistency"
            return output
        }
        if input.Proposal.RevisionIntentDigest == "" {
            output.MissingStage = "revision-intent"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady
        output.PlanSource = input.PlanSource
        output.IntentDigest = input.Proposal.RevisionIntentDigest
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected
    default:
        output.MissingStage = "revision-proposal-status"
        return output
    }
    output.PlanDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(
        output.Status,
        output.ProposalStatus,
        output.ProposalDigest,
        output.PlanSource,
        output.IntentDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown
        output.MissingStage = "revision-plan-evidence"
        output.PlanDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(status, proposalStatus, proposalDigest, planSource, intentDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, proposalStatus, proposalDigest, planSource, intentDigest)))
    return hex.EncodeToString(sum[:])
}