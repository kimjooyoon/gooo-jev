package decision

import "testing"

func TestBindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationBridgeUnknown(t *testing.T) {
    output:=BindJEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationBridge(JEVExternalApplyCapabilityReviewRevisionCandidateGenerationRequestCandidateGateRevisionProposalCandidateMaterializationAdmissionEvaluationBridgeInput{})
    if output.Status!=jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown || output.MissingStage!="capability-boundary" || !output.NonExecuting || !output.NonAuthorizing { t.Fatalf("empty materialization admission evaluation chain was not fail-closed: %+v",output) }
    if err:=output.Validate(); err!=nil { t.Fatalf("unknown materialization admission evaluation chain should validate: %v",err) }
}
