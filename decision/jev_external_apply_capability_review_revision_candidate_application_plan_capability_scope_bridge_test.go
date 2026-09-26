package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeReady(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput{
        ApplicationPlan:       plan,
        SpiffeID:               "spiffe://example.org/service/gooo",
        Audience:              "gooo-jev",
        SandboxPolicyDigest:   "sandbox-policy-digest",
        NetworkAllowlistDigest: "network-allowlist-digest",
        ScopeEvidenceDigest:   "scope-evidence",
        NonAuthorizing:         true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid capability scope bridge, got %v", err)
    }
    if bridge.ScopeStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBound ||
        bridge.SpiffeID != "spiffe://example.org/service/gooo" ||
        bridge.NetworkAllowlistDigest != "network-allowlist-digest" {
        t.Fatalf("capability scope provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeRequiresNetworkAllowlist(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput{
        ApplicationPlan:     plan,
        SpiffeID:             "spiffe://example.org/service/gooo",
        Audience:             "gooo-jev",
        SandboxPolicyDigest: "sandbox-policy-digest",
        ScopeEvidenceDigest: "scope-evidence",
        NonAuthorizing:      true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown ||
        bridge.MissingStage != "network-allowlist" {
        t.Fatalf("missing network allowlist was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete capability scope bridge unexpectedly has digest: %#v", bridge)
    }
}
