package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeUnknown(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeInput{})
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage != "capability-boundary" || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("empty candidate materialization bridge was not fail-closed: %+v",output) }
    if err := output.Validate(); err != nil { t.Fatalf("unknown candidate materialization bridge should validate: %v",err) }
}

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeRejectsUnboundProposal(t *testing.T) {
    output := BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridgeInput{NonAuthorizing:true,RevisionProposal:JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalBridge{NonExecuting:true,NonAuthorizing:true},CandidateMaterialization:JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge{NonExecuting:true,NonAuthorizing:true}})
    if output.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage != "candidate-gate-revision-proposal" { t.Fatalf("unbound proposal was not fail-closed: %+v",output) }
}
