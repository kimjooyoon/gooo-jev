package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPReady(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge{
        Status:               jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeBound,
        CandidateGateStatus:  jevExternalApplyCapabilityReviewRevisionCandidateGated,
        CandidateDecision:    jevExternalApplyCapabilityReviewRevisionCandidateReady,
        CandidateDigest:      "candidate-digest",
        CandidateSource:      "candidate-source",
        CandidateGateDigest:  "candidate-gate-digest",
        ApplicationPlanStatus: jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady,
        PlanDigest:            "plan-digest",
        PlanSource:            "plan-generator",
        PlanEvidenceDigest:   "plan-evidence",
        NonExecuting:         true,
        NonAuthorizing:       true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(
        bridge.Status,
        bridge.CandidateGateStatus,
        bridge.CandidateDecision,
        bridge.CandidateDigest,
        bridge.CandidateSource,
        bridge.CandidateGateDigest,
        bridge.ApplicationPlanStatus,
        bridge.PlanDigest,
        bridge.PlanSource,
        bridge.PlanEvidenceDigest,
    )
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid application plan diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.PlanDigest != "plan-digest" ||
        diagnostic.CandidateDigest != "candidate-digest" {
        t.Fatalf("application plan diagnostic did not preserve evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid application plan bridge was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
