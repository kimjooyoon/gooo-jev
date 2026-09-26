package decision

import "testing"

func candidateGenerationFeedbackDirection() JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge {
    bridge := JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge{
        Status:                    jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound,
        CandidateApplicationStatus: jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady,
        CandidateApplicationDigest: "candidate-application-digest",
        ObservationStatus:          jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationObserved,
        ObservationDigest:           "observation-digest",
        MetricName:                  "reverse-confidence",
        MetricValue:                 0.92,
        MetricDigest:                "metric-digest",
        Direction:                   jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:                      "revision-target",
        CandidateSource:             "candidate-source",
        FeedbackDigest:              "feedback-digest",
        NonExecuting:               true,
        NonAuthorizing:             true,
    }
    bridge.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(
        jevExternalApplyCapabilityReviewImprovementDirectionBound,
        bridge.Direction,
        bridge.Target,
        bridge.CandidateSource,
        bridge.FeedbackDigest,
    )
    bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge(
        bridge.Status,
        bridge.CandidateApplicationStatus,
        bridge.CandidateApplicationDigest,
        bridge.ObservationStatus,
        bridge.ObservationDigest,
        bridge.MetricName,
        bridge.MetricValue,
        bridge.MetricDigest,
        bridge.Direction,
        bridge.Target,
        bridge.CandidateSource,
        bridge.FeedbackDigest,
        bridge.DirectionDigest,
    )
    return bridge
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated(t *testing.T) {
    feedback := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackInput{
        Direction:                 candidateGenerationFeedbackDirection(),
        CandidateDigest:           "generated-candidate-digest",
        GenerationSource:          "candidate-generator",
        GenerationEvidenceDigest: "generation-evidence",
        NonAuthorizing:            true,
    })
    if err := feedback.Validate(); err != nil {
        t.Fatalf("expected valid generated feedback, got %v", err)
    }
    if feedback.GeneratedCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackGenerated ||
        feedback.CandidateDigest != "generated-candidate-digest" ||
        feedback.GenerationEvidenceDigest != "generation-evidence" {
        t.Fatalf("generated candidate evidence was not preserved: %#v", feedback)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackRequiresCandidateEvidence(t *testing.T) {
    feedback := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedback(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackInput{
        Direction:      candidateGenerationFeedbackDirection(),
        NonAuthorizing: true,
    })
    if feedback.Status != jevExternalApplyCapabilityReviewRevisionCandidateGenerationFeedbackUnknown ||
        feedback.MissingStage != "candidate-digest" {
        t.Fatalf("missing candidate evidence was not preserved: %#v", feedback)
    }
    if feedback.BridgeDigest != "" {
        t.Fatalf("incomplete generation feedback unexpectedly has digest: %#v", feedback)
    }
}
