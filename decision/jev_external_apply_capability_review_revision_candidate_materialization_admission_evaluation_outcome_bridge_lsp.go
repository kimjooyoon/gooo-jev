package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPBoundCode = "jev.provenance.candidate-evaluation-outcome-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic
	Outcome               string
	BaselineMetricDigest  string
	ObservedMetricDigest  string
	OutcomeSource         string
	OutcomeEvidenceDigest string
	OutcomeStatus         string
	OutcomeBridgeDigest   string
	ProjectionDigest      string
	Publishable           bool
	NonExecuting          bool
	NonAuthorizing        bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic) Validate() error {
	if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.ProjectionDigest == "" {
		return fmt.Errorf("incomplete JEV candidate evaluation outcome LSP diagnostic")
	}
	if !d.NonExecuting {
		return fmt.Errorf("JEV candidate evaluation outcome LSP diagnostic must be non-executing")
	}
	if !d.NonAuthorizing {
		return fmt.Errorf("JEV candidate evaluation outcome LSP diagnostic must be non-authorizing")
	}
	switch d.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown:
		if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPUnknownCode ||
			d.MissingStage == "" || d.Publishable || d.OutcomeBridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate evaluation outcome LSP diagnostic is inconsistent")
		}
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound:
		if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPBoundCode ||
			!d.Publishable || d.OutcomeBridgeDigest == "" || d.OutcomeStatus == "" {
			return fmt.Errorf("bound JEV candidate evaluation outcome LSP diagnostic is incomplete")
		}
		if err := d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic.Validate(); err != nil {
			return fmt.Errorf("invalid evaluation LSP diagnostic: %w", err)
		}
		switch d.AdmissionDecision {
		case "admit":
			if d.OutcomeStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeBound ||
				d.Outcome == "" || d.BaselineMetricDigest == "" || d.ObservedMetricDigest == "" ||
				d.OutcomeSource == "" || d.OutcomeEvidenceDigest == "" {
				return fmt.Errorf("admitted candidate outcome LSP diagnostic lost metric evidence")
			}
		case "hold", "reject":
			if d.Outcome != "" || d.BaselineMetricDigest != "" || d.ObservedMetricDigest != "" ||
				d.OutcomeSource != "" || d.OutcomeEvidenceDigest != "" {
				return fmt.Errorf("held or rejected candidate outcome LSP diagnostic contains metric evidence")
			}
		default:
			return fmt.Errorf("invalid candidate evaluation outcome admission decision")
		}
	default:
		return fmt.Errorf("invalid candidate evaluation outcome LSP status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic(
		d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic.ProjectionDigest,
		d.Outcome, d.BaselineMetricDigest, d.ObservedMetricDigest, d.OutcomeSource,
		d.OutcomeEvidenceDigest, d.OutcomeStatus, d.OutcomeBridgeDigest, d.Publishable,
	)
	if d.ProjectionDigest != expected {
		return fmt.Errorf("JEV candidate evaluation outcome LSP projection digest mismatch")
	}
	return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic {
	prior := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSP(
		input.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge,
	)
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic: prior,
		Outcome:               input.Outcome,
		BaselineMetricDigest:  input.BaselineMetricDigest,
		ObservedMetricDigest:  input.ObservedMetricDigest,
		OutcomeSource:         input.OutcomeSource,
		OutcomeEvidenceDigest: input.OutcomeEvidenceDigest,
		OutcomeStatus:         input.OutcomeStatus,
		OutcomeBridgeDigest:   input.BridgeDigest,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		input.MissingStage != "" || !prior.Publishable {
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPUnknownCode
		output.Message = "candidate evaluation outcome provenance is incomplete"
		output.Publishable = false
		output.OutcomeBridgeDigest = ""
		if output.MissingStage == "" {
			output.MissingStage = "evaluation-outcome-bridge"
		}
	} else {
		output.Severity = "info"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPBoundCode
		output.Message = "candidate evaluation outcome provenance is bound"
		output.Publishable = true
	}
	output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic.ProjectionDigest,
		output.Outcome, output.BaselineMetricDigest, output.ObservedMetricDigest,
		output.OutcomeSource, output.OutcomeEvidenceDigest, output.OutcomeStatus,
		output.OutcomeBridgeDigest, output.Publishable,
	)
	if err := output.Validate(); err != nil {
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic.Status =
			jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPUnknownCode
		output.Message = "candidate evaluation outcome provenance is incomplete"
		output.MissingStage = "lsp-projection"
		output.Publishable = false
		output.OutcomeBridgeDigest = ""
		output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic(
			output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic.ProjectionDigest,
			output.Outcome, output.BaselineMetricDigest, output.ObservedMetricDigest,
			output.OutcomeSource, output.OutcomeEvidenceDigest, output.OutcomeStatus,
			output.OutcomeBridgeDigest, output.Publishable,
		)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic(
	priorProjectionDigest, outcome, baselineMetricDigest, observedMetricDigest, outcomeSource,
	outcomeEvidenceDigest, outcomeStatus, outcomeBridgeDigest string, publishable bool,
) string {
	values := []string{
		priorProjectionDigest, outcome, baselineMetricDigest, observedMetricDigest,
		outcomeSource, outcomeEvidenceDigest, outcomeStatus, outcomeBridgeDigest,
		strconv.FormatBool(publishable),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

