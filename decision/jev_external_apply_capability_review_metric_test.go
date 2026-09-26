package decision

import (
    "math"
    "testing"
)

func validExternalApplyCapabilityReviewForMetric(status string) JEVExternalApplyCapabilityReview {
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

func TestMeasureJEVExternalApplyCapabilityReviewRecordsEvidenceLinkedMetric(t *testing.T) {
    for _, status := range []string{
        jevExternalApplyCapabilityReviewConfirmed,
        jevExternalApplyCapabilityReviewRefuted,
    } {
        got := MeasureJEVExternalApplyCapabilityReview(JEVExternalApplyCapabilityReviewMetricInput{
            Review:         validExternalApplyCapabilityReviewForMetric(status),
            MetricName:     "review.confidence",
            MetricValue:    0.75,
            NonAuthorizing: true,
        })
        if got.Status != jevExternalApplyCapabilityReviewMetricRecorded || got.ReviewStatus != status || !got.NonExecuting || !got.NonAuthorizing {
            t.Fatalf("status %q got %+v", status, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestMeasureJEVExternalApplyCapabilityReviewRejectsNonFiniteMetric(t *testing.T) {
    got := MeasureJEVExternalApplyCapabilityReview(JEVExternalApplyCapabilityReviewMetricInput{
        Review:         validExternalApplyCapabilityReviewForMetric(jevExternalApplyCapabilityReviewConfirmed),
        MetricName:     "review.confidence",
        MetricValue:    math.NaN(),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewMetricUnknown || got.MissingStage != "metric-value" {
        t.Fatalf("got %+v", got)
    }
}

func TestMeasureJEVExternalApplyCapabilityReviewRejectsMissingMetricName(t *testing.T) {
    got := MeasureJEVExternalApplyCapabilityReview(JEVExternalApplyCapabilityReviewMetricInput{
        Review:         validExternalApplyCapabilityReviewForMetric(jevExternalApplyCapabilityReviewConfirmed),
        MetricValue:    0.75,
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewMetricUnknown || got.MissingStage != "metric-name" {
        t.Fatalf("got %+v", got)
    }
}
