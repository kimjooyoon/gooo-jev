package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPUnknown(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: "generation-request-reverse-observation",
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-reverse-observation" || output.Publishable ||
        !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("unknown generation request candidate gate LSP was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request candidate gate LSP should validate: %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSPRejectsInvalidBridge(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridgeBound,
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-observation-metric-candidate-gate" || output.Publishable {
        t.Fatalf("invalid generation request candidate gate bridge was not fail-closed: %+v", output)
    }
}
