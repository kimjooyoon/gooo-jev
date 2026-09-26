package decision

import "testing"

func originBindingOrigin() JEVImprovementOriginResolution {
    o := JEVImprovementOriginResolution{
        Status: jevImprovementOriginResolved,
        DeclarationDigest: "declaration-digest",
        IRDigest: "ir-digest",
        GenerationDigest: "generation-digest",
        ReverseObservationDigest: "reverse-observation-digest",
        MetricDigest: "metric-digest",
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        ObservationEvidenceDigest: "observation-evidence-digest",
        SelectionEvidenceDigest: "selection-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    o.OriginDigest = digestJEVImprovementOrigin(o.DeclarationDigest, o.IRDigest, o.GenerationDigest, o.ReverseObservationDigest, o.MetricDigest, o.CandidateDigest, o.CandidateSource)
    o.EvidenceDigest = digestJEVImprovementOriginEvidence(o.Status, o.OriginDigest, o.ObservationEvidenceDigest, o.SelectionEvidenceDigest)
    return o
}

func originBindingRequest(candidateDigest string) JEVExternalApplyRequest {
    r := JEVExternalApplyRequest{
        Status: jevExternalApplyRequestReady,
        CandidateDigest: candidateDigest,
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

func TestBindJEVExternalApplyToOrigin(t *testing.T) {
    got := BindJEVExternalApplyToOrigin(JEVExternalApplyOriginBindingInput{
        Origin: originBindingOrigin(),
        Request: originBindingRequest("candidate-digest"),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyOriginBindingReady || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestBindJEVExternalApplyToOriginRejectsCandidateMismatch(t *testing.T) {
    got := BindJEVExternalApplyToOrigin(JEVExternalApplyOriginBindingInput{
        Origin: originBindingOrigin(),
        Request: originBindingRequest("different-candidate"),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyOriginBindingUnknown || got.MissingStage != "origin-candidate-binding" {
        t.Fatalf("got %+v", got)
    }
}
