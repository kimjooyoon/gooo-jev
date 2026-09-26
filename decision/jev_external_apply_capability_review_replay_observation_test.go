package decision

import "testing"

func TestObserveJEVExternalApplyCapabilityReviewReplay(t *testing.T) {
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
    if observation.Status != jevExternalApplyCapabilityReviewReplayObserved {
        t.Fatalf("expected replay observation to be recorded, got %q", observation.Status)
    }
    if err := observation.Validate(); err != nil {
        t.Fatalf("expected valid replay observation, got %v", err)
    }
}

func TestObserveJEVExternalApplyCapabilityReviewReplayPreservesHoldAndRejected(t *testing.T) {
    tests := []struct {
        name           string
        decision       string
        lifecycleState string
        status         string
    }{
        {
            name:           "hold",
            decision:       jevExternalApplyCapabilityReviewRevisionCandidateHold,
            lifecycleState: jevExternalApplyCapabilityReviewCandidateLifecycleHold,
            status:         jevExternalApplyCapabilityReviewReplayObservationHold,
        },
        {
            name:           "rejected",
            decision:       jevExternalApplyCapabilityReviewRevisionCandidateRejected,
            lifecycleState: jevExternalApplyCapabilityReviewCandidateLifecycleRejected,
            status:         jevExternalApplyCapabilityReviewReplayObservationRejected,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
                test.decision,
                test.lifecycleState,
                "",
            )
            preparation := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
                Transition:  transition,
                ReplayScope: "workspace:gooo",
                NonAuthorizing: true,
            })
            observation := ObserveJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayObservationInput{
                Preparation:    preparation,
                NonAuthorizing: true,
            })
            if observation.Status != test.status {
                t.Fatalf("expected status %q, got %q", test.status, observation.Status)
            }
            if err := observation.Validate(); err != nil {
                t.Fatalf("expected valid hold/rejected observation, got %v", err)
            }
        })
    }
}

func TestObserveJEVExternalApplyCapabilityReviewReplayPreservesUnknownStage(t *testing.T) {
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
        NonAuthorizing: true,
    })
    if observation.Status != jevExternalApplyCapabilityReviewReplayObservationUnknown {
        t.Fatalf("expected UNKNOWN replay observation, got %q", observation.Status)
    }
    if observation.MissingStage != "observed-artifact" {
        t.Fatalf("expected first missing stage observed-artifact, got %q", observation.MissingStage)
    }
}
