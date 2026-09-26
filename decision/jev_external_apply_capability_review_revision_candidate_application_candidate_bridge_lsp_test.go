package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSPReady(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge{
        Status:                    jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound,
        CandidateGateStatus:       jevExternalApplyCapabilityReviewRevisionCandidateGated,
        CandidateDecision:         jevExternalApplyCapabilityReviewRevisionCandidateReady,
        CandidateDigest:           "candidate-digest",
        RevisionSource:            "revision-source",
        CandidateGateDigest:       "gate-digest",
        ApplicationCandidateStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        ApplicationTarget:         "application-target",
        ApplicationSource:         "application-source",
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.BridgeDigest != bridge.BridgeDigest {
        t.Fatalf("ready diagnostic did not preserve bridge evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid bridge was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
