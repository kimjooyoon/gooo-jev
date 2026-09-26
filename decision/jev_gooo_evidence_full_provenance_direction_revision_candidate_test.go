package decision

import "testing"

func directionRevisionCandidateFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding {
	t.Helper()
	evaluationInput := evidenceFullProvenanceMaterializationEvaluationInput(t)
	evaluation := BindExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluation(evaluationInput)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput{
		Evaluation:              evaluation,
		CandidateDigest:         evaluation.RevisionCandidateDigest,
		ReplayObservationDigest: "direction-revision-observation",
		NonAuthorizing:          true,
	})
}

func directionRevisionCandidateDirective(feedback ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding) JEVImprovementDirectionDirective {
	direction := JEVImprovementDirectionDirective{
		Status:              jevImprovementDirectiveRevision,
		Directive:           jevImprovementDirectiveRevision,
		CandidateDigest:     feedback.CandidateDigest,
		CandidateSource:     "gooo://candidate/from-feedback",
		InputEvidenceDigest: feedback.AggregationEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	direction.EvidenceDigest = digestJEVImprovementDirectionDirective(
		direction.Status,
		direction.Directive,
		direction.CandidateDigest,
		direction.CandidateSource,
		direction.InputEvidenceDigest,
	)
	return direction
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate(t *testing.T) {
	feedback := directionRevisionCandidateFeedback(t)
	direction := directionRevisionCandidateDirective(feedback)
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/from-feedback",
		RevisionChangeDigest: "revision-change-from-feedback",
		NonAuthorizing:       true,
	})
	if binding.Status != jevGoooEvidenceFullProvenanceDirectionRevisionCandidateBound ||
		binding.RevisionCandidateStatus != jevImprovementRevisionCandidateReady ||
		binding.CandidateDigest != feedback.CandidateDigest ||
		binding.RevisionCandidateDigest == "" {
		t.Fatalf("binding = %#v, want bound revision candidate", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateRejectsReview(t *testing.T) {
	feedback := directionRevisionCandidateFeedback(t)
	direction := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput{
		Feedback:        feedback,
		CandidateSource: "gooo://candidate/from-feedback",
		NonAuthorizing:  true,
	})
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/from-feedback",
		RevisionChangeDigest: "revision-change-from-feedback",
		NonAuthorizing:       true,
	})
	if binding.Status != jevGoooEvidenceFullProvenanceDirectionRevisionCandidateUnknown ||
		binding.MissingStage != "revision-directive" {
		t.Fatalf("binding = %#v, want revision-directive UNKNOWN", binding)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateRejectsCandidateMismatch(t *testing.T) {
	feedback := directionRevisionCandidateFeedback(t)
	direction := directionRevisionCandidateDirective(feedback)
	direction.CandidateDigest = "other-candidate"
	direction.EvidenceDigest = digestJEVImprovementDirectionDirective(
		direction.Status,
		direction.Directive,
		direction.CandidateDigest,
		direction.CandidateSource,
		direction.InputEvidenceDigest,
	)
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/from-feedback",
		RevisionChangeDigest: "revision-change-from-feedback",
		NonAuthorizing:       true,
	})
	if binding.Status != jevGoooEvidenceFullProvenanceDirectionRevisionCandidateUnknown ||
		binding.MissingStage != "direction-feedback-candidate" {
		t.Fatalf("binding = %#v, want direction-feedback-candidate UNKNOWN", binding)
	}
}