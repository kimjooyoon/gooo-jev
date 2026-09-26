package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound = "candidate-observation-metric-direction-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeInput struct {
    Metric         JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge
    Direction      string
    Target         string
    CandidateSource string
    FeedbackDigest string
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge struct {
    Status                   string
    MissingStage             string
    CandidateApplicationStatus string
    CandidateApplicationDigest string
    ObservationStatus        string
    ObservationDigest        string
    MetricName               string
    MetricValue              float64
    MetricDigest              string
    Direction                string
    Target                   string
    CandidateSource          string
    FeedbackDigest           string
    DirectionDigest          string
    BridgeDigest             string
    NonExecuting             bool
    NonAuthorizing           bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge) Validate() error {
    if b.Status == "" || b.CandidateApplicationStatus == "" ||
        b.CandidateApplicationDigest == "" || b.ObservationStatus == "" ||
        b.ObservationDigest == "" || b.MetricName == "" ||
        b.MetricDigest == "" || b.Direction == "" || b.Target == "" ||
        b.CandidateSource == "" || b.FeedbackDigest == "" ||
        b.DirectionDigest == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV candidate observation metric direction bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound {
        return fmt.Errorf("invalid JEV candidate observation metric direction bridge status")
    }
    if b.CandidateApplicationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady {
        return fmt.Errorf("candidate observation direction requires a ready application candidate")
    }
    if b.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate &&
        b.Direction != jevExternalApplyCapabilityReviewImprovementDirectionHold &&
        b.Direction != jevExternalApplyCapabilityReviewImprovementDirectionReject {
        return fmt.Errorf("invalid JEV candidate observation metric direction")
    }
    if b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved &&
        b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch {
        return fmt.Errorf("invalid JEV candidate observation status")
    }
    if b.ObservationStatus == jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch &&
        b.Direction == jevExternalApplyCapabilityReviewImprovementDirectionGenerate {
        return fmt.Errorf("mismatched candidate observation cannot generate a candidate")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV candidate observation metric direction bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate observation metric direction bridge must be non-authorizing")
    }
    expectedDirection := digestJEVExternalApplyCapabilityReviewImprovementDirection(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        b.Direction,
        b.Target,
        b.CandidateSource,
        b.FeedbackDigest,
    )
    if b.DirectionDigest != expectedDirection {
        return fmt.Errorf("JEV candidate observation metric direction digest mismatch")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(
        b.Status,
        b.CandidateApplicationStatus,
        b.CandidateApplicationDigest,
        b.ObservationStatus,
        b.ObservationDigest,
        b.MetricName,
        b.MetricValue,
        b.MetricDigest,
        b.Direction,
        b.Target,
        b.CandidateSource,
        b.FeedbackDigest,
        b.DirectionDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate observation metric direction bridge digest mismatch")
    }
    return nil
}

func DeriveJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Metric.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Metric.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Metric.Validate(); err != nil {
        output.MissingStage = "candidate-observation-metric"
        return output
    }
    if input.Direction == "" {
        output.MissingStage = "improvement-direction"
        return output
    }
    if input.Target == "" {
        output.MissingStage = "improvement-target"
        return output
    }
    if input.CandidateSource == "" {
        output.MissingStage = "candidate-source"
        return output
    }
    if input.FeedbackDigest == "" {
        output.MissingStage = "feedback-digest"
        return output
    }
    if input.Metric.ObservationStatus == jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch &&
        input.Direction == jevExternalApplyCapabilityReviewImprovementDirectionGenerate {
        output.MissingStage = "observation-direction-consistency"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound
    output.CandidateApplicationStatus = input.Metric.CandidateApplicationStatus
    output.CandidateApplicationDigest = input.Metric.CandidateApplicationDigest
    output.ObservationStatus = input.Metric.ObservationStatus
    output.ObservationDigest = input.Metric.ObservationDigest
    output.MetricName = input.Metric.MetricName
    output.MetricValue = input.Metric.MetricValue
    output.MetricDigest = input.Metric.MetricDigest
    output.Direction = input.Direction
    output.Target = input.Target
    output.CandidateSource = input.CandidateSource
    output.FeedbackDigest = input.FeedbackDigest
    output.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        output.Direction,
        output.Target,
        output.CandidateSource,
        output.FeedbackDigest,
    )
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(
        output.Status,
        output.CandidateApplicationStatus,
        output.CandidateApplicationDigest,
        output.ObservationStatus,
        output.ObservationDigest,
        output.MetricName,
        output.MetricValue,
        output.MetricDigest,
        output.Direction,
        output.Target,
        output.CandidateSource,
        output.FeedbackDigest,
        output.DirectionDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown
        output.MissingStage = "candidate-observation-metric-direction-evidence"
        output.BridgeDigest = ""
        output.DirectionDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(status, candidateApplicationStatus, candidateApplicationDigest, observationStatus, observationDigest, metricName string, metricValue float64, metricDigest, direction, target, candidateSource, feedbackDigest, directionDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", status, candidateApplicationStatus, candidateApplicationDigest, observationStatus, observationDigest, metricName, metricValue, metricDigest, direction, target, candidateSource, feedbackDigest, directionDigest)))
    return hex.EncodeToString(sum[:])
}
