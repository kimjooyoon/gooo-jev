package decision

import "testing"

func TestValidateJEVImprovementLifecycleTransitionAllowed(t *testing.T) {
    got := ValidateJEVImprovementLifecycleTransition(JEVImprovementLifecycleTransitionInput{
        FromStatus: jevImprovementCandidateReady,
        ToStatus: jevImprovementCandidateReviewConfirmed,
        InputEvidenceDigest: "candidate-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevLifecycleTransitionAllowed || got.MissingStage != "" {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestValidateJEVImprovementLifecycleTransitionRejectsDirectApply(t *testing.T) {
    got := ValidateJEVImprovementLifecycleTransition(JEVImprovementLifecycleTransitionInput{
        FromStatus: jevImprovementCandidateReviewConfirmed,
        ToStatus: jevExternalApplyObservedApplied,
        InputEvidenceDigest: "review-evidence-digest",
        NonAuthorizing: true,
    })
    if got.Status != jevLifecycleTransitionUnknown || got.MissingStage != "lifecycle-transition" {
        t.Fatalf("got %+v", got)
    }
}

func TestValidateJEVImprovementLifecycleTransitionRequiresEvidence(t *testing.T) {
    got := ValidateJEVImprovementLifecycleTransition(JEVImprovementLifecycleTransitionInput{
        FromStatus: jevImprovementCandidateReady,
        ToStatus: jevImprovementCandidateReviewConfirmed,
        NonAuthorizing: true,
    })
    if got.Status != jevLifecycleTransitionUnknown || got.MissingStage != "transition-evidence" {
        t.Fatalf("got %+v", got)
    }
}
