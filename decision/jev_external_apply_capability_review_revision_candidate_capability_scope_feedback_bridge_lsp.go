package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strconv"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPBoundCode = "jev.provenance.capability-scope-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic struct {
    Severity                         string
    Code                             string
    Message                          string
    Status                           string
    MissingStage                     string
    CandidateDigest                  string
    CandidateSource                  string
    CandidateGateDigest              string
    ApplicationPlanStatus            string
    PlanDigest                      string
    PlanSource                      string
    PlanEvidenceDigest              string
    ApplicationPlanBridgeDigest     string
    ScopeStatus                     string
    SpiffeID                        string
    Audience                        string
    SandboxPolicyDigest             string
    NetworkAllowlistDigest          string
    ScopeEvidenceDigest             string
    CapabilityScopeBridgeDigest     string
    ReverseObservationStatus        string
    ObservationDigest               string
    ObservationSource               string
    ObservationEvidenceDigest       string
    ObservationMetricStatus         string
    ObservationMetricDigest         string
    ObservationMetricSource         string
    ObservationMetricEvidenceDigest string
    ReverseObservationBridgeDigest  string
    FeedbackDirection               string
    FeedbackDigest                  string
    FeedbackSource                  string
    FeedbackEvidenceDigest          string
    FeedbackStatus                  string
    BridgeDigest                    string
    ProjectionDigest                string
    Publishable                     bool
    NonExecuting                    bool
    NonAuthorizing                  bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("incomplete JEV capability scope feedback LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV capability scope feedback LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV capability scope feedback LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown:
        if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPUnknownCode ||
            d.MissingStage == "" || d.Publishable {
            return fmt.Errorf("unknown capability scope feedback LSP diagnostic is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeBound:
        if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPBoundCode ||
            !d.Publishable || d.BridgeDigest == "" || d.ApplicationPlanBridgeDigest == "" ||
            d.CapabilityScopeBridgeDigest == "" || d.ReverseObservationBridgeDigest == "" ||
            d.ApplicationPlanStatus == "" || d.ScopeStatus == "" ||
            d.ReverseObservationStatus == "" || d.FeedbackStatus == "" {
            return fmt.Errorf("bound capability scope feedback LSP diagnostic is incomplete")
        }
        if d.ApplicationPlanStatus == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady &&
            (d.CandidateDigest == "" || d.CandidateSource == "" || d.CandidateGateDigest == "" ||
                d.PlanDigest == "" || d.PlanSource == "" || d.PlanEvidenceDigest == "" ||
                d.SpiffeID == "" || d.Audience == "" || d.SandboxPolicyDigest == "" ||
                d.NetworkAllowlistDigest == "" || d.ScopeEvidenceDigest == "" ||
                d.ObservationDigest == "" || d.ObservationSource == "" ||
                d.ObservationEvidenceDigest == "" || d.ObservationMetricDigest == "" ||
                d.ObservationMetricSource == "" || d.ObservationMetricEvidenceDigest == "" ||
                d.FeedbackDirection == "" || d.FeedbackDigest == "" ||
                d.FeedbackSource == "" || d.FeedbackEvidenceDigest == "") {
            return fmt.Errorf("ready capability scope feedback LSP diagnostic lost evidence")
        }
    default:
        return fmt.Errorf("invalid capability scope feedback LSP status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic(
        d.Severity, d.Code, d.Message, d.Status, d.MissingStage,
        d.CandidateDigest, d.CandidateSource, d.CandidateGateDigest,
        d.ApplicationPlanStatus, d.PlanDigest, d.PlanSource, d.PlanEvidenceDigest,
        d.ApplicationPlanBridgeDigest, d.ScopeStatus, d.SpiffeID, d.Audience,
        d.SandboxPolicyDigest, d.NetworkAllowlistDigest, d.ScopeEvidenceDigest,
        d.CapabilityScopeBridgeDigest, d.ReverseObservationStatus, d.ObservationDigest,
        d.ObservationSource, d.ObservationEvidenceDigest, d.ObservationMetricStatus,
        d.ObservationMetricDigest, d.ObservationMetricSource, d.ObservationMetricEvidenceDigest,
        d.ReverseObservationBridgeDigest, d.FeedbackDirection, d.FeedbackDigest,
        d.FeedbackSource, d.FeedbackEvidenceDigest, d.FeedbackStatus, d.BridgeDigest,
        d.Publishable,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV capability scope feedback LSP projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridge) JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic{
        Status:                          input.Status,
        MissingStage:                    input.MissingStage,
        CandidateDigest:                 input.CandidateDigest,
        CandidateSource:                 input.CandidateSource,
        CandidateGateDigest:             input.CandidateGateDigest,
        ApplicationPlanStatus:           input.ApplicationPlanStatus,
        PlanDigest:                       input.PlanDigest,
        PlanSource:                       input.PlanSource,
        PlanEvidenceDigest:              input.PlanEvidenceDigest,
        ApplicationPlanBridgeDigest:     input.ApplicationPlanBridgeDigest,
        ScopeStatus:                     input.ScopeStatus,
        SpiffeID:                         input.SpiffeID,
        Audience:                        input.Audience,
        SandboxPolicyDigest:              input.SandboxPolicyDigest,
        NetworkAllowlistDigest:           input.NetworkAllowlistDigest,
        ScopeEvidenceDigest:              input.ScopeEvidenceDigest,
        CapabilityScopeBridgeDigest:      input.CapabilityScopeBridgeDigest,
        ReverseObservationStatus:         input.ReverseObservationStatus,
        ObservationDigest:                input.ObservationDigest,
        ObservationSource:                input.ObservationSource,
        ObservationEvidenceDigest:        input.ObservationEvidenceDigest,
        ObservationMetricStatus:          input.ObservationMetricStatus,
        ObservationMetricDigest:           input.ObservationMetricDigest,
        ObservationMetricSource:           input.ObservationMetricSource,
        ObservationMetricEvidenceDigest:   input.ObservationMetricEvidenceDigest,
        ReverseObservationBridgeDigest:   input.ReverseObservationBridgeDigest,
        FeedbackDirection:                 input.FeedbackDirection,
        FeedbackDigest:                    input.FeedbackDigest,
        FeedbackSource:                    input.FeedbackSource,
        FeedbackEvidenceDigest:            input.FeedbackEvidenceDigest,
        FeedbackStatus:                    input.FeedbackStatus,
        BridgeDigest:                      input.BridgeDigest,
        NonExecuting:                      true,
        NonAuthorizing:                    true,
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown || input.MissingStage != "" {
        output.Severity = "warning"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPUnknownCode
        output.Message = "capability scope candidate feedback provenance is incomplete"
        output.Publishable = false
    } else {
        output.Severity = "info"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPBoundCode
        output.Message = "capability scope candidate feedback provenance is bound"
        output.Publishable = true
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic(
        output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
        output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
        output.ApplicationPlanStatus, output.PlanDigest, output.PlanSource, output.PlanEvidenceDigest,
        output.ApplicationPlanBridgeDigest, output.ScopeStatus, output.SpiffeID, output.Audience,
        output.SandboxPolicyDigest, output.NetworkAllowlistDigest, output.ScopeEvidenceDigest,
        output.CapabilityScopeBridgeDigest, output.ReverseObservationStatus, output.ObservationDigest,
        output.ObservationSource, output.ObservationEvidenceDigest, output.ObservationMetricStatus,
        output.ObservationMetricDigest, output.ObservationMetricSource, output.ObservationMetricEvidenceDigest,
        output.ReverseObservationBridgeDigest, output.FeedbackDirection, output.FeedbackDigest,
        output.FeedbackSource, output.FeedbackEvidenceDigest, output.FeedbackStatus, output.BridgeDigest,
        output.Publishable,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeUnknown
        output.MissingStage = "lsp-projection"
        output.Publishable = false
        output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic(
            output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
            output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
            output.ApplicationPlanStatus, output.PlanDigest, output.PlanSource, output.PlanEvidenceDigest,
            output.ApplicationPlanBridgeDigest, output.ScopeStatus, output.SpiffeID, output.Audience,
            output.SandboxPolicyDigest, output.NetworkAllowlistDigest, output.ScopeEvidenceDigest,
            output.CapabilityScopeBridgeDigest, output.ReverseObservationStatus, output.ObservationDigest,
            output.ObservationSource, output.ObservationEvidenceDigest, output.ObservationMetricStatus,
            output.ObservationMetricDigest, output.ObservationMetricSource, output.ObservationMetricEvidenceDigest,
            output.ReverseObservationBridgeDigest, output.FeedbackDirection, output.FeedbackDigest,
            output.FeedbackSource, output.FeedbackEvidenceDigest, output.FeedbackStatus, output.BridgeDigest,
            output.Publishable,
        )
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackBridgeLSPDiagnostic(severity, code, message, status, missingStage, candidateDigest, candidateSource, candidateGateDigest, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, scopeStatus, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest, scopeEvidenceDigest, capabilityScopeBridgeDigest, reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus, observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest, reverseObservationBridgeDigest, feedbackDirection, feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackStatus, bridgeDigest string, publishable bool) string {
    values := []string{
        severity, code, message, status, missingStage, candidateDigest, candidateSource,
        candidateGateDigest, applicationPlanStatus, planDigest, planSource, planEvidenceDigest,
        applicationPlanBridgeDigest, scopeStatus, spiffeID, audience, sandboxPolicyDigest,
        networkAllowlistDigest, scopeEvidenceDigest, capabilityScopeBridgeDigest,
        reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest,
        observationMetricStatus, observationMetricDigest, observationMetricSource,
        observationMetricEvidenceDigest, reverseObservationBridgeDigest, feedbackDirection,
        feedbackDigest, feedbackSource, feedbackEvidenceDigest, feedbackStatus, bridgeDigest,
        strconv.FormatBool(publishable),
    }
    sum := sha256.Sum256([]byte(strings.Join(values, "|")))
    return hex.EncodeToString(sum[:])
}
