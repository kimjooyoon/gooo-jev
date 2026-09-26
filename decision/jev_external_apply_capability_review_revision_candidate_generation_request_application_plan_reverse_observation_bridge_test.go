package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeUnknown(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeInput{},
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "capability-boundary" || !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("empty generation request reverse observation bridge was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request reverse observation bridge should validate: %v", err)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeRejectsDigestDrift(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeInput{
            NonAuthorizing: true,
            GenerationRequestApplicationPlan: JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateFeedbackApplicationPlanBridge{
                NonExecuting: true,
                NonAuthorizing: true,
            },
            ApplicationPlanReverseObservation: JEVExternalApplyCapabilityReviewRevisionCandidateApplicationPlanReverseObservationBridge{
                NonExecuting: true,
                NonAuthorizing: true,
                ApplicationPlanBridgeDigest: "drift",
            },
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-application-plan" {
        t.Fatalf("invalid generation request reverse observation upstream was not fail-closed: %+v", output)
    }
}
