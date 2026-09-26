package decision

import "testing"

func originResolutionObservation() JEVImprovementCycleObservation {
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

func originResolutionSelection(observation JEVImprovementCycleObservation) JEVImprovementCandidateSelection {
    s := JEVImprovementCandidateSelection{
        Status: jevImprovementCandidateReady,
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        ObservationEvidenceDigest: observation.EvidenceDigest,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    s.EvidenceDigest = digestJEVImprovementCandidateSelection(s.Status, s.CandidateDigest, s.CandidateSource, s.ObservationEvidenceDigest)
    return s
}

func TestResolveJEVImprovementOrigin(t *testing.T) {
    observation := originResolutionObservation()
    got := ResolveJEVImprovementOrigin(JEVImprovementOriginResolutionInput{
        Observation: observation,
        Selection: originResolutionSelection(observation),
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementOriginResolved || got.MissingStage != "" || got.OriginDigest == "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestResolveJEVImprovementOriginRejectsUnboundSelection(t *testing.T) {
    observation := originResolutionObservation()
    selection := originResolutionSelection(observation)
    selection.ObservationEvidenceDigest = "different-observation-evidence"
    selection.EvidenceDigest = digestJEVImprovementCandidateSelection(selection.Status, selection.CandidateDigest, selection.CandidateSource, selection.ObservationEvidenceDigest)
    got := ResolveJEVImprovementOrigin(JEVImprovementOriginResolutionInput{
        Observation: observation,
        Selection: selection,
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementOriginUnknown || got.MissingStage != "observation-selection-binding" {
        t.Fatalf("got %+v", got)
    }
}
