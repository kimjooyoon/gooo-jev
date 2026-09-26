package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationInput
// joins the extended revision candidate identity to the existing materialization
// admission and evaluation bridge.
type ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationInput struct {
	RevisionCandidate              ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateBinding
	Admission                      JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge
	EvaluationMode                 string
	ReverseObservationDigest       string
	ReverseObservationSource       string
	ReverseObservationEvidenceDigest string
	EvaluationMetricDigest         string
	EvaluationMetricSource         string
	EvaluationMetricEvidenceDigest string
	NonAuthorizing                 bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding
// preserves candidate, admission, and evaluation evidence without applying code.
type ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding struct {
	Status                    string
	MissingStage              string
	RevisionCandidateDigest   string
	ProvenanceEvidenceDigest  string
	LedgerEvidenceDigest      string
	AdmissionBridgeDigest     string
	EvaluationBridgeDigest    string
	AdmissionDecision          string
	EvaluationStatus           string
	EvidenceDigest             string
	NonExecuting              bool
	NonAuthorizing             bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.RevisionCandidateDigest == "" ||
		b.ProvenanceEvidenceDigest == "" ||
		b.LedgerEvidenceDigest == "" ||
		b.AdmissionBridgeDigest == "" ||
		b.EvaluationBridgeDigest == "" ||
		b.AdmissionDecision == "" ||
		b.EvaluationStatus == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo evidence materialization evaluation binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo evidence materialization evaluation binding must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		RevisionCandidateDigest  string
		ProvenanceEvidenceDigest string
		LedgerEvidenceDigest     string
		AdmissionBridgeDigest    string
		EvaluationBridgeDigest   string
		AdmissionDecision        string
		EvaluationStatus         string
	}{
		RevisionCandidateDigest:  b.RevisionCandidateDigest,
		ProvenanceEvidenceDigest: b.ProvenanceEvidenceDigest,
		LedgerEvidenceDigest:     b.LedgerEvidenceDigest,
		AdmissionBridgeDigest:    b.AdmissionBridgeDigest,
		EvaluationBridgeDigest:   b.EvaluationBridgeDigest,
		AdmissionDecision:        b.AdmissionDecision,
		EvaluationStatus:         b.EvaluationStatus,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo evidence materialization evaluation digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation
// reuses the established admission/evaluation bridge after checking the
// extended candidate identity.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(input ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "materialization-admission-evaluation"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding{
			Status: "UNKNOWN", MissingStage: stage,
			NonExecuting: true, NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.RevisionCandidate.NonAuthorizing ||
		!input.Admission.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.RevisionCandidate.NonExecuting || !input.Admission.NonExecuting {
		return unknown("execution-boundary")
	}
	if input.RevisionCandidate.Status != "bound" {
		return unknown(input.RevisionCandidate.MissingStage)
	}
	if strings.TrimSpace(input.RevisionCandidate.MissingStage) != "" {
		return unknown(input.RevisionCandidate.MissingStage)
	}
	if err := input.Admission.Validate(); err != nil {
		return unknown("candidate-materialization-admission")
	}
	if input.RevisionCandidate.CandidateDigest != input.Admission.CandidateDigest {
		return unknown("candidate-digest-binding")
	}
	evaluation := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeInput{
		Admission:                   input.Admission,
		EvaluationMode:              input.EvaluationMode,
		ReverseObservationDigest:    input.ReverseObservationDigest,
		ReverseObservationSource:    input.ReverseObservationSource,
		ReverseObservationEvidenceDigest: input.ReverseObservationEvidenceDigest,
		EvaluationMetricDigest:      input.EvaluationMetricDigest,
		EvaluationMetricSource:      input.EvaluationMetricSource,
		EvaluationMetricEvidenceDigest: input.EvaluationMetricEvidenceDigest,
		NonAuthorizing:              true,
	})
	if evaluation.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound {
		return unknown(evaluation.MissingStage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding{
		Status:                   "bound",
		RevisionCandidateDigest:  input.RevisionCandidate.CandidateDigest,
		ProvenanceEvidenceDigest: input.RevisionCandidate.ProvenanceEvidenceDigest,
		LedgerEvidenceDigest:     input.RevisionCandidate.LedgerEvidenceDigest,
		AdmissionBridgeDigest:    input.Admission.BridgeDigest,
		EvaluationBridgeDigest:   evaluation.BridgeDigest,
		AdmissionDecision:        evaluation.AdmissionDecision,
		EvaluationStatus:         evaluation.EvaluationStatus,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		RevisionCandidateDigest  string
		ProvenanceEvidenceDigest string
		LedgerEvidenceDigest     string
		AdmissionBridgeDigest    string
		EvaluationBridgeDigest   string
		AdmissionDecision        string
		EvaluationStatus         string
	}{
		RevisionCandidateDigest:  output.RevisionCandidateDigest,
		ProvenanceEvidenceDigest: output.ProvenanceEvidenceDigest,
		LedgerEvidenceDigest:     output.LedgerEvidenceDigest,
		AdmissionBridgeDigest:    output.AdmissionBridgeDigest,
		EvaluationBridgeDigest:   output.EvaluationBridgeDigest,
		AdmissionDecision:        output.AdmissionDecision,
		EvaluationStatus:          output.EvaluationStatus,
	})
	if err := output.Validate(); err != nil {
		return unknown("materialization-evaluation-evidence")
	}
	return output
}
