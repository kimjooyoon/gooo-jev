package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSP(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput{
        ApplicationPlan:        plan,
        SpiffeID:               "spiffe://example.org/service/gooo",
        Audience:               "gooo-jev",
        SandboxPolicyDigest:    "sandbox-policy-digest",
        NetworkAllowlistDigest: "network-allowlist-digest",
        ScopeEvidenceDigest:    "scope-evidence",
        NonAuthorizing:         true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid capability scope LSP projection, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPBoundCode ||
        diagnostic.SpiffeID != "spiffe://example.org/service/gooo" {
        t.Fatalf("capability scope evidence was not published: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge{
            Status:       jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown,
            MissingStage: "network-allowlist",
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN capability scope projection, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.MissingStage != "network-allowlist" {
        t.Fatalf("UNKNOWN capability scope projection lost missing stage: %#v", diagnostic)
    }
}
