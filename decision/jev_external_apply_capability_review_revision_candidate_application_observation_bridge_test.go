package decision

import "testing"

func candidateObservationBridgeApplication(t *testing.T) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge {
    t.Helper()
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge{
        Status:                    jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound,
        CandidateGateStatus:       jevExternalApplyCapabilityReviewRevisionCandidateGated,
        CandidateDecision:         jevExternalApplyCapabilityReviewRevisionCandidateReady,
        CandidateDigest:           "candidate-digest",
        RevisionSource:            "revision-source",
        CandidateGateDigest:        "gate-digest",
        ApplicationCandidateStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        ApplicationTarget:         "application-target",
        ApplicationSource:          "application-source",
        ApplicationEvidenceDigest: "application-evidence",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(
        bridge.Status,
        bridge.CandidateGateStatus,
        bridge.CandidateDecision,
        bridge.CandidateDigest,
        bridge.RevisionSource,
        bridge.CandidateGateDigest,
        bridge.ApplicationCandidateStatus,
        bridge.ApplicationTarget,
        bridge.ApplicationSource,
        bridge.ApplicationEvidenceDigest,
    )
    return bridge
}

func TestObserveJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeRecorded(t *testing.T) {
    bridge := ObserveJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeInput{
        CandidateApplication:      candidateObservationBridgeApplication(t),
        ObservationStatus:         jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationSource:         "reverse-observation-source",
        ObservationEvidenceDigest: "observation-evidence",
        ReverseObservationDigest:  "reverse-observation",
        NonAuthorizing:            true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid observation bridge, got %v", err)
    }
    if bridge.CandidateApplicationDigest == "" || bridge.ObservationDigest == "" ||
        bridge.ObservationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved {
        t.Fatalf("observation bridge did not preserve origin evidence: %#v", bridge)
    }
}

func TestObserveJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeRequiresReverseObservation(t *testing.T) {
    bridge := ObserveJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeInput{
        CandidateApplication: candidateObservationBridgeApplication(t),
        ObservationStatus:    jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationSource:     "reverse-observation-source",
        ObservationEvidenceDigest: "observation-evidence",
        NonAuthorizing:        true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationBridgeUnknown ||
        bridge.MissingStage != "reverse-observation" {
        t.Fatalf("missing reverse observation was not preserved: %#v", bridge)
    }
    if bridge.ObservationDigest != "" {
        t.Fatalf("incomplete observation bridge unexpectedly has digest: %#v", bridge)
    }
}
