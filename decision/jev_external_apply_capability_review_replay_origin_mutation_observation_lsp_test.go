package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSP(t *testing.T) {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        ReverseMutationDigest:  "reverse-mutation-digest",
        NonAuthorizing:         true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSP(observation)
    if !diagnostic.Publishable {
        t.Fatalf("expected mutation observation LSP diagnostic to be publishable")
    }
    if diagnostic.Code != "jev.external-apply.mutation-observation-recorded" {
        t.Fatalf("expected recorded mutation observation LSP code, got %q", diagnostic.Code)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid mutation observation LSP diagnostic, got %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSPPreservesUnknown(t *testing.T) {
    capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(
        jevExternalApplyCapabilityReviewReplayOriginMutationReviewApproved,
    )
    observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
        CapabilityBoundary:     capability,
        ObservedMutationDigest: "observed-mutation-digest",
        NonAuthorizing:         true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationObservationLSP(observation)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN mutation observation LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "reverse-mutation-observation" {
        t.Fatalf("expected reverse-mutation-observation as first missing stage, got %q", diagnostic.MissingStage)
    }
}
