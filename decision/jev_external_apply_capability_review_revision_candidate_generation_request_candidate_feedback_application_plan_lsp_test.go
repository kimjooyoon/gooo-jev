package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPUnknown(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: "generation-request-candidate-feedback",
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-candidate-feedback" || output.Publishable ||
        !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("unknown generation request application plan LSP was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request application plan LSP should validate: %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSPRejectsInvalidBridge(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeBound,
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-candidate-feedback-application-plan" || output.Publishable {
        t.Fatalf("invalid generation request application plan bridge was not fail-closed: %+v", output)
    }
}
