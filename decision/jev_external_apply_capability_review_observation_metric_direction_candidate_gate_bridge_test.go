package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateReady(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge{
        Status:            jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound,
        ObservationStatus: jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded,
        ObservationDigest: "observation-digest",
        MetricName:        "observation-confidence",
        MetricValue:       0.9,
        MetricDigest:      "metric-digest",
        Direction:         jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:            "revision-target",
        CandidateSource:   "candidate-source",
        FeedbackDigest:    "feedback-digest",
        DirectionDigest:   "direction-digest",
        NonExecuting:      true,
        NonAuthorizing:    true,
    }
    bridge.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        bridge.Direction,
        bridge.Target,
        bridge.CandidateSource,
        bridge.FeedbackDigest,
    )
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionBridge(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        bridge.ObservationStatus,
        bridge.ObservationDigest,
        bridge.MetricName,
        bridge.MetricValue,
        bridge.MetricDigest,
        bridge.Direction,
        bridge.Target,
        bridge.CandidateSource,
        bridge.FeedbackDigest,
        bridge.DirectionDigest,
    )

    gateBridge := BindJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGate(JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeInput{
        ObservationMetricDirection: bridge,
        CandidateDigest:             "candidate-digest",
        RevisionSource:              "revision-source",
        NonAuthorizing:              true,
    })
    if err := gateBridge.Validate(); err != nil {
        t.Fatalf("expected valid candidate gate bridge, got %v", err)
    }
    if gateBridge.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
        gateBridge.CandidateDigest != "candidate-digest" ||
        gateBridge.BridgeDigest == "" {
        t.Fatalf("ready candidate gate bridge did not preserve evidence: %#v", gateBridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBlocksMismatchGeneration(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge{
        Status:            jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound,
        ObservationStatus: jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationMismatch,
        ObservationDigest: "observation-digest",
        MetricName:        "observation-confidence",
        MetricValue:       0.1,
        MetricDigest:      "metric-digest",
        Direction:         jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:            "revision-target",
        CandidateSource:   "candidate-source",
        FeedbackDigest:    "feedback-digest",
        DirectionDigest:   "direction-digest",
        NonExecuting:      true,
        NonAuthorizing:    true,
    }
    bridge.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        bridge.Direction,
        bridge.Target,
        bridge.CandidateSource,
        bridge.FeedbackDigest,
    )
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionBridge(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        bridge.ObservationStatus,
        bridge.ObservationDigest,
        bridge.MetricName,
        bridge.MetricValue,
        bridge.MetricDigest,
        bridge.Direction,
        bridge.Target,
        bridge.CandidateSource,
        bridge.FeedbackDigest,
        bridge.DirectionDigest,
    )

    gateBridge := BindJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGate(JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeInput{
        ObservationMetricDirection: bridge,
        CandidateDigest:             "candidate-digest",
        RevisionSource:              "revision-source",
        NonAuthorizing:              true,
    })
    if gateBridge.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown ||
        gateBridge.MissingStage != "observation-direction-consistency" {
        t.Fatalf("mismatch observation was promoted to candidate generation: %#v", gateBridge)
    }
}
