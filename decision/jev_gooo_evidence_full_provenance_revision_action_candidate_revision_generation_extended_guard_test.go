package decision

import (
	"testing"
	"time"
)

func generatedExtendedCandidateRevisionForGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedBinding {
	t.Helper()
	feedback := refutedGeneratedCandidateRevisionFeedbackForDirection(t)
	direction := generatedCandidateRevisionDirectionForExtendedGeneration(t, feedback)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/generated-extended-guard",
		RevisionChangeDigest: "generated-extended-guard-change",
		NonAuthorizing:       true,
	})
}

func validGeneratedExtendedCandidateRevisionGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_690, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardInput{
		Candidate:             generatedExtendedCandidateRevisionForGuard(t),
		CandidateID:           "generated-extended-candidate-revision-1",
		ObservationDigest:     "generated-extended-candidate-observation",
		ObservationAt:         observationAt,
		Now:                   observationAt.Add(10 * time.Second),
		MaxAge:                time.Minute,
		RequestedPath:         "/workspace/project/generated-extended.go",
		AllowedPathPrefix:     "/workspace/project",
		WriteRequested:        true,
		WriteAllowed:          true,
		ConfirmationRequired: true,
		ConfirmationPresent:   true,
		NonAuthorizing:        true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(validGeneratedExtendedCandidateRevisionGuardInput(t))
	if output.Status != "admitted" || output.GuardStatus != "admitted" ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" ||
		output.VerificationEvidenceDigest == "" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want admitted extended generated candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedKeepsStaleReviewable(t *testing.T) {
	input := validGeneratedExtendedCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRejectsTampering(t *testing.T) {
	input := validGeneratedExtendedCandidateRevisionGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-validation" {
		t.Fatalf("output = %#v, want extended generation validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFailsClosedForPath(t *testing.T) {
	input := validGeneratedExtendedCandidateRevisionGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}
