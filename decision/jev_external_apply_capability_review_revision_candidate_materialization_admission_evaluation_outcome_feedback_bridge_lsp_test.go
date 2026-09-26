package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPUnknown(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge{
				JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge: JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridge{
					Status:        jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
					MissingStage:  "outcome-evidence",
					NonExecuting:  true,
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
		t.Fatalf("expected valid UNKNOWN feedback projection, got %v", err)
	}
	if diagnostic.Publishable || diagnostic.MissingStage != "outcome-evidence" ||
		diagnostic.FeedbackBridgeDigest != "" {
		t.Fatalf("UNKNOWN feedback projection lost fail-closed state: %#v", diagnostic)
	}
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDoesNotInfer(t *testing.T) {
	diagnostic := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSP(
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge{
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
	)
	if diagnostic.Publishable || diagnostic.MissingStage == "" || diagnostic.FeedbackBridgeDigest != "" {
		t.Fatalf("missing outcome evidence was not kept fail-closed: %#v", diagnostic)
	}
}

