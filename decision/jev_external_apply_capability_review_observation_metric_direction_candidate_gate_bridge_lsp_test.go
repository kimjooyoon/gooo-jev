package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSPReady(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge{
        Status:                          jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeBound,
        ObservationStatus:                jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded,
        ObservationMetricDirectionDigest: "observation-metric-direction-digest",
        CandidateGateStatus:               jevExternalApplyCapabilityReviewRevisionCandidateGated,
        CandidateDecision:                 jevExternalApplyCapabilityReviewRevisionCandidateReady,
        CandidateDigest:                   "candidate-digest",
        RevisionSource:                    "revision-source",
        GateDigest:                        "gate-digest",
        NonExecuting:                      true,
        NonAuthorizing:                    true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge(
        bridge.Status,
        bridge.ObservationStatus,
        bridge.ObservationMetricDirectionDigest,
        bridge.CandidateGateStatus,
        bridge.CandidateDecision,
        bridge.CandidateDigest,
        bridge.RevisionSource,
        bridge.GateDigest,
    )

    diagnostic := ProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid ready diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
        diagnostic.RevisionSource != "revision-source" || diagnostic.BridgeDigest == "" {
        t.Fatalf("ready diagnostic did not preserve candidate gate evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSPUnknown(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge{
        Status:                          jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeBound,
        ObservationStatus:                jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationObservationRecorded,
        ObservationMetricDirectionDigest: "observation-metric-direction-digest",
        CandidateGateStatus:               jevExternalApplyCapabilityReviewRevisionCandidateGated,
        CandidateDecision:                 jevExternalApplyCapabilityReviewRevisionCandidateReady,
        RevisionSource:                    "revision-source",
        GateDigest:                        "gate-digest",
        NonExecuting:                      true,
        NonAuthorizing:                    true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge(
        bridge.Status,
        bridge.ObservationStatus,
        bridge.ObservationMetricDirectionDigest,
        bridge.CandidateGateStatus,
        bridge.CandidateDecision,
        bridge.CandidateDigest,
        bridge.RevisionSource,
        bridge.GateDigest,
    )

    diagnostic := ProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
    if diagnostic.Publishable || diagnostic.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown ||
        diagnostic.MissingStage != "observation-metric-direction-candidate-gate-evidence" {
        t.Fatalf("incomplete candidate gate bridge was not preserved as UNKNOWN: %#v", diagnostic)
    }
}
