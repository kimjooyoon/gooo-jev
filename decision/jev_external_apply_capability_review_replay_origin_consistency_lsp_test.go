package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginConsistencyLSP(t *testing.T) {
    transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        jevExternalApplyCapabilityReviewCandidateLifecycleReady,
        "candidate-digest",
    )
    preparation := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
        Transition:              transition,
        ReplayScope:             "workspace:gooo",
        ReplayEnvironmentDigest: "environment-digest",
        NonAuthorizing:          true,
    })
    observation := ObserveJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayObservationInput{
        Preparation:              preparation,
        ObservedArtifactDigest:   "observed-artifact-digest",
        ReverseObservationDigest: "reverse-observation-digest",
        NonAuthorizing:            true,
    })
    tests := []struct {
        name   string
        reverse string
        code   string
    }{
        {
            name:   "consistent",
            reverse: "origin-digest",
            code:   "jev.external-apply.origin-consistent",
        },
        {
            name:   "mismatch",
            reverse: "different-origin-digest",
            code:   "jev.external-apply.origin-mismatch",
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := ReconcileJEVExternalApplyCapabilityReviewReplayOrigin(JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput{
                Observation:          observation,
                DeclaredOriginDigest: "origin-digest",
                ReverseOriginDigest:  test.reverse,
                NonAuthorizing:       true,
            })
            diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginConsistencyLSP(result)
            if !diagnostic.Publishable {
                t.Fatalf("expected origin consistency LSP diagnostic to be publishable")
            }
            if diagnostic.Code != test.code {
                t.Fatalf("expected LSP code %q, got %q", test.code, diagnostic.Code)
            }
            if err := diagnostic.Validate(); err != nil {
                t.Fatalf("expected valid origin consistency LSP diagnostic, got %v", err)
            }
        })
    }
}

func TestProjectJEVExternalApplyCapabilityReviewReplayOriginConsistencyLSPPreservesUnknown(t *testing.T) {
    transition := testJEVExternalApplyCapabilityReviewCandidateLifecycleTransition(
        jevExternalApplyCapabilityReviewRevisionCandidateReady,
        jevExternalApplyCapabilityReviewCandidateLifecycleReady,
        "candidate-digest",
    )
    preparation := PrepareJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayPreparationInput{
        Transition:              transition,
        ReplayScope:             "workspace:gooo",
        ReplayEnvironmentDigest: "environment-digest",
        NonAuthorizing:          true,
    })
    observation := ObserveJEVExternalApplyCapabilityReviewReplay(JEVExternalApplyCapabilityReviewReplayObservationInput{
        Preparation:              preparation,
        ObservedArtifactDigest:   "observed-artifact-digest",
        ReverseObservationDigest: "reverse-observation-digest",
        NonAuthorizing:            true,
    })
    result := ReconcileJEVExternalApplyCapabilityReviewReplayOrigin(JEVExternalApplyCapabilityReviewReplayOriginConsistencyInput{
        Observation:          observation,
        DeclaredOriginDigest: "origin-digest",
        NonAuthorizing:       true,
    })
    diagnostic := ProjectJEVExternalApplyCapabilityReviewReplayOriginConsistencyLSP(result)
    if diagnostic.Publishable {
        t.Fatalf("expected UNKNOWN origin consistency LSP diagnostic to be non-publishable")
    }
    if diagnostic.Code != "jev.provenance.unknown" {
        t.Fatalf("expected UNKNOWN LSP code, got %q", diagnostic.Code)
    }
    if diagnostic.MissingStage != "reverse-origin" {
        t.Fatalf("expected reverse-origin as first missing stage, got %q", diagnostic.MissingStage)
    }
}
