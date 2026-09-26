package decision

import (
	"testing"
	"time"
)

func generatedCandidateRevisionForGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationBinding {
	t.Helper()
	metric := refutedActionDerivedCandidateRevisionMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-candidate-guard",
		NonAuthorizing:          true,
	})
	direction := candidateRevisionDirectionForGeneration(t, feedback)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/generated-guard",
		RevisionChangeDigest: "generated-guard-change",
		NonAuthorizing:       true,
	})
}

func validGeneratedCandidateRevisionGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_590, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationGuardInput{
		Candidate:             generatedCandidateRevisionForGuard(t),
		CandidateID:           "generated-candidate-revision-1",
		ObservationDigest:     "generated-candidate-observation",
		ObservationAt:         observationAt,
		Now:                   observationAt.Add(10 * time.Second),
		MaxAge:                time.Minute,
		RequestedPath:         "/workspace/project/generated.go",
		AllowedPathPrefix:     "/workspace/project",
		WriteRequested:        true,
		WriteAllowed:          true,
		ConfirmationRequired: true,
		ConfirmationPresent:   true,
		NonAuthorizing:        true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(validGeneratedCandidateRevisionGuardInput(t))
	if output.Status != "admitted" || output.GuardStatus != "admitted" ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" {
		t.Fatalf("output = %#v, want admitted generated candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationKeepsStaleReviewable(t *testing.T) {
	input := validGeneratedCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationRejectsTampering(t *testing.T) {
	input := validGeneratedCandidateRevisionGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-revision-generation-validation" {
		t.Fatalf("output = %#v, want generation validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFailsClosedForPath(t *testing.T) {
	input := validGeneratedCandidateRevisionGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}