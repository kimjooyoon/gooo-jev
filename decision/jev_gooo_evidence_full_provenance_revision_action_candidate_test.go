package decision

import "testing"

func revisionActionFeedbackForCandidate(t *testing.T, metric ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding, candidate string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding {
	t.Helper()
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         candidate,
		ReplayObservationDigest: "replay-" + candidate,
		NonAuthorizing:          true,
	})
}

func revisionActionDirectionForCandidate(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/action-direction",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(t *testing.T) {
	feedback := revisionActionFeedbackForCandidate(t, counterexampleRevisionActionMetric(t), "candidate-action-revision")
	direction := revisionActionDirectionForCandidate(t, feedback)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-feedback",
		RevisionChangeDigest: "revision-change-action-feedback",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.CandidateDigest == "" || output.BoundRevisionChangeDigest == "" {
		t.Fatalf("output = %#v, want bound revision candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRejectsReview(t *testing.T) {
	feedback := revisionActionFeedbackForCandidate(t, clearRevisionActionMetric(t), "candidate-action-review")
	direction := revisionActionDirectionForCandidate(t, feedback)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-feedback",
		RevisionChangeDigest: "revision-change-action-feedback",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRejectsHold(t *testing.T) {
	feedback := revisionActionFeedbackForCandidate(t, unknownRevisionActionMetric(t), "candidate-action-hold")
	direction := revisionActionDirectionForCandidate(t, feedback)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-feedback",
		RevisionChangeDigest: "revision-change-action-feedback",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRejectsMismatch(t *testing.T) {
	feedback := revisionActionFeedbackForCandidate(t, counterexampleRevisionActionMetric(t), "candidate-action-mismatch")
	direction := revisionActionDirectionForCandidate(t, feedback)
	direction.CandidateDigest = "other-candidate"
	direction.DirectionEvidenceDigest = digestJEVImprovementDirectionDirective(
		direction.DirectiveStatus,
		direction.Directive,
		direction.CandidateDigest,
		direction.CandidateSource,
		direction.InputEvidenceDigest,
	)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-feedback",
		RevisionChangeDigest: "revision-change-action-feedback",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "direction-feedback-candidate" {
		t.Fatalf("output = %#v, want direction-feedback-candidate UNKNOWN", output)
	}
}