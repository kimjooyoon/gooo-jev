package decision

import "testing"

func TestPlanJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionReady(t *testing.T) {
    proposal := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        BridgeStatus:              jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady,
        CandidateGenerationDigest: "generation-digest",
        GateDigest:                "gate-digest",
        CandidateDigest:           "candidate-digest",
        RevisionSource:            "candidate-source",
        RevisionPlanSource:        "revision-plan-source",
        RevisionIntentDigest:      "revision-intent-digest",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    proposal.ProposalDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal(
        proposal.Status,
        proposal.BridgeStatus,
        proposal.CandidateGenerationDigest,
        proposal.GateDigest,
        proposal.CandidateDigest,
        proposal.RevisionSource,
        proposal.RevisionPlanSource,
        proposal.RevisionIntentDigest,
    )

    plan := PlanJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanInput{
        Proposal:       proposal,
        PlanSource:     "revision-plan-source",
        NonAuthorizing: true,
    })
    if err := plan.Validate(); err != nil {
        t.Fatalf("expected valid revision plan, got %v", err)
    }
    if plan.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady ||
        plan.PlanSource != "revision-plan-source" ||
        plan.IntentDigest != "revision-intent-digest" ||
        plan.PlanDigest == "" {
        t.Fatalf("revision plan did not preserve proposal evidence: %#v", plan)
    }
}

func TestPlanJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionUnknownOnSourceMismatch(t *testing.T) {
    proposal := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        BridgeStatus:              jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady,
        CandidateGenerationDigest: "generation-digest",
        GateDigest:                "gate-digest",
        CandidateDigest:           "candidate-digest",
        RevisionSource:            "candidate-source",
        RevisionPlanSource:        "revision-plan-source",
        RevisionIntentDigest:      "revision-intent-digest",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    proposal.ProposalDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal(
        proposal.Status,
        proposal.BridgeStatus,
        proposal.CandidateGenerationDigest,
        proposal.GateDigest,
        proposal.CandidateDigest,
        proposal.RevisionSource,
        proposal.RevisionPlanSource,
        proposal.RevisionIntentDigest,
    )

    plan := PlanJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanInput{
        Proposal:       proposal,
        PlanSource:     "different-plan-source",
        NonAuthorizing: true,
    })
    if plan.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown ||
        plan.MissingStage != "revision-plan-source-consistency" ||
        plan.NonExecuting != true || plan.NonAuthorizing != true {
        t.Fatalf("source mismatch was not preserved as UNKNOWN: %#v", plan)
    }
}
