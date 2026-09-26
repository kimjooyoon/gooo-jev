package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackGenerationRequestChainUnknown(t *testing.T) {
    output:=BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackGenerationRequestChain(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationFeedbackGenerationRequestChainInput{})
    if output.Status!=jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage!="capability-boundary" || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("empty admission feedback generation chain was not fail-closed: %+v",output) }
    if err:=output.Validate(); err!=nil { t.Fatalf("unknown admission feedback generation chain should validate: %v",err) }
}
