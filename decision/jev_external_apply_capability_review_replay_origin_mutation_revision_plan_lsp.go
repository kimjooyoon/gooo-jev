package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPDiagnostic struct {
    Severity        string
    Code            string
    Message         string
    Status          string
    MissingStage    string
    ProposalStatus  string
    ProposalDigest  string
    PlanSource      string
    IntentDigest    string
    PlanDigest      string
    Publishable     bool
    NonExecuting    bool
    NonAuthorizing  bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV revision plan LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV revision plan LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV revision plan LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded:
        if d.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded ||
            d.PlanSource != "" || d.IntentDigest != "" {
            return fmt.Errorf("not-needed revision plan LSP diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady:
        if d.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady ||
            d.ProposalDigest == "" || d.PlanSource == "" ||
            d.IntentDigest == "" || d.PlanDigest == "" {
            return fmt.Errorf("ready revision plan LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold:
        if d.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold ||
            d.ProposalDigest == "" || d.PlanDigest == "" {
            return fmt.Errorf("held revision plan LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected:
        if d.ProposalStatus != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected ||
            d.ProposalDigest == "" || d.PlanDigest == "" {
            return fmt.Errorf("rejected revision plan LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown:
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN revision plan LSP diagnostic must remain non-publishable")
        }
    default:
        return fmt.Errorf("invalid JEV revision plan LSP status")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlan) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanLSPDiagnostic{
        Status:         input.Status,
        MissingStage:   input.MissingStage,
        ProposalStatus: input.ProposalStatus,
        ProposalDigest: input.ProposalDigest,
        PlanSource:     input.PlanSource,
        IntentDigest:   input.IntentDigest,
        PlanDigest:     input.PlanDigest,
        NonExecuting:   input.NonExecuting,
        NonAuthorizing: input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown
        if output.MissingStage == "" {
            output.MissingStage = "revision-plan-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision plan is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-plan-not-needed"
        output.Message = "No revision plan was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-plan-ready"
        output.Message = "Revision plan is ready for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-plan-hold"
        output.Message = "Revision plan is held for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.revision-plan-rejected"
        output.Message = "Revision plan was rejected; no execution or authorization occurred"
    default:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionPlanUnknown
        output.MissingStage = "revision-plan-status"
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision plan status is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
