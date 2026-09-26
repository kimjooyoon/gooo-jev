package decision

import (
    "fmt"
    "math"
)

type JEVExternalApplyCapabilityReviewMetricLSPDiagnostic struct {
    Severity         string
    Code             string
    Message          string
    Status           string
    MissingStage     string
    ReviewStatus     string
    ReviewDigest     string
    MetricName       string
    MetricValue      float64
    EvidenceDigest   string
    MetricDigest     string
    Publishable      bool
    NonExecuting     bool
    NonAuthorizing   bool
}

func (d JEVExternalApplyCapabilityReviewMetricLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review metric LSP diagnostic")
    }
    if math.IsNaN(d.MetricValue) || math.IsInf(d.MetricValue, 0) {
        return fmt.Errorf("external apply capability review metric LSP value must be finite")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review metric LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review metric LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.ReviewStatus == "" || d.ReviewDigest == "" || d.MetricName == "" || d.EvidenceDigest == "" || d.MetricDigest == "") {
        return fmt.Errorf("publishable external apply capability review metric LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewMetricLSP(input JEVExternalApplyCapabilityReviewMetric) JEVExternalApplyCapabilityReviewMetricLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewMetricLSPDiagnostic{
        Status:       input.Status,
        MissingStage: input.MissingStage,
        ReviewStatus: input.ReviewStatus,
        ReviewDigest: input.ReviewDigest,
        MetricName:   input.MetricName,
        MetricValue:  input.MetricValue,
        EvidenceDigest: input.EvidenceDigest,
        MetricDigest: input.MetricDigest,
        NonExecuting: input.NonExecuting,
        NonAuthorizing: input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review metric is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewMetricRecorded:
        output.Severity = "info"
        output.Code = "jev.external-apply.review-metric-recorded"
        output.Message = "Capability review metric is recorded with evidence; no authorization or execution occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV capability review metric is UNKNOWN; evidence must be resolved before use"
        output.Publishable = false
    }
    return output
}
