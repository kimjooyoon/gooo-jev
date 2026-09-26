package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSP(t *testing.T) {
    tests := []struct {
        name     string
        decision string
        status   string
        code     string
    }{
        {
            name:     "approve",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
            code:     "jev.external-apply.origin-mutation-review-approved",
        },
        {
            name:     "hold",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionHold,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold,
            code:     "jev.external-apply.origin-mutation-review-hold",
        },
        {
            name:     "reject",
            decision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionReject,
            status:   jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected,
            code:     "jev.external-apply.origin-mutation-review-rejected",
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
            diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSP(gate)
            if !diagnostic.Publishable {
                t.Fatalf("expected review gate LSP diagnostic to be publishable")
            }
            if diagnostic.Status != test.status || diagnostic.Code != test.code {
                t.Fatalf("expected status %q and code %q, got %q and %q", test.status, test.code, diagnostic.Status, diagnostic.Code)
            }
            if err := diagnostic.Validate(); err != nil {
                t.Fatalf("expected valid review gate LSP diagnostic, got %v", err)
            }
        })
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSPPreservesUnknown(t *testing.T) {
    proposal := testJEVExternalApplyCapabilityReviewReplayOriginMutationProposal(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    gate := GateJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateInput{
        Proposal:       proposal,
        ReviewDecision: jevExternalApplyCapabilityReviewReplayOriginMutationReviewDecisionApprove,
        NonAuthorizing: true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationReviewGateLSP(gate)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN review gate LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "review-evidence" {
        t.Fatalf("expected first missing stage review-evidence, got %q", diagnostic.MissingStage)
    }
}
