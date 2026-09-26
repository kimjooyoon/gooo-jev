package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSP(t *testing.T) {
    tests := []struct {
        name       string
        reverse    string
        proposed   string
        reason     string
        code       string
    }{
        {
            name:     "not-needed",
            reverse:  "origin-digest",
            code:     "jev.external-apply.origin-mutation-not-needed",
        },
        {
            name:     "proposed",
            reverse:  "different-origin-digest",
            proposed: "proposed-origin-digest",
            reason:   "reverse observation is authoritative evidence",
            code:     "jev.external-apply.origin-mutation-proposed",
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            consistencyStatus := jevExternalApplyCapabilityReviewReplayOriginConsistent
            if test.reverse != "origin-digest" {
                consistencyStatus = jevExternalApplyCapabilityReviewReplayOriginMismatch
            }
            consistency := testJEVExternalApplyCapabilityReviewReplayOriginConsistency(
                consistencyStatus,
                test.reverse,
            )
            proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput{
                Consistency:          consistency,
                ProposedOriginDigest: test.proposed,
                Reason:               test.reason,
                NonAuthorizing:       true,
            })
            diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSP(proposal)
            if !diagnostic.Publishable {
                t.Fatalf("expected mutation proposal LSP diagnostic to be publishable")
            }
            if diagnostic.Code != test.code {
                t.Fatalf("expected LSP code %q, got %q", test.code, diagnostic.Code)
            }
            if err := diagnostic.Validate(); err != nil {
                t.Fatalf("expected valid mutation proposal LSP diagnostic, got %v", err)
            }
        })
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSPPreservesUnknown(t *testing.T) {
    consistency := testJEVExternalApplyCapabilityReviewReplayOriginConsistency(
        jevExternalApplyCapabilityReviewReplayOriginMismatch,
        "different-origin-digest",
    )
    proposal := ProposeJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationProposalInput{
        Consistency:          consistency,
        ProposedOriginDigest: "proposed-origin-digest",
        NonAuthorizing:       true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSP(proposal)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN mutation proposal LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "mutation-reason" {
        t.Fatalf("expected mutation-reason as first missing stage, got %q", diagnostic.MissingStage)
    }
}
