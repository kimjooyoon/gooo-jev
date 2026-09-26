package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    CandidateApplicationStatus string
    CandidateApplicationDigest string
    ObservationStatus         string
    ObservationDigest          string
    MetricName                string
    MetricValue               float64
    MetricDigest               string
    BridgeDigest              string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate observation metric LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate observation metric LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate observation metric LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV candidate observation metric diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeBound {
        return fmt.Errorf("invalid JEV candidate observation metric LSP status")
    }
    if !d.Publishable || d.CandidateApplicationStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady ||
        d.CandidateApplicationDigest == "" || d.ObservationStatus == "" ||
        d.ObservationDigest == "" || d.MetricName == "" || d.MetricDigest == "" ||
        d.BridgeDigest == "" {
        return fmt.Errorf("bound JEV candidate observation metric diagnostic is incomplete")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        CandidateApplicationStatus: input.CandidateApplicationStatus,
        CandidateApplicationDigest: input.CandidateApplicationDigest,
        ObservationStatus:          input.ObservationStatus,
        ObservationDigest:          input.ObservationDigest,
        MetricName:                 input.MetricName,
        MetricValue:                input.MetricValue,
        MetricDigest:               input.MetricDigest,
        BridgeDigest:               input.BridgeDigest,
        Publishable:                false,
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMetricBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "candidate-observation-metric-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "Candidate observation metric is UNKNOWN; evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-observation-metric-bound"
    output.Message = "Candidate reverse observation is bound to metric evidence without execution or authorization"
    return output
}
