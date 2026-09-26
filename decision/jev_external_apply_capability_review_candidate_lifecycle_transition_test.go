package decision

import "testing"

func validExternalApplyCapabilityReviewRevisionCandidateGateForLifecycle(decision string) JEVExternalApplyCapabilityReviewRevisionCandidateGate {
    value := JEVExternalApplyCapabilityReviewRevisionCandidateGate{
        Status:          jevExternalApplyCapabilityReviewRevisionCandidateGated,
        Decision:        decision,
        Direction:       jevExternalApplyCapabilityReviewImprovementDirectionGenerate,
        Target:          "decision/provenance",
        CandidateSource: "review-metric-trend",
        DirectionDigest: "direction-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    if decision == jevExternalApplyCapabilityReviewRevisionCandidateReady {
        value.CandidateDigest = "candidate-digest"
    }
    value.GateDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateGate(value.Status, value.Decision, value.Direction, value.Target, value.CandidateSource, value.CandidateDigest, value.DirectionDigest)
    return value
}

func TestTransitionJEVExternalApplyCapabilityReviewCandidateLifecyclePreservesDecisions(t *testing.T) {
    cases := []struct {
        decision string
        toState  string
    }{
        {jevExternalApplyCapabilityReviewRevisionCandidateReady, jevExternalApplyCapabilityReviewCandidateLifecycleReady},
        {jevExternalApplyCapabilityReviewRevisionCandidateHold, jevExternalApplyCapabilityReviewCandidateLifecycleHold},
        {jevExternalApplyCapabilityReviewRevisionCandidateRejected, jevExternalApplyCapabilityReviewCandidateLifecycleRejected},
    }
    for _, testCase := range cases {
        got := TransitionJEVExternalApplyCapabilityReviewCandidateLifecycle(JEVExternalApplyCapabilityReviewCandidateLifecycleTransitionInput{
            Gate:           validExternalApplyCapabilityReviewRevisionCandidateGateForLifecycle(testCase.decision),
            FromState:      "reviewed",
            NonAuthorizing: true,
        })
        if got.ToState != testCase.toState || got.Decision != testCase.decision || !got.NonExecuting || !got.NonAuthorizing {
            t.Fatalf("decision %q got %+v", testCase.decision, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestTransitionJEVExternalApplyCapabilityReviewCandidateLifecycleRejectsMissingSourceState(t *testing.T) {
    got := TransitionJEVExternalApplyCapabilityReviewCandidateLifecycle(JEVExternalApplyCapabilityReviewCandidateLifecycleTransitionInput{
        Gate:           validExternalApplyCapabilityReviewRevisionCandidateGateForLifecycle(jevExternalApplyCapabilityReviewRevisionCandidateReady),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewCandidateLifecycleUnknown || got.MissingStage != "lifecycle-source-state" {
        t.Fatalf("got %+v", got)
    }
}
