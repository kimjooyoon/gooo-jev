package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSP(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
			Status:                          jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
			CandidateDigest:                 "candidate-digest",
			CandidateSource:                 "candidate-source",
			CandidateGateDigest:             "candidate-gate-digest",
			GeneratedCandidateDigest:        "generated-candidate-digest",
			GeneratedCandidateSource:         "candidate-generator",
			GenerationInputDigest:            "generation-input",
			GenerationEvidenceDigest:         "generation-evidence",
			GeneratorIdentity:                "generator-v1",
			MaterializationStatus:            "materialized",
			MaterializationBridgeDigest:      "materialization-bridge",
			AdmissionDecision:                "admit",
			AdmissionStatus:                  "candidate-admission-bound",
			AdmissionDigest:                  "admission-digest",
			AdmissionSource:                  "admission-engine",
			AdmissionEvidenceDigest:          "admission-evidence",
			QualityMetricDigest:              "quality-digest",
			QualityMetricSource:              "quality-metrics",
			QualityMetricEvidenceDigest:      "quality-evidence",
			ReviewDigest:                     "review-digest",
			ReviewSource:                     "review-engine",
			ReviewEvidenceDigest:             "review-evidence",
			AdmissionBridgeDigest:            "admission-bridge",
			EvaluationMode:                   "reverse-observation",
			ReverseObservationDigest:         "reverse-observation-digest",
			ReverseObservationSource:         "reverse-observation-engine",
			ReverseObservationEvidenceDigest: "reverse-observation-evidence",
			EvaluationMetricDigest:           "evaluation-metric-digest",
			EvaluationMetricSource:            "evaluation-metrics",
			EvaluationMetricEvidenceDigest:   "evaluation-metric-evidence",
			EvaluationStatus:                 "candidate-evaluation-bound",
			BridgeDigest:                     "evaluation-bridge",
			NonExecuting:                     true,
			NonAuthorizing:                   true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid evaluation LSP projection, got %v", err)
	}
	if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPBoundCode ||
		diagnostic.ReverseObservationDigest != "reverse-observation-digest" ||
		diagnostic.EvaluationMetricEvidenceDigest != "evaluation-metric-evidence" {
		t.Fatalf("evaluation evidence was not published: %#v", diagnostic)
	}
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSPUnknown(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
			Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
			MissingStage:  "reverse-observation-evidence",
			NonExecuting:  true,
			NonAuthorizing: true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid UNKNOWN evaluation projection, got %v", err)
	}
	if diagnostic.Publishable || diagnostic.MissingStage != "reverse-observation-evidence" ||
		diagnostic.MaterializationBridgeDigest != "" || diagnostic.AdmissionBridgeDigest != "" ||
		diagnostic.BridgeDigest != "" {
		t.Fatalf("UNKNOWN evaluation projection lost fail-closed state: %#v", diagnostic)
	}
}

