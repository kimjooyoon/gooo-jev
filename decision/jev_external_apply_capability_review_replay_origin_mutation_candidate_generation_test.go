package decision

import "testing"

func TestGenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(t *testing.T) {
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
    if candidate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated {
        t.Fatalf("expected mutation candidate to be generated, got %q", candidate.Status)
    }
    if err := candidate.Validate(); err != nil {
        t.Fatalf("expected valid mutation candidate generation, got %v", err)
    }
}

func TestGenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidatePreservesHoldAndRejected(t *testing.T) {
    tests := []struct {
        name   string
        review string
        status string
    }{
        {
            name:   "hold",
            review: jevExternalApplyCapabilityReviewReplayOriginMutationReviewHold,
            status: jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold,
        },
        {
            name:   "rejected",
            review: jevExternalApplyCapabilityReviewReplayOriginMutationReviewRejected,
            status: jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            capability := testJEVExternalApplyCapabilityReviewReplayOriginMutationCapability(test.review)
            observation := ObserveJEVExternalApplyCapabilityReviewReplayOriginMutation(JEVExternalApplyCapabilityReviewReplayOriginMutationObservationInput{
                CapabilityBoundary: capability,
                NonAuthorizing:     true,
            })
            candidate := GenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidate(JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerationInput{
                Observation:    observation,
                NonAuthorizing: true,
            })
            if candidate.Status != test.status {
                t.Fatalf("expected status %q, got %q", test.status, candidate.Status)
            }
            if err := candidate.Validate(); err != nil {
                t.Fatalf("expected valid candidate generation, got %v", err)
            }
        })
    }
}

func TestGenerateJEVExternalApplyCapabilityReviewReplayOriginMutationCandidatePreservesUnknownStage(t *testing.T) {
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
    if candidate.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateUnknown {
        t.Fatalf("expected UNKNOWN mutation candidate, got %q", candidate.Status)
    }
    if candidate.MissingStage != "candidate-digest" {
        t.Fatalf("expected first missing stage candidate-digest, got %q", candidate.MissingStage)
    }
}
