package decision

import (
	"testing"
	"time"
)

func validJEVActionCandidateInput() JEVActionCandidateGuardInput {
	observationAt := time.Unix(1_800_000_190, 0).UTC()
	return JEVActionCandidateGuardInput{
		CandidateID:          "candidate-1",
		ObservationDigest:    "observation-digest",
		ObservationAt:        observationAt,
		Now:                  observationAt.Add(10 * time.Second),
		MaxAge:               time.Minute,
		RequestedPath:        "/workspace/project/main.go",
		AllowedPathPrefix:    "/workspace/project",
		WriteRequested:       true,
		WriteAllowed:         true,
		ConfirmationRequired: true,
		ConfirmationPresent:  true,
		NonAuthorizing:      true,
	}
}

func TestGuardJEVActionCandidateAdmitsFreshBoundedCandidate(t *testing.T) {
	output := GuardJEVActionCandidate(validJEVActionCandidateInput())
	if output.Status != "admitted" || output.MissingStage != "" {
		t.Fatalf("unexpected admitted output: %+v", output)
	}
	if output.EvidenceDigest == "" || !output.NonExecuting || !output.NonAuthorizing {
		t.Fatalf("candidate guard lost evidence or safety boundary: %+v", output)
	}
}

func TestGuardJEVActionCandidateKeepsFreshnessAndPathReviewable(t *testing.T) {
	input := validJEVActionCandidateInput()
	input.Now = input.ObservationAt.Add(2 * time.Minute)
	output := GuardJEVActionCandidate(input)
	if output.Status != "review" || output.MissingStage != "observation-stale" {
		t.Fatalf("unexpected stale output: %+v", output)
	}

	input = validJEVActionCandidateInput()
	input.RequestedPath = "/workspace/project/../secret.txt"
	output = GuardJEVActionCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("unexpected path output: %+v", output)
	}

	input = validJEVActionCandidateInput()
	input.RequestedPath = "/workspace/project-other/main.go"
	output = GuardJEVActionCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "path-scope" {
		t.Fatalf("unexpected path-boundary output: %+v", output)
	}
}

func TestGuardJEVActionCandidateRejectsWriteAndRequiresConfirmation(t *testing.T) {
	input := validJEVActionCandidateInput()
	input.WriteAllowed = false
	output := GuardJEVActionCandidate(input)
	if output.Status != "rejected" || output.MissingStage != "write-permission" {
		t.Fatalf("unexpected write output: %+v", output)
	}

	input = validJEVActionCandidateInput()
	input.ConfirmationPresent = false
	output = GuardJEVActionCandidate(input)
	if output.Status != "review" || output.MissingStage != "confirmation" {
		t.Fatalf("unexpected confirmation output: %+v", output)
	}
}

func TestGuardJEVActionCandidateFailsClosedForFutureObservationAndAuthorization(t *testing.T) {
	input := validJEVActionCandidateInput()
	input.ObservationAt = input.Now.Add(time.Second)
	output := GuardJEVActionCandidate(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "observation-future" {
		t.Fatalf("unexpected future observation output: %+v", output)
	}

	input = validJEVActionCandidateInput()
	input.NonAuthorizing = false
	output = GuardJEVActionCandidate(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
