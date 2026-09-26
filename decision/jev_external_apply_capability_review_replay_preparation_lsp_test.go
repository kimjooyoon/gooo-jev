package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayPreparationLSP(t *testing.T) {
    transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        jevExternalApplyCapabilityReviewCandidateLifecycleReady,
        "candidate-digest",
    )
    preparation := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
        Transition:              transition,
        ReplayScope:             "workspace:gooo",
        ReplayEnvironmentDigest: "environment-digest",
        NonAuthorizing:          true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayPreparationLSP(preparation)
    if !diagnostic.Publishable {
        t.Fatalf("expected replay preparation LSP diagnostic to be publishable")
    }
    if diagnostic.Code != "jev.external-apply.replay-ready" {
        t.Fatalf("expected replay-ready LSP code, got %q", diagnostic.Code)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid replay preparation LSP diagnostic, got %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayPreparationLSPPreservesUnknown(t *testing.T) {
    transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        jevExternalApplyCapabilityReviewCandidateLifecycleReady,
        "candidate-digest",
    )
    preparation := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
        Transition:     transition,
        ReplayScope:    "",
        NonAuthorizing: true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayPreparationLSP(preparation)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN replay preparation LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "replay-scope" {
        t.Fatalf("expected replay-scope as first missing stage, got %q", diagnostic.MissingStage)
    }
}
