package decision

import "testing"

func testJEVExternalApplyCapabilityReviewImprovementDirection(direction string) JEVExternalApplyCapabilityReviewImprovementDirection {
    output := JEVExternalApplyCapabilityReviewImprovementDirection{
        Status:          jevExternalApplyCapabilityReviewImprovementDirectionBound,
        Direction:       direction,
        Target:          "gooo-replay",
        CandidateSource: "candidate-source",
        FeedbackDigest:  "feedback-digest",
        NonExecuting:    true,
        NonAuthorizing:  true,
    }
    output.DirectionDigest = digestJEVExternalApplyCapabilityReviewImprovementDirection(
        output.Status,
        output.Direction,
        output.Target,
        output.CandidateSource,
        output.FeedbackDigest,
    )
    return output
}

func TestGateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(t *testing.T) {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        ReverseMutationDigest:  "reverse-mutation-digest",
        NonAuthorizing:         true,
    })
    candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput{
        Observation:     observation,
        CandidateDigest: "candidate-digest",
        CandidateSource: "candidate-source",
        NonAuthorizing:  true,
    })
    bridge := GateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeInput{
        Direction:      testJEVExternalApplyCapabilityReviewImprovementDirection(jevExternalApplyCapabilityReviewImprovementDirectionGenerate),
        Candidate:      candidate,
        RevisionSource: "candidate-source",
        NonAuthorizing: true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady {
        t.Fatalf("expected ready revision gate bridge, got %q", bridge.Status)
    }
    if err := bridge.Validate(); err != nil {
        t.Fatalf("expected valid revision gate bridge, got %v", err)
    }
}

func TestGateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidatePreservesHoldAndRejected(t *testing.T) {
    tests := []struct {
        name      string
        review    string
        direction string
        status    string
    }{
        {
            name:      "hold",
            review:    jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold,
            direction: jevExternalApplyCapabilityReviewImprovementDirectionHold,
            status:    jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold,
        },
        {
            name:      "rejected",
            review:    jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected,
            direction: jevExternalApplyCapabilityReviewImprovementDirectionReject,
            status:    jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(test.review)
            observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
                CapabilityBoundary: capability,
                NonAuthorizing:     true,
            })
            candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput{
                Observation:    observation,
                NonAuthorizing: true,
            })
            bridge := GateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeInput{
                Direction:      testJEVExternalApplyCapabilityReviewImprovementDirection(test.direction),
                Candidate:      candidate,
                NonAuthorizing: true,
            })
            if bridge.Status != test.status {
                t.Fatalf("expected bridge status %q, got %q", test.status, bridge.Status)
            }
            if err := bridge.Validate(); err != nil {
                t.Fatalf("expected valid revision gate bridge, got %v", err)
            }
        })
    }
}

func TestGateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidatePreservesUnknownStage(t *testing.T) {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        ReverseMutationDigest:  "reverse-mutation-digest",
        NonAuthorizing:         true,
    })
    candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput{
        Observation:     observation,
        CandidateDigest: "candidate-digest",
        CandidateSource: "candidate-source",
        NonAuthorizing:  true,
    })
    bridge := GateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeInput{
        Direction:      testJEVExternalApplyCapabilityReviewImprovementDirection(jevExternalApplyCapabilityReviewImprovementDirectionGenerate),
        Candidate:      candidate,
        NonAuthorizing: true,
    })
    if bridge.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateUnknown {
        t.Fatalf("expected UNKNOWN bridge, got %q", bridge.Status)
    }
    if bridge.MissingStage != "revision-source" {
        t.Fatalf("expected first missing stage revision-source, got %q", bridge.MissingStage)
    }
}
