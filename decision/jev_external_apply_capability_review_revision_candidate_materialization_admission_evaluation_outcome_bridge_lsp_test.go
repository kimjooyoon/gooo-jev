package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSP(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
				Status:                          jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
				CandidateDigest:                 "candidate-digest",
				CandidateSource:                 "candidate-source",
				CandidateGateDigest:             "candidate-gate-digest",
				GeneratedCandidateDigest:        "generated-candidate-digest",
				GeneratedCandidateSource:        "candidate-generator",
				GenerationInputDigest:           "generation-input",
				GenerationEvidenceDigest:        "generation-evidence",
				GeneratorIdentity:               "generator-v1",
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
				EvaluationMetricDigest:            "evaluation-metric-digest",
				EvaluationMetricSource:            "evaluation-metrics",
				EvaluationMetricEvidenceDigest:   "evaluation-metric-evidence",
				EvaluationStatus:                 "candidate-evaluation-bound",
				BridgeDigest:                     "evaluation-bridge",
				NonExecuting:                     true,
				NonAuthorizing:                   true,
			},
			Outcome:               "improved",
			BaselineMetricDigest:  "baseline-metric",
			ObservedMetricDigest:  "observed-metric",
			OutcomeSource:         "outcome-engine",
			OutcomeEvidenceDigest: "outcome-evidence",
			OutcomeStatus:         jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeBound,
			BridgeDigest:          "outcome-bridge",
			NonExecuting:          true,
			NonAuthorizing:        true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid outcome LSP projection, got %v", err)
	}
	if !diagnostic.Publishable || diagnostic.Outcome != "improved" ||
		diagnostic.BaselineMetricDigest != "baseline-metric" ||
		diagnostic.OutcomeEvidenceDigest != "outcome-evidence" {
		t.Fatalf("outcome evidence was not published: %#v", diagnostic)
	}
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPUnknown(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
				Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
				MissingStage:  "reverse-observation-evidence",
				NonExecuting:  true,
				NonAuthorizing: true,
			},
			NonExecuting:   true,
			NonAuthorizing: true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid UNKNOWN outcome projection, got %v", err)
	}
	if diagnostic.Publishable || diagnostic.MissingStage != "reverse-observation-evidence" ||
		diagnostic.OutcomeBridgeDigest != "" {
		t.Fatalf("UNKNOWN outcome projection lost fail-closed state: %#v", diagnostic)
	}
}

