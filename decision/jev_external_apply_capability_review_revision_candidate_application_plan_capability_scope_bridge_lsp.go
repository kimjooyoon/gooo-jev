package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPBoundCode = "jev.provenance.application-plan-capability-scope-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic struct {
    Severity                   string
    Code                       string
    Message                    string
    Status                     string
    MissingStage               string
    ApplicationPlanStatus      string
    PlanDigest                 string
    PlanSource                 string
    PlanEvidenceDigest         string
    ApplicationPlanBridgeDigest string
    SpiffeID                   string
    Audience                   string
    SandboxPolicyDigest        string
    NetworkAllowlistDigest     string
    ScopeEvidenceDigest        string
    ScopeStatus                string
    BridgeDigest               string
    ProjectionDigest           string
    Publishable                bool
    NonExecuting               bool
    NonAuthorizing             bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" ||
        d.ProjectionDigest == "" {
        return fmt.Errorf("incomplete JEV capability scope LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV capability scope LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV capability scope LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown:
        if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPUnknownCode ||
            d.MissingStage == "" || d.Publishable {
            return fmt.Errorf("unknown capability scope LSP diagnostic is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeBound:
        if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPBoundCode ||
            !d.Publishable || d.BridgeDigest == "" || d.ApplicationPlanBridgeDigest == "" ||
            d.ApplicationPlanStatus == "" || d.ScopeStatus == "" {
            return fmt.Errorf("bound capability scope LSP diagnostic is incomplete")
        }
        if d.ApplicationPlanStatus == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady &&
            (d.PlanDigest == "" || d.PlanSource == "" || d.PlanEvidenceDigest == "" ||
                d.SpiffeID == "" || d.Audience == "" || d.SandboxPolicyDigest == "" ||
                d.NetworkAllowlistDigest == "" || d.ScopeEvidenceDigest == "") {
            return fmt.Errorf("ready capability scope LSP diagnostic lost evidence")
        }
    default:
        return fmt.Errorf("invalid capability scope LSP status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic(
        d.Severity,
        d.Code,
        d.Message,
        d.Status,
        d.MissingStage,
        d.ApplicationPlanStatus,
        d.PlanDigest,
        d.PlanSource,
        d.PlanEvidenceDigest,
        d.ApplicationPlanBridgeDigest,
        d.SpiffeID,
        d.Audience,
        d.SandboxPolicyDigest,
        d.NetworkAllowlistDigest,
        d.ScopeEvidenceDigest,
        d.ScopeStatus,
        d.BridgeDigest,
        d.Publishable,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV capability scope LSP projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridge) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic{
        Status:                      input.Status,
        MissingStage:                input.MissingStage,
        ApplicationPlanStatus:       input.ApplicationPlanStatus,
        PlanDigest:                  input.PlanDigest,
        PlanSource:                  input.PlanSource,
        PlanEvidenceDigest:          input.PlanEvidenceDigest,
        ApplicationPlanBridgeDigest: input.ApplicationPlanBridgeDigest,
        SpiffeID:                    input.SpiffeID,
        Audience:                    input.Audience,
        SandboxPolicyDigest:         input.SandboxPolicyDigest,
        NetworkAllowlistDigest:      input.NetworkAllowlistDigest,
        ScopeEvidenceDigest:         input.ScopeEvidenceDigest,
        ScopeStatus:                 input.ScopeStatus,
        BridgeDigest:                input.BridgeDigest,
        NonExecuting:                true,
        NonAuthorizing:              true,
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown || input.MissingStage != "" {
        output.Severity = "warning"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPUnknownCode
        output.Message = "application plan capability scope provenance is incomplete"
        output.Publishable = false
    } else {
        output.Severity = "info"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPBoundCode
        output.Message = "application plan capability scope provenance is bound"
        output.Publishable = true
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic(
        output.Severity,
        output.Code,
        output.Message,
        output.Status,
        output.MissingStage,
        output.ApplicationPlanStatus,
        output.PlanDigest,
        output.PlanSource,
        output.PlanEvidenceDigest,
        output.ApplicationPlanBridgeDigest,
        output.SpiffeID,
        output.Audience,
        output.SandboxPolicyDigest,
        output.NetworkAllowlistDigest,
        output.ScopeEvidenceDigest,
        output.ScopeStatus,
        output.BridgeDigest,
        output.Publishable,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeUnknown
        output.MissingStage = "lsp-projection"
        output.Publishable = false
        output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic(
            output.Severity,
            output.Code,
            output.Message,
            output.Status,
            output.MissingStage,
            output.ApplicationPlanStatus,
            output.PlanDigest,
            output.PlanSource,
            output.PlanEvidenceDigest,
            output.ApplicationPlanBridgeDigest,
            output.SpiffeID,
            output.Audience,
            output.SandboxPolicyDigest,
            output.NetworkAllowlistDigest,
            output.ScopeEvidenceDigest,
            output.ScopeStatus,
            output.BridgeDigest,
            output.Publishable,
        )
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanCapabilityScopeBridgeLSPDiagnostic(severity, code, message, status, missingStage, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest, scopeEvidenceDigest, scopeStatus, bridgeDigest string, publishable bool) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t", severity, code, message, status, missingStage, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, spiffeID, audience, sandboxPolicyDigest, networkAllowlistDigest, scopeEvidenceDigest, scopeStatus, bridgeDigest, publishable)))
    return hex.EncodeToString(sum[:])
}
