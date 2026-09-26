package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSPBound(t *testing.T) {
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
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionBridge(
        bridge.Status,
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

    diagnostic := ProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid bound diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.MetricName != "observation-confidence" ||
        diagnostic.Direction != jevExternalApplyCapabilityReviewImprovementDirectionGenerate ||
        diagnostic.BridgeDigest == "" {
        t.Fatalf("bound diagnostic did not preserve bridge evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSPUnknown(t *testing.T) {
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

    diagnostic := ProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown ||
        diagnostic.MissingStage != "observation-metric-direction-evidence" {
        t.Fatalf("incomplete bridge was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
