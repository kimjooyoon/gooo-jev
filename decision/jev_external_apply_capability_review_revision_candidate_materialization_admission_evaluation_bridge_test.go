package decision

import "testing"

func TestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeValidate(t *testing.T) {
	bridge := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
		Status:                          jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
		CandidateDigest:                 "candidate-digest",
		CandidateSource:                 "candidate-source",
		CandidateGateDigest:             "candidate-gate-digest",
		GeneratedCandidateDigest:        "generated-candidate-digest",
		GeneratedCandidateSource:        "candidate-generator",
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
		EvaluationMetricSource:           "evaluation-metrics",
		EvaluationMetricEvidenceDigest:   "evaluation-metric-evidence",
		EvaluationStatus:                 jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeEvaluationBound,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	bridge.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge(
		bridge.Status, bridge.MissingStage, bridge.CandidateDigest, bridge.CandidateSource,
		bridge.CandidateGateDigest, bridge.GeneratedCandidateDigest, bridge.GeneratedCandidateSource,
		bridge.GenerationInputDigest, bridge.GenerationEvidenceDigest, bridge.GeneratorIdentity,
		bridge.MaterializationStatus, bridge.MaterializationBridgeDigest, bridge.AdmissionDecision,
		bridge.AdmissionStatus, bridge.AdmissionDigest, bridge.AdmissionSource,
		bridge.AdmissionEvidenceDigest, bridge.QualityMetricDigest, bridge.QualityMetricSource,
		bridge.QualityMetricEvidenceDigest, bridge.ReviewDigest, bridge.ReviewSource,
		bridge.ReviewEvidenceDigest, bridge.AdmissionBridgeDigest, bridge.EvaluationMode,
		bridge.ReverseObservationDigest, bridge.ReverseObservationSource,
		bridge.ReverseObservationEvidenceDigest, bridge.EvaluationMetricDigest,
		bridge.EvaluationMetricSource, bridge.EvaluationMetricEvidenceDigest, bridge.EvaluationStatus,
		bridge.NonExecuting, bridge.NonAuthorizing,
	)
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected valid evaluation bridge, got %v", err)
	}
}

func TestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown(t *testing.T) {
	bridge := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
		Status:         jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
		MissingStage:   "reverse-observation-evidence",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected UNKNOWN bridge to remain valid without digest evidence, got %v", err)
	}
}

