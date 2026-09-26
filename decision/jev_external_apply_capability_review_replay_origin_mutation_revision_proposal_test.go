package decision

import "testing"

func testJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(status string) JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        ReverseMutationDigest:  "reverse-mutation-digest",
        NonAuthorizing:         true,
    })
    candidateDigest := "candidate-digest"
    candidateSource := "candidate-source"
    candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput{
        Observation:     observation,
        CandidateDigest: candidateDigest,
        CandidateSource: candidateSource,
        NonAuthorizing:  true,
    })
    direction := testJEVExternalApplyCapabilityReviewImprovementDirection(jevExternalApplyCapabilityReviewImprovementDirectionGenerate)
    if status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady {
        review := jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold
        directionName := jevExternalApplyCapabilityReviewImprovementDirectionHold
        if status == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected {
            review = jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected
            directionName = jevExternalApplyCapabilityReviewImprovementDirectionReject
        }
        capability = testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(review)
        observation = ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
            CapabilityBoundary: capability,
            NonAuthorizing:     true,
        })
        candidate = GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput{
            Observation:    observation,
            NonAuthorizing: true,
        })
        direction = testJEVExternalApplyCapabilityReviewImprovementDirection(directionName)
    }
    bridge := GateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeInput{
        Direction:      direction,
        Candidate:      candidate,
        RevisionSource: candidateSource,
        NonAuthorizing: true,
    })
    if bridge.Status != status {
        panic("test bridge status mismatch")
    }
    return bridge
}

func TestProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(t *testing.T) {
    bridge := testJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(
        jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady,
    )
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalInput{
        Bridge:             bridge,
        RevisionPlanSource: "revision-plan-source",
        RevisionIntentDigest: "revision-intent-digest",
        NonAuthorizing:     true,
    })
    if proposal.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady {
        t.Fatalf("expected ready revision proposal, got %q", proposal.Status)
    }
    if err := proposal.Validate(); err != nil {
        t.Fatalf("expected valid revision proposal, got %v", err)
    }
}

func TestProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPreservesHoldAndRejected(t *testing.T) {
    tests := []struct {
        name   string
        status string
        want   string
    }{
        {
            name:   "hold",
            status: jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold,
            want:   jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold,
        },
        {
            name:   "rejected",
            status: jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected,
            want:   jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            bridge := testJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(test.status)
            proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalInput{
                Bridge:         bridge,
                NonAuthorizing: true,
            })
            if proposal.Status != test.want {
                t.Fatalf("expected proposal status %q, got %q", test.want, proposal.Status)
            }
            if err := proposal.Validate(); err != nil {
                t.Fatalf("expected valid revision proposal, got %v", err)
            }
        })
    }
}

func TestProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPreservesUnknownStage(t *testing.T) {
    bridge := testJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge(
        jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady,
    )
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutationRevision(JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalInput{
        Bridge:         bridge,
        RevisionPlanSource: "revision-plan-source",
        NonAuthorizing: true,
    })
    if proposal.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown {
        t.Fatalf("expected UNKNOWN revision proposal, got %q", proposal.Status)
    }
    if proposal.MissingStage != "revision-intent" {
        t.Fatalf("expected first missing stage revision-intent, got %q", proposal.MissingStage)
    }
}
