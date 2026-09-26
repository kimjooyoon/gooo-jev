package decision

import "testing"

func generationCandidateGateFeedback() JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback {
    feedback := JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback{
        Status:                   jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackBound,
        DirectionStatus:          jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound,
        Direction:                jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:                   "revision-target",
        FeedbackDigest:           "feedback-digest",
        DirectionDigest:          "direction-digest",
        GeneratedCandidateStatus: jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated,
        CandidateDigest:           "generated-candidate-digest",
        GenerationSource:          "candidate-generator",
        GenerationEvidenceDigest: "generation-evidence",
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    feedback.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(
        feedback.Status,
        feedback.DirectionStatus,
        feedback.Direction,
        feedback.Target,
        feedback.FeedbackDigest,
        feedback.DirectionDigest,
        feedback.GeneratedCandidateStatus,
        feedback.CandidateDigest,
        feedback.GenerationSource,
        feedback.GenerationEvidenceDigest,
    )
    return feedback
}

func generationCandidateGate() JEVExternalApplyCapabilityReviewRevisionCandidateGate {
    gate := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:          jevExternalApplyCapabilityReviewRevisionCandidateGated,
        Decision:        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "revision-target",
        CandidateSource: "candidate-generator",
        CandidateDigest: "generated-candidate-digest",
        DirectionDigest: "direction-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
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

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeConsistent(t *testing.T) {
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeInput{
        GenerationFeedback: generationCandidateGateFeedback(),
        CandidateGate:      generationCandidateGate(),
        NonAuthorizing:     true,
    })
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid generation gate bridge, got %v", err)
    }
    if bridge.CandidateDigest != "generated-candidate-digest" ||
        bridge.GenerationSource != "candidate-generator" ||
        bridge.CandidateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady {
        t.Fatalf("generation gate bridge did not preserve consistency evidence: %#v", bridge)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeRejectsMismatch(t *testing.T) {
    gate := generationCandidateGate()
    gate.CandidateDigest = "different-candidate-digest"
    gate.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(
        gate.Status,
        gate.Decision,
        gate.Direction,
        gate.Target,
        gate.CandidateSource,
        gate.CandidateDigest,
        gate.DirectionDigest,
    )
    bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeInput{
        GenerationFeedback: generationCandidateGateFeedback(),
        CandidateGate:      gate,
        NonAuthorizing:     true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackCandidateGateBridgeUnknown ||
        bridge.MissingStage != "candidate-digest-consistency" {
        t.Fatalf("candidate digest mismatch was not preserved: %#v", bridge)
    }
}
