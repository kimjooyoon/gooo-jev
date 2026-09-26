package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSP(t *testing.T) {
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSP(bridge)
    if !diagnostic.Publishable {
        t.Fatalf("expected candidate revision gate bridge LSP diagnostic to be publishable")
    }
    if diagnostic.Code != "jev.external-apply.mutation-candidate-revision-ready" {
        t.Fatalf("expected ready bridge LSP code, got %q", diagnostic.Code)
    }
    if err := diagnostic.Validate(); err != nil {
        t.Fatalf("expected valid bridge LSP diagnostic, got %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSPPreservesUnknown(t *testing.T) {
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
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSP(bridge)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN bridge LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "revision-source" {
        t.Fatalf("expected revision-source as first missing stage, got %q", diagnostic.MissingStage)
    }
}
