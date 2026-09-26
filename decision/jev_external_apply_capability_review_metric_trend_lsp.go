package decision

import (
    "fmt"
    "math"
)

type JEVExternalApplyCapabilityReviewMetricTrendLSPDiagnostic struct {
    Severity        string
    Code            string
    Message         string
    Status          string
    MissingStage    string
    MetricName      string
    SampleCount     int
    ConfirmedCount  int
    RefutedCount    int
    MeanValue       float64
    EvidenceDigest  string
    TrendDigest     string
    Publishable     bool
    NonExecuting    bool
    NonAuthorizing  bool
}

func (d JEVExternalApplyCapabilityReviewMetricTrendLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review metric trend LSP diagnostic")
    }
    if d.SampleCount <= 0 || d.ConfirmedCount < 0 || d.RefutedCount < 0 || d.ConfirmedCount+d.RefutedCount != d.SampleCount {
        return fmt.Errorf("invalid external apply capability review metric trend LSP counts")
    }
    if math.IsNaN(d.MeanValue) || math.IsInf(d.MeanValue, 0) {
        return fmt.Errorf("external apply capability review metric trend LSP mean must be finite")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review metric trend LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review metric trend LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.MetricName == "" || d.EvidenceDigest == "" || d.TrendDigest == "") {
        return fmt.Errorf("publishable external apply capability review metric trend LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewMetricTrendLSP(input JEVExternalApplyCapabilityReviewMetricTrend) JEVExternalApplyCapabilityReviewMetricTrendLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewMetricTrendLSPDiagnostic{
        Status:         input.Status,
        MissingStage:   input.MissingStage,
        MetricName:     input.MetricName,
        SampleCount:    input.SampleCount,
        ConfirmedCount: input.ConfirmedCount,
        RefutedCount:   input.RefutedCount,
        MeanValue:      input.MeanValue,
        EvidenceDigest: input.EvidenceDigest,
        TrendDigest:    input.TrendDigest,
        NonExecuting:   input.NonExecuting,
        NonAuthorizing: input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review metric trend is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewMetricTrendRecorded:
        output.Severity = "info"
        output.Code = "jev.external-apply.review-metric-trend-recorded"
        output.Message = "Capability review metric trend is recorded with evidence; no authorization or execution occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review metric trend is UNKNOWN; evidence must be resolved before use"
        output.Publishable = false
    }
    return output
}
