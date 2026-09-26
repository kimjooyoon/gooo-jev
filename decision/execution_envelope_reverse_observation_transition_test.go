package decision

import (
	"testing"
	"time"
)

func makeReverseObservationTransitionTestValue(
	t *testing.T,
	status ProvenanceObservationStatus,
	missingStage string,
	outputDigest string,
	verifierDigest string,
) ReverseObservation {
	t.Helper()
	observation := ReverseObservation{
		Schema:                 ReverseObservationSchemaV1,
		ExecutionReceiptDigest: "receipt-digest",
		ObservedOutputDigest:   outputDigest,
		VerifierDigest:         verifierDigest,
		Status:                 status,
		MissingStage:           missingStage,
		ObservedAt:             time.Unix(1_800_000_000, 0).UTC(),
	}
	digest, err := observation.computeDigest()
	if err != nil {
		t.Fatalf("computeDigest() error = %v", err)
	}
	observation.ObservationDigest = digest
	if err := observation.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	return observation
}

func TestCompareExecutionEnvelopeReverseObservationTransition(t *testing.T) {
	previous := makeReverseObservationTransitionTestValue(
		t, ProvenanceObserved, "", "output-a", "verifier-a",
	)
	current := makeReverseObservationTransitionTestValue(
		t, ProvenanceObserved, "", "output-a", "verifier-a",
	)
	output := CompareExecutionEnvelopeReverseObservationTransition(previous, current)
	if output.Status != "unchanged" || output.Changed ||
		output.Direction != "stable" || output.ObservationDigest != current.ObservationDigest {
		t.Fatalf("unexpected stable transition: %#v", output)
	}

	current = makeReverseObservationTransitionTestValue(
		t, ProvenanceObserved, "", "output-b", "verifier-a",
	)
	output = CompareExecutionEnvelopeReverseObservationTransition(previous, current)
	if output.Status != "changed" || !output.Changed ||
		output.Direction != "output-transition" ||
		output.FirstMismatch != "reverse-observation-output" {
		t.Fatalf("unexpected output transition: %#v", output)
	}

	current = makeReverseObservationTransitionTestValue(
		t, ProvenanceReview, "review", "", "",
	)
	output = CompareExecutionEnvelopeReverseObservationTransition(previous, current)
	if output.Status != "changed" || output.Direction != "status-transition" ||
		output.FirstMismatch != "reverse-observation-status" {
		t.Fatalf("unexpected status transition: %#v", output)
	}

	invalid := current
	invalid.MissingStage = "tampered"
	output = CompareExecutionEnvelopeReverseObservationTransition(previous, invalid)
	if output.Status != "UNKNOWN" || output.Changed ||
		output.FirstMismatch != "current-reverse-observation" {
		t.Fatalf("unexpected invalid transition: %#v", output)
	}

	invalid = previous
	invalid.ObservationDigest = "tampered"
	output = CompareExecutionEnvelopeReverseObservationTransition(invalid, current)
	if output.Status != "UNKNOWN" || output.Changed ||
		output.FirstMismatch != "previous-reverse-observation" {
		t.Fatalf("unexpected previous invalid transition: %#v", output)
	}
}