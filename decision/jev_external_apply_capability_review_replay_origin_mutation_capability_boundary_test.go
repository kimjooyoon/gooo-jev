package decision

import "testing"

func TestRecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(t *testing.T) {
    proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
        Proposal:             proposal,
        ReviewDecision:       jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove,
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing:       true,
    })
    capability := RecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryInput{
        ReviewGate:             gate,
        Workspace:              "workspace:gooo",
        PrincipalURI:           "spiffe://gooo.example/reviewer",
        Audience:               "gooo-replay",
        NetworkAllowlistDigest: "network-allowlist-digest",
        NonAuthorizing:         true,
    })
    if capability.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRecorded {
        t.Fatalf("expected recorded capability boundary, got %q", capability.Status)
    }
    if err := capability.Validate(); err != nil {
        t.Fatalf("expected valid capability boundary, got %v", err)
    }
}

func TestRecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityPreservesHoldAndRejected(t *testing.T) {
    tests := []struct {
        name       string
        decision   string
        status     string
    }{
        {
            name:     "hold",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionHold,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityHold,
        },
        {
            name:     "rejected",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionReject,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityRejected,
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
            capability := RecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryInput{
                ReviewGate:     gate,
                NonAuthorizing: true,
            })
            if capability.Status != test.status {
                t.Fatalf("expected status %q, got %q", test.status, capability.Status)
            }
            if err := capability.Validate(); err != nil {
                t.Fatalf("expected valid capability boundary, got %v", err)
            }
        })
    }
}

func TestRecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityPreservesUnknownStage(t *testing.T) {
    proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
        Proposal:             proposal,
        ReviewDecision:       jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove,
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing:       true,
    })
    capability := RecordJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(JEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryInput{
        ReviewGate:     gate,
        Workspace:      "workspace:gooo",
        PrincipalURI:   "https://not-spiffe",
        Audience:       "gooo-replay",
        NonAuthorizing: true,
    })
    if capability.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCapabilityUnknown {
        t.Fatalf("expected UNKNOWN capability boundary, got %q", capability.Status)
    }
    if capability.MissingStage != "spiffe-principal" {
        t.Fatalf("expected first missing stage spiffe-principal, got %q", capability.MissingStage)
    }
}
