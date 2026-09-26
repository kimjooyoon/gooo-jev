package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPOProjectionUnknown(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPOProjection(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge{
            Status:         jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage:   "generation-request",
            NonExecuting:   true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request" || output.Publishable ||
        !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("unknown generation request projection was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request projection should validate: %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPOProjectionRejectsInvalidBridge(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackLSPOProjection(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridge{
            Status:         jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGenerationFeedbackBridgeBound,
            NonExecuting:   true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-candidate-generation-feedback" || output.Publishable {
        t.Fatalf("invalid generation request bridge was not fail-closed: %+v", output)
    }
}
