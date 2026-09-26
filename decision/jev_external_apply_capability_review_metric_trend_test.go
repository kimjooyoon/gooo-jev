package decision

import (
    "fmt"
    "testing"
)

func validExternalApplyCapabilityReviewMetricForTrend(status string, value float64, suffix string) JEVExternalApplyCapabilityReviewMetric {
    metric := JEVExternalApplyCapabilityReviewMetric{
        Status:         jevExternalApplyCapabilityReviewMetricRecorded,
        ReviewStatus:   status,
        ReviewDigest:   fmt.Sprintf("review-digest-%s", suffix),
        MetricName:     "review.confidence",
        MetricValue:    value,
        EvidenceDigest: fmt.Sprintf("review-evidence-%s", suffix),
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    metric.MetricDigest = digestJEVExternalApplyCapabilityReviewMetric(metric.ReviewStatus, metric.ReviewDigest, metric.MetricName, metric.MetricValue, metric.EvidenceDigest)
    return metric
}

func TestAggregateJEVExternalApplyCapabilityReviewMetricTrendPreservesCounts(t *testing.T) {
    got := AggregateJEVExternalApplyCapabilityReviewMetricTrend(JEVExternalApplyCapabilityReviewMetricTrendInput{
        Metrics: []JEVExternalApplyCapabilityReviewMetric{
            validExternalApplyCapabilityReviewMetricForTrend(jevExternalApplyCapabilityReviewConfirmed, 0.75, "one"),
            validExternalApplyCapabilityReviewMetricForTrend(jevExternalApplyCapabilityReviewRefuted, 0.25, "two"),
        },
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewMetricTrendRecorded || got.SampleCount != 2 || got.ConfirmedCount != 1 || got.RefutedCount != 1 || got.MeanValue != 0.5 {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestAggregateJEVExternalApplyCapabilityReviewMetricTrendRejectsMixedMetricNames(t *testing.T) {
    first := validExternalApplyCapabilityReviewMetricForTrend(jevExternalApplyCapabilityReviewConfirmed, 0.75, "one")
    second := validExternalApplyCapabilityReviewMetricForTrend(jevExternalApplyCapabilityReviewRefuted, 0.25, "two")
    second.MetricName = "review.latency"
    second.MetricDigest = digestJEVExternalApplyCapabilityReviewMetric(second.ReviewStatus, second.ReviewDigest, second.MetricName, second.MetricValue, second.EvidenceDigest)
    got := AggregateJEVExternalApplyCapabilityReviewMetricTrend(JEVExternalApplyCapabilityReviewMetricTrendInput{
        Metrics:        []JEVExternalApplyCapabilityReviewMetric{first, second},
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewMetricTrendUnknown || got.MissingStage != "metric-name-consistency" {
        t.Fatalf("got %+v", got)
    }
}

func TestAggregateJEVExternalApplyCapabilityReviewMetricTrendRejectsEmptySamples(t *testing.T) {
    got := AggregateJEVExternalApplyCapabilityReviewMetricTrend(JEVExternalApplyCapabilityReviewMetricTrendInput{
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewMetricTrendUnknown || got.MissingStage != "metric-samples" {
        t.Fatalf("got %+v", got)
    }
}
