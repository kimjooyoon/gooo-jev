package decision

import "testing"

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(t *testing.T) {
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:              evaluation,
		CandidateDigest:         evaluation.RevisionCandidateDigest,
		ReplayObservationDigest: "evaluation-replay-observation",
		NonAuthorizing:          true,
	})
	if binding.Status != "bound" || binding.FeedbackStatus != jevReplayFeedbackConfirmed ||
		binding.AggregationStatus != jevImprovementFeedbackStableForReview ||
		binding.Total != 1 || binding.Confirmed != 1 || binding.EvidenceDigest == "" {
		t.Fatalf("binding = %#v, want stable-for-review feedback", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackRequiresObservation(t *testing.T) {
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:      evaluation,
		CandidateDigest: evaluation.RevisionCandidateDigest,
		NonAuthorizing:  true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "replay-observation" {
		t.Fatalf("binding = %#v, want replay-observation UNKNOWN", binding)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackRejectsTampering(t *testing.T) {
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	evaluation.EvidenceDigest = "tampered"
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:              evaluation,
		CandidateDigest:         evaluation.RevisionCandidateDigest,
		ReplayObservationDigest: "evaluation-replay-observation",
		NonAuthorizing:          true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "evaluation-validation" {
		t.Fatalf("binding = %#v, want evaluation-validation UNKNOWN", binding)
	}
}
