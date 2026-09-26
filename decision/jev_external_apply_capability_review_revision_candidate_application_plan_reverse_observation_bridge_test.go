package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeReady(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeInput{
        ApplicationPlan:                 plan,
        ObservationDigest:               "observation-digest",
        ObservationSource:               "reverse-observer",
        ObservationEvidenceDigest:       "observation-evidence",
        ObservationMetricDigest:         "metric-digest",
        ObservationMetricSource:         "observation-metric",
        ObservationMetricEvidenceDigest: "metric-evidence",
        NonAuthorizing:                   true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid application plan reverse observation bridge, got %v", err)
    }
    if bridge.ReverseObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBound ||
        bridge.ObservationMetricStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanObservationMetricBound ||
        bridge.PlanDigest != "plan-digest" {
        t.Fatalf("application plan reverse observation provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeRequiresEvidence(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeInput{
        ApplicationPlan:           plan,
        ObservationDigest:         "observation-digest",
        ObservationSource:         "reverse-observer",
        ObservationEvidenceDigest: "observation-evidence",
        ObservationMetricDigest:   "metric-digest",
        ObservationMetricSource:   "observation-metric",
        NonAuthorizing:             true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown ||
        bridge.MissingStage != "observation-metric-evidence" {
        t.Fatalf("missing observation metric evidence was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete observation bridge unexpectedly has digest: %#v", bridge)
    }
}
