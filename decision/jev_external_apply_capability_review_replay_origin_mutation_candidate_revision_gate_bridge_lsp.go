package decision

import "fmt"

type JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSPDiagnostic struct {
    Severity                    string
    Code                       string
    Message                    string
    Status                     string
    MissingStage               string
    CandidateStatus             string
    CandidateGenerationDigest  string
    GateStatus                  string
    GateDecision                string
    GateDigest                  string
    BridgeDigest               string
    Publishable                 bool
    NonExecuting                bool
    NonAuthorizing              bool
}

func (d JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete JEV mutation candidate revision gate bridge LSP diagnostic")
    }
    if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold &&
        d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected {
        return fmt.Errorf("invalid mutation candidate revision gate bridge LSP status")
    }
    if d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded &&
        d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated &&
        d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold &&
        d.CandidateStatus != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        return fmt.Errorf("invalid mutation candidate generation status")
    }
    if d.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateNotNeeded {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded || d.GateStatus != "" {
            return fmt.Errorf("candidate-not-needed LSP diagnostic requires empty gate status")
        }
    }
    if d.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateGenerated {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady ||
            d.GateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated ||
            d.GateDecision != jevExternalApplyCapabilityReviewRevisionCandidateReady ||
            d.GateDigest == "" {
            return fmt.Errorf("generated candidate LSP diagnostic requires ready gate evidence")
        }
    }
    if d.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateHold {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold ||
            d.GateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated ||
            d.GateDecision != jevExternalApplyCapabilityReviewRevisionCandidateHold ||
            d.GateDigest == "" {
            return fmt.Errorf("held candidate LSP diagnostic requires hold gate evidence")
        }
    }
    if d.CandidateStatus == jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRejected {
        if d.Status != jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected ||
            d.GateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated ||
            d.GateDecision != jevExternalApplyCapabilityReviewRevisionCandidateRejected ||
            d.GateDigest == "" {
            return fmt.Errorf("rejected candidate LSP diagnostic requires rejected gate evidence")
        }
    }
    if d.CandidateGenerationDigest == "" || d.BridgeDigest == "" {
        return fmt.Errorf("mutation candidate revision gate bridge LSP diagnostic requires digest evidence")
    }
    if !d.NonExecuting {
        return fmt.Errorf("mutation candidate revision gate bridge LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("mutation candidate revision gate bridge LSP diagnostic must be non-authorizing")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSP(input JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridge) JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateBridgeLSPDiagnostic{
        Status:                    input.Status,
        MissingStage:              input.MissingStage,
        CandidateStatus:            input.CandidateStatus,
        CandidateGenerationDigest: input.CandidateGenerationDigest,
        GateStatus:                 input.Gate.Status,
        GateDecision:               input.Gate.Decision,
        GateDigest:                 input.Gate.GateDigest,
        BridgeDigest:               input.BridgeDigest,
        NonExecuting:              input.NonExecuting,
        NonAuthorizing:            input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV candidate revision gate bridge is not publishable: candidate or gate evidence is incomplete"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateNotNeeded:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-candidate-revision-gate-not-needed"
        output.Message = "No revision candidate gate was needed; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateReady:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-candidate-revision-ready"
        output.Message = "Generated candidate passed the revision gate for review; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateHold:
        output.Severity = "info"
        output.Code = "jev.external-apply.mutation-candidate-revision-hold"
        output.Message = "Generated candidate is held by the revision gate; no execution or authorization occurred"
    case jevExternalApplyCapabilityReviewReplayOriginMutationCandidateRevisionGateRejected:
        output.Severity = "warning"
        output.Code = "jev.external-apply.mutation-candidate-revision-rejected"
        output.Message = "Generated candidate was rejected by the revision gate; no execution or authorization occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV candidate revision gate bridge is UNKNOWN; evidence must be resolved"
        output.Publishable = false
    }
    return output
}
