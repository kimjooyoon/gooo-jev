package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationLSPUnknown(t *testing.T) {
    output:=ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationLSP(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationBridge{Status:jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,MissingStage:"candidate-materialization",NonExecuting:true,NonAuthorizing:true})
    if output.Status!=jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage!="candidate-materialization" || output.Publishable || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("unknown candidate materialization projection was not preserved: %+v",output) }
    if err:=output.Validate(); err!=nil { t.Fatalf("unknown candidate materialization projection should validate: %v",err) }
}
