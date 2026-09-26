package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSPDiagnostic struct {
    Severity                    string
    Code                        string
    Message                     string
    Status                     string
    MissingStage                string
    CandidateApplicationStatus  string
    CandidateApplicationDigest  string
    ObservationStatus           string
    ObservationDigest           string
    MetricName                  string
    MetricValue                 float64
    MetricDigest                string
    Direction                   string
    Target                      string
    CandidateSource             string
    FeedbackDigest              string
    DirectionDigest             string
    BridgeDigest               string
    Publishable                 bool
    NonExecuting                bool
    NonAuthorizing              bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate observation metric direction LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate observation metric direction LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate observation metric direction LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV candidate observation direction diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeBound {
        return fmt.Errorf("invalid JEV candidate observation metric direction LSP status")
    }
    if !d.Publishable || d.CandidateApplicationDigest == "" ||
        d.ObservationDigest == "" || d.MetricDigest == "" ||
        d.Direction == "" || d.DirectionDigest == "" ||
        d.BridgeDigest == "" {
        return fmt.Errorf("bound JEV candidate observation metric direction diagnostic is incomplete")
    }
    if d.ObservationStatus == jevExternalApplyCapabilityReviewRevisionCandidateApplicationObservationMismatch &&
        d.Direction == jevExternalApplyCapabilityReviewImprovementDirectionGenerate {
        return fmt.Errorf("mismatched candidate observation was exposed as generation")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        CandidateApplicationStatus: input.CandidateApplicationStatus,
        CandidateApplicationDigest: input.CandidateApplicationDigest,
        ObservationStatus:           input.ObservationStatus,
        ObservationDigest:           input.ObservationDigest,
        MetricName:                 input.MetricName,
        MetricValue:                input.MetricValue,
        MetricDigest:               input.MetricDigest,
        Direction:                  input.Direction,
        Target:                    input.Target,
        CandidateSource:           input.CandidateSource,
        FeedbackDigest:            input.FeedbackDigest,
        DirectionDigest:           input.DirectionDigest,
        BridgeDigest:               input.BridgeDigest,
        Publishable:               false,
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "candidate-observation-metric-direction-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "Candidate observation metric direction is UNKNOWN; evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-observation-metric-direction-bound"
    output.Message = "Candidate observation metric is bound to an improvement direction without execution or authorization"
    if input.Direction == jevExternalApplyCapabilityReviewImprovementDirectionReject {
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-observation-metric-direction-rejected"
        output.Message = "Candidate observation metric direction rejected further generation without execution or authorization"
    }
    return output
}
