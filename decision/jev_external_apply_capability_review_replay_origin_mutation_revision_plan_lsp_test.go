package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPReady(t *testing.T) {
    plan := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        ProposalStatus:   jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        CandidateDigest:  "candidate-digest",
        ProposalDigest:   "proposal-digest",
        RevisionSource:   "candidate-source",
        PlanSource:       "revision-plan-source",
        IntentDigest:     "revision-intent-digest",
        PlanDigest:       digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan(
            jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
            jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
            "proposal-digest",
            "revision-plan-source",
            "revision-intent-digest",
        ),
        NonExecuting:   true,
        NonAuthorizing: true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSP(plan)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.PlanSource != "revision-plan-source" ||
        diagnostic.IntentDigest != "revision-intent-digest" || diagnostic.PlanDigest == "" {
        t.Fatalf("ready diagnostic did not preserve plan evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPUnknown(t *testing.T) {
    plan := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:         jevExternalApplyCapabilityReviewOriginMutationRevisionPlanReady,
        ProposalStatus: jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        NonExecuting:   true,
        NonAuthorizing: true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSP(plan)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown ||
        diagnostic.MissingStage != "revision-plan-evidence" {
        t.Fatalf("incomplete plan was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
