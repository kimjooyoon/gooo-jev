package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSPDiagnostic struct {
    Severity                  string
    Code                      string
    Message                   string
    Status                    string
    MissingStage              string
    BridgeStatus              string
    CandidateGenerationDigest string
    GateDigest                string
    CandidateDigest           string
    RevisionSource            string
    RevisionPlanSource        string
    RevisionIntentDigest     string
    ProposalDigest            string
    Publishable               bool
    NonExecuting              bool
    NonAuthorizing            bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV revision proposal LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV revision proposal LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV revision proposal LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded:
        if d.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded ||
            d.RevisionPlanSource != "" ||
            d.RevisionIntentDigest != "" {
            return fmt.Errorf("not-needed revision proposal LSP diagnostic has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady:
        if d.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady ||
            d.CandidateGenerationDigest == "" ||
            d.GateDigest == "" ||
            d.CandidateDigest == "" ||
            d.RevisionSource == "" ||
            d.RevisionPlanSource == "" ||
            d.RevisionIntentDigest == "" ||
            d.ProposalDigest == "" {
            return fmt.Errorf("ready revision proposal LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold:
        if d.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold ||
            d.ProposalDigest == "" {
            return fmt.Errorf("held revision proposal LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected:
        if d.BridgeStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected ||
            d.ProposalDigest == "" {
            return fmt.Errorf("rejected revision proposal LSP diagnostic is incomplete")
        }
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown:
        if d.Publishable || d.MissingStage == "" {
            return fmt.Errorf("UNKNOWN revision proposal LSP diagnostic must remain non-publishable")
        }
    default:
        return fmt.Errorf("invalid JEV revision proposal LSP status")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown &&
        (d.CandidateGenerationDigest == "" || d.ProposalDigest == "") {
        return fmt.Errorf("revision proposal LSP diagnostic requires digest evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposal) JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        BridgeStatus:              input.BridgeStatus,
        CandidateGenerationDigest: input.CandidateGenerationDigest,
        GateDigest:                input.GateDigest,
        CandidateDigest:           input.CandidateDigest,
        RevisionSource:            input.RevisionSource,
        RevisionPlanSource:        input.RevisionPlanSource,
        RevisionIntentDigest:      input.RevisionIntentDigest,
        ProposalDigest:            input.ProposalDigest,
        NonExecuting:              input.NonExecuting,
        NonAuthorizing:            input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown
        if output.MissingStage == "" {
            output.MissingStage = "revision-proposal-evidence"
        }
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision proposal is not publishable; evidence must be resolved"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-proposal-not-needed"
        output.Message = "No revision proposal was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-proposal-ready"
        output.Message = "Revision proposal is ready for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.revision-proposal-hold"
        output.Message = "Revision proposal is held for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.revision-proposal-rejected"
        output.Message = "Revision proposal was rejected; no execution or authorization occurred"
    default:
        output.Status = jevExternalApplyCapabilityReviewReplayOriginMutationRevisionProposalUnknown
        output.MissingStage = "revision-proposal-status"
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV revision proposal status is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
