package decision

import "testing"

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection(t *testing.T) {
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:              evaluation,
		CandidateDigest:         evaluation.RevisionCandidateDigest,
		ReplayObservationDigest: "evaluation-replay-observation",
		NonAuthorizing:          true,
	})
	direction := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/from-evaluation",
		NonAuthorizing: true,
	})
	if direction.Status != jevImprovementDirectiveExternalReview ||
		direction.Directive != jevImprovementDirectiveExternalReview ||
		direction.CandidateDigest != feedback.CandidateDigest {
		t.Fatalf("direction = %#v, want external review direction", direction)
	}
	if err := direction.Validate(); err != nil {
		t.Fatalf("direction should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRequiresCandidateSource(t *testing.T) {
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:              evaluation,
		CandidateDigest:         evaluation.RevisionCandidateDigest,
		ReplayObservationDigest: "evaluation-replay-observation",
		NonAuthorizing:          true,
	})
	direction := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput{
		Feedback:       feedback,
		NonAuthorizing: true,
	})
	if direction.Status != "UNKNOWN" || direction.MissingStage != "candidate-source" {
		t.Fatalf("direction = %#v, want candidate-source UNKNOWN", direction)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRejectsTampering(t *testing.T) {
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:              evaluation,
		CandidateDigest:         evaluation.RevisionCandidateDigest,
		ReplayObservationDigest: "evaluation-replay-observation",
		NonAuthorizing:          true,
	})
	feedback.AggregationEvidenceDigest = "tampered"
	direction := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/from-evaluation",
		NonAuthorizing: true,
	})
	if direction.Status != "UNKNOWN" || direction.MissingStage != "evaluation-feedback-validation" {
		t.Fatalf("direction = %#v, want evaluation-feedback-validation UNKNOWN", direction)
	}
}
