package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSPReady(t *testing.T) {
    candidate := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady,
        ReviewGateStatus:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady,
        Decision:                  jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        PlanDigest:                "plan-digest",
        GateDigest:                "gate-digest",
        ApplicationTarget:         "target",
        ApplicationSource:         "application-source",
        ApplicationEvidenceDigest: "application-evidence",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    candidate.CandidateDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(
        candidate.Status,
        candidate.ReviewGateStatus,
        candidate.Decision,
        candidate.PlanDigest,
        candidate.GateDigest,
        candidate.ApplicationTarget,
        candidate.ApplicationSource,
        candidate.ApplicationEvidenceDigest,
    )

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSP(candidate)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.ApplicationTarget != "target" ||
        diagnostic.ApplicationSource != "application-source" || diagnostic.CandidateDigest == "" {
        t.Fatalf("ready diagnostic did not preserve application evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSPUnknown(t *testing.T) {
    candidate := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate{
        Status:             jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady,
        ReviewGateStatus:   jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady,
        Decision:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        PlanDigest:         "plan-digest",
        GateDigest:         "gate-digest",
        NonExecuting:       true,
        NonAuthorizing:     true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSP(candidate)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown ||
        diagnostic.MissingStage != "application-candidate-evidence" {
        t.Fatalf("incomplete application candidate was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
