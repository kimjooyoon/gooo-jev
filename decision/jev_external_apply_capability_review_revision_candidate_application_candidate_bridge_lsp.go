package decision

import "fmt"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    CandidateGateStatus       string
    CandidateDecision         string
    CandidateDigest           string
    RevisionSource            string
    CandidateGateDigest       string
    ApplicationCandidateStatus string
    ApplicationTarget         string
    ApplicationSource         string
    ApplicationEvidenceDigest string
    BridgeDigest              string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV candidate application candidate LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV candidate application candidate LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV candidate application candidate LSP diagnostic must be non-authorizing")
    }
    if d.Status == jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown {
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN JEV candidate application candidate diagnostic must remain non-publishable")
        }
        return nil
    }
    if d.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound {
        return fmt.Errorf("invalid JEV candidate application candidate LSP status")
    }
    if d.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated ||
        d.CandidateDecision == "" || d.CandidateGateDigest == "" ||
        d.ApplicationCandidateStatus == "" || d.BridgeDigest == "" || !d.Publishable {
        return fmt.Errorf("bound JEV candidate application candidate diagnostic is incomplete")
    }
    switch d.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if d.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady ||
            d.CandidateDigest == "" || d.RevisionSource == "" || d.ApplicationTarget == "" ||
            d.ApplicationSource == "" || d.ApplicationEvidenceDigest == "" {
            return fmt.Errorf("ready JEV candidate application candidate diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if d.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold ||
            d.ApplicationTarget != "" || d.ApplicationSource != "" || d.ApplicationEvidenceDigest != "" {
            return fmt.Errorf("held JEV candidate application candidate diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if d.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected ||
            d.ApplicationTarget != "" || d.ApplicationSource != "" || d.ApplicationEvidenceDigest != "" {
            return fmt.Errorf("rejected JEV candidate application candidate diagnostic has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid JEV candidate gate decision")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSP(
    input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge,
) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        CandidateGateStatus:       input.CandidateGateStatus,
        CandidateDecision:         input.CandidateDecision,
        CandidateDigest:           input.CandidateDigest,
        RevisionSource:             input.RevisionSource,
        CandidateGateDigest:        input.CandidateGateDigest,
        ApplicationCandidateStatus: input.ApplicationCandidateStatus,
        ApplicationTarget:         input.ApplicationTarget,
        ApplicationSource:          input.ApplicationSource,
        ApplicationEvidenceDigest: input.ApplicationEvidenceDigest,
        BridgeDigest:              input.BridgeDigest,
        Publishable:               false,
        NonExecuting:              true,
        NonAuthorizing:            true,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown
        if output.MissingStage == "" {
            output.MissingStage = "candidate-application-candidate-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV candidate application candidate is UNKNOWN; evidence must be resolved"
        return output
    }
    output.Publishable = true
    output.Severity = "info"
    output.Code = "jev.external-apply.candidate-application-candidate-bound"
    output.Message = "Candidate gate evidence is bound to an application candidate without execution or authorization"
    if input.ApplicationCandidateStatus == jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected {
        output.Severity = "warning"
        output.Code = "jev.external-apply.candidate-application-candidate-rejected"
        output.Message = "Candidate gate rejected the application candidate without execution or authorization"
    }
    return output
}
