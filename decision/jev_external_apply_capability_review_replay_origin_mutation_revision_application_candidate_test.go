package decision

import "testing"

func TestGenerateJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady(t *testing.T) {
    gate := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate{
        Status:               jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady,
        PlanStatus:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        Decision:             jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        PlanDigest:           "plan-digest",
        ReviewSource:         "review-source",
        ReviewEvidenceDigest: "review-evidence",
        NonExecuting:         true,
        NonAuthorizing:       true,
    }
    gate.GateDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate(
        gate.Status,
        gate.PlanStatus,
        gate.Decision,
        gate.PlanDigest,
        gate.ReviewSource,
        gate.ReviewEvidenceDigest,
    )

    candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateInput{
        ReviewGate:                gate,
        ApplicationTarget:         "target",
        ApplicationSource:         "application-source",
        ApplicationEvidenceDigest: "application-evidence",
        NonAuthorizing:            true,
    })
    if err := candidate.Validate(); err != nil {
        t.Fatalf("expected valid application candidate, got %v", err)
    }
    if candidate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady ||
        candidate.ApplicationTarget != "target" || candidate.CandidateDigest == "" {
        t.Fatalf("ready application candidate did not preserve evidence: %#v", candidate)
    }
}

func TestGenerateJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknownWithoutSource(t *testing.T) {
    gate := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGate{
        Status:               jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady,
        PlanStatus:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady,
        Decision:             jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove,
        PlanDigest:           "plan-digest",
        GateDigest:           "gate-digest",
        ReviewSource:         "review-source",
        ReviewEvidenceDigest: "review-evidence",
        NonExecuting:         true,
        NonAuthorizing:       true,
    }

    candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateInput{
        ReviewGate:      gate,
        ApplicationTarget: "target",
        NonAuthorizing:  true,
    })
    if candidate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown ||
        candidate.MissingStage != "application-source" ||
        !candidate.NonExecuting || !candidate.NonAuthorizing {
        t.Fatalf("missing application source was not preserved as UNKNOWN: %#v", candidate)
    }
}
