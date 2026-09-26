package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeUnknown(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeInput{},
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "capability-boundary" || !output.NonExecuting || !output.NonAuthorizing {
        t.Fatalf("empty candidate gate revision proposal bridge was not fail-closed: %+v", output)
    }
    if err := output.Validate(); err != nil {
        t.Fatalf("unknown candidate gate revision proposal bridge should validate: %v", err)
    }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeRejectsDigestDrift(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge(
        JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridgeInput{
            NonAuthorizing: true,
            ObservationMetricCandidateGate: JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestObservationMetricCandidateGateBridge{
                NonExecuting: true,
                NonAuthorizing: true,
            },
            FeedbackRevisionProposal: JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalBridge{
                NonExecuting: true,
                NonAuthorizing: true,
                CandidateGateDigest: "drift",
            },
        },
    )
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
        output.MissingStage != "observation-metric-candidate-gate" {
        t.Fatalf("invalid candidate gate upstream was not fail-closed: %+v", output)
    }
}
