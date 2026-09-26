package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeBound = "application-plan-reverse-observation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound = "reverse-observation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationHeld = "reverse-observation-held"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationRejected = "reverse-observation-rejected"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound = "observation-metric-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricHeld = "observation-metric-held"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricRejected = "observation-metric-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeInput struct {
    ApplicationPlan                    JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge
    ObservationDigest                   string
    ObservationSource                   string
    ObservationEvidenceDigest          string
    ObservationMetricDigest            string
    ObservationMetricSource            string
    ObservationMetricEvidenceDigest    string
    NonAuthorizing                     bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge struct {
    Status                         string
    MissingStage                   string
    ApplicationPlanStatus          string
    PlanDigest                     string
    PlanSource                     string
    PlanEvidenceDigest             string
    ApplicationPlanBridgeDigest    string
    ReverseObservationStatus       string
    ObservationDigest              string
    ObservationSource              string
    ObservationEvidenceDigest      string
    ObservationMetricStatus        string
    ObservationMetricDigest        string
    ObservationMetricSource        string
    ObservationMetricEvidenceDigest string
    BridgeDigest                   string
    NonExecuting                   bool
    NonAuthorizing                 bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge) Validate() error {
    if b.Status == "" || b.ApplicationPlanStatus == "" ||
        b.ApplicationPlanBridgeDigest == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV application plan reverse observation bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeBound {
        return fmt.Errorf("invalid JEV application plan reverse observation bridge status")
    }
    switch b.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if b.PlanDigest == "" || b.PlanSource == "" || b.PlanEvidenceDigest == "" ||
            b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound ||
            b.ObservationDigest == "" || b.ObservationSource == "" || b.ObservationEvidenceDigest == "" ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound ||
            b.ObservationMetricDigest == "" || b.ObservationMetricSource == "" ||
            b.ObservationMetricEvidenceDigest == "" {
            return fmt.Errorf("ready application plan reverse observation bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationHeld ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricHeld ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" ||
            b.ObservationDigest != "" || b.ObservationSource != "" || b.ObservationEvidenceDigest != "" ||
            b.ObservationMetricDigest != "" || b.ObservationMetricSource != "" ||
            b.ObservationMetricEvidenceDigest != "" {
            return fmt.Errorf("held application plan reverse observation bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        if b.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationRejected ||
            b.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricRejected ||
            b.PlanDigest != "" || b.PlanSource != "" || b.PlanEvidenceDigest != "" ||
            b.ObservationDigest != "" || b.ObservationSource != "" || b.ObservationEvidenceDigest != "" ||
            b.ObservationMetricDigest != "" || b.ObservationMetricSource != "" ||
            b.ObservationMetricEvidenceDigest != "" {
            return fmt.Errorf("rejected application plan reverse observation bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid application plan status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV application plan reverse observation bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV application plan reverse observation bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(
        b.Status,
        b.ApplicationPlanStatus,
        b.PlanDigest,
        b.PlanSource,
        b.PlanEvidenceDigest,
        b.ApplicationPlanBridgeDigest,
        b.ReverseObservationStatus,
        b.ObservationDigest,
        b.ObservationSource,
        b.ObservationEvidenceDigest,
        b.ObservationMetricStatus,
        b.ObservationMetricDigest,
        b.ObservationMetricSource,
        b.ObservationMetricEvidenceDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV application plan reverse observation bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.ApplicationPlan.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.ApplicationPlan.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.ApplicationPlan.Validate(); err != nil {
        output.MissingStage = "application-plan"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeBound
    output.ApplicationPlanStatus = input.ApplicationPlan.ApplicationPlanStatus
    output.PlanDigest = input.ApplicationPlan.PlanDigest
    output.PlanSource = input.ApplicationPlan.PlanSource
    output.PlanEvidenceDigest = input.ApplicationPlan.PlanEvidenceDigest
    output.ApplicationPlanBridgeDigest = input.ApplicationPlan.BridgeDigest
    switch input.ApplicationPlan.ApplicationPlanStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady:
        if input.ObservationDigest == "" {
            output.MissingStage = "reverse-observation-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
            return output
        }
        if input.ObservationSource == "" {
            output.MissingStage = "reverse-observation-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
            return output
        }
        if input.ObservationEvidenceDigest == "" {
            output.MissingStage = "reverse-observation-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
            return output
        }
        if input.ObservationMetricDigest == "" {
            output.MissingStage = "observation-metric-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
            return output
        }
        if input.ObservationMetricSource == "" {
            output.MissingStage = "observation-metric-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
            return output
        }
        if input.ObservationMetricEvidenceDigest == "" {
            output.MissingStage = "observation-metric-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
            return output
        }
        output.ReverseObservationStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound
        output.ObservationDigest = input.ObservationDigest
        output.ObservationSource = input.ObservationSource
        output.ObservationEvidenceDigest = input.ObservationEvidenceDigest
        output.ObservationMetricStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound
        output.ObservationMetricDigest = input.ObservationMetricDigest
        output.ObservationMetricSource = input.ObservationMetricSource
        output.ObservationMetricEvidenceDigest = input.ObservationMetricEvidenceDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld:
        output.ReverseObservationStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationHeld
        output.ObservationMetricStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected:
        output.ReverseObservationStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationRejected
        output.ObservationMetricStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricRejected
    default:
        output.MissingStage = "application-plan-status"
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(
        output.Status,
        output.ApplicationPlanStatus,
        output.PlanDigest,
        output.PlanSource,
        output.PlanEvidenceDigest,
        output.ApplicationPlanBridgeDigest,
        output.ReverseObservationStatus,
        output.ObservationDigest,
        output.ObservationSource,
        output.ObservationEvidenceDigest,
        output.ObservationMetricStatus,
        output.ObservationMetricDigest,
        output.ObservationMetricSource,
        output.ObservationMetricEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
        output.MissingStage = "application-plan-reverse-observation-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(status, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus, observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", status, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus, observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
