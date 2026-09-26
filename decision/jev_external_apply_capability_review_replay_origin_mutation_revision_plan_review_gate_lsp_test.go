package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSPReady(t *testing.T) {
    gate := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate{
        Status:              jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady,
        PlanStatus:          jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        Decision:            jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        PlanDigest:          "plan-digest",
        ReviewSource:        "review-source",
        ReviewEvidenceDigest: "review-evidence",
        NonExecuting:        true,
        NonAuthorizing:      true,
    }
    gate.GateDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate(
        gate.Status,
        gate.PlanStatus,
        gate.Decision,
        gate.PlanDigest,
        gate.ReviewSource,
        gate.ReviewEvidenceDigest,
    )

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSP(gate)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove ||
        diagnostic.ReviewSource != "review-source" || diagnostic.GateDigest == "" {
        t.Fatalf("ready diagnostic did not preserve review evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSPUnknown(t *testing.T) {
    gate := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate{
        Status:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady,
        PlanStatus:     jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        Decision:       jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        PlanDigest:     "plan-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateLSP(gate)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateUnknown ||
        diagnostic.MissingStage != "revision-plan-review-evidence" {
        t.Fatalf("incomplete review gate was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
