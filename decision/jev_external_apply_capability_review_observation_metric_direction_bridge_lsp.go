package decision

import "fmt"

type JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSPDiagnostic struct {
    Severity           string
    Code               string
    Message            string
    Status             string
    MissingStage       string
    ObservationStatus  string
    ObservationDigest  string
    MetricName         string
    MetricValue        float64
    MetricDigest       string
    Direction          string
    Target             string
    CandidateSource    string
    FeedbackDigest     string
    DirectionDigest    string
    BridgeDigest       string
    Publishable        bool
    NonExecuting       bool
    NonAuthorizing     bool
}

func (d JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV observation metric direction bridge LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV observation metric direction bridge LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV observation metric direction bridge LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN observation metric direction bridge diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeBound {
        return fmt.Errorf("invalid JEV observation metric direction bridge LSP status")
    }
    if !d.Publishable || d.ObservationStatus == "" || d.ObservationDigest == "" ||
        d.MetricName == "" || d.MetricDigest == "" || d.Direction == "" ||
        d.Target == "" || d.CandidateSource == "" || d.FeedbackDigest == "" ||
        d.DirectionDigest == "" || d.BridgeDigest == "" {
        return fmt.Errorf("bound observation metric direction bridge diagnostic is incomplete")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSP(input JEVExternalApplyCapabilityReviewObservationMetricDirectionBridge) JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewObservationMetricDirectionBridgeLSPDiagnostic{
        Status:            input.Status,
        MissingStage:      input.MissingStage,
        ObservationStatus: input.ObservationStatus,
        ObservationDigest: input.ObservationDigest,
        MetricName:        input.MetricName,
        MetricValue:       input.MetricValue,
        MetricDigest:      input.MetricDigest,
        Direction:         input.Direction,
        Target:            input.Target,
        CandidateSource:   input.CandidateSource,
        FeedbackDigest:    input.FeedbackDigest,
        DirectionDigest:   input.DirectionDigest,
        BridgeDigest:      input.BridgeDigest,
        NonExecuting:      input.NonExecuting,
        NonAuthorizing:    input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewObservationMetricDirectionBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "observation-metric-direction-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV observation metric direction bridge is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.observation-metric-direction-bound"
    output.Message = "Observation, metric, and improvement direction are provenance-bound; no execution or authorization occurred"
    return output
}
