package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPUnknown(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge{
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
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid UNKNOWN generation request projection, got %v", err)
	}
	if diagnostic.Publishable || diagnostic.MissingStage != "proposal-evidence" ||
		diagnostic.GenerationRequestBridgeDigest != "" {
		t.Fatalf("UNKNOWN generation request projection lost fail-closed state: %#v", diagnostic)
	}
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionProposalGenerationRequestLSPDoesNotInfer(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge{
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
				ProposalDecision: "revise",
			},
			NonExecuting:   true,
			NonAuthorizing: true,
		},
	)
	if diagnostic.Publishable || diagnostic.MissingStage == "" ||
		diagnostic.GenerationRequestBridgeDigest != "" {
		t.Fatalf("missing generation request evidence was not kept fail-closed: %#v", diagnostic)
	}
}

