package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSPBound(t *testing.T) {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge{
        Status:                   jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeBound,
        GenerationFeedbackStatus: jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound,
        GeneratedCandidateStatus: jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated,
        CandidateDigest:           "generated-candidate-digest",
        GenerationSource:          "candidate-generator",
        GenerationEvidenceDigest: "generation-evidence",
        CandidateGateStatus:       jevExternalApplyCapabilityReviewRevisionCandidateGated,
        CandidateDecision:         jevExternalApplyCapabilityReviewRevisionCandidateReady,
        CandidateGateDigest:       "candidate-gate-digest",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(
        bridge.Status,
        bridge.GenerationFeedbackStatus,
        bridge.GeneratedCandidateStatus,
        bridge.CandidateDigest,
        bridge.GenerationSource,
        bridge.GenerationEvidenceDigest,
        bridge.CandidateGateStatus,
        bridge.CandidateDecision,
        bridge.CandidateGateDigest,
    )
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSP(bridge)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid consistency diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.CandidateDigest != "generated-candidate-digest" ||
        diagnostic.CandidateGateDigest != "candidate-gate-digest" {
        t.Fatalf("consistency diagnostic did not preserve evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid consistency bridge was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
