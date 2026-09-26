package decision

import "testing"

func replayFeedbackObservation() JEVImprovementReplayObservation {
    o := JEVImprovementReplayObservation{
        Status: jevImprovementReplayReproduced,
        CandidateDigest: "candidate-digest",
        PreparationEvidenceDigest: "preparation-evidence-digest",
        ObservationEvidenceDigest: "observation-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    o.EvidenceDigest = digestJEVImprovementReplayObservation(o.Status, o.CandidateDigest, o.PreparationEvidenceDigest, o.ObservationEvidenceDigest)
    return o
}

func TestDeriveJEVImprovementReplayFeedbackConfirmed(t *testing.T) {
    got := DeriveJEVImprovementReplayFeedback(JEVImprovementReplayFeedbackInput{
        ReplayObservation: replayFeedbackObservation(),
        MetricDigest: "metric-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevReplayFeedbackConfirmed || got.FeedbackKind != jevReplayFeedbackConfirmed || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestDeriveJEVImprovementReplayFeedbackCounterexample(t *testing.T) {
    observation := replayFeedbackObservation()
    observation.Status = jevImprovementReplayCounterexample
    observation.EvidenceDigest = digestJEVImprovementReplayObservation(observation.Status, observation.CandidateDigest, observation.PreparationEvidenceDigest, observation.ObservationEvidenceDigest)
    got := DeriveJEVImprovementReplayFeedback(JEVImprovementReplayFeedbackInput{
        ReplayObservation: observation,
        MetricDigest: "metric-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevReplayFeedbackRefuted || got.FeedbackKind != jevReplayFeedbackRefuted {
        t.Fatalf("got %+v", got)
    }
}

func TestDeriveJEVImprovementReplayFeedbackRequiresMetric(t *testing.T) {
    got := DeriveJEVImprovementReplayFeedback(JEVImprovementReplayFeedbackInput{
        ReplayObservation: replayFeedbackObservation(),
        NonAuthorizing: true,
    })
    if got.Status != jevReplayFeedbackUnknown || got.MissingStage != "metric" {
        t.Fatalf("got %+v", got)
    }
}
