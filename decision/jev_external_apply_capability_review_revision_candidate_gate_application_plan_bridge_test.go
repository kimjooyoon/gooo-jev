package decision

import "testing"

func applicationPlanCandidateGate() JEVExternalApplyCapabilityReviewRevisionCandidateGate {
    gate := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:          jevExternalApplyCapabilityReviewRevisionCandidateGated,
        Decision:        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "candidate-source",
        CandidateDigest: "candidate-digest",
        DirectionDigest: "direction-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    gate.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(
        gate.Status,
        gate.Decision,
        gate.Direction,
        gate.Target,
        gate.CandidateSource,
        gate.CandidateDigest,
        gate.DirectionDigest,
    )
    return gate
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeReady(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid application plan bridge, got %v", err)
    }
    if bridge.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady ||
        bridge.PlanDigest != "plan-digest" ||
        bridge.CandidateDigest != "candidate-digest" {
        t.Fatalf("application plan provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeRequiresPlanDigest(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:  applicationPlanCandidateGate(),
        PlanSource:     "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing: true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown ||
        bridge.MissingStage != "plan-digest" {
        t.Fatalf("missing plan digest was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete application plan bridge unexpectedly has digest: %#v", bridge)
    }
}
