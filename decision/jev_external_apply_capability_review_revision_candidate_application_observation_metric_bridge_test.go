package decision

import "testing"

func candidateObservationMetricBridgeObservation() JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge {
    observation := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge{
        Status:                     jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeBound,
        CandidateApplicationStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        CandidateApplicationDigest: "candidate-application-digest",
        ObservationStatus:          jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationSource:           "reverse-observation-source",
        ObservationEvidenceDigest:  "observation-evidence",
        ReverseObservationDigest:   "reverse-observation",
        NonExecuting:               true,
        NonAuthorizing:             true,
    }
    observation.ObservationDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(
        observation.Status,
        observation.CandidateApplicationStatus,
        observation.CandidateApplicationDigest,
        observation.ObservationStatus,
        observation.ObservationSource,
        observation.ObservationEvidenceDigest,
        observation.ReverseObservationDigest,
    )
    return observation
}

func TestMeasureJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeRecorded(t *testing.T) {
    bridge := MeasureJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeInput{
        Observation:  candidateObservationMetricBridgeObservation(),
        MetricName:   "reverse-confidence",
        MetricValue:  0.92,
        MetricDigest: "metric-digest",
        NonAuthorizing: true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid candidate observation metric bridge, got %v", err)
    }
    if bridge.CandidateApplicationDigest != "candidate-application-digest" ||
        bridge.ObservationDigest == "" || bridge.MetricName != "reverse-confidence" {
        t.Fatalf("candidate observation metric provenance was not preserved: %#v", bridge)
    }
}

func TestMeasureJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeRequiresMetricDigest(t *testing.T) {
    bridge := MeasureJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeInput{
        Observation: candidateObservationMetricBridgeObservation(),
        MetricName:  "reverse-confidence",
        MetricValue: 0.92,
        NonAuthorizing: true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown ||
        bridge.MissingStage != "metric-digest" {
        t.Fatalf("missing metric digest was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete metric bridge unexpectedly has digest: %#v", bridge)
    }
}
