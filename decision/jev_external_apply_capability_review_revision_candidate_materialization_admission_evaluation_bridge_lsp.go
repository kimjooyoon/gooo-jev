package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPBoundCode = "jev.provenance.candidate-materialization-admission-evaluation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic struct {
	Severity                          string
	Code                              string
	Message                           string
	Status                            string
	MissingStage                      string
	CandidateDigest                   string
	CandidateSource                   string
	CandidateGateDigest               string
	GeneratedCandidateDigest          string
	GeneratedCandidateSource          string
	GenerationInputDigest             string
	GenerationEvidenceDigest          string
	GeneratorIdentity                 string
	MaterializationStatus             string
	MaterializationBridgeDigest       string
	AdmissionDecision                 string
	AdmissionStatus                   string
	AdmissionDigest                   string
	AdmissionSource                   string
	AdmissionEvidenceDigest           string
	QualityMetricDigest               string
	QualityMetricSource               string
	QualityMetricEvidenceDigest       string
	ReviewDigest                      string
	ReviewSource                      string
	ReviewEvidenceDigest              string
	AdmissionBridgeDigest             string
	EvaluationMode                    string
	ReverseObservationDigest          string
	ReverseObservationSource          string
	ReverseObservationEvidenceDigest  string
	EvaluationMetricDigest            string
	EvaluationMetricSource            string
	EvaluationMetricEvidenceDigest    string
	EvaluationStatus                  string
	BridgeDigest                      string
	ProjectionDigest                  string
	Publishable                       bool
	NonExecuting                      bool
	NonAuthorizing                    bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic) Validate() error {
	if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.ProjectionDigest == "" {
		return fmt.Errorf("incomplete JEV candidate materialization admission evaluation LSP diagnostic")
	}
	if !d.NonExecuting {
		return fmt.Errorf("JEV candidate materialization admission evaluation LSP diagnostic must be non-executing")
	}
	if !d.NonAuthorizing {
		return fmt.Errorf("JEV candidate materialization admission evaluation LSP diagnostic must be non-authorizing")
	}
	switch d.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown:
		if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPUnknownCode ||
			d.MissingStage == "" || d.Publishable || d.MaterializationBridgeDigest != "" ||
			d.AdmissionBridgeDigest != "" || d.BridgeDigest != "" {
			return fmt.Errorf("unknown candidate materialization admission evaluation LSP diagnostic is inconsistent")
		}
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound:
		if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPBoundCode ||
			!d.Publishable || d.MaterializationBridgeDigest == "" || d.AdmissionBridgeDigest == "" ||
			d.BridgeDigest == "" || d.EvaluationStatus == "" {
			return fmt.Errorf("bound candidate materialization admission evaluation LSP diagnostic is incomplete")
		}
		switch d.AdmissionDecision {
		case "admit":
			if d.EvaluationMode == "" || d.ReverseObservationDigest == "" ||
				d.ReverseObservationSource == "" || d.ReverseObservationEvidenceDigest == "" ||
				d.EvaluationMetricDigest == "" || d.EvaluationMetricSource == "" ||
				d.EvaluationMetricEvidenceDigest == "" {
				return fmt.Errorf("admitted candidate evaluation LSP diagnostic lost observation or metric evidence")
			}
		case "hold", "reject":
			if d.EvaluationMode != "" || d.ReverseObservationDigest != "" ||
				d.ReverseObservationSource != "" || d.ReverseObservationEvidenceDigest != "" ||
				d.EvaluationMetricDigest != "" || d.EvaluationMetricSource != "" ||
				d.EvaluationMetricEvidenceDigest != "" {
				return fmt.Errorf("held or rejected candidate evaluation LSP diagnostic contains evaluation evidence")
			}
		default:
			return fmt.Errorf("invalid candidate materialization admission evaluation decision")
		}
	default:
		return fmt.Errorf("invalid candidate materialization admission evaluation LSP status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic(
		d.Severity, d.Code, d.Message, d.Status, d.MissingStage,
		d.CandidateDigest, d.CandidateSource, d.CandidateGateDigest,
		d.GeneratedCandidateDigest, d.GeneratedCandidateSource, d.GenerationInputDigest,
		d.GenerationEvidenceDigest, d.GeneratorIdentity, d.MaterializationStatus,
		d.MaterializationBridgeDigest, d.AdmissionDecision, d.AdmissionStatus,
		d.AdmissionDigest, d.AdmissionSource, d.AdmissionEvidenceDigest,
		d.QualityMetricDigest, d.QualityMetricSource, d.QualityMetricEvidenceDigest,
		d.ReviewDigest, d.ReviewSource, d.ReviewEvidenceDigest, d.AdmissionBridgeDigest,
		d.EvaluationMode, d.ReverseObservationDigest, d.ReverseObservationSource,
		d.ReverseObservationEvidenceDigest, d.EvaluationMetricDigest, d.EvaluationMetricSource,
		d.EvaluationMetricEvidenceDigest, d.EvaluationStatus, d.BridgeDigest, d.Publishable,
	)
	if d.ProjectionDigest != expected {
		return fmt.Errorf("JEV candidate materialization admission evaluation LSP projection digest mismatch")
	}
	return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic {
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic{
		Status:                          input.Status,
		MissingStage:                    input.MissingStage,
		CandidateDigest:                 input.CandidateDigest,
		CandidateSource:                 input.CandidateSource,
		CandidateGateDigest:             input.CandidateGateDigest,
		GeneratedCandidateDigest:         input.GeneratedCandidateDigest,
		GeneratedCandidateSource:         input.GeneratedCandidateSource,
		GenerationInputDigest:            input.GenerationInputDigest,
		GenerationEvidenceDigest:         input.GenerationEvidenceDigest,
		GeneratorIdentity:                input.GeneratorIdentity,
		MaterializationStatus:            input.MaterializationStatus,
		MaterializationBridgeDigest:      input.MaterializationBridgeDigest,
		AdmissionDecision:                input.AdmissionDecision,
		AdmissionStatus:                  input.AdmissionStatus,
		AdmissionDigest:                  input.AdmissionDigest,
		AdmissionSource:                  input.AdmissionSource,
		AdmissionEvidenceDigest:          input.AdmissionEvidenceDigest,
		QualityMetricDigest:              input.QualityMetricDigest,
		QualityMetricSource:              input.QualityMetricSource,
		QualityMetricEvidenceDigest:      input.QualityMetricEvidenceDigest,
		ReviewDigest:                     input.ReviewDigest,
		ReviewSource:                     input.ReviewSource,
		ReviewEvidenceDigest:             input.ReviewEvidenceDigest,
		AdmissionBridgeDigest:            input.AdmissionBridgeDigest,
		EvaluationMode:                   input.EvaluationMode,
		ReverseObservationDigest:         input.ReverseObservationDigest,
		ReverseObservationSource:         input.ReverseObservationSource,
		ReverseObservationEvidenceDigest: input.ReverseObservationEvidenceDigest,
		EvaluationMetricDigest:           input.EvaluationMetricDigest,
		EvaluationMetricSource:           input.EvaluationMetricSource,
		EvaluationMetricEvidenceDigest:   input.EvaluationMetricEvidenceDigest,
		EvaluationStatus:                 input.EvaluationStatus,
		BridgeDigest:                     input.BridgeDigest,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || input.MissingStage != "" {
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPUnknownCode
		output.Message = "candidate materialization admission evaluation provenance is incomplete"
		output.Publishable = false
		output.MaterializationBridgeDigest = ""
		output.AdmissionBridgeDigest = ""
		output.BridgeDigest = ""
		if output.MissingStage == "" {
			output.MissingStage = "evaluation-bridge"
		}
	} else {
		output.Severity = "info"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPBoundCode
		output.Message = "candidate materialization admission evaluation provenance is bound"
		output.Publishable = true
	}
	output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic(
		output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
		output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
		output.GeneratedCandidateDigest, output.GeneratedCandidateSource, output.GenerationInputDigest,
		output.GenerationEvidenceDigest, output.GeneratorIdentity, output.MaterializationStatus,
		output.MaterializationBridgeDigest, output.AdmissionDecision, output.AdmissionStatus,
		output.AdmissionDigest, output.AdmissionSource, output.AdmissionEvidenceDigest,
		output.QualityMetricDigest, output.QualityMetricSource, output.QualityMetricEvidenceDigest,
		output.ReviewDigest, output.ReviewSource, output.ReviewEvidenceDigest, output.AdmissionBridgeDigest,
		output.EvaluationMode, output.ReverseObservationDigest, output.ReverseObservationSource,
		output.ReverseObservationEvidenceDigest, output.EvaluationMetricDigest, output.EvaluationMetricSource,
		output.EvaluationMetricEvidenceDigest, output.EvaluationStatus, output.BridgeDigest, output.Publishable,
	)
	if err := output.Validate(); err != nil {
		output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		output.MissingStage = "lsp-projection"
		output.Publishable = false
		output.MaterializationBridgeDigest = ""
		output.AdmissionBridgeDigest = ""
		output.BridgeDigest = ""
		output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic(
			output.Severity, output.Code, output.Message, output.Status, output.MissingStage,
			output.CandidateDigest, output.CandidateSource, output.CandidateGateDigest,
			output.GeneratedCandidateDigest, output.GeneratedCandidateSource, output.GenerationInputDigest,
			output.GenerationEvidenceDigest, output.GeneratorIdentity, output.MaterializationStatus,
			output.MaterializationBridgeDigest, output.AdmissionDecision, output.AdmissionStatus,
			output.AdmissionDigest, output.AdmissionSource, output.AdmissionEvidenceDigest,
			output.QualityMetricDigest, output.QualityMetricSource, output.QualityMetricEvidenceDigest,
			output.ReviewDigest, output.ReviewSource, output.ReviewEvidenceDigest, output.AdmissionBridgeDigest,
			output.EvaluationMode, output.ReverseObservationDigest, output.ReverseObservationSource,
			output.ReverseObservationEvidenceDigest, output.EvaluationMetricDigest, output.EvaluationMetricSource,
			output.EvaluationMetricEvidenceDigest, output.EvaluationStatus, output.BridgeDigest, output.Publishable,
		)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPDiagnostic(
	severity, code, message, status, missingStage, candidateDigest, candidateSource, candidateGateDigest,
	generatedCandidateDigest, generatedCandidateSource, generationInputDigest, generationEvidenceDigest,
	generatorIdentity, materializationStatus, materializationBridgeDigest, admissionDecision, admissionStatus,
	admissionDigest, admissionSource, admissionEvidenceDigest, qualityMetricDigest, qualityMetricSource,
	qualityMetricEvidenceDigest, reviewDigest, reviewSource, reviewEvidenceDigest, admissionBridgeDigest,
	evaluationMode, reverseObservationDigest, reverseObservationSource, reverseObservationEvidenceDigest,
	evaluationMetricDigest, evaluationMetricSource, evaluationMetricEvidenceDigest, evaluationStatus,
	bridgeDigest string, publishable bool,
) string {
	values := []string{
		severity, code, message, status, missingStage, candidateDigest, candidateSource,
		candidateGateDigest, generatedCandidateDigest, generatedCandidateSource, generationInputDigest,
		generationEvidenceDigest, generatorIdentity, materializationStatus, materializationBridgeDigest,
		admissionDecision, admissionStatus, admissionDigest, admissionSource, admissionEvidenceDigest,
		qualityMetricDigest, qualityMetricSource, qualityMetricEvidenceDigest, reviewDigest, reviewSource,
		reviewEvidenceDigest, admissionBridgeDigest, evaluationMode, reverseObservationDigest,
		reverseObservationSource, reverseObservationEvidenceDigest, evaluationMetricDigest,
		evaluationMetricSource, evaluationMetricEvidenceDigest, evaluationStatus, bridgeDigest,
		strconv.FormatBool(publishable),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

