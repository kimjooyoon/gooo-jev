package decision

import "testing"

func testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(decision, lifecycleState, candidateDigest string) JEVExternalApplyCapabilityReviewCandidateLifecycleTransition {
    gate := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:          jevExternalApplyCapabilityReviewRevisionCandidateGated,
        Decision:        decision,
        Direction:       "candidate-generation",
        Target:          "gooo-replay",
        CandidateSource: "candidate-source",
        CandidateDigest: candidateDigest,
        DirectionDigest: "direction-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    gate.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(
        gate.Status,
        gate.Decision,
        gate.Direction,
        gate.Target,
        gate.CandidateSource,
        gate.CandidateDigest,
        gate.DirectionDigest,
    )
    transition := JEVExternalApplyCapabilityReviewCandidateLifecycleTransition{
        Status:          jevExternalApplyCapabilityReviewCandidateLifecycleTransitionRecorded,
        FromState:       "candidate-source-state",
        ToState:         lifecycleState,
        Decision:        decision,
        GateDigest:      gate.GateDigest,
        CandidateDigest: gate.CandidateDigest,
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    transition.TransitionDigest = digestJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        transition.Status,
        transition.FromState,
        transition.ToState,
        transition.Decision,
        transition.GateDigest,
        transition.CandidateDigest,
    )
    return transition
}

func TestPrepareJEVExternalApplyCapabilityReviewReplay(t *testing.T) {
    tests := []struct {
        name            string
        decision        string
        lifecycleState  string
        candidateDigest string
        environment     string
        replayDecision  string
    }{
        {
            name:            "ready",
            decision:        jevExternalApplyCapabilityReviewRevisionCandidateReady,
            lifecycleState:  jevExternalApplyCapabilityReviewCandidateLifecycleReady,
            candidateDigest: "candidate-digest",
            environment:     "environment-digest",
            replayDecision:  jevExternalApplyCapabilityReviewReplayReady,
        },
        {
            name:            "hold",
            decision:        jevExternalApplyCapabilityReviewRevisionCandidateHold,
            lifecycleState:  jevExternalApplyCapabilityReviewCandidateLifecycleHold,
            replayDecision:  jevExternalApplyCapabilityReviewReplayHold,
        },
        {
            name:            "rejected",
            decision:        jevExternalApplyCapabilityReviewRevisionCandidateRejected,
            lifecycleState:  jevExternalApplyCapabilityReviewCandidateLifecycleRejected,
            replayDecision:  jevExternalApplyCapabilityReviewReplayRejected,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
                test.decision,
                test.lifecycleState,
                test.candidateDigest,
            )
            output := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
                Transition:              transition,
                ReplayScope:             "workspace:gooo",
                ReplayEnvironmentDigest: test.environment,
                NonAuthorizing:          true,
            })
            if output.Status != jevExternalApplyCapabilityReviewReplayPreparationRecorded {
                t.Fatalf("expected replay preparation to be recorded, got %q", output.Status)
            }
            if output.Decision != test.replayDecision {
                t.Fatalf("expected replay decision %q, got %q", test.replayDecision, output.Decision)
            }
            if err := output.Validate(); err != nil {
                t.Fatalf("expected valid replay preparation, got %v", err)
            }
        })
    }
}

func TestPrepareJEVExternalApplyCapabilityReviewReplayPreservesUnknownStage(t *testing.T) {
    transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        jevExternalApplyCapabilityReviewCandidateLifecycleReady,
        "candidate-digest",
    )
    output := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
        Transition:     transition,
        ReplayScope:    "",
        NonAuthorizing: true,
    })
    if output.Status != jevExternalApplyCapabilityReviewReplayUnknown {
        t.Fatalf("expected UNKNOWN replay preparation, got %q", output.Status)
    }
    if output.MissingStage != "replay-scope" {
        t.Fatalf("expected first missing stage replay-scope, got %q", output.MissingStage)
    }
}

func TestPrepareJEVExternalApplyCapabilityReviewReplayRequiresEnvironmentForReady(t *testing.T) {
    transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        jevExternalApplyCapabilityReviewCandidateLifecycleReady,
        "candidate-digest",
    )
    output := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
        Transition:  transition,
        ReplayScope: "workspace:gooo",
        NonAuthorizing: true,
    })
    if output.Status != jevExternalApplyCapabilityReviewReplayUnknown {
        t.Fatalf("expected UNKNOWN replay preparation, got %q", output.Status)
    }
    if output.MissingStage != "replay-environment" {
        t.Fatalf("expected first missing stage replay-environment, got %q", output.MissingStage)
    }
}
