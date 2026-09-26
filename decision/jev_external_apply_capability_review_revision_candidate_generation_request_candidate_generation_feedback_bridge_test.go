package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeUnknown(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeInput{
			GenerationRequest: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
					JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
						JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
							JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
								Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
								MissingStage:  "proposal-evidence",
								NonExecuting:  true,
								NonAuthorizing: true,
							},
							NonExecuting:   true,
							NonAuthorizing: true,
						},
						NonExecuting:   true,
						NonAuthorizing: true,
					},
					NonExecuting:   true,
					NonAuthorizing: true,
				},
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			Direction: JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge{
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "proposal-evidence" || bridge.BridgeDigest != "" {
		t.Fatalf("generation feedback bridge did not preserve UNKNOWN provenance: %#v", bridge)
	}
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected fail-closed UNKNOWN bridge to validate, got %v", err)
	}
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeRejectsDirectionInference(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeInput{
			GenerationRequest: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
					JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
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
						BridgeDigest:   "feedback-bridge",
					},
					NonExecuting:   true,
					NonAuthorizing: true,
					BridgeDigest:   "proposal-bridge",
				},
				NonExecuting:   true,
				NonAuthorizing: true,
				BridgeDigest:   "request-bridge",
			},
			Direction: JEVExternalApplyCapabilityReviewRevisionCandidateObservationMetricDirectionBridge{
				Direction:     jevExternalApplyCapabilityReviewImprovementDirectionHold,
				NonExecuting:   true,
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "generation-request" {
		t.Fatalf("invalid request provenance was not fail-closed: %#v", bridge)
	}
}

