package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeReady(t *testing.T) {
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
        CandidateDigest:        "candidate-digest",
        CandidateSource:        "candidate-source",
        CandidateGateDigest:    "candidate-gate-digest",
        FeedbackDirection:      jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackDigest:         "feedback-digest",
        FeedbackSource:         "feedback-engine",
        FeedbackEvidenceDigest: "feedback-evidence",
        NonAuthorizing:         true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeInput{
        FeedbackBridge:         feedback,
        ProposalDecision:        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise,
        ProposalTarget:          "revision-target",
        ProposalDigest:          "proposal-digest",
        ProposalSource:          "proposal-generator",
        ProposalEvidenceDigest: "proposal-evidence",
        NonAuthorizing:          true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid feedback revision proposal bridge, got %v", err)
    }
    if bridge.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBound ||
        bridge.ProposalDecision != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise ||
        bridge.CandidateDigest != "candidate-digest" {
        t.Fatalf("revision proposal provenance was not preserved: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeRequiresEvidence(t *testing.T) {
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
        CandidateDigest:        "candidate-digest",
        CandidateSource:        "candidate-source",
        CandidateGateDigest:    "candidate-gate-digest",
        FeedbackDirection:      jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackImprove,
        FeedbackDigest:         "feedback-digest",
        FeedbackSource:         "feedback-engine",
        FeedbackEvidenceDigest: "feedback-evidence",
        NonAuthorizing:         true,
    })
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge(JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeInput{
        FeedbackBridge:         feedback,
        ProposalDecision:        jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalRevise,
        ProposalTarget:          "revision-target",
        ProposalSource:          "proposal-generator",
        ProposalEvidenceDigest: "proposal-evidence",
        NonAuthorizing:          true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridgeUnknown ||
        bridge.MissingStage != "proposal-digest" {
        t.Fatalf("missing proposal digest was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete revision proposal bridge unexpectedly has digest: %#v", bridge)
    }
}
