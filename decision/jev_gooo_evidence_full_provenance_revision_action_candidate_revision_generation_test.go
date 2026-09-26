package decision

import "testing"

func candidateRevisionDirectionForGeneration(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/revision-generation",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(t *testing.T) {
	metric := refutedActionDerivedCandidateRevisionMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-generation",
		NonAuthorizing:          true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationInput{
		Feedback:             feedback,
		Direction:            candidateRevisionDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/action-derived-next",
		RevisionChangeDigest: "action-derived-next-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.SourceRevisionCandidateDigest == "" || output.CandidateDigest == "" ||
		output.BoundRevisionChangeDigest == "" {
		t.Fatalf("output = %#v, want bound next revision candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationRejectsReview(t *testing.T) {
	metric := confirmedActionDerivedCandidateRevisionMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-generation-review",
		NonAuthorizing:          true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationInput{
		Feedback:             feedback,
		Direction:            candidateRevisionDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/action-derived-next",
		RevisionChangeDigest: "action-derived-next-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationRejectsHold(t *testing.T) {
	metric := unknownActionDerivedCandidateRevisionMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-generation-hold",
		NonAuthorizing:          true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationInput{
		Feedback:             feedback,
		Direction:            candidateRevisionDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/action-derived-next",
		RevisionChangeDigest: "action-derived-next-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationRequiresChange(t *testing.T) {
	metric := refutedActionDerivedCandidateRevisionMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-generation-change",
		NonAuthorizing:          true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationInput{
		Feedback:       feedback,
		Direction:      candidateRevisionDirectionForGeneration(t, feedback),
		RevisionSource: "gooo://revision/action-derived-next",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-change" {
		t.Fatalf("output = %#v, want revision-change UNKNOWN", output)
	}
}