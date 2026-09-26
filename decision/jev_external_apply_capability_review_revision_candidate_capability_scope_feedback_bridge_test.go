package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeReady(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    scope := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput{
        ApplicationPlan:        plan,
        SpiffeID:               "spiffe://example.org/service/gooo",
        Audience:               "gooo-jev",
        SandboxPolicyDigest:    "sandbox-policy-digest",
        NetworkAllowlistDigest: "network-allowlist-digest",
        ScopeEvidenceDigest:    "scope-evidence",
        NonAuthorizing:         true,
    })
    observation := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeInput{
        ApplicationPlan:                 plan,
        ObservationDigest:               "observation-digest",
        ObservationSource:               "reverse-observer",
        ObservationEvidenceDigest:       "observation-evidence",
        ObservationMetricDigest:         "metric-digest",
        ObservationMetricSource:         "observation-metric",
        ObservationMetricEvidenceDigest: "metric-evidence",
        NonAuthorizing:                   true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeInput{
        CapabilityScope:        scope,
        ReverseObservation:     observation,
        CandidateDigest:         "candidate-digest",
        CandidateSource:         "candidate-source",
        CandidateGateDigest:     "candidate-gate-digest",
        FeedbackDirection:      jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackDigest:         "feedback-digest",
        FeedbackSource:         "feedback-engine",
        FeedbackEvidenceDigest: "feedback-evidence",
        NonAuthorizing:         true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid capability scope feedback bridge, got %v", err)
    }
    if bridge.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBound ||
        bridge.FeedbackDirection != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove ||
        bridge.CandidateDigest != "candidate-digest" {
        t.Fatalf("candidate feedback provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeRequiresEvidence(t *testing.T) {
    plan := BindJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeInput{
        CandidateGate:      applicationPlanCandidateGate(),
        PlanDigest:         "plan-digest",
        PlanSource:         "plan-generator",
        PlanEvidenceDigest: "plan-evidence",
        NonAuthorizing:     true,
    })
    scope := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeInput{
        ApplicationPlan:        plan,
        SpiffeID:               "spiffe://example.org/service/gooo",
        Audience:               "gooo-jev",
        SandboxPolicyDigest:    "sandbox-policy-digest",
        NetworkAllowlistDigest: "network-allowlist-digest",
        ScopeEvidenceDigest:    "scope-evidence",
        NonAuthorizing:         true,
    })
    observation := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeInput{
        ApplicationPlan:                 plan,
        ObservationDigest:               "observation-digest",
        ObservationSource:               "reverse-observer",
        ObservationEvidenceDigest:       "observation-evidence",
        ObservationMetricDigest:         "metric-digest",
        ObservationMetricSource:         "observation-metric",
        ObservationMetricEvidenceDigest: "metric-evidence",
        NonAuthorizing:                   true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeInput{
        CapabilityScope:      scope,
        ReverseObservation:   observation,
        CandidateDigest:       "candidate-digest",
        CandidateSource:       "candidate-source",
        CandidateGateDigest:   "candidate-gate-digest",
        FeedbackDirection:    jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackSource:       "feedback-engine",
        FeedbackEvidenceDigest: "feedback-evidence",
        NonAuthorizing:       true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown ||
        bridge.MissingStage != "feedback-digest" {
        t.Fatalf("missing feedback digest was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete feedback bridge unexpectedly has digest: %#v", bridge)
    }
}
