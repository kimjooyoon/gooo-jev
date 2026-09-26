package decision

import "testing"

func externalApplyObservationRequest() JEVExternalApplyRequest {
    r := JEVExternalApplyRequest{
        Status: jevExternalApplyRequestReady,
        CandidateDigest: "candidate-digest",
        ReviewEvidenceDigest: "review-evidence-digest",
        AggregationEvidenceDigest: "aggregation-evidence-digest",
        TargetReference: "gooo://external/apply/review/one",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    r.RequestEvidenceDigest = digestJEVExternalApplyRequestEvidence(r.CandidateDigest, r.ReviewEvidenceDigest, r.AggregationEvidenceDigest, r.TargetReference)
    r.EvidenceDigest = digestJEVExternalApplyRequest(r.Status, r.CandidateDigest, r.RequestEvidenceDigest)
    return r
}

func TestObserveJEVExternalApplyApplied(t *testing.T) {
    got := ObserveJEVExternalApply(JEVExternalApplyObservationInput{
        Request: externalApplyObservationRequest(),
        AppliedTargetReference: "gooo://external/apply/review/one",
        Outcome: jevExternalApplyObservedApplied,
        ApplyEvidenceDigest: "apply-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyObservedApplied || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestObserveJEVExternalApplyRejected(t *testing.T) {
    got := ObserveJEVExternalApply(JEVExternalApplyObservationInput{
        Request: externalApplyObservationRequest(),
        AppliedTargetReference: "gooo://external/apply/review/one",
        Outcome: jevExternalApplyObservedRejected,
        ApplyEvidenceDigest: "rejection-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyObservedRejected || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
}

func TestObserveJEVExternalApplyRejectsTargetMismatch(t *testing.T) {
    got := ObserveJEVExternalApply(JEVExternalApplyObservationInput{
        Request: externalApplyObservationRequest(),
        AppliedTargetReference: "gooo://external/apply/other",
        Outcome: jevExternalApplyObservedApplied,
        ApplyEvidenceDigest: "apply-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyObservedUnknown || got.MissingStage != "target-binding" {
        t.Fatalf("got %+v", got)
    }
}
