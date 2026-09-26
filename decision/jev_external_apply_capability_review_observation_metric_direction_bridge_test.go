package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewObservationMetricDirectionRecorded(t *testing.T) {
    observation := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded,
        CandidateStatus:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady,
        CandidateDigest:           "candidate-digest",
        ObservationStatus:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved,
        ObservationSource:          "observation-source",
        ObservationEvidenceDigest: "observation-evidence",
        ReverseObservationDigest:  "reverse-observation",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    observation.ObservationDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation(
        observation.Status,
        observation.CandidateStatus,
        observation.CandidateDigest,
        observation.ObservationStatus,
        observation.ObservationSource,
        observation.ObservationEvidenceDigest,
        observation.ReverseObservationDigest,
    )

    metric := JEVExternalApplyCapabilityReviewMetric{
        Status:         jevExternalApplyCapabilityReviewMetricRecorded,
        ReviewStatus:   jevExternalApplyCapabilityReviewConfirmed,
        ReviewDigest:   "review-digest",
        MetricName:     "observation-confidence",
        MetricValue:    0.9,
        EvidenceDigest: "metric-evidence",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    metric.MetricDigest = digestJEVExternalApplyCapabilityReviewMetric(
        metric.ReviewStatus,
        metric.ReviewDigest,
        metric.MetricName,
        metric.MetricValue,
        metric.EvidenceDigest,
    )

    direction := JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:          jevExternalApplyCapabilityReviewImprovementDirectionBound,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "candidate-source",
        FeedbackDigest:  "feedback-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    direction.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(
        direction.Status,
        direction.Direction,
        direction.Target,
        direction.CandidateSource,
        direction.FeedbackDigest,
    )

    bridge := BindJEVExternalApplyCapabilityReviewObservationMetricDirection(JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeInput{
        Observation:    observation,
        Metric:         metric,
        Direction:      direction,
        NonAuthorizing: true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid observation metric direction bridge, got %v", err)
    }
    if bridge.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound ||
        bridge.MetricName != "observation-confidence" ||
        bridge.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate ||
        bridge.BridgeDigest == "" {
        t.Fatalf("bridge did not preserve improvement evidence: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewObservationMetricDirectionMismatchBlocksGeneration(t *testing.T) {
    observation := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch,
        CandidateStatus:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady,
        CandidateDigest:           "candidate-digest",
        ObservationStatus:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatched,
        ObservationSource:          "observation-source",
        ObservationEvidenceDigest: "observation-evidence",
        ReverseObservationDigest:  "reverse-observation",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    observation.ObservationDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation(
        observation.Status,
        observation.CandidateStatus,
        observation.CandidateDigest,
        observation.ObservationStatus,
        observation.ObservationSource,
        observation.ObservationEvidenceDigest,
        observation.ReverseObservationDigest,
    )

    metric := JEVExternalApplyCapabilityReviewMetric{
        Status:         jevExternalApplyCapabilityReviewMetricRecorded,
        ReviewStatus:   jevExternalApplyCapabilityReviewConfirmed,
        ReviewDigest:   "review-digest",
        MetricName:     "observation-confidence",
        MetricValue:    0.1,
        EvidenceDigest: "metric-evidence",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    metric.MetricDigest = digestJEVExternalApplyCapabilityReviewMetric(metric.ReviewStatus, metric.ReviewDigest, metric.MetricName, metric.MetricValue, metric.EvidenceDigest)
    direction := JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:          jevExternalApplyCapabilityReviewImprovementDirectionBound,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "candidate-source",
        FeedbackDigest:  "feedback-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    direction.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(direction.Status, direction.Direction, direction.Target, direction.CandidateSource, direction.FeedbackDigest)

    bridge := BindJEVExternalApplyCapabilityReviewObservationMetricDirection(JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeInput{
        Observation:    observation,
        Metric:         metric,
        Direction:      direction,
        NonAuthorizing: true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown ||
        bridge.MissingStage != "observation-direction-consistency" {
        t.Fatalf("mismatched observation was not blocked from generation: %#v", bridge)
    }
}
