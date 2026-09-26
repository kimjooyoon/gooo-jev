package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPUnknown(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,
            MissingStage: "application-plan-reverse-observation",
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "application-plan-reverse-observation" || output.Publishable ||
        !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("unknown generation request reverse observation LSP was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown generation request reverse observation LSP should validate: %v", err)
    }
}

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSPRejectsInvalidBridge(t *testing.T) {
    output := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationLSP(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridge{
            Status: jevExternalApplyCapabilityReviewRevisionCandidateGenerationRequestApplicationPlanReverseObservationBridgeBound,
            NonExecuting: true,
            NonAuthorizing: true,
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "generation-request-application-plan-reverse-observation" || output.Publishable {
        t.Fatalf("invalid generation request reverse observation bridge was not fail-closed: %+v", output)
    }
}
