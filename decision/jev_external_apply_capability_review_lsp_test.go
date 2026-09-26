package decision

import "testing"

func validExternalApplyCapabilityReviewForLSP(status string) JEVExternalApplyCapabilityReview {
    review := JEVExternalApplyCapabilityReview{
        Status:               status,
        CapabilityDigest:     "capability-digest",
        ReviewerPrincipal:    "spiffe://example.org/ns/prod/sa/jev-reviewer",
        ReviewEvidenceDigest: "review-evidence-digest",
        NonExecuting:         true,
        NonAuthorizing:       true,
    }
    review.ReviewDigest = digestJEVExternalApplyCapabilityReview(review.Status, review.CapabilityDigest, review.ReviewerPrincipal, review.ReviewEvidenceDigest)
    return review
}

func TestProjectJEVExternalApplyCapabilityReviewLSPConfirmedIsInformational(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewLSP(validExternalApplyCapabilityReviewForLSP(jevExternalApplyCapabilityReviewConfirmed))
    if got.Severity != "info" || got.Code != "jev.external-apply.review-confirmed" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewLSPRefutedIsBlockingWarning(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewLSP(validExternalApplyCapabilityReviewForLSP(jevExternalApplyCapabilityReviewRefuted))
    if got.Severity != "warning" || got.Code != "jev.external-apply.review-refuted" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewLSP(JEVExternalApplyCapabilityReview{
        Status:         jevExternalApplyCapabilityReviewUnknown,
        MissingStage:   "review-evidence",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
