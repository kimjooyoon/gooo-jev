package decision

import "testing"

func validExternalApplyCapabilityForReview() JEVExternalApplyCapabilityBoundary {
    scopeDigest, ok := digestJEVExternalApplyCapabilityScopes([]string{"review"})
    if !ok {
        panic("test scope digest must be valid")
    }
    capability := JEVExternalApplyCapabilityBoundary{
        Status:                 jevExternalApplyCapabilityDescribed,
        BindingDigest:          "binding-evidence-digest",
        Principal:              "spiffe://example.org/ns/prod/sa/jev-reviewer",
        Audience:               "gooo-jev",
        Workspace:              "gooo-jev",
        NetworkAllowlistDigest: "network-allowlist-digest",
        ScopeDigest:             scopeDigest,
        NonExecuting:            true,
        NonAuthorizing:          true,
    }
    capability.CapabilityDigest = digestJEVExternalApplyCapability(capability.BindingDigest, capability.Principal, capability.Audience, capability.Workspace, capability.NetworkAllowlistDigest, capability.ScopeDigest)
    return capability
}

func TestObserveJEVExternalApplyCapabilityReviewPreservesOutcome(t *testing.T) {
    for _, outcome := range []string{
        jevExternalApplyCapabilityReviewConfirmed,
        jevExternalApplyCapabilityReviewRefuted,
    } {
        got := ObserveJEVExternalApplyCapabilityReview(JEVExternalApplyCapabilityReviewInput{
            Capability:           validExternalApplyCapabilityForReview(),
            Outcome:              outcome,
            ReviewerPrincipal:    "spiffe://example.org/ns/prod/sa/jev-reviewer",
            ReviewEvidenceDigest: "review-evidence-digest",
            NonAuthorizing:       true,
        })
        if got.Status != outcome || !got.NonExecuting || !got.NonAuthorizing {
            t.Fatalf("outcome %q got %+v", outcome, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestObserveJEVExternalApplyCapabilityReviewRejectsUnboundReviewer(t *testing.T) {
    got := ObserveJEVExternalApplyCapabilityReview(JEVExternalApplyCapabilityReviewInput{
        Capability:           validExternalApplyCapabilityForReview(),
        Outcome:              jevExternalApplyCapabilityReviewConfirmed,
        ReviewerPrincipal:    "service-account:jev-reviewer",
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing:       true,
    })
    if got.Status != jevExternalApplyCapabilityReviewUnknown || got.MissingStage != "reviewer-principal" {
        t.Fatalf("got %+v", got)
    }
}

func TestObserveJEVExternalApplyCapabilityReviewPreservesUnknownOutcome(t *testing.T) {
    got := ObserveJEVExternalApplyCapabilityReview(JEVExternalApplyCapabilityReviewInput{
        Capability:           validExternalApplyCapabilityForReview(),
        Outcome:              "pending",
        ReviewerPrincipal:    "spiffe://example.org/ns/prod/sa/jev-reviewer",
        ReviewEvidenceDigest: "review-evidence-digest",
        NonAuthorizing:       true,
    })
    if got.Status != jevExternalApplyCapabilityReviewUnknown || got.MissingStage != "review-outcome" {
        t.Fatalf("got %+v", got)
    }
}
