package decision

import "testing"

func aggregateReplayFeedback(kind string) JEVImprovementReplayFeedback {
    f := JEVImprovementReplayFeedback{
        Status: kind,
        FeedbackKind: kind,
        CandidateDigest: "candidate-" + kind,
        ReplayObservationDigest: "observation-" + kind,
        MetricDigest: "metric-" + kind,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    f.FeedbackEvidenceDigest = digestJEVImprovementReplayFeedbackEvidence(f.ReplayObservationDigest, f.MetricDigest, f.FeedbackKind)
    f.EvidenceDigest = digestJEVImprovementReplayFeedback(f.Status, f.FeedbackKind, f.CandidateDigest, f.ReplayObservationDigest, f.MetricDigest, f.FeedbackEvidenceDigest)
    return f
}

func TestAggregateJEVImprovementFeedbackStable(t *testing.T) {
    got := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
        Feedback: []JEVImprovementReplayFeedback{
            aggregateReplayFeedback(jevReplayFeedbackConfirmed),
            aggregateReplayFeedback(jevReplayFeedbackConfirmed),
        },
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementFeedbackStableForReview || got.Total != 2 || got.Confirmed != 2 || got.Refuted != 0 || got.Unknown != 0 {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestAggregateJEVImprovementFeedbackRefutedNeedsRevision(t *testing.T) {
    got := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
        Feedback: []JEVImprovementReplayFeedback{
            aggregateReplayFeedback(jevReplayFeedbackConfirmed),
            aggregateReplayFeedback(jevReplayFeedbackRefuted),
        },
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementFeedbackNeedsRevision || got.Refuted != 1 {
        t.Fatalf("got %+v", got)
    }
}

func TestAggregateJEVImprovementFeedbackUnknownHolds(t *testing.T) {
    got := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
        Feedback: []JEVImprovementReplayFeedback{
            aggregateReplayFeedback(jevReplayFeedbackUnknown),
        },
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementFeedbackHold || got.Unknown != 1 {
        t.Fatalf("got %+v", got)
    }
}

func TestAggregateJEVImprovementFeedbackRequiresHistory(t *testing.T) {
    got := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{NonAuthorizing: true})
    if got.Status != jevImprovementFeedbackAggregateUnknown || got.MissingStage != "feedback-history" {
        t.Fatalf("got %+v", got)
    }
}
