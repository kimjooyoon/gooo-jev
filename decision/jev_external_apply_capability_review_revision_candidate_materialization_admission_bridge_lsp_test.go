package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSP(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge{
			Status:                     jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeBound,
			CandidateDigest:            "candidate-digest",
			CandidateSource:             "candidate-source",
			CandidateGateDigest:         "candidate-gate-digest",
			GeneratedCandidateDigest:   "generated-candidate-digest",
			GeneratedCandidateSource:   "candidate-generator",
			GenerationInputDigest:       "generation-input",
			GenerationEvidenceDigest:    "generation-evidence",
			GeneratorIdentity:            "generator-v1",
			MaterializationStatus:        "materialized",
			MaterializationBridgeDigest: "materialization-bridge-digest",
			AdmissionDecision:            "admit",
			AdmissionStatus:              "candidate-admission-bound",
			AdmissionDigest:              "admission-digest",
			AdmissionSource:              "admission-engine",
			AdmissionEvidenceDigest:      "admission-evidence",
			QualityMetricDigest:          "quality-digest",
			QualityMetricSource:          "quality-metrics",
			QualityMetricEvidenceDigest:  "quality-evidence",
			ReviewDigest:                 "review-digest",
			ReviewSource:                 "review-engine",
			ReviewEvidenceDigest:         "review-evidence",
			BridgeDigest:                 "admission-bridge-digest",
			NonExecuting:                 true,
			NonAuthorizing:               true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid candidate materialization admission LSP projection, got %v", err)
	}
	if !diagnostic.Publishable || diagnostic.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPBoundCode ||
		diagnostic.AdmissionDigest != "admission-digest" || diagnostic.QualityMetricEvidenceDigest != "quality-evidence" ||
		diagnostic.ReviewEvidenceDigest != "review-evidence" {
		t.Fatalf("candidate materialization admission evidence was not published: %#v", diagnostic)
	}
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSPUnknown(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge{
			Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown,
			MissingStage:  "quality-metric-evidence",
			NonExecuting: true,
			NonAuthorizing: true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid UNKNOWN candidate materialization admission projection, got %v", err)
	}
	if diagnostic.Publishable || diagnostic.MissingStage != "quality-metric-evidence" ||
		diagnostic.BridgeDigest != "" || diagnostic.MaterializationBridgeDigest != "" {
		t.Fatalf("UNKNOWN candidate materialization admission projection lost fail-closed state: %#v", diagnostic)
	}
}

