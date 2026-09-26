package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSP(t *testing.T) {
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSP(candidate)
    if !diagnostic.Publishable {
        t.Fatalf("expected candidate generation LSP diagnostic to be publishable")
    }
    if diagnostic.Code != "jev.external-apply.mutation-candidate-generated" {
        t.Fatalf("expected generated candidate LSP code, got %q", diagnostic.Code)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid candidate generation LSP diagnostic, got %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSPPreservesUnknown(t *testing.T) {
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
        Observation:    observation,
        CandidateSource: "candidate-source",
        NonAuthorizing: true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationLSP(candidate)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN candidate generation LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "candidate-digest" {
        t.Fatalf("expected candidate-digest as first missing stage, got %q", diagnostic.MissingStage)
    }
}
