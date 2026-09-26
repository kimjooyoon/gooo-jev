package decision

import "testing"

func externalApplyReview() JEVImprovementCandidateReview {
    r := JEVImprovementCandidateReview{
        Status: jevImprovementCandidateReviewConfirmed,
        CandidateDigest: "candidate-digest",
        ReviewerReference: "reviewer://external/one",
        Decision: jevImprovementCandidateReviewAccepted,
        SelectionEvidenceDigest: "selection-evidence-digest",
        ReviewEvidenceDigest: "review-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    r.EvidenceDigest = digestJEVImprovementCandidateReview(r.Status, r.CandidateDigest, r.ReviewerReference, r.Decision, r.SelectionEvidenceDigest, r.ReviewEvidenceDigest)
    return r
}

func externalApplyAggregation(status string) JEVImprovementFeedbackAggregation {
    a := JEVImprovementFeedbackAggregation{
        Status: status,
        Total: 1,
        Confirmed: 1,
        InputEvidenceDigest: "input-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    a.EvidenceDigest = digestJEVImprovementFeedbackAggregation(a.Status, a.Total, a.Confirmed, a.Refuted, a.Unknown, a.InputEvidenceDigest)
    return a
}

func TestCreateJEVExternalApplyRequest(t *testing.T) {
    got := CreateJEVExternalApplyRequest(JEVExternalApplyRequestInput{
        Review: externalApplyReview(),
        Aggregation: externalApplyAggregation(jevImprovementFeedbackStableForReview),
        TargetReference: "gooo://external/apply/review/one",
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyRequestReady || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestCreateJEVExternalApplyRequestRejectsNeedsRevision(t *testing.T) {
    aggregation := externalApplyAggregation(jevImprovementFeedbackNeedsRevision)
    aggregation.Confirmed = 0
    aggregation.Refuted = 1
    aggregation.EvidenceDigest = digestJEVImprovementFeedbackAggregation(aggregation.Status, aggregation.Total, aggregation.Confirmed, aggregation.Refuted, aggregation.Unknown, aggregation.InputEvidenceDigest)
    got := CreateJEVExternalApplyRequest(JEVExternalApplyRequestInput{
        Review: externalApplyReview(),
        Aggregation: aggregation,
        TargetReference: "gooo://external/apply/review/one",
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyRequestUnknown || got.MissingStage != "stable-feedback-aggregation" {
        t.Fatalf("got %+v", got)
    }
}

func TestCreateJEVExternalApplyRequestRequiresTarget(t *testing.T) {
    got := CreateJEVExternalApplyRequest(JEVExternalApplyRequestInput{
        Review: externalApplyReview(),
        Aggregation: externalApplyAggregation(jevImprovementFeedbackStableForReview),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyRequestUnknown || got.MissingStage != "target-reference" {
        t.Fatalf("got %+v", got)
    }
}
