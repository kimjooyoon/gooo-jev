package decision

import "testing"

func validExternalApplyCapabilityReviewMetricTrendForLSP() JEVExternalApplyCapabilityReviewMetricTrend {
    trend := JEVExternalApplyCapabilityReviewMetricTrend{
        Status:         jevExternalApplyCapabilityReviewMetricTrendRecorded,
        MetricName:     "review.confidence",
        SampleCount:    2,
        ConfirmedCount: 1,
        RefutedCount:   1,
        MeanValue:      0.5,
        EvidenceDigest: "trend-evidence-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    trend.TrendDigest = digestJEVExternalApplyCapabilityReviewMetricTrend(trend.Status, trend.MetricName, trend.SampleCount, trend.ConfirmedCount, trend.RefutedCount, trend.MeanValue, trend.EvidenceDigest)
    return trend
}

func TestProjectJEVExternalApplyCapabilityReviewMetricTrendLSPIsPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewMetricTrendLSP(validExternalApplyCapabilityReviewMetricTrendForLSP())
    if got.Severity != "info" || got.Code != "jev.external-apply.review-metric-trend-recorded" || !got.Publishable {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewMetricTrendLSPUnknownIsNotPublishable(t *testing.T) {
    got := ProjectJEVExternalApplyCapabilityReviewMetricTrendLSP(JEVExternalApplyCapabilityReviewMetricTrend{
        Status:         jevExternalApplyCapabilityReviewMetricTrendUnknown,
        MissingStage:   "metric-name-consistency",
        SampleCount:    0,
        NonExecuting:   true,
        NonAuthorizing: true,
    })
    if got.Publishable || got.Severity != "error" || got.Code != "jev.provenance.unknown" {
        t.Fatalf("got %+v", got)
    }
}
