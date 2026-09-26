package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound = "candidate-materialization-admission-evaluation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationBound = "candidate-evaluation-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationHeld = "candidate-evaluation-held"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationRejected = "candidate-evaluation-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeInput struct {
	Admission                   JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge
	EvaluationMode              string
	ReverseObservationDigest    string
	ReverseObservationSource    string
	ReverseObservationEvidenceDigest string
	EvaluationMetricDigest      string
	EvaluationMetricSource      string
	EvaluationMetricEvidenceDigest string
	NonAuthorizing              bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge struct {
	Status                          string
	MissingStage                    string
	CandidateDigest                 string
	CandidateSource                 string
	CandidateGateDigest             string
	GeneratedCandidateDigest        string
	GeneratedCandidateSource        string
	GenerationInputDigest           string
	GenerationEvidenceDigest        string
	GeneratorIdentity               string
	MaterializationStatus            string
	MaterializationBridgeDigest      string
	AdmissionDecision                string
	AdmissionStatus                 string
	AdmissionDigest                 string
	AdmissionSource                 string
	AdmissionEvidenceDigest         string
	QualityMetricDigest              string
	QualityMetricSource              string
	QualityMetricEvidenceDigest     string
	ReviewDigest                     string
	ReviewSource                     string
	ReviewEvidenceDigest             string
	AdmissionBridgeDigest             string
	EvaluationMode                   string
	ReverseObservationDigest         string
	ReverseObservationSource         string
	ReverseObservationEvidenceDigest string
	EvaluationMetricDigest           string
	EvaluationMetricSource           string
	EvaluationMetricEvidenceDigest   string
	EvaluationStatus                 string
	BridgeDigest                     string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge) Validate() error {
	if b.Status == "" {
		return fmt.Errorf("incomplete JEV candidate materialization admission evaluation bridge")
	}
	if !b.NonExecuting {
		return fmt.Errorf("JEV candidate materialization admission evaluation bridge must be non-executing")
	}
	if !b.NonAuthorizing {
		return fmt.Errorf("JEV candidate materialization admission evaluation bridge must be non-authorizing")
	}
	switch b.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown:
		if b.MissingStage == "" || b.BridgeDigest != "" || b.AdmissionBridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate materialization admission evaluation bridge is inconsistent")
		}
		return nil
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound:
		if b.MissingStage != "" || b.EvaluationStatus == "" || b.AdmissionDecision == "" {
			return fmt.Errorf("bound JEV candidate materialization admission evaluation bridge is incomplete")
		}
		switch b.AdmissionDecision {
		case "admit":
			if b.EvaluationStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationBound ||
				b.EvaluationMode == "" || b.ReverseObservationDigest == "" ||
				b.ReverseObservationSource == "" || b.ReverseObservationEvidenceDigest == "" ||
				b.EvaluationMetricDigest == "" || b.EvaluationMetricSource == "" ||
				b.EvaluationMetricEvidenceDigest == "" {
				return fmt.Errorf("admitted JEV candidate evaluation lost reverse observation or metric evidence")
			}
		case "hold":
			if b.EvaluationStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationHeld ||
				b.EvaluationMode != "" || b.ReverseObservationDigest != "" ||
				b.ReverseObservationSource != "" || b.ReverseObservationEvidenceDigest != "" ||
				b.EvaluationMetricDigest != "" || b.EvaluationMetricSource != "" ||
				b.EvaluationMetricEvidenceDigest != "" {
				return fmt.Errorf("held JEV candidate evaluation contains execution evidence")
			}
		case "reject":
			if b.EvaluationStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationRejected ||
				b.EvaluationMode != "" || b.ReverseObservationDigest != "" ||
				b.ReverseObservationSource != "" || b.ReverseObservationEvidenceDigest != "" ||
				b.EvaluationMetricDigest != "" || b.EvaluationMetricSource != "" ||
				b.EvaluationMetricEvidenceDigest != "" {
				return fmt.Errorf("rejected JEV candidate evaluation contains execution evidence")
			}
		default:
			return fmt.Errorf("invalid JEV candidate materialization admission decision")
		}
	default:
		return fmt.Errorf("invalid JEV candidate materialization admission evaluation bridge status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge(
		b.Status, b.MissingStage, b.CandidateDigest, b.CandidateSource, b.CandidateGateDigest,
		b.GeneratedCandidateDigest, b.GeneratedCandidateSource, b.GenerationInputDigest,
		b.GenerationEvidenceDigest, b.GeneratorIdentity, b.MaterializationStatus,
		b.MaterializationBridgeDigest, b.AdmissionDecision, b.AdmissionStatus,
		b.AdmissionDigest, b.AdmissionSource, b.AdmissionEvidenceDigest,
		b.QualityMetricDigest, b.QualityMetricSource, b.QualityMetricEvidenceDigest,
		b.ReviewDigest, b.ReviewSource, b.ReviewEvidenceDigest, b.AdmissionBridgeDigest,
		b.EvaluationMode, b.ReverseObservationDigest, b.ReverseObservationSource,
		b.ReverseObservationEvidenceDigest, b.EvaluationMetricDigest, b.EvaluationMetricSource,
		b.EvaluationMetricEvidenceDigest, b.EvaluationStatus, b.NonExecuting, b.NonAuthorizing,
	)
	if b.BridgeDigest != expected {
		return fmt.Errorf("JEV candidate materialization admission evaluation bridge digest mismatch")
	}
	return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge {
	unknown := func(stage string) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge {
		if stage == "" {
			stage = "admission-evaluation-bridge"
		}
		output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
			Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
			MissingStage:  stage,
			NonExecuting:  true,
			NonAuthorizing: true,
		}
		return output
	}
	if !input.NonAuthorizing || !input.Admission.NonExecuting || !input.Admission.NonAuthorizing {
		return unknown("capability-boundary")
	}
	if err := input.Admission.Validate(); err != nil {
		return unknown("admission-bridge")
	}
	if input.Admission.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown {
		return unknown(input.Admission.MissingStage)
	}
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
		Status:                          jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
		CandidateDigest:                 input.Admission.CandidateDigest,
		CandidateSource:                 input.Admission.CandidateSource,
		CandidateGateDigest:             input.Admission.CandidateGateDigest,
		GeneratedCandidateDigest:        input.Admission.GeneratedCandidateDigest,
		GeneratedCandidateSource:        input.Admission.GeneratedCandidateSource,
		GenerationInputDigest:           input.Admission.GenerationInputDigest,
		GenerationEvidenceDigest:        input.Admission.GenerationEvidenceDigest,
		GeneratorIdentity:               input.Admission.GeneratorIdentity,
		MaterializationStatus:            input.Admission.MaterializationStatus,
		MaterializationBridgeDigest:      input.Admission.MaterializationBridgeDigest,
		AdmissionDecision:                input.Admission.AdmissionDecision,
		AdmissionStatus:                  input.Admission.AdmissionStatus,
		AdmissionDigest:                 input.Admission.AdmissionDigest,
		AdmissionSource:                 input.Admission.AdmissionSource,
		AdmissionEvidenceDigest:          input.Admission.AdmissionEvidenceDigest,
		QualityMetricDigest:              input.Admission.QualityMetricDigest,
		QualityMetricSource:              input.Admission.QualityMetricSource,
		QualityMetricEvidenceDigest:      input.Admission.QualityMetricEvidenceDigest,
		ReviewDigest:                     input.Admission.ReviewDigest,
		ReviewSource:                     input.Admission.ReviewSource,
		ReviewEvidenceDigest:             input.Admission.ReviewEvidenceDigest,
		AdmissionBridgeDigest:             input.Admission.BridgeDigest,
		EvaluationMode:                   input.EvaluationMode,
		ReverseObservationDigest:         input.ReverseObservationDigest,
		ReverseObservationSource:         input.ReverseObservationSource,
		ReverseObservationEvidenceDigest: input.ReverseObservationEvidenceDigest,
		EvaluationMetricDigest:            input.EvaluationMetricDigest,
		EvaluationMetricSource:            input.EvaluationMetricSource,
		EvaluationMetricEvidenceDigest:   input.EvaluationMetricEvidenceDigest,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	switch output.AdmissionDecision {
	case "admit":
		output.EvaluationStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationBound
		stages := []struct {
			value string
			stage string
		}{
			{output.EvaluationMode, "evaluation-mode"},
			{output.ReverseObservationDigest, "reverse-observation-digest"},
			{output.ReverseObservationSource, "reverse-observation-source"},
			{output.ReverseObservationEvidenceDigest, "reverse-observation-evidence"},
			{output.EvaluationMetricDigest, "evaluation-metric-digest"},
			{output.EvaluationMetricSource, "evaluation-metric-source"},
			{output.EvaluationMetricEvidenceDigest, "evaluation-metric-evidence"},
		}
		for _, item := range stages {
			if item.value == "" {
				return unknown(item.stage)
			}
		}
	case "hold":
		output.EvaluationStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationHeld
		if input.EvaluationMode != "" || input.ReverseObservationDigest != "" ||
			input.ReverseObservationSource != "" || input.ReverseObservationEvidenceDigest != "" ||
			input.EvaluationMetricDigest != "" || input.EvaluationMetricSource != "" ||
			input.EvaluationMetricEvidenceDigest != "" {
			return unknown("held-evaluation-evidence")
		}
	case "reject":
		output.EvaluationStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationRejected
		if input.EvaluationMode != "" || input.ReverseObservationDigest != "" ||
			input.ReverseObservationSource != "" || input.ReverseObservationEvidenceDigest != "" ||
			input.EvaluationMetricDigest != "" || input.EvaluationMetricSource != "" ||
			input.EvaluationMetricEvidenceDigest != "" {
			return unknown("rejected-evaluation-evidence")
		}
	default:
		return unknown("admission-decision")
	}
	output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge(
		output.Status, output.MissingStage, output.CandidateDigest, output.CandidateSource,
		output.CandidateGateDigest, output.GeneratedCandidateDigest, output.GeneratedCandidateSource,
		output.GenerationInputDigest, output.GenerationEvidenceDigest, output.GeneratorIdentity,
		output.MaterializationStatus, output.MaterializationBridgeDigest, output.AdmissionDecision,
		output.AdmissionStatus, output.AdmissionDigest, output.AdmissionSource,
		output.AdmissionEvidenceDigest, output.QualityMetricDigest, output.QualityMetricSource,
		output.QualityMetricEvidenceDigest, output.ReviewDigest, output.ReviewSource,
		output.ReviewEvidenceDigest, output.AdmissionBridgeDigest, output.EvaluationMode,
		output.ReverseObservationDigest, output.ReverseObservationSource,
		output.ReverseObservationEvidenceDigest, output.EvaluationMetricDigest,
		output.EvaluationMetricSource, output.EvaluationMetricEvidenceDigest, output.EvaluationStatus,
		output.NonExecuting, output.NonAuthorizing,
	)
	if err := output.Validate(); err != nil {
		return unknown("evaluation-bridge")
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge(
	status, missingStage, candidateDigest, candidateSource, candidateGateDigest,
	generatedCandidateDigest, generatedCandidateSource, generationInputDigest,
	generationEvidenceDigest, generatorIdentity, materializationStatus,
	materializationBridgeDigest, admissionDecision, admissionStatus, admissionDigest,
	admissionSource, admissionEvidenceDigest, qualityMetricDigest, qualityMetricSource,
	qualityMetricEvidenceDigest, reviewDigest, reviewSource, reviewEvidenceDigest,
	admissionBridgeDigest, evaluationMode, reverseObservationDigest, reverseObservationSource,
	reverseObservationEvidenceDigest, evaluationMetricDigest, evaluationMetricSource,
	evaluationMetricEvidenceDigest, evaluationStatus string, nonExecuting, nonAuthorizing bool,
) string {
	values := []string{
		status, missingStage, candidateDigest, candidateSource, candidateGateDigest,
		generatedCandidateDigest, generatedCandidateSource, generationInputDigest,
		generationEvidenceDigest, generatorIdentity, materializationStatus,
		materializationBridgeDigest, admissionDecision, admissionStatus, admissionDigest,
		admissionSource, admissionEvidenceDigest, qualityMetricDigest, qualityMetricSource,
		qualityMetricEvidenceDigest, reviewDigest, reviewSource, reviewEvidenceDigest,
		admissionBridgeDigest, evaluationMode, reverseObservationDigest, reverseObservationSource,
		reverseObservationEvidenceDigest, evaluationMetricDigest, evaluationMetricSource,
		evaluationMetricEvidenceDigest, evaluationStatus, strconv.FormatBool(nonExecuting),
		strconv.FormatBool(nonAuthorizing),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

