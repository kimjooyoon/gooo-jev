package decision

import "testing"

func validExternalApplyCapabilityReviewTrendSignalForFeedback() JEVExternalApplyCapabilityReviewTrendSignal {
    signal := JEVExternalApplyCapabilityReviewTrendSignal{
        Status:         jevExternalApplyCapabilityReviewTrendSignalRecorded,
        Signal:         jevExternalApplyCapabilityReviewTrendSignalImprove,
        MetricName:     "review.confidence",
        TrendDigest:    "trend-digest",
        MeanValue:      0.75,
        Threshold:      0.7,
        EvidenceDigest: "trend-evidence-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    signal.SignalDigest = digestJEVExternalApplyCapabilityReviewTrendSignal(signal.Status, signal.Signal, signal.MetricName, signal.TrendDigest, signal.MeanValue, signal.Threshold, signal.EvidenceDigest)
    return signal
}

func TestBindJEVExternalApplyCapabilityReviewFeedbackPreservesSignal(t *testing.T) {
    got := BindJEVExternalApplyCapabilityReviewFeedback(JEVExternalApplyCapabilityReviewFeedbackBridgeInput{
        Signal:                 validExternalApplyCapabilityReviewTrendSignalForFeedback(),
        FeedbackEvidenceDigest: "feedback-evidence-digest",
        NonAuthorizing:         true,
    })
    if got.Status != jevExternalApplyCapabilityReviewFeedbackBound || got.Signal != jevExternalApplyCapabilityReviewTrendSignalImprove || !got.NonExecuting || !got.NonAuthorizing {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestBindJEVExternalApplyCapabilityReviewFeedbackRejectsMissingEvidence(t *testing.T) {
    got := BindJEVExternalApplyCapabilityReviewFeedback(JEVExternalApplyCapabilityReviewFeedbackBridgeInput{
        Signal:         validExternalApplyCapabilityReviewTrendSignalForFeedback(),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewFeedbackUnknown || got.MissingStage != "feedback-evidence" {
        t.Fatalf("got %+v", got)
    }
}

func TestBindJEVExternalApplyCapabilityReviewFeedbackPreservesUnknownBoundary(t *testing.T) {
    got := BindJEVExternalApplyCapabilityReviewFeedback(JEVExternalApplyCapabilityReviewFeedbackBridgeInput{
        Signal: JEVExternalApplyCapabilityReviewTrendSignal{
            Status:         jevExternalApplyCapabilityReviewTrendSignalUnknown,
            MissingStage:   "signal-threshold",
            NonExecuting:   true,
            NonAuthorizing: true,
        },
        FeedbackEvidenceDigest: "feedback-evidence-digest",
        NonAuthorizing:         true,
    })
    if got.Status != jevExternalApplyCapabilityReviewFeedbackUnknown || got.MissingStage != "trend-signal" {
        t.Fatalf("got %+v", got)
    }
}
