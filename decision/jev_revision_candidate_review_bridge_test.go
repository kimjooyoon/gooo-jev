package decision

import "testing"

func bridgeRevisionCandidate() JEVImprovementRevisionCandidate {
    c := JEVImprovementRevisionCandidate{
        Status: jevImprovementRevisionCandidateReady,
        ParentCandidateDigest: "parent-candidate-digest",
        RevisionSource: "gooo://revision/source/one",
        RevisionChangeDigest: "revision-change-digest",
        DirectiveEvidenceDigest: "directive-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    c.CandidateDigest = digestJEVImprovementRevisionCandidate(c.ParentCandidateDigest, c.RevisionSource, c.RevisionChangeDigest, c.DirectiveEvidenceDigest)
    c.EvidenceDigest = digestJEVImprovementRevisionCandidateEvidence(c.Status, c.CandidateDigest, c.DirectiveEvidenceDigest, c.RevisionChangeDigest)
    return c
}

func TestBridgeJEVImprovementRevisionCandidateForReview(t *testing.T) {
    got := BridgeJEVImprovementRevisionCandidateForReview(JEVRevisionCandidateReviewBridgeInput{
        Candidate: bridgeRevisionCandidate(),
        NonAuthorizing: true,
    })
    if got.Status != jevRevisionCandidateReviewBridgeReady || got.MissingStage != "" || got.Selection.Status != jevImprovementCandidateReady {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestBridgeJEVImprovementRevisionCandidateRejectsTamperedCandidate(t *testing.T) {
    candidate := bridgeRevisionCandidate()
    candidate.RevisionChangeDigest = "tampered-revision-change"
    got := BridgeJEVImprovementRevisionCandidateForReview(JEVRevisionCandidateReviewBridgeInput{
        Candidate: candidate,
        NonAuthorizing: true,
    })
    if got.Status != jevRevisionCandidateReviewBridgeUnknown || got.MissingStage != "revision-candidate" {
        t.Fatalf("got %+v", got)
    }
}
