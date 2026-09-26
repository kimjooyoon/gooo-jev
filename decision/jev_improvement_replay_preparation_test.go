package decision

import "testing"

func replayPreparationReview() JEVImprovementCandidateReview {
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

func TestPrepareJEVImprovementReplay(t *testing.T) {
    got := PrepareJEVImprovementReplay(JEVImprovementReplayPreparationInput{
        Review: replayPreparationReview(),
        ReplayScope: "sandbox://candidate-one",
        ReplayEnvironment: "gooo://replay/environment/v1",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementReplayReady || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestPrepareJEVImprovementReplayDoesNotPromoteHold(t *testing.T) {
    review := replayPreparationReview()
    review.Status = jevImprovementCandidateReviewHold
    review.EvidenceDigest = digestJEVImprovementCandidateReview(review.Status, review.CandidateDigest, review.ReviewerReference, review.Decision, review.SelectionEvidenceDigest, review.ReviewEvidenceDigest)
    got := PrepareJEVImprovementReplay(JEVImprovementReplayPreparationInput{
        Review: review,
        ReplayScope: "sandbox://candidate-one",
        ReplayEnvironment: "gooo://replay/environment/v1",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementReplayUnknown || got.MissingStage != "review-confirmation" {
        t.Fatalf("got %+v", got)
    }
}

func TestPrepareJEVImprovementReplayRequiresScope(t *testing.T) {
    got := PrepareJEVImprovementReplay(JEVImprovementReplayPreparationInput{
        Review: replayPreparationReview(),
        ReplayEnvironment: "gooo://replay/environment/v1",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementReplayUnknown || got.MissingStage != "replay-scope" {
        t.Fatalf("got %+v", got)
    }
}
