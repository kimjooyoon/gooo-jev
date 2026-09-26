package decision

import "testing"

func candidateDirectionForRevision(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/action-derived-direction",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(t *testing.T) {
	metric := refutedActionDerivedCandidateMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-replay",
		NonAuthorizing:          true,
	})
	direction := candidateDirectionForRevision(t, feedback)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-derived",
		RevisionChangeDigest: "action-derived-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.SourceRevisionCandidateDigest == "" || output.CandidateDigest == "" ||
		output.BoundRevisionChangeDigest == "" {
		t.Fatalf("output = %#v, want bound candidate revision", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRejectsReview(t *testing.T) {
	metric := confirmedActionDerivedCandidateMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-review",
		NonAuthorizing:          true,
	})
	direction := candidateDirectionForRevision(t, feedback)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-derived",
		RevisionChangeDigest: "action-derived-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRejectsHold(t *testing.T) {
	metric := unknownActionDerivedCandidateMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-hold",
		NonAuthorizing:          true,
	})
	direction := candidateDirectionForRevision(t, feedback)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-derived",
		RevisionChangeDigest: "action-derived-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRequiresSource(t *testing.T) {
	metric := refutedActionDerivedCandidateMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-source",
		NonAuthorizing:          true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput{
		Feedback:       feedback,
		Direction:      candidateDirectionForRevision(t, feedback),
		RevisionChangeDigest: "action-derived-revision-change",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-source" {
		t.Fatalf("output = %#v, want revision-source UNKNOWN", output)
	}
}