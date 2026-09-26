package decision

import (
	"testing"
	"time"
)

func actionDerivedRevisionCandidateForGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding {
	t.Helper()
	feedback := revisionActionFeedbackForCandidate(t, counterexampleRevisionActionMetric(t), "candidate-guard")
	direction := revisionActionDirectionForCandidate(t, feedback)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-guard",
		RevisionChangeDigest: "revision-change-action-guard",
		NonAuthorizing:       true,
	})
}

func validActionDerivedRevisionCandidateGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_190, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardInput{
		Candidate:            actionDerivedRevisionCandidateForGuard(t),
		CandidateID:          "action-derived-revision-1",
		ObservationDigest:    "action-derived-observation",
		ObservationAt:        observationAt,
		Now:                  observationAt.Add(10 * time.Second),
		MaxAge:               time.Minute,
		RequestedPath:        "/workspace/project/revision.go",
		AllowedPathPrefix:    "/workspace/project",
		WriteRequested:       true,
		WriteAllowed:         true,
		ConfirmationRequired: true,
		ConfirmationPresent:  true,
		NonAuthorizing:       true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(validActionDerivedRevisionCandidateGuardInput(t))
	if output.Status != "admitted" || output.GuardStatus != "admitted" || output.MissingStage != "" {
		t.Fatalf("output = %#v, want admitted action-derived candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateKeepsStaleReviewable(t *testing.T) {
	input := validActionDerivedRevisionCandidateGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRejectsTampering(t *testing.T) {
	input := validActionDerivedRevisionCandidateGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-validation" {
		t.Fatalf("output = %#v, want candidate validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFailsClosedForPath(t *testing.T) {
	input := validActionDerivedRevisionCandidateGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}