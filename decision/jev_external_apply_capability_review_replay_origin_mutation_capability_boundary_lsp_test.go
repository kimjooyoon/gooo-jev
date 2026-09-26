package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSP(t *testing.T) {
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
        NonAuthorizing:          true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSP(capability)
    if !diagnostic.Publishable {
        t.Fatalf("expected capability boundary LSP diagnostic to be publishable")
    }
    if diagnostic.Code != "jev.external-apply.mutation-capability-recorded" {
        t.Fatalf("expected recorded capability LSP code, got %q", diagnostic.Code)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid capability boundary LSP diagnostic, got %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSPPreservesUnknown(t *testing.T) {
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCapabilityBoundaryLSP(capability)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN capability boundary LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "spiffe-principal" {
        t.Fatalf("expected first missing stage spiffe-principal, got %q", diagnostic.MissingStage)
    }
}
