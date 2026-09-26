package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeUnknown(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeInput{
			Outcome: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
					Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
					MissingStage:  "outcome-evidence",
					NonExecuting:  true,
					NonAuthorizing: true,
				},
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "outcome-evidence" || bridge.BridgeDigest != "" {
		t.Fatalf("feedback bridge did not preserve UNKNOWN provenance: %#v", bridge)
	}
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected fail-closed UNKNOWN feedback bridge to validate, got %v", err)
	}
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeDoesNotInfer(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeInput{
			Outcome: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
					Status:             jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
					AdmissionDecision:  "admit",
					BridgeDigest:       "outcome-bridge",
					NonExecuting:       true,
					NonAuthorizing:     true,
				},
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "evaluation-outcome-bridge" {
		t.Fatalf("missing outcome evidence was not fail-closed: %#v", bridge)
	}
}

