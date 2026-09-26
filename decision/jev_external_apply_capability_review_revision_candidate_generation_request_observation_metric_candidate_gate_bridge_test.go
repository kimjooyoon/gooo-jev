package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeUnknown(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeInput{},
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "capability-boundary" || !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("empty generation request candidate gate bridge was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request candidate gate bridge should validate: %v", err)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeRejectsInvalidDirectionGate(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeInput{
            NonAuthorizing: true,
            GenerationRequestReverseObservation: JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge{
                NonExecuting: true,
                NonAuthorizing: true,
            },
            ObservationMetricCandidateGate: JEVExternalApplyCapabilityReviewObservationMetricDirectionCandidateGateBridge{
                NonExecuting: true,
                NonAuthorizing: true,
            },
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-reverse-observation" {
        t.Fatalf("invalid reverse observation upstream was not fail-closed: %+v", output)
    }
}
