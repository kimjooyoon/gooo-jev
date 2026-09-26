package decision

import (
	"testing"
	"time"
)

func extendedLineageCandidateRevisionForRevisionGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionBinding {
	t.Helper()
	feedback := candidateFeedbackForRevisionGeneration(t)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            candidateDirectionForRevisionGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-candidate-revision-guard",
		RevisionChangeDigest: "extended-lineage-candidate-revision-guard-change",
		NonAuthorizing:       true,
	})
}

func validExtendedLineageCandidateRevisionCandidateGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_910, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionGuardInput{
		Candidate:            extendedLineageCandidateRevisionForRevisionGuard(t),
		CandidateID:          "extended-lineage-candidate-revision-guard-1",
		ObservationDigest:    "extended-lineage-candidate-revision-guard-observation",
		ObservationAt:        observationAt,
		Now:                  observationAt.Add(10 * time.Second),
		MaxAge:               time.Minute,
		RequestedPath:        "/workspace/project/generated-extended-lineage-candidate-revision.go",
		AllowedPathPrefix:    "/workspace/project",
		WriteRequested:       true,
		WriteAllowed:         true,
		ConfirmationRequired: true,
		ConfirmationPresent:  true,
		NonAuthorizing:      true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(validExtendedLineageCandidateRevisionCandidateGuardInput(t))
	if output.Status != "admitted" || output.SourceCandidateID == "" ||
		output.ParentCandidateDigest == "" || output.CandidateDigest == "" ||
		output.GuardEvidenceDigest == "" {
		t.Fatalf("output = %#v, want admitted extended lineage candidate revision guard", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionKeepsStaleReviewable(t *testing.T) {
	input := validExtendedLineageCandidateRevisionCandidateGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionRejectsTampering(t *testing.T) {
	input := validExtendedLineageCandidateRevisionCandidateGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-validation" {
		t.Fatalf("output = %#v, want candidate-revision validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFailsClosedForPath(t *testing.T) {
	input := validExtendedLineageCandidateRevisionCandidateGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}