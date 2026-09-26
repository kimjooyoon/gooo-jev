package decision

import "testing"

func testJEVExternalApplyCapabilityReviewReplayOriginConsistency(status, reverseOriginDigest string) JEVExternalApplyCapabilityReviewReplayOriginConsistency {
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
    consistency := ReconcileJEVExternalApplyCapabilityReviewReplayOrigin(JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput{
        Observation:          observation,
        DeclaredOriginDigest: "origin-digest",
        ReverseOriginDigest:  reverseOriginDigest,
        NonAuthorizing:       true,
    })
    if consistency.Status != status {
        panic("test origin consistency status mismatch")
    }
    return consistency
}

func TestProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(t *testing.T) {
    mismatch := testJEVExternalApplyCapabilityReviewReplayOriginConsistency(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput{
        Consistency:          mismatch,
        ProposedOriginDigest: "proposed-origin-digest",
        Reason:               "reverse observation is authoritative evidence",
        NonAuthorizing:       true,
    })
    if proposal.Status != jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        t.Fatalf("expected origin mutation proposal, got %q", proposal.Status)
    }
    if err := proposal.Validate(); err != nil {
        t.Fatalf("expected valid origin mutation proposal, got %v", err)
    }
}

func TestProposeJEVExternalApplyCapabilityReviewReplayOriginMutationNotNeeded(t *testing.T) {
    consistent := testJEVExternalApplyCapabilityReviewReplayOriginConsistency(
        jevExternalApplyCapabilityReviewReplayOriginConsistent,
        "origin-digest",
    )
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput{
        Consistency:    consistent,
        NonAuthorizing: true,
    })
    if proposal.Status != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded {
        t.Fatalf("expected mutation-not-needed, got %q", proposal.Status)
    }
    if err := proposal.Validate(); err != nil {
        t.Fatalf("expected valid mutation-not-needed result, got %v", err)
    }
}

func TestProposeJEVExternalApplyCapabilityReviewReplayOriginMutationPreservesUnknownStage(t *testing.T) {
    mismatch := testJEVExternalApplyCapabilityReviewReplayOriginConsistency(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput{
        Consistency:          mismatch,
        ProposedOriginDigest: "proposed-origin-digest",
        NonAuthorizing:       true,
    })
    if proposal.Status != jevExternalApplyCapabilityReviewReplayOriginMutationUnknown {
        t.Fatalf("expected UNKNOWN mutation proposal, got %q", proposal.Status)
    }
    if proposal.MissingStage != "mutation-reason" {
        t.Fatalf("expected first missing stage mutation-reason, got %q", proposal.MissingStage)
    }
}
