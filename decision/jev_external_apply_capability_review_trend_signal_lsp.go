package decision

import (
    "fmt"
    "math"
)

type JEVExternalApplyCapabilityReviewTrendSignalLSPDiagnostic struct {
    Severity        string
    Code            string
    Message         string
    Status          string
    MissingStage    string
    Signal          string
    MetricName      string
    TrendDigest     string
    MeanValue       float64
    Threshold       float64
    EvidenceDigest  string
    SignalDigest    string
    Publishable     bool
    NonExecuting    bool
    NonAuthorizing  bool
}

func (d JEVExternalApplyCapabilityReviewTrendSignalLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review trend signal LSP diagnostic")
    }
    if d.Signal != jevExternalApplyCapabilityReviewTrendSignalImprove && d.Signal != jevExternalApplyCapabilityReviewTrendSignalHold && d.Signal != jevExternalApplyCapabilityReviewTrendSignalRollback {
        return fmt.Errorf("invalid external apply capability review trend signal LSP signal")
    }
    if math.IsNaN(d.MeanValue) || math.IsInf(d.MeanValue, 0) || math.IsNaN(d.Threshold) || math.IsInf(d.Threshold, 0) {
        return fmt.Errorf("external apply capability review trend signal LSP values must be finite")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review trend signal LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review trend signal LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.MetricName == "" || d.TrendDigest == "" || d.EvidenceDigest == "" || d.SignalDigest == "") {
        return fmt.Errorf("publishable external apply capability review trend signal LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewTrendSignalLSP(input JEVExternalApplyCapabilityReviewTrendSignal) JEVExternalApplyCapabilityReviewTrendSignalLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewTrendSignalLSPDiagnostic{
        Status:         input.Status,
        MissingStage:   input.MissingStage,
        Signal:         input.Signal,
        MetricName:     input.MetricName,
        TrendDigest:    input.TrendDigest,
        MeanValue:      input.MeanValue,
        Threshold:      input.Threshold,
        EvidenceDigest: input.EvidenceDigest,
        SignalDigest:   input.SignalDigest,
        NonExecuting:   input.NonExecuting,
        NonAuthorizing: input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV trend signal is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Signal {
    case jevExternalApplyCapabilityReviewTrendSignalImprove:
        output.Severity = "info"
        output.Code = "jev.external-apply.trend.improve"
        output.Message = "Evidence supports an improvement direction for external review; no automatic change occurred"
    case jevExternalApplyCapabilityReviewTrendSignalHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.trend.hold"
        output.Message = "Evidence is inconclusive for change; external review remains required"
    case jevExternalApplyCapabilityReviewTrendSignalRollback:
        output.Severity = "warning"
        output.Code = "jev.external-apply.trend.rollback"
        output.Message = "Evidence supports rollback review; no automatic change occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV trend signal is UNKNOWN; evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
