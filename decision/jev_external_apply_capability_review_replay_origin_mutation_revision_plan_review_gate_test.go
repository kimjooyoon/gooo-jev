package decision

import "testing"

func TestReviewJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady(t *testing.T) {
    plan := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:          jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        ProposalStatus:  jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        CandidateDigest: "candidate-digest",
        ProposalDigest:  "proposal-digest",
        RevisionSource:  "candidate-source",
        PlanSource:      "revision-plan-source",
        IntentDigest:    "revision-intent-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    plan.PlanDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(
        plan.Status,
        plan.ProposalStatus,
        plan.ProposalDigest,
        plan.PlanSource,
        plan.IntentDigest,
    )

    gate := ReviewJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateInput{
        Plan:                plan,
        ReviewDecision:      jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        ReviewSource:        "review-source",
        ReviewEvidenceDigest: "review-evidence",
        NonAuthorizing:      true,
    })
    if err := gate.Validate(); err != nil {
        t.Fatalf("expected valid revision plan review gate, got %v", err)
    }
    if gate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady ||
        gate.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove ||
        gate.GateDigest == "" {
        t.Fatalf("ready review gate did not preserve evidence: %#v", gate)
    }
}

func TestReviewJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknownWithoutSource(t *testing.T) {
    plan := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:          jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        ProposalStatus:  jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        CandidateDigest: "candidate-digest",
        ProposalDigest:  "proposal-digest",
        RevisionSource:  "candidate-source",
        PlanSource:      "revision-plan-source",
        IntentDigest:    "revision-intent-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    plan.PlanDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(
        plan.Status,
        plan.ProposalStatus,
        plan.ProposalDigest,
        plan.PlanSource,
        plan.IntentDigest,
    )

    gate := ReviewJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateInput{
        Plan:           plan,
        ReviewDecision: jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        NonAuthorizing: true,
    })
    if gate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown ||
        gate.MissingStage != "revision-review-source" ||
        !gate.NonExecuting || !gate.NonAuthorizing {
        t.Fatalf("missing review source was not preserved as UNKNOWN: %#v", gate)
    }
}
