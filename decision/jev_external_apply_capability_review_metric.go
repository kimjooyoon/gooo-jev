package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math"
)

const jevExternalApplyCapabilityReviewMetricRecorded = "capability-review-metric-recorded"
const jevExternalApplyCapabilityReviewMetricUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewMetricInput struct {
    Review         JEVExternalApplyCapabilityReview
    MetricName     string
    MetricValue    float64
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewMetric struct {
    Status          string
    MissingStage    string
    ReviewStatus    string
    ReviewDigest    string
    MetricName      string
    MetricValue     float64
    EvidenceDigest  string
    MetricDigest    string
    NonExecuting    bool
    NonAuthorizing  bool
}

func (m JEVExternalApplyCapabilityReviewMetric) Validate() error {
    if m.Status == "" || m.ReviewStatus == "" || m.ReviewDigest == "" || m.MetricName == "" || m.EvidenceDigest == "" || m.MetricDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review metric")
    }
    if m.Status != jevExternalApplyCapabilityReviewMetricRecorded {
        return fmt.Errorf("invalid JEV external apply capability review metric status")
    }
    if m.ReviewStatus != jevExternalApplyCapabilityReviewConfirmed && m.ReviewStatus != jevExternalApplyCapabilityReviewRefuted {
        return fmt.Errorf("invalid JEV external apply capability review metric review status")
    }
    if math.IsNaN(m.MetricValue) || math.IsInf(m.MetricValue, 0) {
        return fmt.Errorf("JEV external apply capability review metric value must be finite")
    }
    if !m.NonExecuting {
        return fmt.Errorf("JEV external apply capability review metric must be non-executing")
    }
    if !m.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review metric must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewMetric(m.ReviewStatus, m.ReviewDigest, m.MetricName, m.MetricValue, m.EvidenceDigest)
    if m.MetricDigest != expected {
        return fmt.Errorf("JEV external apply capability review metric digest mismatch")
    }
    return nil
}

func MeasureJEVExternalApplyCapabilityReview(input JEVExternalApplyCapabilityReviewMetricInput) JEVExternalApplyCapabilityReviewMetric {
    output := JEVExternalApplyCapabilityReviewMetric{
        Status:         jevExternalApplyCapabilityReviewMetricUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Review.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Review.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Review.Validate(); err != nil {
        output.MissingStage = "capability-review"
        return output
    }
    if input.MetricName == "" {
        output.MissingStage = "metric-name"
        return output
    }
    if math.IsNaN(input.MetricValue) || math.IsInf(input.MetricValue, 0) {
        output.MissingStage = "metric-value"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewMetricRecorded
    output.ReviewStatus = input.Review.Status
    output.ReviewDigest = input.Review.ReviewDigest
    output.MetricName = input.MetricName
    output.MetricValue = input.MetricValue
    output.EvidenceDigest = input.Review.ReviewEvidenceDigest
    output.MetricDigest = digestJEVExternalApplyCapabilityReviewMetric(output.ReviewStatus, output.ReviewDigest, output.MetricName, output.MetricValue, output.EvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewMetricUnknown
        output.MissingStage = "capability-review-metric-evidence"
        output.ReviewDigest = ""
        output.MetricDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewMetric(reviewStatus, reviewDigest, metricName string, metricValue float64, evidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%0.17g|%s", reviewStatus, reviewDigest, metricName, metricValue, evidenceDigest)))
    return hex.EncodeToString(sum[:])
}
