package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeUnknown(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeInput{},
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "capability-boundary" || !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("empty generation request application plan bridge was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request application plan bridge should validate: %v", err)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeRejectsInvalidUpstream(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridgeInput{
            NonAuthorizing: true,
            GenerationRequestFeedback: JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge{
                NonExecuting: true,
                NonAuthorizing: true,
            },
            CandidateGateApplicationPlan: JEVExternalApplyCapabilityReviewRevisionCandidateGateApplicationPlanBridge{
                NonExecuting: true,
                NonAuthorizing: true,
            },
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-candidate-feedback" {
        t.Fatalf("invalid upstream generation request feedback was not fail-closed: %+v", output)
    }
}
