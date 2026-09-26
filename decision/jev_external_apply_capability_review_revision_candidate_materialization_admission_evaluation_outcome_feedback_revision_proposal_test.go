package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeUnknown(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeInput{
			FeedbackBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
					JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
						Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
						MissingStage:  "feedback-evidence",
						NonExecuting:  true,
						NonAuthorizing: true,
					},
					NonExecuting:   true,
					NonAuthorizing: true,
				},
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "feedback-evidence" || bridge.BridgeDigest != "" {
		t.Fatalf("revision proposal bridge did not preserve UNKNOWN provenance: %#v", bridge)
	}
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected fail-closed UNKNOWN revision proposal bridge to validate, got %v", err)
	}
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeDoesNotInfer(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeInput{
			FeedbackBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
					JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
						Status:             jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound,
						AdmissionDecision:  "admit",
						NonExecuting:       true,
						NonAuthorizing:     true,
						BridgeDigest:       "evaluation-bridge",
					},
					NonExecuting:   true,
					NonAuthorizing: true,
					BridgeDigest:   "outcome-bridge",
				},
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "evaluation-outcome-feedback" {
		t.Fatalf("invalid prior evidence was not fail-closed: %#v", bridge)
	}
}

