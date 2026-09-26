package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound = "observation-metric-direction-bound"
const jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeInput struct {
    Observation    JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation
    Metric         JEVExternalApplyCapabilityReviewMetric
    Direction      JEVExternalApplyCapabilityReviewImprovementDirection
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge struct {
    Status             string
    MissingStage       string
    ObservationStatus  string
    ObservationDigest  string
    MetricName         string
    MetricValue        float64
    MetricDigest       string
    Direction          string
    Target             string
    CandidateSource    string
    FeedbackDigest     string
    DirectionDigest    string
    BridgeDigest       string
    NonExecuting       bool
    NonAuthorizing     bool
}

func (b JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge) Validate() error {
    if b.Status == "" || b.ObservationStatus == "" || b.ObservationDigest == "" ||
        b.MetricName == "" || b.MetricDigest == "" || b.Direction == "" ||
        b.Target == "" || b.CandidateSource == "" || b.FeedbackDigest == "" ||
        b.DirectionDigest == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV observation metric direction bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound {
        return fmt.Errorf("invalid JEV observation metric direction bridge status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV observation metric direction bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV observation metric direction bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewObservationMetricDirectionBridge(
        b.Status,
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
        return fmt.Errorf("JEV observation metric direction bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewObservationMetricDirection(input JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeInput) JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge {
    output := JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge{
        Status:         jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Observation.NonAuthorizing ||
        !input.Metric.NonAuthorizing || !input.Direction.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Observation.NonExecuting || !input.Metric.NonExecuting || !input.Direction.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Observation.Validate(); err != nil {
        output.MissingStage = "application-observation"
        return output
    }
    if err := input.Metric.Validate(); err != nil {
        output.MissingStage = "capability-review-metric"
        return output
    }
    if err := input.Direction.Validate(); err != nil {
        output.MissingStage = "improvement-direction"
        return output
    }
    if input.Observation.Status == jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch &&
        input.Direction.Direction == jevExternalApplyCapabilityReviewImprovementDirectionGenerate {
        output.MissingStage = "observation-direction-consistency"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound
    output.ObservationStatus = input.Observation.Status
    output.ObservationDigest = input.Observation.ObservationDigest
    output.MetricName = input.Metric.MetricName
    output.MetricValue = input.Metric.MetricValue
    output.MetricDigest = input.Metric.MetricDigest
    output.Direction = input.Direction.Direction
    output.Target = input.Direction.Target
    output.CandidateSource = input.Direction.CandidateSource
    output.FeedbackDigest = input.Direction.FeedbackDigest
    output.DirectionDigest = input.Direction.DirectionDigest
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionBridge(
        output.Status,
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
        output.Status = jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown
        output.MissingStage = "observation-metric-direction-evidence"
        output.ObservationDigest = ""
        output.MetricDigest = ""
        output.DirectionDigest = ""
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewObservationMetricDirectionBridge(status, observationStatus, observationDigest, metricName string, metricValue float64, metricDigest, direction, target, candidateSource, feedbackDigest, directionDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%0.17g|%s|%s|%s|%s|%s|%s", status, observationStatus, observationDigest, metricName, metricValue, metricDigest, direction, target, candidateSource, feedbackDigest, directionDigest)))
    return hex.EncodeToString(sum[:])
}
