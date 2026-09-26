package decision

import (
	"testing"
	"time"
)

func extendedLineageCandidateRevisionForGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding {
	t.Helper()
	feedback := refutedGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	direction := extendedLineageDirectionForGeneration(t, feedback)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/generated-extended-lineage-guard",
		RevisionChangeDigest: "generated-extended-lineage-guard-change",
		NonAuthorizing:       true,
	})
}

func validExtendedLineageCandidateRevisionGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_790, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardInput{
		Candidate:             extendedLineageCandidateRevisionForGuard(t),
		CandidateID:           "extended-lineage-candidate-revision-1",
		ObservationDigest:     "extended-lineage-candidate-observation",
		ObservationAt:         observationAt,
		Now:                   observationAt.Add(10 * time.Second),
		MaxAge:                time.Minute,
		RequestedPath:         "/workspace/project/generated-extended-lineage.go",
		AllowedPathPrefix:     "/workspace/project",
		WriteRequested:        true,
		WriteAllowed:          true,
		ConfirmationRequired: true,
		ConfirmationPresent:   true,
		NonAuthorizing:        true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(validExtendedLineageCandidateRevisionGuardInput(t))
	if output.Status != "admitted" || output.SourceCandidateID == "" ||
		output.RevisionCandidateDigest == "" || output.GuardEvidenceDigest == "" {
		t.Fatalf("output = %#v, want admitted extended lineage candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageKeepsStaleReviewable(t *testing.T) {
	input := validExtendedLineageCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRejectsTampering(t *testing.T) {
	input := validExtendedLineageCandidateRevisionGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-validation" {
		t.Fatalf("output = %#v, want lineage generation validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFailsClosedForPath(t *testing.T) {
	input := validExtendedLineageCandidateRevisionGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}
