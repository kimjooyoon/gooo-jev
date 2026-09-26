package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPReady(t *testing.T) {
    plan := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:          jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        ProposalStatus:  jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady,
        ProposalDigest:  "proposal-digest",
        PlanSource:      "revision-plan-source",
        IntentDigest:    "revision-intent-digest",
        PlanDigest:      "plan-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSP(plan)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.PlanSource != "revision-plan-source" ||
        diagnostic.IntentDigest != "revision-intent-digest" || diagnostic.PlanDigest != "plan-digest" {
        t.Fatalf("ready diagnostic did not preserve plan evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPUnknown(t *testing.T) {
    plan := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
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
