package decision

import "testing"

func replayObservationPreparation() JEVImprovementReplayPreparation {
    p := JEVImprovementReplayPreparation{
        Status: jevImprovementReplayReady,
        CandidateDigest: "candidate-digest",
        ReviewEvidenceDigest: "review-evidence-digest",
        ReplayScope: "sandbox://candidate-one",
        ReplayEnvironment: "gooo://replay/environment/v1",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    p.EvidenceDigest = digestJEVImprovementReplayPreparation(p.Status, p.CandidateDigest, p.ReviewEvidenceDigest, p.ReplayScope, p.ReplayEnvironment)
    return p
}

func TestObserveJEVImprovementReplayReproduced(t *testing.T) {
    got := ObserveJEVImprovementReplay(JEVImprovementReplayObservationInput{
        Preparation: replayObservationPreparation(),
        ObservedCandidateDigest: "candidate-digest",
        Outcome: jevImprovementReplayReproduced,
        ObservationEvidenceDigest: "observation-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementReplayReproduced || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestObserveJEVImprovementReplayCounterexample(t *testing.T) {
    got := ObserveJEVImprovementReplay(JEVImprovementReplayObservationInput{
        Preparation: replayObservationPreparation(),
        ObservedCandidateDigest: "candidate-digest",
        Outcome: jevImprovementReplayCounterexample,
        ObservationEvidenceDigest: "counterexample-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementReplayCounterexample || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
}

func TestObserveJEVImprovementReplayRejectsCandidateMismatch(t *testing.T) {
    got := ObserveJEVImprovementReplay(JEVImprovementReplayObservationInput{
        Preparation: replayObservationPreparation(),
        ObservedCandidateDigest: "different-candidate",
        Outcome: jevImprovementReplayReproduced,
        ObservationEvidenceDigest: "observation-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementReplayObservationUnknown || got.MissingStage != "candidate-binding" {
        t.Fatalf("got %+v", got)
    }
}
