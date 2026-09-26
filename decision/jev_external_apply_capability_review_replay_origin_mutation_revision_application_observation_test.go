package decision

import "testing"

func TestObserveJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationRecorded(t *testing.T) {
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

    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplication(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationInput{
        Candidate:                 candidate,
        ObservationStatus:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved,
        ObservationSource:         "runtime-observation-source",
        ObservationEvidenceDigest: "runtime-observation-evidence",
        ReverseObservationDigest:  "reverse-observation",
        NonAuthorizing:            true,
    })
    if err := observation.Validate(); err != nil {
        t.Fatalf("expected valid application observation, got %v", err)
    }
    if observation.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded ||
        observation.ReverseObservationDigest != "reverse-observation" ||
        observation.ObservationDigest == "" {
        t.Fatalf("recorded observation did not preserve evidence: %#v", observation)
    }
}

func TestObserveJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationUnknownWithoutReverseObservation(t *testing.T) {
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

    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplication(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationInput{
        Candidate:                 candidate,
        ObservationStatus:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved,
        ObservationSource:         "runtime-observation-source",
        ObservationEvidenceDigest: "runtime-observation-evidence",
        NonAuthorizing:            true,
    })
    if observation.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown ||
        observation.MissingStage != "reverse-observation" ||
        !observation.NonExecuting || !observation.NonAuthorizing {
        t.Fatalf("missing reverse observation was not preserved as UNKNOWN: %#v", observation)
    }
}
