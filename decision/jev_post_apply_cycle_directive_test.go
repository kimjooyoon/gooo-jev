package decision

import "testing"

func postApplyCycleObservation(status string) JEVExternalApplyObservation {
    o := JEVExternalApplyObservation{
        Status: status,
        CandidateDigest: "candidate-digest",
        RequestEvidenceDigest: "request-evidence-digest",
        AppliedTargetReference: "gooo://external/apply/review/one",
        ApplyEvidenceDigest: "apply-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    o.EvidenceDigest = digestJEVExternalApplyObservation(o.Status, o.CandidateDigest, o.RequestEvidenceDigest, o.AppliedTargetReference, o.ApplyEvidenceDigest)
    return o
}

func TestDeriveJEVPostApplyCycleDirectiveApplied(t *testing.T) {
    got := DeriveJEVPostApplyCycleDirective(JEVPostApplyCycleDirectiveInput{
        Observation: postApplyCycleObservation(jevExternalApplyObservedApplied),
        NonAuthorizing: true,
    })
    if got.Status != jevPostApplyCycleReplayRequired || got.NextStage != "replay-observation" || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestDeriveJEVPostApplyCycleDirectiveRejected(t *testing.T) {
    got := DeriveJEVPostApplyCycleDirective(JEVPostApplyCycleDirectiveInput{
        Observation: postApplyCycleObservation(jevExternalApplyObservedRejected),
        NonAuthorizing: true,
    })
    if got.Status != jevPostApplyCycleRevisionRequired || got.NextStage != "revision-candidate" {
        t.Fatalf("got %+v", got)
    }
}

func TestDeriveJEVPostApplyCycleDirectiveUnknown(t *testing.T) {
    got := DeriveJEVPostApplyCycleDirective(JEVPostApplyCycleDirectiveInput{
        Observation: postApplyCycleObservation(jevExternalApplyObservedUnknown),
        NonAuthorizing: true,
    })
    if got.Status != jevPostApplyCycleEvidenceRequired || got.NextStage != "apply-evidence" {
        t.Fatalf("got %+v", got)
    }
}
