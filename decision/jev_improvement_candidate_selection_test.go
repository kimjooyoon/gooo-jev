package decision

import "testing"

func candidateSelectionObservation() JEVImprovementCycleObservation {
    o := JEVImprovementCycleObservation{
        Status: changePlanFeedbackDispositionReady,
        DeclarationDigest: "declaration-digest",
        IRDigest: "ir-digest",
        GenerationDigest: "generation-digest",
        ReverseObservationDigest: "reverse-observation-digest",
        MetricDigest: "metric-digest",
        ChangePlanDigest: "plan-digest",
        DispositionStatus: changePlanFeedbackDispositionReady,
        LSPCode: "jev.change-plan.ready-for-external-apply-review",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    o.EvidenceDigest = digestJEVImprovementCycleObservation(o.DeclarationDigest, o.IRDigest, o.GenerationDigest, o.ReverseObservationDigest, o.MetricDigest, o.ChangePlanDigest, o.DispositionStatus, o.LSPCode)
    return o
}

func TestSelectJEVImprovementCandidate(t *testing.T) {
    observation := candidateSelectionObservation()
    got := SelectJEVImprovementCandidate(JEVImprovementCandidateSelectionInput{
        Observation: observation,
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementCandidateReady || got.MissingStage != "" || got.EvidenceDigest == "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestSelectJEVImprovementCandidateRequiresCandidateDigest(t *testing.T) {
    observation := candidateSelectionObservation()
    got := SelectJEVImprovementCandidate(JEVImprovementCandidateSelectionInput{
        Observation: observation,
        CandidateSource: "gooo://candidate/one",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementCandidateUnknown || got.MissingStage != "candidate-digest" {
        t.Fatalf("got %+v", got)
    }
}
