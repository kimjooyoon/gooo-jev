package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math"
)

const jevExternalApplyCapabilityReviewTrendSignalRecorded = "capability-review-trend-signal-recorded"
const jevExternalApplyCapabilityReviewTrendSignalImprove = "improve"
const jevExternalApplyCapabilityReviewTrendSignalHold = "hold"
const jevExternalApplyCapabilityReviewTrendSignalRollback = "rollback"
const jevExternalApplyCapabilityReviewTrendSignalUnknown = "UNKNOWN"

type JEVExternalApplyCapabilityReviewTrendSignalInput struct {
    Trend          JEVExternalApplyCapabilityReviewMetricTrend
    Threshold      float64
    NonAuthorizing bool
}

type JEVExternalApplyCapabilityReviewTrendSignal struct {
    Status          string
    MissingStage    string
    Signal          string
    MetricName      string
    TrendDigest     string
    MeanValue       float64
    Threshold       float64
    EvidenceDigest  string
    SignalDigest    string
    NonExecuting    bool
    NonAuthorizing  bool
}

func (s JEVExternalApplyCapabilityReviewTrendSignal) Validate() error {
    if s.Status == "" || s.Signal == "" || s.MetricName == "" || s.TrendDigest == "" || s.EvidenceDigest == "" || s.SignalDigest == "" {
        return fmt.Errorf("incomplete JEV external apply capability review trend signal")
    }
    if s.Status != jevExternalApplyCapabilityReviewTrendSignalRecorded {
        return fmt.Errorf("invalid JEV external apply capability review trend signal status")
    }
    if s.Signal != jevExternalApplyCapabilityReviewTrendSignalImprove && s.Signal != jevExternalApplyCapabilityReviewTrendSignalHold && s.Signal != jevExternalApplyCapabilityReviewTrendSignalRollback {
        return fmt.Errorf("invalid JEV external apply capability review trend signal")
    }
    if math.IsNaN(s.MeanValue) || math.IsInf(s.MeanValue, 0) || math.IsNaN(s.Threshold) || math.IsInf(s.Threshold, 0) {
        return fmt.Errorf("JEV external apply capability review trend signal values must be finite")
    }
    if !s.NonExecuting {
        return fmt.Errorf("JEV external apply capability review trend signal must be non-executing")
    }
    if !s.NonAuthorizing {
        return fmt.Errorf("JEV external apply capability review trend signal must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewTrendSignal(s.Status, s.Signal, s.MetricName, s.TrendDigest, s.MeanValue, s.Threshold, s.EvidenceDigest)
    if s.SignalDigest != expected {
        return fmt.Errorf("JEV external apply capability review trend signal digest mismatch")
    }
    return nil
}

func DeriveJEVExternalApplyCapabilityReviewTrendSignal(input JEVExternalApplyCapabilityReviewTrendSignalInput) JEVExternalApplyCapabilityReviewTrendSignal {
    output := JEVExternalApplyCapabilityReviewTrendSignal{
        Status:         jevExternalApplyCapabilityReviewTrendSignalUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Trend.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Trend.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Trend.Validate(); err != nil {
        output.MissingStage = "metric-trend"
        return output
    }
    if math.IsNaN(input.Threshold) || math.IsInf(input.Threshold, 0) {
        output.MissingStage = "signal-threshold"
        return output
    }
    signal := jevExternalApplyCapabilityReviewTrendSignalHold
    if input.Trend.ConfirmedCount > input.Trend.RefutedCount && input.Trend.MeanValue >= input.Threshold {
        signal = jevExternalApplyCapabilityReviewTrendSignalImprove
    } else if input.Trend.RefutedCount > input.Trend.ConfirmedCount || input.Trend.MeanValue < input.Threshold {
        signal = jevExternalApplyCapabilityReviewTrendSignalRollback
    }
    output.Status = jevExternalApplyCapabilityReviewTrendSignalRecorded
    output.Signal = signal
    output.MetricName = input.Trend.MetricName
    output.TrendDigest = input.Trend.TrendDigest
    output.MeanValue = input.Trend.MeanValue
    output.Threshold = input.Threshold
    output.EvidenceDigest = input.Trend.EvidenceDigest
    output.SignalDigest = digestJEVExternalApplyCapabilityReviewTrendSignal(output.Status, output.Signal, output.MetricName, output.TrendDigest, output.MeanValue, output.Threshold, output.EvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewTrendSignalUnknown
        output.MissingStage = "trend-signal-evidence"
        output.TrendDigest = ""
        output.SignalDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewTrendSignal(status, signal, metricName, trendDigest string, meanValue, threshold float64, evidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%0.17g|%0.17g|%s", status, signal, metricName, trendDigest, meanValue, threshold, evidenceDigest)))
    return hex.EncodeToString(sum[:])
}
