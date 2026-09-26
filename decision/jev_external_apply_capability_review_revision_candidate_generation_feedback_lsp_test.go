package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPGenerated(t *testing.T) {
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSP(feedback)
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid generation diagnostic, got %v", err)
    }
    if !diagnostic.Publishable || diagnostic.CandidateDigest != "generated-candidate-digest" ||
        diagnostic.GenerationEvidenceDigest != "generation-evidence" {
        t.Fatalf("generation diagnostic did not preserve evidence: %#v", diagnostic)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSPUnknown(t *testing.T) {
    diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback{},
    )
    if diagnostic.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown ||
        diagnostic.Publishable || diagnostic.MissingStage == "" {
        t.Fatalf("invalid generation feedback was published: %#v", diagnostic)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid UNKNOWN diagnostic, got %v", err)
    }
}
