package decision

import "testing"

func testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(status string) JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundary {
    proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    decision := jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove
    if status == jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold {
        decision = jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionHold
    }
    if status == jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected {
        decision = jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionReject
    }
    gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
        Proposal:             proposal,
        ReviewDecision:       decision,
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing:       true,
    })
    return RecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryInput{
        ReviewGate:             gate,
        Workspace:              "workspace:gooo",
        PrincipalURI:           "spiffe://gooo.example/reviewer",
        Audience:               "gooo-replay",
        NetworkAllowlistDigest: "network-allowlist-digest",
        NonAuthorizing:         true,
    })
}

func TestObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(t *testing.T) {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        ReverseMutationDigest:  "reverse-mutation-digest",
        NonAuthorizing:         true,
    })
    if observation.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationRecorded {
        t.Fatalf("expected mutation observation to be recorded, got %q", observation.Status)
    }
    if err := observation.Validate(); err != nil {
        t.Fatalf("expected valid mutation observation, got %v", err)
    }
}

func TestObserveJEVExternalApplyCapabilityReviewReplayOriginMutationPreservesHoldAndRejected(t *testing.T) {
    tests := []struct {
        name   string
        review string
        status string
    }{
        {
            name:   "hold",
            review: jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold,
            status: jevExternalApplyCapabilityReviewReplayOriginMutationObservationHold,
        },
        {
            name:   "rejected",
            review: jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected,
            status: jevExternalApplyCapabilityReviewReplayOriginMutationObservationRejected,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(test.review)
            observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
                CapabilityBoundary: capability,
                NonAuthorizing:     true,
            })
            if observation.Status != test.status {
                t.Fatalf("expected status %q, got %q", test.status, observation.Status)
            }
            if err := observation.Validate(); err != nil {
                t.Fatalf("expected valid mutation observation, got %v", err)
            }
        })
    }
}

func TestObserveJEVExternalApplyCapabilityReviewReplayOriginMutationPreservesUnknownStage(t *testing.T) {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        NonAuthorizing:         true,
    })
    if observation.Status != jevExternalApplyCapabilityReviewReplayOriginMutationObservationUnknown {
        t.Fatalf("expected UNKNOWN mutation observation, got %q", observation.Status)
    }
    if observation.MissingStage != "reverse-mutation-observation" {
        t.Fatalf("expected first missing stage reverse-mutation-observation, got %q", observation.MissingStage)
    }
}
