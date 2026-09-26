package decision

import "testing"

func revisionCandidateDirective() JEVImprovementDirectionDirective {
    d := JEVImprovementDirectionDirective{
        Status: jevImprovementDirectiveRevision,
        Directive: jevImprovementDirectiveRevision,
        CandidateDigest: "parent-candidate-digest",
        CandidateSource: "gooo://candidate/parent",
        InputEvidenceDigest: "aggregation-evidence-digest",
        NonExecuting: true,
        NonAuthorizing: true,
    }
    d.EvidenceDigest = digestJEVImprovementDirectionDirective(d.Status, d.Directive, d.CandidateDigest, d.CandidateSource, d.InputEvidenceDigest)
    return d
}

func TestGenerateJEVImprovementRevisionCandidate(t *testing.T) {
    got := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
        Directive: revisionCandidateDirective(),
        RevisionSource: "gooo://revision/source/one",
        RevisionChangeDigest: "revision-change-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementRevisionCandidateReady || got.MissingStage != "" || got.CandidateDigest == "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestGenerateJEVImprovementRevisionCandidateRequiresRevisionChange(t *testing.T) {
    got := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
        Directive: revisionCandidateDirective(),
        RevisionSource: "gooo://revision/source/one",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementRevisionCandidateUnknown || got.MissingStage != "revision-change" {
        t.Fatalf("got %+v", got)
    }
}

func TestGenerateJEVImprovementRevisionCandidateRejectsExternalReviewDirective(t *testing.T) {
    directive := revisionCandidateDirective()
    directive.Status = jevImprovementDirectiveExternalReview
    directive.Directive = jevImprovementDirectiveExternalReview
    directive.EvidenceDigest = digestJEVImprovementDirectionDirective(directive.Status, directive.Directive, directive.CandidateDigest, directive.CandidateSource, directive.InputEvidenceDigest)
    got := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
        Directive: directive,
        RevisionSource: "gooo://revision/source/one",
        RevisionChangeDigest: "revision-change-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevImprovementRevisionCandidateUnknown || got.MissingStage != "revision-directive" {
        t.Fatalf("got %+v", got)
    }
}
