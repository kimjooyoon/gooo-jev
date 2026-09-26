package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSPRecorded(t *testing.T) {
    observation := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation{
        Status:                    jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded,
        CandidateStatus:           jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady,
        CandidateDigest:           "candidate-digest",
        ObservationStatus:         jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved,
        ObservationSource:          "observation-source",
        ObservationEvidenceDigest: "observation-evidence",
        ReverseObservationDigest:  "reverse-observation",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    observation.ObservationDigest = digestJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation(
        observation.Status,
        observation.CandidateStatus,
        observation.CandidateDigest,
        observation.ObservationStatus,
        observation.ObservationSource,
        observation.ObservationEvidenceDigest,
        observation.ReverseObservationDigest,
    )

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSP(observation)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid recorded diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.ReverseObservationDigest != "reverse-observation" ||
        diagnostic.ObservationDigest == "" {
        t.Fatalf("recorded diagnostic did not preserve reverse evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSPUnknown(t *testing.T) {
    observation := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservation{
        Status:            jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded,
        CandidateStatus:   jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady,
        CandidateDigest:   "candidate-digest",
        ObservationStatus: jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationObserved,
        ObservationSource: "observation-source",
        NonExecuting:      true,
        NonAuthorizing:    true,
    }

    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationLSP(observation)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationUnknown ||
        diagnostic.MissingStage != "application-observation-evidence" {
        t.Fatalf("incomplete observation was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
