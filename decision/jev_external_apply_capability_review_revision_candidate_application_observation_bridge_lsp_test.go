package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSPRecorded(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge{
        Status:                     jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeBound,
        CandidateApplicationStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        CandidateApplicationDigest: "candidate-application-digest",
        ObservationStatus:          jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationSource:           "reverse-observation-source",
        ObservationEvidenceDigest:  "observation-evidence",
        ReverseObservationDigest:   "reverse-observation",
        NonExecuting:               true,
        NonAuthorizing:             true,
    }
    bridge.ObservationDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(
        bridge.Status,
        bridge.CandidateApplicationStatus,
        bridge.CandidateApplicationDigest,
        bridge.ObservationStatus,
        bridge.ObservationSource,
        bridge.ObservationEvidenceDigest,
        bridge.ReverseObservationDigest,
    )
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid recorded diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.ObservationDigest != bridge.ObservationDigest {
        t.Fatalf("recorded diagnostic did not preserve observation evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid observation bridge was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
