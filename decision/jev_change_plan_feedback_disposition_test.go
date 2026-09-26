package decision

import "testing"

func TestDecisionConfidenceChangePlanFeedbackDispositionRejectsAuthorization(t *testing.T) {
    got := DispositionDecisionConfidenceChangePlanWithFeedback(DecisionConfidenceChangePlanFeedbackDispositionInput{})
    if got.Status != changePlanFeedbackDispositionUnknown || got.MissingStage != "authorization-boundary" {
        t.Fatalf("got %+v", got)
    }
}

func TestDecisionConfidenceChangePlanFeedbackDispositionRejectsIncompleteEvidence(t *testing.T) {
    input := DecisionConfidenceChangePlanFeedbackDispositionInput{
        NonAuthorizing: true,
        FeedbackBinding: DecisionConfidenceChangePlanReplayFeedbackBinding{
            NonAuthorizing: true,
            NonExecuting: true,
            FeedbackStatus: changePlanFeedbackStatusConfirmed,
        },
    }
    got := DispositionDecisionConfidenceChangePlanWithFeedback(input)
    if got.Status != changePlanFeedbackDispositionUnknown || got.MissingStage != "change-plan" {
        t.Fatalf("got %+v", got)
    }
}
