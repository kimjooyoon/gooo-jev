package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeBound = "candidate-application-observation-metric-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeInput struct {
    Observation    JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge
    MetricName     string
    MetricValue   float64
    MetricDigest   string
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge struct {
    Status                    string
    MissingStage              string
    CandidateApplicationStatus string
    CandidateApplicationDigest string
    ObservationStatus         string
    ObservationDigest         string
    MetricName                string
    MetricValue               float64
    MetricDigest              string
    BridgeDigest              string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge) Validate() error {
    if b.Status == "" || b.CandidateApplicationStatus == "" ||
        b.CandidateApplicationDigest == "" || b.ObservationStatus == "" ||
        b.ObservationDigest == "" || b.MetricName == "" ||
        b.MetricDigest == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV candidate observation metric bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeBound {
        return fmt.Errorf("invalid JEV candidate observation metric bridge status")
    }
    if b.CandidateApplicationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady {
        return fmt.Errorf("candidate observation metric requires a ready application candidate")
    }
    if b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved &&
        b.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch {
        return fmt.Errorf("candidate observation metric has invalid observation status")
    }
    if math.IsNaN(b.MetricValue) || math.IsInf(b.MetricValue, 0) {
        return fmt.Errorf("candidate observation metric value must be finite")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV candidate observation metric bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate observation metric bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(
        b.Status,
        b.CandidateApplicationStatus,
        b.CandidateApplicationDigest,
        b.ObservationStatus,
        b.ObservationDigest,
        b.MetricName,
        b.MetricValue,
        b.MetricDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate observation metric bridge digest mismatch")
    }
    return nil
}

func MeasureJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Observation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Observation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Observation.Validate(); err != nil {
        output.MissingStage = "candidate-application-observation"
        return output
    }
    if input.Observation.CandidateApplicationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady {
        output.MissingStage = "ready-application-candidate"
        return output
    }
    if input.MetricName == "" {
        output.MissingStage = "metric-name"
        return output
    }
    if math.IsNaN(input.MetricValue) || math.IsInf(input.MetricValue, 0) {
        output.MissingStage = "metric-value"
        return output
    }
    if input.MetricDigest == "" {
        output.MissingStage = "metric-digest"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeBound
    output.CandidateApplicationStatus = input.Observation.CandidateApplicationStatus
    output.CandidateApplicationDigest = input.Observation.CandidateApplicationDigest
    output.ObservationStatus = input.Observation.ObservationStatus
    output.ObservationDigest = input.Observation.ObservationDigest
    output.MetricName = input.MetricName
    output.MetricValue = input.MetricValue
    output.MetricDigest = input.MetricDigest
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(
        output.Status,
        output.CandidateApplicationStatus,
        output.CandidateApplicationDigest,
        output.ObservationStatus,
        output.ObservationDigest,
        output.MetricName,
        output.MetricValue,
        output.MetricDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown
        output.MissingStage = "candidate-observation-metric-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(status, candidateApplicationStatus, candidateApplicationDigest, observationStatus, observationDigest, metricName string, metricValue float64, metricDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%0.17g|%s", status, candidateApplicationStatus, candidateApplicationDigest, observationStatus, observationDigest, metricName, metricValue, metricDigest)))
    return hex.EncodeToString(sum[:])
}
