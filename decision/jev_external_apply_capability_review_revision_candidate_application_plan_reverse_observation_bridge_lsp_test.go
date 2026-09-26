package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSP(t *testing.T) {
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid LSP projection, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPBoundCode ||
        diagnostic.ObservationMetricEvidenceDigest != "metric-evidence" {
        t.Fatalf("LSP projection did not preserve bound evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge{
            Status:       jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown,
            MissingStage: "reverse-observation-evidence",
            BridgeDigest: "bridge-digest",
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN LSP projection, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPUnknownCode {
        t.Fatalf("UNKNOWN projection was published: %#v", diagnostic)
    }
}
