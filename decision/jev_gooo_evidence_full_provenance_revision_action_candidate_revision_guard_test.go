package decision

import (
	"testing"
	"time"
)

func actionDerivedCandidateRevisionForGuard(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionBinding {
	t.Helper()
	metric := refutedActionDerivedCandidateMetric(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-guard",
		NonAuthorizing:          true,
	})
	direction := candidateDirectionForRevision(t, feedback)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/action-derived-guard",
		RevisionChangeDigest: "action-derived-guard-change",
		NonAuthorizing:       true,
	})
}

func validActionDerivedCandidateRevisionGuardInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardInput {
	t.Helper()
	observationAt := time.Unix(1_800_000_390, 0).UTC()
	return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardInput{
		Candidate:             actionDerivedCandidateRevisionForGuard(t),
		CandidateID:           "action-derived-candidate-revision-1",
		ObservationDigest:     "action-derived-candidate-revision-observation",
		ObservationAt:         observationAt,
		Now:                   observationAt.Add(10 * time.Second),
		MaxAge:                time.Minute,
		RequestedPath:         "/workspace/project/revision.go",
		AllowedPathPrefix:     "/workspace/project",
		WriteRequested:        true,
		WriteAllowed:          true,
		ConfirmationRequired: true,
		ConfirmationPresent:  true,
		NonAuthorizing:       true,
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(t *testing.T) {
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(validActionDerivedCandidateRevisionGuardInput(t))
	if output.Status != "admitted" || output.GuardStatus != "admitted" ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" {
		t.Fatalf("output = %#v, want admitted candidate revision", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionKeepsStaleReviewable(t *testing.T) {
	input := validActionDerivedCandidateRevisionGuardInput(t)
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("output = %#v, want observation-stale review", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("review output should validate: %v", err)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionRejectsTampering(t *testing.T) {
	input := validActionDerivedCandidateRevisionGuardInput(t)
	input.Candidate.CandidateEvidenceDigest = "tampered"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-revision-validation" {
		t.Fatalf("output = %#v, want candidate revision validation UNKNOWN", output)
	}
}

func TestGuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFailsClosedForPath(t *testing.T) {
	input := validActionDerivedCandidateRevisionGuardInput(t)
	input.RequestedPath = "/workspace/project/../secret.go"
	output := GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("output = %#v, want path-scope UNKNOWN", output)
	}
}