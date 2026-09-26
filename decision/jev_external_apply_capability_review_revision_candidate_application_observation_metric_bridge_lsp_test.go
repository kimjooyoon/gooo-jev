package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSPBound(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge{
        Status:                    jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeBound,
        CandidateApplicationStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        CandidateApplicationDigest: "candidate-application-digest",
        ObservationStatus:          jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationDigest:          "observation-digest",
        MetricName:                 "reverse-confidence",
        MetricValue:                0.92,
        MetricDigest:               "metric-digest",
        NonExecuting:               true,
        NonAuthorizing:             true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(
        bridge.Status,
        bridge.CandidateApplicationStatus,
        bridge.CandidateApplicationDigest,
        bridge.ObservationStatus,
        bridge.ObservationDigest,
        bridge.MetricName,
        bridge.MetricValue,
        bridge.MetricDigest,
    )
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid metric diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.MetricDigest != "metric-digest" ||
        diagnostic.CandidateApplicationDigest != "candidate-application-digest" {
        t.Fatalf("metric diagnostic did not preserve provenance: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid metric bridge was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
