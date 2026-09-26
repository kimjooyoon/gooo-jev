package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackCandidateGenerationLoopUnknown(t *testing.T) {
    output:=ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackCandidateGenerationLoopLSP(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackCandidateGenerationLoopBridge{Status:jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,MissingStage:"feedback-candidate-generation-loop",NonExecuting:true,NonAuthorizing:true})
    if output.Status!=jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage!="feedback-candidate-generation-loop" || output.Publishable || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("unknown feedback candidate generation loop projection was not preserved: %+v",output) }
    if err:=output.Validate(); err!=nil { t.Fatalf("unknown feedback candidate generation loop projection should validate: %v",err) }
}
