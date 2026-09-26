package decision

import "testing"

func validExternalApplyCapabilityReviewFeedbackForDirection() JEVExternalApplyCapabilityReviewFeedbackBridge {
    bridge := JEVExternalApplyCapabilityReviewFeedbackBridge{
        Status:                jevExternalApplyCapabilityReviewFeedbackBound,
        Signal:                jevExternalApplyCapabilityReviewTrendSignalImprove,
        SignalDigest:          "signal-digest",
        TrendDigest:           "trend-digest",
        FeedbackEvidenceDigest: "feedback-evidence-digest",
        NonExecuting:          true,
        NonAuthorizing:        true,
    }
    bridge.FeedbackDigest = digestJEVExternalApplyCapabilityReviewFeedback(bridge.Signal, bridge.SignalDigest, bridge.TrendDigest, bridge.FeedbackEvidenceDigest)
    return bridge
}

func TestDeriveJEVExternalApplyCapabilityReviewImprovementDirectionMapsSignals(t *testing.T) {
    cases := []struct {
        signal    string
        direction string
    }{
        {jevExternalApplyCapabilityReviewTrendSignalImprove, jevExternalApplyCapabilityReviewImprovementDirectionGenerate},
        {jevExternalApplyCapabilityReviewTrendSignalHold, jevExternalApplyCapabilityReviewImprovementDirectionHold},
        {jevExternalApplyCapabilityReviewTrendSignalRollback, jevExternalApplyCapabilityReviewImprovementDirectionReject},
    }
    for _, testCase := range cases {
        feedback := validExternalApplyCapabilityReviewFeedbackForDirection()
        feedback.Signal = testCase.signal
        feedback.FeedbackDigest = digestJEVExternalApplyCapabilityReviewFeedback(feedback.Signal, feedback.SignalDigest, feedback.TrendDigest, feedback.FeedbackEvidenceDigest)
        got := DeriveJEVExternalApplyCapabilityReviewImprovementDirection(JEVExternalApplyCapabilityReviewImprovementDirectionInput{
            Feedback:        feedback,
            Target:          "decision/provenance",
            CandidateSource: "review-metric-trend",
            NonAuthorizing:  true,
        })
        if got.Direction != testCase.direction || !got.NonExecuting || !got.NonAuthorizing {
            t.Fatalf("signal %q got %+v", testCase.signal, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestDeriveJEVExternalApplyCapabilityReviewImprovementDirectionRejectsMissingTarget(t *testing.T) {
    got := DeriveJEVExternalApplyCapabilityReviewImprovementDirection(JEVExternalApplyCapabilityReviewImprovementDirectionInput{
        Feedback:        validExternalApplyCapabilityReviewFeedbackForDirection(),
        CandidateSource: "review-metric-trend",
        NonAuthorizing:  true,
    })
    if got.Status != jevExternalApplyCapabilityReviewImprovementDirectionUnknown || got.MissingStage != "improvement-target" {
        t.Fatalf("got %+v", got)
    }
}

func TestDeriveJEVExternalApplyCapabilityReviewImprovementDirectionPreservesUnknownFeedback(t *testing.T) {
    got := DeriveJEVExternalApplyCapabilityReviewImprovementDirection(JEVExternalApplyCapabilityReviewImprovementDirectionInput{
        Feedback: JEVExternalApplyCapabilityReviewFeedbackBridge{
            Status:         jevExternalApplyCapabilityReviewFeedbackUnknown,
            MissingStage:   "feedback-evidence",
            NonExecuting:   true,
            NonAuthorizing: true,
        },
        Target:          "decision/provenance",
        CandidateSource: "review-metric-trend",
        NonAuthorizing:  true,
    })
    if got.Status != jevExternalApplyCapabilityReviewImprovementDirectionUnknown || got.MissingStage != "feedback-bridge" {
        t.Fatalf("got %+v", got)
    }
}
