package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPUnknown(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
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
			NonExecuting:   true,
			NonAuthorizing: true,
		},
	)
	if err := diagnostic.Validate(); err != nil {
		t.Fatalf("expected valid UNKNOWN revision proposal projection, got %v", err)
	}
	if diagnostic.Publishable || diagnostic.MissingStage != "feedback-evidence" ||
		diagnostic.ProposalBridgeDigest != "" {
		t.Fatalf("UNKNOWN revision proposal projection lost fail-closed state: %#v", diagnostic)
	}
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDoesNotInfer(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge{
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
			},
			NonExecuting:   true,
			NonAuthorizing: true,
		},
	)
	if diagnostic.Publishable || diagnostic.MissingStage == "" || diagnostic.ProposalBridgeDigest != "" {
		t.Fatalf("missing feedback evidence was not kept fail-closed: %#v", diagnostic)
	}
}

