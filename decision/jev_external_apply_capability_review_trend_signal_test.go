package decision

import "testing"

func validExternalApplyCapabilityReviewMetricTrendForSignal() JEVExternalApplyCapabilityReviewMetricTrend {
    trend := JEVExternalApplyCapabilityReviewMetricTrend{
        Status:         jevExternalApplyCapabilityReviewMetricTrendRecorded,
        MetricName:     "review.confidence",
        SampleCount:    4,
        ConfirmedCount: 3,
        RefutedCount:   1,
        MeanValue:      0.8,
        EvidenceDigest: "trend-evidence-digest",
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    trend.TrendDigest = digestJEVExternalApplyCapabilityReviewMetricTrend(trend.Status, trend.MetricName, trend.SampleCount, trend.ConfirmedCount, trend.RefutedCount, trend.MeanValue, trend.EvidenceDigest)
    return trend
}

func TestDeriveJEVExternalApplyCapabilityReviewTrendSignalImproves(t *testing.T) {
    got := DeriveJEVExternalApplyCapabilityReviewTrendSignal(JEVExternalApplyCapabilityReviewTrendSignalInput{
        Trend:          validExternalApplyCapabilityReviewMetricTrendForSignal(),
        Threshold:      0.7,
        NonAuthorizing: true,
    })
    if got.Signal != jevExternalApplyCapabilityReviewTrendSignalImprove || !got.NonExecuting || !got.NonAuthorizing {
        t.Fatalf("got %+v", got)
    }
    if err := got.Validate(); err != nil {
        t.Fatal(err)
    }
}

func TestDeriveJEVExternalApplyCapabilityReviewTrendSignalRollsBack(t *testing.T) {
    trend := validExternalApplyCapabilityReviewMetricTrendForSignal()
    trend.ConfirmedCount = 1
    trend.RefutedCount = 3
    trend.MeanValue = 0.4
    trend.TrendDigest = digestJEVExternalApplyCapabilityReviewMetricTrend(trend.Status, trend.MetricName, trend.SampleCount, trend.ConfirmedCount, trend.RefutedCount, trend.MeanValue, trend.EvidenceDigest)
    got := DeriveJEVExternalApplyCapabilityReviewTrendSignal(JEVExternalApplyCapabilityReviewTrendSignalInput{
        Trend:          trend,
        Threshold:      0.7,
        NonAuthorizing: true,
    })
    if got.Signal != jevExternalApplyCapabilityReviewTrendSignalRollback {
        t.Fatalf("got %+v", got)
    }
}

func TestDeriveJEVExternalApplyCapabilityReviewTrendSignalHoldsTie(t *testing.T) {
    trend := validExternalApplyCapabilityReviewMetricTrendForSignal()
    trend.ConfirmedCount = 2
    trend.RefutedCount = 2
    trend.MeanValue = 0.7
    trend.TrendDigest = digestJEVExternalApplyCapabilityReviewMetricTrend(trend.Status, trend.MetricName, trend.SampleCount, trend.ConfirmedCount, trend.RefutedCount, trend.MeanValue, trend.EvidenceDigest)
    got := DeriveJEVExternalApplyCapabilityReviewTrendSignal(JEVExternalApplyCapabilityReviewTrendSignalInput{
        Trend:          trend,
        Threshold:      0.7,
        NonAuthorizing: true,
    })
    if got.Signal != jevExternalApplyCapabilityReviewTrendSignalHold {
        t.Fatalf("got %+v", got)
    }
}

func TestDeriveJEVExternalApplyCapabilityReviewTrendSignalPreservesUnknownThreshold(t *testing.T) {
    got := DeriveJEVExternalApplyCapabilityReviewTrendSignal(JEVExternalApplyCapabilityReviewTrendSignalInput{
        Trend:          validExternalApplyCapabilityReviewMetricTrendForSignal(),
        NonAuthorizing: true,
    })
    if got.Status != jevExternalApplyCapabilityReviewTrendSignalRecorded || got.Signal == "" {
        t.Fatalf("got %+v", got)
    }
}
