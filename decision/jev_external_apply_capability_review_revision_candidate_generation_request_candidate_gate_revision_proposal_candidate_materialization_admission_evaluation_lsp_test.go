package decision

import "testing"

func TestProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationLSPUnknown(t *testing.T) {
    output:=ProjectJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationLSP(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationBridge{Status:jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown,MissingStage:"admission-evaluation",NonExecuting:true,NonAuthorizing:true})
    if output.Status!=jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage!="admission-evaluation" || output.Publishable || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("unknown materialization admission evaluation projection was not preserved: %+v",output) }
    if err:=output.Validate(); err!=nil { t.Fatalf("unknown materialization admission evaluation projection should validate: %v",err) }
}
