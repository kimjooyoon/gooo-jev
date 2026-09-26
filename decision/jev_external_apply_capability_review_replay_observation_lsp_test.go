package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayObservationLSP(t *testing.T) {
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
    observation := ObserveJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayObservationInput{
        Preparation:              preparation,
        ObservedArtifactDigest:   "observed-artifact-digest",
        ReverseObservationDigest: "reverse-observation-digest",
        NonAuthorizing:            true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayObservationLSP(observation)
    if !diagnostic.Publishable {
        t.Fatalf("expected replay observation LSP diagnostic to be publishable")
    }
    if diagnostic.Code != "jev.external-apply.replay-observed" {
        t.Fatalf("expected replay-observed LSP code, got %q", diagnostic.Code)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid replay observation LSP diagnostic, got %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayObservationLSPPreservesUnknown(t *testing.T) {
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
    observation := ObserveJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayObservationInput{
        Preparation:    preparation,
        ObservedArtifactDigest: "observed-artifact-digest",
        NonAuthorizing: true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayObservationLSP(observation)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN replay observation LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "reverse-observation" {
        t.Fatalf("expected reverse-observation as first missing stage, got %q", diagnostic.MissingStage)
    }
}
