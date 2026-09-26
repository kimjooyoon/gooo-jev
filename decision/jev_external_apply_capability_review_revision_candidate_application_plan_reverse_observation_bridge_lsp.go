package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPBoundCode = "jev.provenance.application-plan-reverse-observation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic struct {
    Severity                         string
    Code                             string
    Message                          string
    Status                           string
    MissingStage                     string
    ApplicationPlanStatus            string
    PlanDigest                       string
    PlanSource                       string
    PlanEvidenceDigest               string
    ApplicationPlanBridgeDigest      string
    ReverseObservationStatus         string
    ObservationDigest                string
    ObservationSource                string
    ObservationEvidenceDigest        string
    ObservationMetricStatus          string
    ObservationMetricDigest          string
    ObservationMetricSource          string
    ObservationMetricEvidenceDigest  string
    BridgeDigest                     string
    ProjectionDigest                 string
    Publishable                      bool
    NonExecuting                     bool
    NonAuthorizing                   bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic) Validate() error {
    if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV application plan reverse observation LSP diagnostic")
    }
    if !d.NonExecuting {
        return fmt.Errorf("JEV application plan reverse observation LSP diagnostic must be non-executing")
    }
    if !d.NonAuthorizing {
        return fmt.Errorf("JEV application plan reverse observation LSP diagnostic must be non-authorizing")
    }
    switch d.Status {
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown:
        if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPUnknownCode ||
            d.MissingStage == "" || d.Publishable {
            return fmt.Errorf("unknown application plan reverse observation LSP diagnostic is inconsistent")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeBound:
        if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPBoundCode ||
            !d.Publishable || d.ApplicationPlanStatus == "" || d.ApplicationPlanBridgeDigest == "" ||
            d.ReverseObservationStatus == "" || d.ObservationMetricStatus == "" {
            return fmt.Errorf("bound application plan reverse observation LSP diagnostic is incomplete")
        }
        if d.ApplicationPlanStatus == jevExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanReady &&
            (d.PlanDigest == "" || d.PlanSource == "" || d.PlanEvidenceDigest == "" ||
                d.ObservationDigest == "" || d.ObservationSource == "" || d.ObservationEvidenceDigest == "" ||
                d.ObservationMetricDigest == "" || d.ObservationMetricSource == "" ||
                d.ObservationMetricEvidenceDigest == "") {
            return fmt.Errorf("ready application plan reverse observation LSP diagnostic lost evidence")
        }
    default:
        return fmt.Errorf("invalid application plan reverse observation LSP status")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic(
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
        d.ReverseObservationStatus,
        d.ObservationDigest,
        d.ObservationSource,
        d.ObservationEvidenceDigest,
        d.ObservationMetricStatus,
        d.ObservationMetricDigest,
        d.ObservationMetricSource,
        d.ObservationMetricEvidenceDigest,
        d.BridgeDigest,
        d.Publishable,
    )
    if d.ProjectionDigest != expected {
        return fmt.Errorf("JEV application plan reverse observation LSP projection digest mismatch")
    }
    return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic{
        Status:                         input.Status,
        MissingStage:                   input.MissingStage,
        ApplicationPlanStatus:          input.ApplicationPlanStatus,
        PlanDigest:                     input.PlanDigest,
        PlanSource:                     input.PlanSource,
        PlanEvidenceDigest:             input.PlanEvidenceDigest,
        ApplicationPlanBridgeDigest:    input.ApplicationPlanBridgeDigest,
        ReverseObservationStatus:       input.ReverseObservationStatus,
        ObservationDigest:               input.ObservationDigest,
        ObservationSource:               input.ObservationSource,
        ObservationEvidenceDigest:      input.ObservationEvidenceDigest,
        ObservationMetricStatus:        input.ObservationMetricStatus,
        ObservationMetricDigest:        input.ObservationMetricDigest,
        ObservationMetricSource:        input.ObservationMetricSource,
        ObservationMetricEvidenceDigest: input.ObservationMetricEvidenceDigest,
        BridgeDigest:                   input.BridgeDigest,
        NonExecuting:                   true,
        NonAuthorizing:                 true,
    }
    if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown || input.MissingStage != "" {
        output.Severity = "warning"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPUnknownCode
        output.Message = "application plan reverse observation provenance is incomplete"
        output.Publishable = false
    } else {
        output.Severity = "info"
        output.Code = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPBoundCode
        output.Message = "application plan reverse observation provenance is bound"
        output.Publishable = true
    }
    output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic(
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
        output.ReverseObservationStatus,
        output.ObservationDigest,
        output.ObservationSource,
        output.ObservationEvidenceDigest,
        output.ObservationMetricStatus,
        output.ObservationMetricDigest,
        output.ObservationMetricSource,
        output.ObservationMetricEvidenceDigest,
        output.BridgeDigest,
        output.Publishable,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeUnknown
        output.MissingStage = "lsp-projection"
        output.Publishable = false
        output.ProjectionDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridgeLSPDiagnostic(severity, code, message, status, missingStage, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus, observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest, bridgeDigest string, publishable bool) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t", severity, code, message, status, missingStage, applicationPlanStatus, planDigest, planSource, planEvidenceDigest, applicationPlanBridgeDigest, reverseObservationStatus, observationDigest, observationSource, observationEvidenceDigest, observationMetricStatus, observationMetricDigest, observationMetricSource, observationMetricEvidenceDigest, bridgeDigest, publishable)))
    return hex.EncodeToString(sum[:])
}
