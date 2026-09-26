package decision

import "fmt"

type JEVExternalApplyCapabilityBoundaryLSPDiagnostic struct {
    Severity                 string
    Code                     string
    Message                  string
    Status                   string
    MissingStage             string
    BindingDigest            string
    Principal                string
    Audience                 string
    Workspace                string
    NetworkAllowlistDigest   string
    ScopeDigest              string
    CapabilityDigest         string
    Publishable              bool
    NonExecuting             bool
    NonAuthorizing           bool
}

func (d JEVExternalApplyCapabilityBoundaryLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" {
        return fmt.Errorf("incomplete external apply capability boundary LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("external apply capability boundary LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("external apply capability boundary LSP diagnostic must be non-authorizing")
    }
    if d.Publishable && (d.BindingDigest == "" || d.Principal == "" || d.Audience == "" || d.Workspace == "" || d.NetworkAllowlistDigest == "" || d.ScopeDigest == "" || d.CapabilityDigest == "") {
        return fmt.Errorf("publishable external apply capability boundary LSP diagnostic requires evidence")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityBoundaryLSP(input JEVExternalApplyCapabilityBoundary) JEVExternalApplyCapabilityBoundaryLSPDiagnostic {
    output := JEVExternalApplyCapabilityBoundaryLSPDiagnostic{
        Status:                 input.Status,
        MissingStage:           input.MissingStage,
        BindingDigest:          input.BindingDigest,
        Principal:              input.Principal,
        Audience:               input.Audience,
        Workspace:              input.Workspace,
        NetworkAllowlistDigest: input.NetworkAllowlistDigest,
        ScopeDigest:            input.ScopeDigest,
        CapabilityDigest:       input.CapabilityDigest,
        NonExecuting:           input.NonExecuting,
        NonAuthorizing:         input.NonAuthorizing,
    }
    if err := input.Validate(); err != nil {
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV external apply capability boundary is not publishable: evidence is incomplete or invalid"
        output.Publishable = false
        return output
    }
    output.Publishable = true
    switch input.Status {
    case jevExternalApplyCapabilityDescribed:
        output.Severity = "info"
        output.Code = "jev.external-apply.capability-described"
        output.Message = "Capability boundary is described for external review; no authorization or execution occurred"
    default:
        output.Severity = "error"
        output.Code = "jev.provenance.unknown"
        output.Message = "JEV external apply capability boundary is UNKNOWN; evidence must be resolved before review"
        output.Publishable = false
    }
    return output
}
