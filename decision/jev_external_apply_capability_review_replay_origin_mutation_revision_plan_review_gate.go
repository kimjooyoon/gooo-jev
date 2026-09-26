package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded = "revision-plan-review-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady = "revision-plan-review-ready"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold = "revision-plan-review-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected = "revision-plan-review-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown = "UNKNOWN"

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove = "approve"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold = "hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject = "reject"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewNotNeeded = "not-needed"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateInput struct {
    Plan              JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan
    ReviewDecision    string
    ReviewSource      string
    ReviewEvidenceDigest string
    NonAuthorizing    bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate struct {
    Status             string
    MissingStage       string
    PlanStatus         string
    Decision           string
    PlanDigest         string
    ReviewSource       string
    ReviewEvidenceDigest string
    GateDigest         string
    NonExecuting       bool
    NonAuthorizing     bool
}

func (g JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate) Validate() error {
    if g.Status == "" || g.PlanStatus == "" || g.Decision == "" || g.PlanDigest == "" || g.GateDigest == "" {
        return fmt.Errorf("incomplete JEV revision plan review gate")
    }
    if g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold &&
        g.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected {
        return fmt.Errorf("invalid JEV revision plan review gate status")
    }
    switch g.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded:
        if g.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded ||
            g.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewNotNeeded ||
            g.ReviewSource != "" || g.ReviewEvidenceDigest != "" {
            return fmt.Errorf("not-needed revision plan review gate has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady:
        if g.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady ||
            g.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove ||
            g.ReviewSource == "" || g.ReviewEvidenceDigest == "" {
            return fmt.Errorf("ready revision plan review gate is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold:
        if (g.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady &&
            g.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold) ||
            g.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold ||
            g.ReviewSource == "" || g.ReviewEvidenceDigest == "" {
            return fmt.Errorf("held revision plan review gate is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected:
        if (g.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady &&
            g.PlanStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected) ||
            g.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject ||
            g.ReviewSource == "" || g.ReviewEvidenceDigest == "" {
            return fmt.Errorf("rejected revision plan review gate is incomplete")
        }
    }
    if !g.NonExecuting {
        return fmt.Errorf("JEV revision plan review gate must be non-executing")
    }
    if !g.NonAuthorizing {
        return fmt.Errorf("JEV revision plan review gate must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate(
        g.Status,
        g.PlanStatus,
        g.Decision,
        g.PlanDigest,
        g.ReviewSource,
        g.ReviewEvidenceDigest,
    )
    if g.GateDigest != expected {
        return fmt.Errorf("JEV revision plan review gate digest mismatch")
    }
    return nil
}

func ReviewJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateInput) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown,
        PlanStatus:     input.Plan.Status,
        PlanDigest:     input.Plan.PlanDigest,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Plan.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Plan.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Plan.Validate(); err != nil {
        output.MissingStage = "revision-plan"
        return output
    }
    switch input.Plan.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded
        output.Decision = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady:
        if input.ReviewDecision == "" {
            output.MissingStage = "revision-review-decision"
            return output
        }
        if input.ReviewSource == "" {
            output.MissingStage = "revision-review-source"
            return output
        }
        if input.ReviewEvidenceDigest == "" {
            output.MissingStage = "revision-review-evidence"
            return output
        }
        output.ReviewSource = input.ReviewSource
        output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
        switch input.ReviewDecision {
        case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove:
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady
        case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold:
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold
        case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject:
            output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected
        default:
            output.MissingStage = "revision-review-decision"
            return output
        }
        output.Decision = input.ReviewDecision
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold:
        if input.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold {
            output.MissingStage = "revision-review-decision"
            return output
        }
        if input.ReviewSource == "" {
            output.MissingStage = "revision-review-source"
            return output
        }
        if input.ReviewEvidenceDigest == "" {
            output.MissingStage = "revision-review-evidence"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold
        output.Decision = input.ReviewDecision
        output.ReviewSource = input.ReviewSource
        output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected:
        if input.ReviewDecision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject {
            output.MissingStage = "revision-review-decision"
            return output
        }
        if input.ReviewSource == "" {
            output.MissingStage = "revision-review-source"
            return output
        }
        if input.ReviewEvidenceDigest == "" {
            output.MissingStage = "revision-review-evidence"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected
        output.Decision = input.ReviewDecision
        output.ReviewSource = input.ReviewSource
        output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
    default:
        output.MissingStage = "revision-plan-status"
        return output
    }
    output.GateDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate(
        output.Status,
        output.PlanStatus,
        output.Decision,
        output.PlanDigest,
        output.ReviewSource,
        output.ReviewEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown
        output.MissingStage = "revision-plan-review-evidence"
        output.GateDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate(status, planStatus, decision, planDigest, reviewSource, reviewEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s", status, planStatus, decision, planDigest, reviewSource, reviewEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
