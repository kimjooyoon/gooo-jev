package decision

import (
	"testing"
	"time"
)

func revisionActionGuardCandidate(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding {
	t.Helper()
	feedback := directionRevisionCandidateFeedback(t)
	direction := directionRevisionCandidateDirective(feedback)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/from-feedback",
		RevisionChangeDigest: "revision-change-from-feedback",
		NonAuthorizing:       true,
	})
}

func validRevisionActionGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardInput {
	observationAt := time.Unix(1_800_000_190, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardInput{
		RevisionCandidate:    revisionActionGuardCandidate(t),
		CandidateID:          "revision-action-1",
		ObservationDigest:    "revision-observation-digest",
		ObservationAt:        observationAt,
		Now:                  observationAt.Add(10 * time.Second),
		MaxAge:               time.Minute,
		RequestedPath:        "/workspace/project/main.go",
		AllowedPathPrefix:    "/workspace/project",
		WriteRequested:       true,
		WriteAllowed:         true,
		ConfirmationRequired: true,
		ConfirmationPresent:  true,
		NonAuthorizing:       true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionAdmitsTypedCandidate(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(validRevisionActionGuardInput(t))
	if output.Status != "admitted" || output.GuardStatus != "admitted" || output.MissingStage != "" {
		t.Fatalf("output = %#v, want admitted", output)
	}
	if output.RevisionCandidateDigest == "" || output.GuardEvidenceDigest == "" {
		t.Fatalf("output lost revision or guard provenance: %#v", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionKeepsStaleCandidateReviewable(t *testing.T) {
	input := validRevisionActionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should preserve evidence: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionRejectsTamperedRevision(t *testing.T) {
	input := validRevisionActionGuardInput(t)
	input.RevisionCandidate.RevisionCandidateDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-candidate-validation" {
		t.Fatalf("output = %#v, want revision-candidate-validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFailsClosedForPathAndAuthorization(t *testing.T) {
	input := validRevisionActionGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.txt"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}

	input = validRevisionActionGuardInput(t)
	input.NonAuthorizing = false
	output = GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "authorization-boundary" {
		t.Fatalf("output = %#v, want authorization-boundary UNKNOWN", output)
	}
}