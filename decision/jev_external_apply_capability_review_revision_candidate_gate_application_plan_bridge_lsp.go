package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPDiagnostic struct {
    Severity                string
    Code                    string
    Message                 string
    Status                  string
    MissingStage            string
    CandidateGateStatus     string
    CandidateDecision       string
    CandidateDigest         string
    CandidateSource         string
    CandidateGateDigest     string
    ApplicationPlanStatus   string
    PlanDigest              string
    PlanSource              string
    PlanEvidenceDigest      string
    BridgeDigest            string
    Publishable             bool
    NonExecuting            bool
    NonAuthorizing          bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate gate application plan LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate gate application plan LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate gate application plan LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV application plan diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeBound {
        return fmt.Errorf("invalid JEV candidate gate application plan LSP status")
    }
    if !d.Publishable || d.CandidateGateStatus == "" ||
        d.CandidateDecision == "" || d.CandidateGateDigest == "" ||
        d.ApplicationPlanStatus == "" || d.BridgeDigest == "" {
        return fmt.Errorf("bound JEV application plan diagnostic is incomplete")
    }
    switch d.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if d.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady ||
            d.CandidateDigest == "" || d.CandidateSource == "" ||
            d.PlanDigest == "" || d.PlanSource == "" ||
            d.PlanEvidenceDigest == "" {
            return fmt.Errorf("ready JEV application plan diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if d.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanHeld ||
            d.PlanDigest != "" || d.PlanSource != "" || d.PlanEvidenceDigest != "" {
            return fmt.Errorf("held JEV application plan diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if d.ApplicationPlanStatus != jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanRejected ||
            d.PlanDigest != "" || d.PlanSource != "" || d.PlanEvidenceDigest != "" {
            return fmt.Errorf("rejected JEV application plan diagnostic has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid candidate gate decision")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeLSPDiagnostic{
        Status:                input.Status,
        MissingStage:          input.MissingStage,
        CandidateGateStatus:   input.CandidateGateStatus,
        CandidateDecision:     input.CandidateDecision,
        CandidateDigest:       input.CandidateDigest,
        CandidateSource:       input.CandidateSource,
        CandidateGateDigest:   input.CandidateGateDigest,
        ApplicationPlanStatus: input.ApplicationPlanStatus,
        PlanDigest:            input.PlanDigest,
        PlanSource:            input.PlanSource,
        PlanEvidenceDigest:    input.PlanEvidenceDigest,
        BridgeDigest:          input.BridgeDigest,
        Publishable:           false,
        NonExecuting:          true,
        NonAuthorizing:        true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "candidate-gate-application-plan-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "Candidate gate application plan evidence is UNKNOWN; evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-gate-application-plan-bound"
    output.Message = "Candidate gate is bound to an application plan without execution or authorization"
    if input.CandidateDecision == jevExternalApplyCapabilityReviewRevisionCandidateRejected {
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-gate-application-plan-rejected"
        output.Message = "Candidate gate rejected the application plan without execution or authorization"
    }
    return output
}
