package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    ReviewGateStatus          string
    Decision                  string
    PlanDigest                string
    GateDigest                string
    ApplicationTarget         string
    ApplicationSource         string
    ApplicationEvidenceDigest string
    CandidateDigest           string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV revision application candidate LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV revision application candidate LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV revision application candidate LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded:
        if d.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateNotNeeded ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewNotNeeded ||
            d.ApplicationTarget != "" || d.ApplicationSource != "" ||
            d.ApplicationEvidenceDigest != "" || d.CandidateDigest == "" {
            return fmt.Errorf("not-needed application candidate LSP diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady:
        if d.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateReady ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewApprove ||
            d.PlanDigest == "" || d.GateDigest == "" || d.ApplicationTarget == "" ||
            d.ApplicationSource == "" || d.ApplicationEvidenceDigest == "" ||
            d.CandidateDigest == "" {
            return fmt.Errorf("ready application candidate LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold:
        if d.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateHold ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewHold ||
            d.PlanDigest == "" || d.GateDigest == "" || d.CandidateDigest == "" {
            return fmt.Errorf("held application candidate LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected:
        if d.ReviewGateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewGateRejected ||
            d.Decision != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReviewReject ||
            d.PlanDigest == "" || d.GateDigest == "" || d.CandidateDigest == "" {
            return fmt.Errorf("rejected application candidate LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown:
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN application candidate LSP diagnostic must remain non-publishable")
        }
    default:
        return fmt.Errorf("invalid JEV revision application candidate LSP status")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidate) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        ReviewGateStatus:          input.ReviewGateStatus,
        Decision:                  input.Decision,
        PlanDigest:                input.PlanDigest,
        GateDigest:                input.GateDigest,
        ApplicationTarget:         input.ApplicationTarget,
        ApplicationSource:         input.ApplicationSource,
        ApplicationEvidenceDigest: input.ApplicationEvidenceDigest,
        CandidateDigest:           input.CandidateDigest,
        NonExecuting:              input.NonExecuting,
        NonAuthorizing:            input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown
        if output.MissingStage == "" {
            output.MissingStage = "application-candidate-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision application candidate is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.application-candidate-not-needed"
        output.Message = "No application candidate was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.application-candidate-ready"
        output.Message = "Application candidate is ready for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.application-candidate-hold"
        output.Message = "Application candidate is held; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.application-candidate-rejected"
        output.Message = "Application candidate was rejected; no execution or authorization occurred"
    default:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionApplicationCandidateUnknown
        output.MissingStage = "application-candidate-status"
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision application candidate status is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
