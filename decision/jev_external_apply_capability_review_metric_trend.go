package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math"
    "strings"
)

const jevExternalApplyCapabilityReviewMetricTrendRecorded = "capability-review-metric-trend-recorded"
const jevExternalApplyCapabilityReviewMetricTrendUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewMetricTrendInput struct {
    Metrics        []JEVExternalApplyCapabilityReviewMetric
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewMetricTrend struct {
    Status          string
    MissingStage    string
    MetricName      string
    SampleCount     int
    ConfirmedCount  int
    RefutedCount    int
    MeanValue       float64
    EvidenceDigest  string
    TrendDigest     string
    NonExecuting    bool
    NonAuthorizing  bool
}

func (t JEVExternalApplyCapabilityReviewMetricTrend) Validate() error {
    if t.Status == "" || t.MetricName == "" || t.EvidenceDigest == "" || t.TrendDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review metric trend")
    }
    if t.Status != jevExternalApplyCapabilityReviewMetricTrendRecorded {
        return fmt.Errorf("invalid JEV external apply capability review metric trend status")
    }
    if t.SampleCount <= 0 || t.ConfirmedCount < 0 || t.RefutedCount < 0 || t.ConfirmedCount+t.RefutedCount != t.SampleCount {
        return fmt.Errorf("invalid JEV external apply capability review metric trend counts")
    }
    if math.IsNaN(t.MeanValue) || math.IsInf(t.MeanValue, 0) {
        return fmt.Errorf("JEV external apply capability review metric trend mean must be finite")
    }
    if !t.NonExecuting {
        return fmt.Errorf("JEV external apply capability review metric trend must be non-executing")
    }
    if !t.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review metric trend must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewMetricTrend(t.Status, t.MetricName, t.SampleCount, t.ConfirmedCount, t.RefutedCount, t.MeanValue, t.EvidenceDigest)
    if t.TrendDigest != expected {
        return fmt.Errorf("JEV external apply capability review metric trend digest mismatch")
    }
    return nil
}

func AggregateJEVExternalApplyCapabilityReviewMetricTrend(input JEVExternalApplyCapabilityReviewMetricTrendInput) JEVExternalApplyCapabilityReviewMetricTrend {
    output := JEVExternalApplyCapabilityReviewMetricTrend{
        Status:         jevExternalApplyCapabilityReviewMetricTrendUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if len(input.Metrics) == 0 {
        output.MissingStage = "metric-samples"
        return output
    }
    metricName := input.Metrics[0].MetricName
    evidenceParts := make([]string, 0, len(input.Metrics))
    total := 0.0
    confirmedCount := 0
    refutedCount := 0
    for _, metric := range input.Metrics {
        if err := metric.Validate(); err != nil {
            output.MissingStage = "metric-sample"
            return output
        }
        if metric.MetricName != metricName {
            output.MissingStage = "metric-name-consistency"
            return output
        }
        if !metric.NonExecuting {
            output.MissingStage = "execution-boundary"
            return output
        }
        if !metric.NonAuthorizing {
            output.MissingStage = "authorization-boundary"
            return output
        }
        total += metric.MetricValue
        evidenceParts = append(evidenceParts, metric.MetricDigest, metric.EvidenceDigest)
        switch metric.ReviewStatus {
        case jevExternalApplyCapabilityReviewConfirmed:
            confirmedCount++
        case jevExternalApplyCapabilityReviewRefuted:
            refutedCount++
        default:
            output.MissingStage = "review-status"
            return output
        }
    }
    output.Status = jevExternalApplyCapabilityReviewMetricTrendRecorded
    output.MetricName = metricName
    output.SampleCount = len(input.Metrics)
    output.ConfirmedCount = confirmedCount
    output.RefutedCount = refutedCount
    output.MeanValue = total / float64(output.SampleCount)
    output.EvidenceDigest = digestJEVExternalApplyCapabilityReviewMetricTrendEvidence(strings.Join(evidenceParts, "|"))
    output.TrendDigest = digestJEVExternalApplyCapabilityReviewMetricTrend(output.Status, output.MetricName, output.SampleCount, output.ConfirmedCount, output.RefutedCount, output.MeanValue, output.EvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewMetricTrendUnknown
        output.MissingStage = "metric-trend-evidence"
        output.EvidenceDigest = ""
        output.TrendDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewMetricTrendEvidence(parts string) string {
    sum := sha256.Sum256([]byte(parts))
    return hex.EncodeToString(sum[:])
}

func digestJEVExternalApplyCapabilityReviewMetricTrend(status, metricName string, sampleCount, confirmedCount, refutedCount int, meanValue float64, evidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d|%d|%d|%0.17g|%s", status, metricName, sampleCount, confirmedCount, refutedCount, meanValue, evidenceDigest)))
    return hex.EncodeToString(sum[:])
}
