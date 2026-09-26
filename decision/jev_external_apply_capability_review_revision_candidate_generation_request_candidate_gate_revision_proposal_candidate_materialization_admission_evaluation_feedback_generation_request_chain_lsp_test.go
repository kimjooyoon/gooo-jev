package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackGenerationRequestChainUnknown(t *testing.T) {
    output:=ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackGenerationRequestChainLSP(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackGenerationRequestChain{Status:jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,MissingStage:"feedback-generation",NonExecuting:true,NonAuthorizing:true})
    if output.Status!=jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage!="feedback-generation" || output.Publishable || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("unknown admission feedback generation projection was not preserved: %+v",output) }
    if err:=output.Validate(); err!=nil { t.Fatalf("unknown admission feedback generation projection should validate: %v",err) }
}
