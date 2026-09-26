package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded = "revision-proposal-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady = "revision-proposal-ready"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold = "revision-proposal-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected = "revision-proposal-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalInput struct {
    Bridge              JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge
    RevisionPlanSource  string
    RevisionIntentDigest string
    NonAuthorizing      bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal struct {
    Status                    string
    MissingStage              string
    BridgeStatus              string
    CandidateGenerationDigest string
    GateDigest                string
    CandidateDigest           string
    RevisionSource            string
    RevisionPlanSource        string
    RevisionIntentDigest      string
    ProposalDigest            string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (p JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal) Validate() error {
    if p.Status == "" || p.BridgeStatus == "" || p.CandidateGenerationDigest == "" || p.ProposalDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review replay origin mutation revision proposal")
    }
    if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected {
        return fmt.Errorf("invalid JEV revision proposal status")
    }
    if p.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded &&
        p.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady &&
        p.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold &&
        p.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected {
        return fmt.Errorf("invalid candidate revision gate bridge status")
    }
    if p.BridgeStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded {
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded ||
            p.RevisionPlanSource != "" ||
            p.RevisionIntentDigest != "" {
            return fmt.Errorf("not-needed bridge requires revision-proposal-not-needed status")
        }
    }
    if p.BridgeStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady {
        if p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady ||
            p.CandidateDigest == "" ||
            p.RevisionSource == "" ||
            p.RevisionPlanSource == "" ||
            p.RevisionIntentDigest == "" {
            return fmt.Errorf("ready bridge requires complete revision proposal evidence")
        }
    }
    if p.BridgeStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold {
        return fmt.Errorf("held bridge requires revision-proposal-hold status")
    }
    if p.BridgeStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected &&
        p.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected {
        return fmt.Errorf("rejected bridge requires revision-proposal-rejected status")
    }
    if !p.NonExecuting {
        return fmt.Errorf("JEV revision proposal must be non-executing")
    }
    if !p.NonAuthorizing {
        return fmt.Errorf("JEV revision proposal must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal(
        p.Status,
        p.BridgeStatus,
        p.CandidateGenerationDigest,
        p.GateDigest,
        p.CandidateDigest,
        p.RevisionSource,
        p.RevisionPlanSource,
        p.RevisionIntentDigest,
    )
    if p.ProposalDigest != expected {
        return fmt.Errorf("JEV revision proposal digest mismatch")
    }
    return nil
}

func ProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalInput) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Bridge.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Bridge.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Bridge.Validate(); err != nil {
        output.MissingStage = "candidate-revision-gate-bridge"
        return output
    }
    output.BridgeStatus = input.Bridge.Status
    output.CandidateGenerationDigest = input.Bridge.CandidateGenerationDigest
    output.GateDigest = input.Bridge.Gate.GateDigest
    output.CandidateDigest = input.Bridge.Gate.CandidateDigest
    output.RevisionSource = input.Bridge.Gate.CandidateSource
    switch input.Bridge.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady:
        if input.RevisionPlanSource == "" {
            output.MissingStage = "revision-plan-source"
            return output
        }
        if input.RevisionIntentDigest == "" {
            output.MissingStage = "revision-intent"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady
        output.RevisionPlanSource = input.RevisionPlanSource
        output.RevisionIntentDigest = input.RevisionIntentDigest
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected
    default:
        output.MissingStage = "candidate-revision-gate-status"
        return output
    }
    output.ProposalDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal(
        output.Status,
        output.BridgeStatus,
        output.CandidateGenerationDigest,
        output.GateDigest,
        output.CandidateDigest,
        output.RevisionSource,
        output.RevisionPlanSource,
        output.RevisionIntentDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown
        output.MissingStage = "revision-proposal-evidence"
        output.ProposalDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal(status, bridgeStatus, candidateGenerationDigest, gateDigest, candidateDigest, revisionSource, revisionPlanSource, revisionIntentDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", status, bridgeStatus, candidateGenerationDigest, gateDigest, candidateDigest, revisionSource, revisionPlanSource, revisionIntentDigest)))
    return hex.EncodeToString(sum[:])
}
