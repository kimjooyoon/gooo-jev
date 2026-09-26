package decision

import "testing"

func testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(status, reverseOriginDigest string) JEVExternalApplyCapabilityReviewReplayOriginMutationProposal {
    consistency := testJEVExternalApplyCapabilityReviewReplayOriginConsistency(status, reverseOriginDigest)
    proposalStatus := jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded
    proposedOriginDigest := ""
    reason := ""
    if status == jevExternalApplyCapabilityReviewReplayOriginMismatch {
        proposalStatus = jevExternalApplyCapabilityReviewReplayOriginMutationProposed
        proposedOriginDigest = "proposed-origin-digest"
        reason = "reverse observation is authoritative evidence"
    }
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput{
        Consistency:          consistency,
        ProposedOriginDigest: proposedOriginDigest,
        Reason:               reason,
        NonAuthorizing:       true,
    })
    if proposal.Status != proposalStatus {
        panic("test mutation proposal status mismatch")
    }
    return proposal
}

func TestGateJEVExternalApplyCapabilityReviewReplayOriginMutation(t *testing.T) {
    tests := []struct {
        name     string
        decision string
        status   string
    }{
        {
            name:     "approve",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
        },
        {
            name:     "hold",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionHold,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold,
        },
        {
            name:     "reject",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionReject,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
                jevExternalApplyCapabilityReviewReplayOriginMismatch,
                "different-origin-digest",
            )
            gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
                Proposal:             proposal,
                ReviewDecision:       test.decision,
                ReviewEvidenceDigest: "review-evidence-digest",
                NonAuthorizing:       true,
            })
            if gate.Status != test.status {
                t.Fatalf("expected review status %q, got %q", test.status, gate.Status)
            }
            if err := gate.Validate(); err != nil {
                t.Fatalf("expected valid review gate, got %v", err)
            }
        })
    }
}

func TestGateJEVExternalApplyCapabilityReviewReplayOriginMutationNotNeeded(t *testing.T) {
    proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        jevExternalApplyCapabilityReviewReplayOriginConsistent,
        "origin-digest",
    )
    gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
        Proposal:       proposal,
        NonAuthorizing: true,
    })
    if gate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewNotNeeded {
        t.Fatalf("expected mutation-not-needed review status, got %q", gate.Status)
    }
    if err := gate.Validate(); err != nil {
        t.Fatalf("expected valid mutation-not-needed review gate, got %v", err)
    }
}

func TestGateJEVExternalApplyCapabilityReviewReplayOriginMutationPreservesUnknownStage(t *testing.T) {
    proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
        Proposal:       proposal,
        ReviewDecision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove,
        NonAuthorizing: true,
    })
    if gate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationReviewUnknown {
        t.Fatalf("expected UNKNOWN review gate, got %q", gate.Status)
    }
    if gate.MissingStage != "review-evidence" {
        t.Fatalf("expected first missing stage review-evidence, got %q", gate.MissingStage)
    }
}
