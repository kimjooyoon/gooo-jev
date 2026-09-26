package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSPDiagnostic struct {
    Severity             string
    Code                 string
    Message              string
    Status               string
    MissingStage         string
    ConsistencyStatus    string
    DeclaredOriginDigest string
    ReverseOriginDigest  string
    ProposedOriginDigest string
    Reason               string
    ConsistencyDigest    string
    MutationProposalDigest string
    Publishable          bool
    NonExecuting         bool
    NonAuthorizing       bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability review replay origin mutation proposal LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationProposed {
        return fmt.Errorf("invalid external apply capability review replay origin mutation proposal LSP status")
    }
    if d.ConsistencyStatus != jevExternalApplyCapabilityReviewReplayOriginConsistent &&
        d.ConsistencyStatus != jevExternalApplyCapabilityReviewReplayOriginMismatch {
        return fmt.Errorf("invalid external apply capability review replay origin mutation proposal LSP consistency status")
    }
    if d.ConsistencyStatus == jevExternalApplyCapabilityReviewReplayOriginConsistent {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded ||
            d.ProposedOriginDigest != "" ||
            d.Reason != "" {
            return fmt.Errorf("consistent mutation-not-needed LSP diagnostic contains a proposal")
        }
    }
    if d.ConsistencyStatus == jevExternalApplyCapabilityReviewReplayOriginMismatch {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationProposed ||
            d.ProposedOriginDigest == "" ||
            d.Reason == "" {
            return fmt.Errorf("mismatch mutation proposal LSP diagnostic requires proposed digest and reason")
        }
    }
    if d.DeclaredOriginDigest == "" || d.ReverseOriginDigest == "" || d.ConsistencyDigest == "" || d.MutationProposalDigest == "" {
        return fmt.Errorf("origin mutation proposal LSP diagnostic requires complete evidence")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability review replay origin mutation proposal LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability review replay origin mutation proposal LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationProposal) JEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationProposalLSPDiagnostic{
        Status:                 input.Status,
        MissingStage:           input.MissingStage,
        ConsistencyStatus:      input.ConsistencyStatus,
        DeclaredOriginDigest:   input.DeclaredOriginDigest,
        ReverseOriginDigest:    input.ReverseOriginDigest,
        ProposedOriginDigest:   input.ProposedOriginDigest,
        Reason:                 input.Reason,
        ConsistencyDigest:      input.ConsistencyDigest,
        MutationProposalDigest: input.MutationProposalDigest,
        NonExecuting:           input.NonExecuting,
        NonAuthorizing:         input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV origin mutation proposal is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.origin-mutation-not-needed"
        output.Message = "Origin evidence is consistent; no mutation proposal was generated"
    case jevExternalApplyCapabilityReviewReplayOriginMutationProposed:
        output.Severity = "warning"
        output.Code = "jev.external-apply.origin-mutation-proposed"
        output.Message = "Origin mutation is proposed for review; no automatic change or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV origin mutation proposal is UNKNOWN; evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
