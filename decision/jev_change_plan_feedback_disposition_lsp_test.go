package decision

import "testing"

func validFeedbackDispositionForLSP(status string) DecisionConfidenceChangePlanFeedbackDisposition {
    d := DecisionConfidenceChangePlanFeedbackDisposition{
        Status: status,
        FeedbackStatus: changePlanFeedbackStatusConfirmed,
        ChangePlanDigest: "plan-digest",
        VerificationDigest: "verification-digest",
        DispositionDigest: "disposition-digest",
        FeedbackEvidenceDigest: "feedback-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    d.EvidenceDigest = digestChangePlanFeedbackDisposition(d.ChangePlanDigest, d.VerificationDigest, d.DispositionDigest, d.FeedbackEvidenceDigest, d.FeedbackStatus)
    return d
}

func TestProjectDecisionConfidenceChangePlanFeedbackDispositionLSPReadyIsReviewOnly(t *testing.T) {
    got := ProjectDecisionConfidenceChangePlanFeedbackDispositionLSP(validFeedbackDispositionForLSP(changePlanFeedbackDispositionReady))
    if got.Severity != "info" || got.Code != "jev.change-plan.ready-for-external-apply-review" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectDecisionConfidenceChangePlanFeedbackDispositionLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectDecisionConfidenceChangePlanFeedbackDispositionLSP(DecisionConfidenceChangePlanFeedbackDisposition{
        Status: changePlanFeedbackDispositionUnknown,
        MissingStage: "feedback-binding",
        NonAuthorizing: true,
        NonExecuting: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
