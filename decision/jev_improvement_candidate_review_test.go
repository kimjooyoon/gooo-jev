package decision

import "testing"

func reviewCandidateSelection() JEVImprovementCandidateSelection {
    s := JEVImprovementCandidateSelection{
        Status: jevImprovementCandidateReady,
        CandidateDigest: "candidate-digest",
        CandidateSource: "gooo://candidate/one",
        ObservationEvidenceDigest: "observation-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    s.EvidenceDigest = digestJEVImprovementCandidateSelection(s.Status, s.CandidateDigest, s.CandidateSource, s.ObservationEvidenceDigest)
    return s
}

func TestReviewJEVImprovementCandidateAccepted(t *testing.T) {
    got := ReviewJEVImprovementCandidate(JEVImprovementCandidateReviewInput{
        Selection: reviewCandidateSelection(),
        ReviewerReference: "reviewer://external/one",
        Decision: jevImprovementCandidateReviewAccepted,
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementCandidateReviewConfirmed || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestReviewJEVImprovementCandidateUnknownIsHold(t *testing.T) {
    got := ReviewJEVImprovementCandidate(JEVImprovementCandidateReviewInput{
        Selection: reviewCandidateSelection(),
        ReviewerReference: "reviewer://external/one",
        Decision: jevImprovementCandidateReviewUnknown,
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementCandidateReviewHold || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
}

func TestReviewJEVImprovementCandidateRequiresReviewerEvidence(t *testing.T) {
    got := ReviewJEVImprovementCandidate(JEVImprovementCandidateReviewInput{
        Selection: reviewCandidateSelection(),
        Decision: jevImprovementCandidateReviewAccepted,
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementCandidateReviewStatusUnknown || got.MissingStage != "reviewer-reference" {
        t.Fatalf("got %+v", got)
    }
}
