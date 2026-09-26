package decision

import "testing"

func candidateObservationMetricDirectionMetric(observationStatus string) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge {
    metric := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge{
        Status:                    jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeBound,
        CandidateApplicationStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        CandidateApplicationDigest: "candidate-application-digest",
        ObservationStatus:          observationStatus,
        ObservationDigest:          "observation-digest",
        MetricName:                 "reverse-confidence",
        MetricValue:                0.92,
        MetricDigest:               "metric-digest",
        NonExecuting:               true,
        NonAuthorizing:             true,
    }
    metric.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(
        metric.Status,
        metric.CandidateApplicationStatus,
        metric.CandidateApplicationDigest,
        metric.ObservationStatus,
        metric.ObservationDigest,
        metric.MetricName,
        metric.MetricValue,
        metric.MetricDigest,
    )
    return metric
}

func TestDeriveJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeRecorded(t *testing.T) {
    bridge := DeriveJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeInput{
        Metric:          candidateObservationMetricDirectionMetric(jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved),
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "candidate-source",
        FeedbackDigest:  "feedback-digest",
        NonAuthorizing:  true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid direction bridge, got %v", err)
    }
    if bridge.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate ||
        bridge.CandidateApplicationDigest != "candidate-application-digest" ||
        bridge.DirectionDigest == "" {
        t.Fatalf("direction bridge did not preserve metric origin: %#v", bridge)
    }
}

func TestDeriveJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBlocksMismatchGeneration(t *testing.T) {
    bridge := DeriveJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeInput{
        Metric:          candidateObservationMetricDirectionMetric(jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch),
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "candidate-source",
        FeedbackDigest:  "feedback-digest",
        NonAuthorizing:  true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown ||
        bridge.MissingStage != "observation-direction-consistency" {
        t.Fatalf("mismatched observation generated direction: %#v", bridge)
    }
}
