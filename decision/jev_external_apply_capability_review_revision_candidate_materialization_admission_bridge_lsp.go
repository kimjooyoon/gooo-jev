package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPBoundCode = "jev.provenance.candidate-materialization-admission-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic struct {
	Severity                    string
	Code                        string
	Message                     string
	Status                     string
	MissingStage                string
	CandidateDigest             string
	CandidateSource             string
	CandidateGateDigest         string
	GeneratedCandidateDigest    string
	GeneratedCandidateSource    string
	GenerationInputDigest       string
	GenerationEvidenceDigest    string
	GeneratorIdentity            string
	MaterializationStatus       string
	MaterializationBridgeDigest string
	AdmissionDecision            string
	AdmissionStatus              string
	AdmissionDigest              string
	AdmissionSource              string
	AdmissionEvidenceDigest     string
	QualityMetricDigest          string
	QualityMetricSource          string
	QualityMetricEvidenceDigest string
	ReviewDigest                 string
	ReviewSource                 string
	ReviewEvidenceDigest         string
	BridgeDigest                 string
	ProjectionDigest             string
	Publishable                  bool
	NonExecuting                 bool
	NonAuthorizing               bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic) Validate() error {
	if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.ProjectionDigest == "" {
		return fmt.Errorf("incomplete JEV candidate materialization admission LSP diagnostic")
	}
	if !d.NonExecuting {
		return fmt.Errorf("JEV candidate materialization admission LSP diagnostic must be non-executing")
	}
	if !d.NonAuthorizing {
		return fmt.Errorf("JEV candidate materialization admission LSP diagnostic must be non-authorizing")
	}
	switch d.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown:
		if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPUnknownCode ||
			d.MissingStage == "" || d.Publishable || d.BridgeDigest != "" || d.MaterializationBridgeDigest != "" {
			return fmt.Errorf("unknown candidate materialization admission LSP diagnostic is inconsistent")
		}
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeBound:
		if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPBoundCode ||
			!d.Publishable || d.BridgeDigest == "" || d.MaterializationBridgeDigest == "" || d.AdmissionStatus == "" {
			return fmt.Errorf("bound candidate materialization admission LSP diagnostic is incomplete")
		}
		switch d.AdmissionDecision {
		case "admit":
			if d.CandidateDigest == "" || d.CandidateSource == "" || d.CandidateGateDigest == "" ||
				d.GeneratedCandidateDigest == "" || d.GeneratedCandidateSource == "" ||
				d.GenerationInputDigest == "" || d.GenerationEvidenceDigest == "" ||
				d.GeneratorIdentity == "" || d.AdmissionDigest == "" || d.AdmissionSource == "" ||
				d.AdmissionEvidenceDigest == "" || d.QualityMetricDigest == "" ||
				d.QualityMetricSource == "" || d.QualityMetricEvidenceDigest == "" ||
				d.ReviewDigest == "" || d.ReviewSource == "" || d.ReviewEvidenceDigest == "" {
				return fmt.Errorf("admitted candidate materialization admission LSP diagnostic lost evidence")
			}
		case "hold", "reject":
			if d.AdmissionDigest != "" || d.AdmissionSource != "" || d.AdmissionEvidenceDigest != "" ||
				d.QualityMetricDigest != "" || d.QualityMetricSource != "" || d.QualityMetricEvidenceDigest != "" ||
				d.ReviewDigest != "" || d.ReviewSource != "" || d.ReviewEvidenceDigest != "" {
				return fmt.Errorf("held or rejected candidate materialization admission LSP diagnostic contains admission evidence")
			}
		default:
			return fmt.Errorf("invalid candidate materialization admission decision")
		}
	default:
		return fmt.Errorf("invalid candidate materialization admission LSP status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic(
		d.Severity, d.Code, d.Message, d.Status, d.MissingStage,
		d.CandidateDigest, d.CandidateSource, d.CandidateGateDigest,
		d.GeneratedCandidateDigest, d.GeneratedCandidateSource, d.GenerationInputDigest,
		d.GenerationEvidenceDigest, d.GeneratorIdentity, d.MaterializationStatus,
		d.MaterializationBridgeDigest, d.AdmissionDecision, d.AdmissionStatus,
		d.AdmissionDigest, d.AdmissionSource, d.AdmissionEvidenceDigest,
		d.QualityMetricDigest, d.QualityMetricSource, d.QualityMetricEvidenceDigest,
		d.ReviewDigest, d.ReviewSource, d.ReviewEvidenceDigest, d.BridgeDigest,
		d.Publishable,
	)
	if d.ProjectionDigest != expected {
		return fmt.Errorf("JEV candidate materialization admission LSP projection digest mismatch")
	}
	return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic {
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic{
		Status:                     input.Status,
		MissingStage:               input.MissingStage,
		CandidateDigest:            input.CandidateDigest,
		CandidateSource:             input.CandidateSource,
		CandidateGateDigest:        input.CandidateGateDigest,
		GeneratedCandidateDigest:   input.GeneratedCandidateDigest,
		GeneratedCandidateSource:   input.GeneratedCandidateSource,
		GenerationInputDigest:      input.GenerationInputDigest,
		GenerationEvidenceDigest:   input.GenerationEvidenceDigest,
		GeneratorIdentity:           input.GeneratorIdentity,
		MaterializationStatus:      input.MaterializationStatus,
		MaterializationBridgeDigest: input.MaterializationBridgeDigest,
		AdmissionDecision:           input.AdmissionDecision,
		AdmissionStatus:             input.AdmissionStatus,
		AdmissionDigest:             input.AdmissionDigest,
		AdmissionSource:             input.AdmissionSource,
		AdmissionEvidenceDigest:     input.AdmissionEvidenceDigest,
		QualityMetricDigest:         input.QualityMetricDigest,
		QualityMetricSource:         input.QualityMetricSource,
		QualityMetricEvidenceDigest: input.QualityMetricEvidenceDigest,
		ReviewDigest:                input.ReviewDigest,
		ReviewSource:                input.ReviewSource,
		ReviewEvidenceDigest:        input.ReviewEvidenceDigest,
		BridgeDigest:                input.BridgeDigest,
		NonExecuting:                true,
		NonAuthorizing:              true,
	}
	if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown || input.MissingStage != "" {
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPUnknownCode
		output.Message = "candidate materialization admission provenance is incomplete"
		output.Publishable = false
		output.BridgeDigest = ""
		output.MaterializationBridgeDigest = ""
		if output.MissingStage == "" {
			output.MissingStage = "admission-bridge"
		}
	} else {
		output.Severity = "info"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPBoundCode
		output.Message = "candidate materialization admission provenance is bound"
		output.Publishable = true
	}
	output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic(
		output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
		output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
		output.GeneratedCandidateDigest, output.GeneratedCandidateSource, output.GenerationInputDigest,
		output.GenerationEvidenceDigest, output.GeneratorIdentity, output.MaterializationStatus,
		output.MaterializationBridgeDigest, output.AdmissionDecision, output.AdmissionStatus,
		output.AdmissionDigest, output.AdmissionSource, output.AdmissionEvidenceDigest,
		output.QualityMetricDigest, output.QualityMetricSource, output.QualityMetricEvidenceDigest,
		output.ReviewDigest, output.ReviewSource, output.ReviewEvidenceDigest, output.BridgeDigest,
		output.Publishable,
	)
	if err := output.Validate(); err != nil {
		output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
		output.MissingStage = "lsp-projection"
		output.Publishable = false
		output.BridgeDigest = ""
		output.MaterializationBridgeDigest = ""
		output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic(
			output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
			output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
			output.GeneratedCandidateDigest, output.GeneratedCandidateSource, output.GenerationInputDigest,
			output.GenerationEvidenceDigest, output.GeneratorIdentity, output.MaterializationStatus,
			output.MaterializationBridgeDigest, output.AdmissionDecision, output.AdmissionStatus,
			output.AdmissionDigest, output.AdmissionSource, output.AdmissionEvidenceDigest,
			output.QualityMetricDigest, output.QualityMetricSource, output.QualityMetricEvidenceDigest,
			output.ReviewDigest, output.ReviewSource, output.ReviewEvidenceDigest, output.BridgeDigest,
			output.Publishable,
		)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPDiagnostic(
	severity, code, message, status, missingStage, candidateDigest, candidateSource, candidateGateDigest,
	generatedCandidateDigest, generatedCandidateSource, generationInputDigest, generationEvidenceDigest,
	generatorIdentity, materializationStatus, materializationBridgeDigest, admissionDecision, admissionStatus,
	admissionDigest, admissionSource, admissionEvidenceDigest, qualityMetricDigest, qualityMetricSource,
	qualityMetricEvidenceDigest, reviewDigest, reviewSource, reviewEvidenceDigest, bridgeDigest string,
	publishable bool,
) string {
	values := []string{
		severity, code, message, status, missingStage, candidateDigest, candidateSource,
		candidateGateDigest, generatedCandidateDigest, generatedCandidateSource, generationInputDigest,
		generationEvidenceDigest, generatorIdentity, materializationStatus, materializationBridgeDigest,
		admissionDecision, admissionStatus, admissionDigest, admissionSource, admissionEvidenceDigest,
		qualityMetricDigest, qualityMetricSource, qualityMetricEvidenceDigest, reviewDigest, reviewSource,
		reviewEvidenceDigest, bridgeDigest, strconv.FormatBool(publishable),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

