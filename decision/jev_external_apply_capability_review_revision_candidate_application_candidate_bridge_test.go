package decision

import "testing"

func candidateApplicationBridgeGate(decision string) JEVExternalApplyCapabilityReviewRevisionCandidateGate {
    gate := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:          jevExternalApplyCapabilityReviewRevisionCandidateGated,
        Decision:        decision,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "revision-source",
        CandidateDigest: "candidate-digest",
        DirectionDigest: "direction-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    if decision != jevExternalApplyCapabilityReviewRevisionCandidateReady {
        gate.CandidateDigest = ""
    }
    gate.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(
        gate.Status,
        gate.Decision,
        gate.Direction,
        gate.Target,
        gate.CandidateSource,
        gate.CandidateDigest,
        gate.DirectionDigest,
    )
    return gate
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeReady(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeInput{
        CandidateGate:             candidateApplicationBridgeGate(jevExternalApplyCapabilityReviewRevisionCandidateReady),
        ApplicationTarget:         "application-target",
        ApplicationSource:         "application-source",
        ApplicationEvidenceDigest: "application-evidence",
        NonAuthorizing:            true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid ready bridge, got %v", err)
    }
    if bridge.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady ||
        bridge.ApplicationTarget != "application-target" ||
        bridge.CandidateDigest != "candidate-digest" {
        t.Fatalf("ready bridge did not preserve candidate and application evidence: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeHold(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeInput{
        CandidateGate:     candidateApplicationBridgeGate(jevExternalApplyCapabilityReviewRevisionCandidateHold),
        NonAuthorizing:    true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid hold bridge, got %v", err)
    }
    if bridge.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold ||
        bridge.ApplicationTarget != "" || bridge.ApplicationSource != "" {
        t.Fatalf("hold bridge exposed application evidence: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeRequiresEvidence(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeInput{
        CandidateGate:     candidateApplicationBridgeGate(jevExternalApplyCapabilityReviewRevisionCandidateReady),
        NonAuthorizing:    true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound ||
        bridge.MissingStage != "application-target" {
        t.Fatalf("missing application evidence was not preserved: %#v", bridge)
    }
    if bridge.BridgeDigest != "" {
        t.Fatalf("incomplete bridge unexpectedly has digest: %#v", bridge)
    }
}
