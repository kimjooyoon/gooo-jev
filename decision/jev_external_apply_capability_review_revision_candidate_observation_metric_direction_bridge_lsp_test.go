package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSPBound(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge{
        Status:                    jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound,
        CandidateApplicationStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        CandidateApplicationDigest: "candidate-application-digest",
        ObservationStatus:          jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationDigest:           "observation-digest",
        MetricName:                  "reverse-confidence",
        MetricValue:                 0.92,
        MetricDigest:                "metric-digest",
        Direction:                   jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:                      "revision-target",
        CandidateSource:             "candidate-source",
        FeedbackDigest:              "feedback-digest",
        DirectionDigest:             digestJEVExternalApplyCapabilityReviewImprovementDirection(
            jevExternalApplyCapabilityReviewImprovementDirectionBound,
            jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
            "revision-target",
            "candidate-source",
            "feedback-digest",
        ),
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(
        bridge.Status,
        bridge.CandidateApplicationStatus,
        bridge.CandidateApplicationDigest,
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid direction diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.DirectionDigest != bridge.DirectionDigest {
        t.Fatalf("direction diagnostic did not preserve provenance: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid direction bridge was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
