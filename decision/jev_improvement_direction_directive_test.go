package decision

import "testing"

func directionAggregation(status string) JEVImprovementFeedbackAggregation {
    a := JEVImprovementFeedbackAggregation{
        Status: status,
        Total: 1,
        InputEvidenceDigest: "input-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    switch status {
    case jevImprovementFeedbackStableForReview:
        a.Confirmed = 1
    case jevImprovementFeedbackNeedsRevision:
        a.Refuted = 1
    case jevImprovementFeedbackHold:
        a.Unknown = 1
    }
    a.EvidenceDigest = digestJEVImprovementFeedbackAggregation(a.Status, a.Total, a.Confirmed, a.Refuted, a.Unknown, a.InputEvidenceDigest)
    return a
}

func TestDeriveJEVImprovementDirectionDirectiveStable(t *testing.T) {
    got := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
        Aggregation: directionAggregation(jevImprovementFeedbackStableForReview),
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        NonAuthorizing: true,
    })
    if got.Directive != jevImprovementDirectiveExternalReview || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestDeriveJEVImprovementDirectionDirectiveRevision(t *testing.T) {
    got := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
        Aggregation: directionAggregation(jevImprovementFeedbackNeedsRevision),
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        NonAuthorizing: true,
    })
    if got.Directive != jevImprovementDirectiveRevision {
        t.Fatalf("got %+v", got)
    }
}

func TestDeriveJEVImprovementDirectionDirectiveEvidenceHold(t *testing.T) {
    got := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
        Aggregation: directionAggregation(jevImprovementFeedbackHold),
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        NonAuthorizing: true,
    })
    if got.Directive != jevImprovementDirectiveEvidence {
        t.Fatalf("got %+v", got)
    }
}
