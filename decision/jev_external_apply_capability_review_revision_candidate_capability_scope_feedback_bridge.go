package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeBound = "capability-scope-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound = "candidate-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld = "candidate-feedback-held"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected = "candidate-feedback-rejected"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove = "improve"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRetain = "retain"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRollback = "rollback"

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeInput struct {
    CapabilityScope     JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge
    ReverseObservation  JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge
    FeedbackDirection   string
    FeedbackDigest      string
    FeedbackSource      string
    FeedbackEvidenceDigest string
    NonAuthorizing      bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge struct {
    Status                         string
    MissingStage                   string
    CandidateDigest                string
    CandidateSource                string
    CandidateGateDigest            string
    ApplicationPlanStatus          string
    PlanDigest                     string
    PlanSource                     string
    PlanEvidenceDigest             string
    ApplicationPlanBridgeDigest    string
    ScopeStatus                    string
    SpiffeID                       string
    Audience                       string
    SandboxPolicyDigest            string
    NetworkAllowlistDigest         string
    ScopeEvidenceDigest            string
    CapabilityScopeBridgeDigest    string
    ReverseObservationStatus       string
    ObservationDigest              string
    ObservationSource              string
    ObservationEvidenceDigest      string
    ObservationMetricStatus        string
    ObservationMetricDigest        string
    ObservationMetricSource        string
    ObservationMetricEvidenceDigest string
    ReverseObservationBridgeDigest string
    FeedbackDirection              string
    FeedbackDigest                 string
    FeedbackSource                 string
    FeedbackEvidenceDigest         string
    FeedbackStatus                 string
    BridgeDigest                   string
    NonExecuting                   bool
    NonAuthorizing                 bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge) Validate() error {
    if b.Status == "" || b.ApplicationPlanStatus == "" ||
        b.ApplicationPlanBridgeDigest == "" || b.ScopeStatus == "" ||
        b.CapabilityScopeBridgeDigest == "" || b.ReverseObservationStatus == "" ||
        b.ReverseObservationBridgeDigest == "" || b.FeedbackStatus == "" ||
        b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV capability scope feedback bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeBound {
        return fmt.Errorf("invalid JEV capability scope feedback bridge status")
    }
    switch b.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if b.CandidateDigest == "" || b.CandidateSource == "" || b.CandidateGateDigest == "" ||
            b.PlanDigest == "" || b.PlanSource == "" || b.PlanEvidenceDigest == "" ||
            b.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBound ||
            b.SpiffeID == "" || b.Audience == "" || b.SandboxPolicyDigest == "" ||
            b.NetworkAllowlistDigest == "" || b.ScopeEvidenceDigest == "" ||
            b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound ||
            b.ObservationDigest == "" || b.ObservationSource == "" || b.ObservationEvidenceDigest == "" ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound ||
            b.ObservationMetricDigest == "" || b.ObservationMetricSource == "" ||
            b.ObservationMetricEvidenceDigest == "" ||
            b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound ||
            !validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackDirection(b.FeedbackDirection) ||
            b.FeedbackDigest == "" || b.FeedbackSource == "" ||
            b.FeedbackEvidenceDigest == "" {
            return fmt.Errorf("ready capability scope feedback bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        if b.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeHeld ||
            b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationHeld ||
            b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld ||
            b.FeedbackDirection != "" || b.FeedbackDigest != "" || b.FeedbackSource != "" ||
            b.FeedbackEvidenceDigest != "" {
            return fmt.Errorf("held capability scope feedback bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        if b.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeRejected ||
            b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationRejected ||
            b.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected ||
            b.FeedbackDirection != "" || b.FeedbackDigest != "" || b.FeedbackSource != "" ||
            b.FeedbackEvidenceDigest != "" {
            return fmt.Errorf("rejected capability scope feedback bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid capability scope feedback application plan status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV capability scope feedback bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV capability scope feedback bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(
        b.Status,
        b.CandidateDigest,
        b.CandidateSource,
        b.CandidateGateDigest,
        b.ApplicationPlanStatus,
        b.PlanDigest,
        b.PlanSource,
        b.PlanEvidenceDigest,
        b.ApplicationPlanBridgeDigest,
        b.ScopeStatus,
        b.SpiffeID,
        b.Audience,
        b.SandboxPolicyDigest,
        b.NetworkAllowlistDigest,
        b.ScopeEvidenceDigest,
        b.CapabilityScopeBridgeDigest,
        b.ReverseObservationStatus,
        b.ObservationDigest,
        b.ObservationSource,
        b.ObservationEvidenceDigest,
        b.ObservationMetricStatus,
        b.ObservationMetricDigest,
        b.ObservationMetricSource,
        b.ObservationMetricEvidenceDigest,
        b.ReverseObservationBridgeDigest,
        b.FeedbackDirection,
        b.FeedbackDigest,
        b.FeedbackSource,
        b.FeedbackEvidenceDigest,
        b.FeedbackStatus,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV capability scope feedback bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.CapabilityScope.NonAuthorizing || !input.ReverseObservation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.CapabilityScope.NonExecuting || !input.ReverseObservation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.CapabilityScope.Validate(); err != nil {
        output.MissingStage = "capability-scope"
        return output
    }
    if err := input.ReverseObservation.Validate(); err != nil {
        output.MissingStage = "reverse-observation"
        return output
    }
    if input.CapabilityScope.ApplicationPlanBridgeDigest != input.ReverseObservation.ApplicationPlanBridgeDigest ||
        input.CapabilityScope.ApplicationPlanStatus != input.ReverseObservation.ApplicationPlanStatus ||
        input.CapabilityScope.PlanDigest != input.ReverseObservation.PlanDigest ||
        input.CapabilityScope.PlanSource != input.ReverseObservation.PlanSource ||
        input.CapabilityScope.PlanEvidenceDigest != input.ReverseObservation.PlanEvidenceDigest {
        output.MissingStage = "provenance-consistency"
        return output
    }

    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeBound
    output.CandidateDigest = input.CapabilityScope.ApplicationPlan.CandidateDigest
    output.CandidateSource = input.CapabilityScope.ApplicationPlan.CandidateSource
    output.CandidateGateDigest = input.CapabilityScope.ApplicationPlan.CandidateGateDigest
    output.ApplicationPlanStatus = input.CapabilityScope.ApplicationPlanStatus
    output.PlanDigest = input.CapabilityScope.PlanDigest
    output.PlanSource = input.CapabilityScope.PlanSource
    output.PlanEvidenceDigest = input.CapabilityScope.PlanEvidenceDigest
    output.ApplicationPlanBridgeDigest = input.CapabilityScope.ApplicationPlanBridgeDigest
    output.ScopeStatus = input.CapabilityScope.ScopeStatus
    output.SpiffeID = input.CapabilityScope.SpiffeID
    output.Audience = input.CapabilityScope.Audience
    output.SandboxPolicyDigest = input.CapabilityScope.SandboxPolicyDigest
    output.NetworkAllowlistDigest = input.CapabilityScope.NetworkAllowlistDigest
    output.ScopeEvidenceDigest = input.CapabilityScope.ScopeEvidenceDigest
    output.CapabilityScopeBridgeDigest = input.CapabilityScope.BridgeDigest
    output.ReverseObservationStatus = input.ReverseObservation.ReverseObservationStatus
    output.ObservationDigest = input.ReverseObservation.ObservationDigest
    output.ObservationSource = input.ReverseObservation.ObservationSource
    output.ObservationEvidenceDigest = input.ReverseObservation.ObservationEvidenceDigest
    output.ObservationMetricStatus = input.ReverseObservation.ObservationMetricStatus
    output.ObservationMetricDigest = input.ReverseObservation.ObservationMetricDigest
    output.ObservationMetricSource = input.ReverseObservation.ObservationMetricSource
    output.ObservationMetricEvidenceDigest = input.ReverseObservation.ObservationMetricEvidenceDigest
    output.ReverseObservationBridgeDigest = input.ReverseObservation.BridgeDigest

    switch input.CapabilityScope.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if !validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackDirection(input.FeedbackDirection) {
            output.MissingStage = "feedback-direction"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
            return output
        }
        if input.FeedbackDigest == "" {
            output.MissingStage = "feedback-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
            return output
        }
        if input.FeedbackSource == "" {
            output.MissingStage = "feedback-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
            return output
        }
        if input.FeedbackEvidenceDigest == "" {
            output.MissingStage = "feedback-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
            return output
        }
        output.FeedbackDirection = input.FeedbackDirection
        output.FeedbackDigest = input.FeedbackDigest
        output.FeedbackSource = input.FeedbackSource
        output.FeedbackEvidenceDigest = input.FeedbackEvidenceDigest
        output.FeedbackStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        output.FeedbackStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        output.FeedbackStatus = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRejected
    default:
        output.MissingStage = "application-plan-status"
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(
        output.Status,
        output.CandidateDigest,
        output.CandidateSource,
        output.CandidateGateDigest,
        output.ApplicationPlanStatus,
        output.PlanDigest,
        output.PlanSource,
        output.PlanEvidenceDigest,
        output.ApplicationPlanBridgeDigest,
        output.ScopeStatus,
        output.SpiffeID,
        output.Audience,
        output.SandboxPolicyDigest,
        output.NetworkAllowlistDigest,
        output.ScopeEvidenceDigest,
        output.CapabilityScopeBridgeDigest,
        output.ReverseObservationStatus,
        output.ObservationDigest,
        output.ObservationSource,
        output.ObservationEvidenceDigest,
        output.ObservationMetricStatus,
        output.ObservationMetricDigest,
        output.ObservationMetricSource,
        output.ObservationMetricEvidenceDigest,
        output.ReverseObservationBridgeDigest,
        output.FeedbackDirection,
        output.FeedbackDigest,
        output.FeedbackSource,
        output.FeedbackEvidenceDigest,
        output.FeedbackStatus,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
        output.MissingStage = "capability-scope-feedback-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func validJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackDirection(direction string) bool {
    switch direction {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRetain,
        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRollback:
        return true
    default:
        return false
    }
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(status, candidateDigest, candidateSource, candidateGateDigest, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, scopeStatus, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest, scopeEvidenceDigest, capabilityScopeBridgeDigest, reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus, observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest, reverseObservationBridgeDigest, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackStatus string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{
        status, candidateDigest, candidateSource, candidateGateDigest, applicationPlanStatus,
        planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest,
        scopeStatus, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest,
        scopeEvidenceDigest, capabilityScopeBridgeDigest, reverseObservationStatus,
        observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus,
        observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest,
        reverseObservationBridgeDigest, feedbackDirection, feedbackDigest, feedbackSource,
        feedbackEvidenceDigest, feedbackStatus,
    }, "|")))
    return hex.EncodeToString(sum[:])
}
