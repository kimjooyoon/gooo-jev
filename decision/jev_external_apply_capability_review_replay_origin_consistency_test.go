package decision

import "testing"

func TestReconcileJEVExternalApplyCapabilityReviewReplayOrigin(t *testing.T) {
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
    tests := []struct {
        name   string
        reverse string
        status string
    }{
        {
            name:   "consistent",
            reverse: "origin-digest",
            status: jevExternalApplyCapabilityReviewReplayOriginConsistent,
        },
        {
            name:   "mismatch",
            reverse: "different-origin-digest",
            status: jevExternalApplyCapabilityReviewReplayOriginMismatch,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := ReconcileJEVExternalApplyCapabilityReviewReplayOrigin(JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput{
                Observation:          observation,
                DeclaredOriginDigest: "origin-digest",
                ReverseOriginDigest:  test.reverse,
                NonAuthorizing:       true,
            })
            if result.Status != test.status {
                t.Fatalf("expected status %q, got %q", test.status, result.Status)
            }
            if err := result.Validate(); err != nil {
                t.Fatalf("expected valid origin consistency, got %v", err)
            }
        })
    }
}

func TestReconcileJEVExternalApplyCapabilityReviewReplayOriginPreservesUnknownStage(t *testing.T) {
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
    result := ReconcileJEVExternalApplyCapabilityReviewReplayOrigin(JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput{
        Observation:          observation,
        DeclaredOriginDigest: "origin-digest",
        ReverseOriginDigest:  "",
        NonAuthorizing:       true,
    })
    if result.Status != jevExternalApplyCapabilityReviewReplayOriginUnknown {
        t.Fatalf("expected UNKNOWN origin consistency, got %q", result.Status)
    }
    if result.MissingStage != "reverse-origin" {
        t.Fatalf("expected first missing stage reverse-origin, got %q", result.MissingStage)
    }
}
