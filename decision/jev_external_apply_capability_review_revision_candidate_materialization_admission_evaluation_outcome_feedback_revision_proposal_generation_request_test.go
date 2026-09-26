package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeUnknown(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeInput{
			RevisionProposal: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
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
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "proposal-evidence" || bridge.BridgeDigest != "" {
		t.Fatalf("generation request bridge did not preserve UNKNOWN provenance: %#v", bridge)
	}
	if err := bridge.Validate(); err != nil {
		t.Fatalf("expected fail-closed UNKNOWN generation request bridge to validate, got %v", err)
	}
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeRequiresEvidence(t *testing.T) {
	bridge := BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeInput{
			RevisionProposal: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
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
				ProposalDecision: "revise",
			},
			NonAuthorizing: true,
		},
	)
	if bridge.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		bridge.MissingStage != "revision-proposal-bridge" {
		t.Fatalf("invalid proposal provenance was not fail-closed: %#v", bridge)
	}
}

