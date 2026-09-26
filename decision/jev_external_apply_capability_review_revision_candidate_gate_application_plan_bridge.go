package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeBound = "candidate-gate-application-plan-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady = "application-plan-ready"
const jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld = "application-plan-held"
const jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected = "application-plan-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput struct {
    CandidateGate     JEVExternalApplyCapabilityReviewRevisionCandidateGate
    PlanDigest        string
    PlanSource        string
    PlanEvidenceDigest string
    NonAuthorizing    bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge struct {
    Status             string
    MissingStage       string
    CandidateGateStatus string
    CandidateDecision   string
    CandidateDigest     string
    CandidateSource     string
    CandidateGateDigest string
    ApplicationPlanStatus string
    PlanDigest          string
    PlanSource          string
    PlanEvidenceDigest  string
    BridgeDigest        string
    NonExecuting       bool
    NonAuthorizing     bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge) Validate() error {
    if b.Status == "" || b.CandidateGateStatus == "" ||
        b.CandidateDecision == "" || b.CandidateGateDigest == "" ||
        b.ApplicationPlanStatus == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV candidate gate application plan bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeBound {
        return fmt.Errorf("invalid JEV candidate gate application plan bridge status")
    }
    if b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("invalid JEV candidate gate status")
    }
    switch b.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if b.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady ||
            b.CandidateDigest == "" || b.CandidateSource == "" ||
            b.PlanDigest == "" || b.PlanSource == "" ||
            b.PlanEvidenceDigest == "" {
            return fmt.Errorf("ready application plan bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if b.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" {
            return fmt.Errorf("held application plan bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if b.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" {
            return fmt.Errorf("rejected application plan bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid JEV candidate gate decision")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV candidate gate application plan bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate gate application plan bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(
        b.Status,
        b.CandidateGateStatus,
        b.CandidateDecision,
        b.CandidateDigest,
        b.CandidateSource,
        b.CandidateGateDigest,
        b.ApplicationPlanStatus,
        b.PlanDigest,
        b.PlanSource,
        b.PlanEvidenceDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate gate application plan bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.CandidateGate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.CandidateGate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.CandidateGate.Validate(); err != nil {
        output.MissingStage = "revision-candidate-gate"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeBound
    output.CandidateGateStatus = input.CandidateGate.Status
    output.CandidateDecision = input.CandidateGate.Decision
    output.CandidateDigest = input.CandidateGate.CandidateDigest
    output.CandidateSource = input.CandidateGate.CandidateSource
    output.CandidateGateDigest = input.CandidateGate.GateDigest
    switch input.CandidateGate.Decision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if input.PlanDigest == "" {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown
            output.MissingStage = "plan-digest"
            return output
        }
        if input.PlanSource == "" {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown
            output.MissingStage = "plan-source"
            return output
        }
        if input.PlanEvidenceDigest == "" {
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown
            output.MissingStage = "plan-evidence"
            return output
        }
        output.ApplicationPlanStatus = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady
        output.PlanDigest = input.PlanDigest
        output.PlanSource = input.PlanSource
        output.PlanEvidenceDigest = input.PlanEvidenceDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        output.ApplicationPlanStatus = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        output.ApplicationPlanStatus = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected
    default:
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown
        output.MissingStage = "candidate-gate-decision"
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(
        output.Status,
        output.CandidateGateStatus,
        output.CandidateDecision,
        output.CandidateDigest,
        output.CandidateSource,
        output.CandidateGateDigest,
        output.ApplicationPlanStatus,
        output.PlanDigest,
        output.PlanSource,
        output.PlanEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown
        output.MissingStage = "candidate-gate-application-plan-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(status, candidateGateStatus, candidateDecision, candidateDigest, candidateSource, candidateGateDigest, applicationPlanStatus, planDigest, planSource, planEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", status, candidateGateStatus, candidateDecision, candidateDigest, candidateSource, candidateGateDigest, applicationPlanStatus, planDigest, planSource, planEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
