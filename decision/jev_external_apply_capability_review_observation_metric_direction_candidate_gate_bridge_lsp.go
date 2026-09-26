package decision

import "fmt"

type JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSPDiagnostic struct {
    Severity                          string
    Code                              string
    Message                           string
    Status                            string
    MissingStage                      string
    ObservationStatus                 string
    ObservationMetricDirectionDigest  string
    CandidateGateStatus               string
    CandidateDecision                 string
    CandidateDigest                   string
    RevisionSource                    string
    GateDigest                        string
    BridgeDigest                      string
    Publishable                       bool
    NonExecuting                      bool
    NonAuthorizing                    bool
}

func (d JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV observation metric direction candidate gate bridge LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV observation metric direction candidate gate bridge LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV observation metric direction candidate gate bridge LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN candidate gate bridge LSP diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeBound {
        return fmt.Errorf("invalid JEV observation metric direction candidate gate bridge LSP status")
    }
    if !d.Publishable || d.ObservationStatus == "" || d.ObservationMetricDirectionDigest == "" ||
        d.CandidateGateStatus == "" || d.CandidateDecision == "" ||
        d.GateDigest == "" || d.BridgeDigest == "" {
        return fmt.Errorf("bound candidate gate bridge LSP diagnostic is incomplete")
    }
    if d.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateReady &&
        (d.CandidateDigest == "" || d.RevisionSource == "") {
        return fmt.Errorf("ready candidate gate bridge LSP diagnostic requires candidate evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSP(input JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge) JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeLSPDiagnostic{
        Status:                          input.Status,
        MissingStage:                    input.MissingStage,
        ObservationStatus:                input.ObservationStatus,
        ObservationMetricDirectionDigest: input.ObservationMetricDirectionDigest,
        CandidateGateStatus:               input.CandidateGateStatus,
        CandidateDecision:                 input.CandidateDecision,
        CandidateDigest:                   input.CandidateDigest,
        RevisionSource:                   input.RevisionSource,
        GateDigest:                       input.GateDigest,
        BridgeDigest:                     input.BridgeDigest,
        NonExecuting:                     input.NonExecuting,
        NonAuthorizing:                   input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "observation-metric-direction-candidate-gate-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV candidate gate bridge is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.observation-metric-direction-candidate-gate-bound"
    output.Message = "Observation metric direction is bound to the revision candidate gate; no execution or authorization occurred"
    return output
}
