package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSP(t *testing.T) {
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
    feedback := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeInput{
        CapabilityScope:        scope,
        ReverseObservation:     observation,
        CandidateDigest:         "candidate-digest",
        CandidateSource:         "candidate-source",
        CandidateGateDigest:     "candidate-gate-digest",
        FeedbackDirection:       jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackDigest:          "feedback-digest",
        FeedbackSource:          "feedback-engine",
        FeedbackEvidenceDigest:  "feedback-evidence",
        NonAuthorizing:          true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSP(feedback)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid feedback LSP projection, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPBoundCode ||
        diagnostic.FeedbackDirection != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove {
        t.Fatalf("candidate feedback evidence was not published: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge{
            Status:        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown,
            MissingStage:  "feedback-evidence",
            NonExecuting:  true,
            NonAuthorizing: true,
        },
    )
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN feedback LSP projection, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.MissingStage != "feedback-evidence" {
        t.Fatalf("UNKNOWN feedback projection lost missing stage: %#v", diagnostic)
    }
}
