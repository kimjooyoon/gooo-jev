package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded = "application-candidate-not-needed"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady = "application-candidate-ready"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold = "application-candidate-hold"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected = "application-candidate-rejected"
const jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateInput struct {
    ReviewGate               JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate
    ApplicationTarget        string
    ApplicationSource        string
    ApplicationEvidenceDigest string
    NonAuthorizing           bool
}

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate struct {
    Status                    string
    MissingStage              string
    ReviewGateStatus          string
    Decision                  string
    PlanDigest                string
    GateDigest                string
    ApplicationTarget         string
    ApplicationSource         string
    ApplicationEvidenceDigest string
    CandidateDigest           string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (c JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate) Validate() error {
    if c.Status == "" || c.ReviewGateStatus == "" || c.Decision == "" ||
        c.PlanDigest == "" || c.GateDigest == "" || c.CandidateDigest == "" {
        return fmt.Errorf("incomplete JEV revision application candidate")
    }
    if c.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded &&
        c.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady &&
        c.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold &&
        c.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected {
        return fmt.Errorf("invalid JEV revision application candidate status")
    }
    switch c.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded:
        if c.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded ||
            c.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewNotNeeded ||
            c.ApplicationTarget != "" || c.ApplicationSource != "" ||
            c.ApplicationEvidenceDigest != "" {
            return fmt.Errorf("not-needed application candidate has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady:
        if c.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady ||
            c.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove ||
            c.ApplicationTarget == "" || c.ApplicationSource == "" ||
            c.ApplicationEvidenceDigest == "" {
            return fmt.Errorf("ready application candidate is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold:
        if c.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold ||
            c.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold {
            return fmt.Errorf("held application candidate has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected:
        if c.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected ||
            c.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject {
            return fmt.Errorf("rejected application candidate has inconsistent evidence")
        }
    }
    if !c.NonExecuting {
        return fmt.Errorf("JEV revision application candidate must be non-executing")
    }
    if !c.NonAuthorizing {
        return fmt.Errorf("JEV revision application candidate must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(
        c.Status,
        c.ReviewGateStatus,
        c.Decision,
        c.PlanDigest,
        c.GateDigest,
        c.ApplicationTarget,
        c.ApplicationSource,
        c.ApplicationEvidenceDigest,
    )
    if c.CandidateDigest != expected {
        return fmt.Errorf("JEV revision application candidate digest mismatch")
    }
    return nil
}

func GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateInput) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate{
        Status:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown,
        ReviewGateStatus: input.ReviewGate.Status,
        Decision:         input.ReviewGate.Decision,
        PlanDigest:       input.ReviewGate.PlanDigest,
        GateDigest:       input.ReviewGate.GateDigest,
        NonExecuting:     true,
        NonAuthorizing:   true,
    }
    if !input.NonAuthorizing || !input.ReviewGate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.ReviewGate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.ReviewGate.Validate(); err != nil {
        output.MissingStage = "revision-plan-review-gate"
        return output
    }
    switch input.ReviewGate.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady:
        if input.ApplicationTarget == "" {
            output.MissingStage = "application-target"
            return output
        }
        if input.ApplicationSource == "" {
            output.MissingStage = "application-source"
            return output
        }
        if input.ApplicationEvidenceDigest == "" {
            output.MissingStage = "application-evidence"
            return output
        }
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady
        output.ApplicationTarget = input.ApplicationTarget
        output.ApplicationSource = input.ApplicationSource
        output.ApplicationEvidenceDigest = input.ApplicationEvidenceDigest
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected
    default:
        output.MissingStage = "revision-plan-review-status"
        return output
    }
    output.CandidateDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(
        output.Status,
        output.ReviewGateStatus,
        output.Decision,
        output.PlanDigest,
        output.GateDigest,
        output.ApplicationTarget,
        output.ApplicationSource,
        output.ApplicationEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown
        output.MissingStage = "application-candidate-evidence"
        output.CandidateDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(status, reviewGateStatus, decision, planDigest, gateDigest, applicationTarget, applicationSource, applicationEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", status, reviewGateStatus, decision, planDigest, gateDigest, applicationTarget, applicationSource, applicationEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
