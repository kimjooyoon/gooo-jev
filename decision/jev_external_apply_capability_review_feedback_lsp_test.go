package decision

import "testing"

func validExternalApplyCapabilityReviewFeedbackForLSP() JEVExternalApplyCapabilityReviewFeedbackBridge {
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

func TestProjectJEVExternalApplyCapabilityReviewFeedbackLSPIsPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewFeedbackLSP(validExternalApplyCapabilityReviewFeedbackForLSP())
    if got.Severity != "info" || got.Code != "jev.external-apply.feedback-bound" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewFeedbackLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewFeedbackLSP(JEVExternalApplyCapabilityReviewFeedbackBridge{
        Status:         jevExternalApplyCapabilityReviewFeedbackUnknown,
        MissingStage:   "feedback-evidence",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
