package decision

import (
	"testing"
	"time"
)

func extendedLineageCandidateForGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateBinding {
	t.Helper()
	feedback := refutedExtendedLineageCandidateRevisionFeedbackForDirection(t)
	direction := extendedLineageCandidateDirectionForGeneration(t, feedback)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/extended-lineage-candidate-guard",
		RevisionChangeDigest: "extended-lineage-candidate-guard-change",
		NonAuthorizing:       true,
	})
}

func validExtendedLineageCandidateGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_850, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardInput{
		Candidate:            extendedLineageCandidateForGuard(t),
		CandidateID:          "extended-lineage-candidate-guard-1",
		ObservationDigest:    "extended-lineage-candidate-guard-observation",
		ObservationAt:        observationAt,
		Now:                  observationAt.Add(10 * time.Second),
		MaxAge:               time.Minute,
		RequestedPath:        "/workspace/project/generated-extended-lineage-candidate.go",
		AllowedPathPrefix:    "/workspace/project",
		WriteRequested:       true,
		WriteAllowed:         true,
		ConfirmationRequired: true,
		ConfirmationPresent:  true,
		NonAuthorizing:       true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(validExtendedLineageCandidateGuardInput(t))
	if output.Status != "admitted" || output.SourceCandidateID == "" ||
		output.CandidateDigest == "" || output.CandidateEvidenceDigest == "" ||
		output.GuardEvidenceDigest == "" {
		t.Fatalf("output = %#v, want admitted extended lineage candidate guard", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateKeepsStaleReviewable(t *testing.T) {
	input := validExtendedLineageCandidateGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRejectsTampering(t *testing.T) {
	input := validExtendedLineageCandidateGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-validation" {
		t.Fatalf("output = %#v, want candidate-validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFailsClosedForPath(t *testing.T) {
	input := validExtendedLineageCandidateGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}