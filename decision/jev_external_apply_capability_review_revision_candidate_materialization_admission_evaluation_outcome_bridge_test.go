package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeUnknown(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeInput{
			Evaluation: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
				Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
				MissingStage:  "reverse-observation-evidence",
				NonExecuting:  true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "reverse-observation-evidence" || bridge.BridgeDigest != "" {
		t.Fatalf("outcome bridge did not preserve UNKNOWN provenance: %#v", bridge)
	}
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected fail-closed UNKNOWN outcome bridge to validate, got %v", err)
	}
}

func TestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeRejectsInference(t *testing.T) {
	bridge := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
			Status:             jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
			AdmissionDecision:  "admit",
			BridgeDigest:       "evaluation-bridge",
			NonExecuting:       true,
			NonAuthorizing:     true,
		},
		OutcomeStatus:   jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeOutcomeBound,
		Outcome:         "improved",
		NonExecuting:    true,
		NonAuthorizing:  true,
		BridgeDigest:    "outcome-bridge",
	}
	if err := bridge.Validate(); err == nil {
		t.Fatalf("expected missing metric evidence to reject inferred improvement")
	}
}

