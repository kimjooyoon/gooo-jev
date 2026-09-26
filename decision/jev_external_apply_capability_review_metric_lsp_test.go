package decision

import "testing"

func validExternalApplyCapabilityReviewMetricForLSP(status string) JEVExternalApplyCapabilityReviewMetric {
    metric := JEVExternalApplyCapabilityReviewMetric{
        Status:         jevExternalApplyCapabilityReviewMetricRecorded,
        ReviewStatus:   status,
        ReviewDigest:   "review-digest",
        MetricName:     "review.confidence",
        MetricValue:    0.75,
        EvidenceDigest: "review-evidence-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    metric.MetricDigest = digestJEVExternalApplyCapabilityReviewMetric(metric.ReviewStatus, metric.ReviewDigest, metric.MetricName, metric.MetricValue, metric.EvidenceDigest)
    return metric
}

func TestProjectJEVExternalApplyCapabilityReviewMetricLSPIsPublishable(t *testing.T) {
    for _, status := range []string{
        jevExternalApplyCapabilityReviewConfirmed,
        jevExternalApplyCapabilityReviewRefuted,
    } {
        got := ProjectJEVExternalApplyCapabilityReviewMetricLSP(validExternalApplyCapabilityReviewMetricForLSP(status))
        if got.Severity != "info" || got.Code != "jev.external-apply.review-metric-recorded" || !got.Publishable {
            t.Fatalf("status %q got %+v", status, got)
        }
        if err := got.Validate(); err != nil {
            t.Fatal(err)
        }
    }
}

func TestProjectJEVExternalApplyCapabilityReviewMetricLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewMetricLSP(JEVExternalApplyCapabilityReviewMetric{
        Status:         jevExternalApplyCapabilityReviewMetricUnknown,
        MissingStage:   "metric-value",
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
