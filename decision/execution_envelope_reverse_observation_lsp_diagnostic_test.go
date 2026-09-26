package decision

import (
	"testing"
	"time"
)

func validReverseObservationForDiagnostic(t *testing.T, status ProvenanceObservationStatus, missingStage string) ReverseObservation {
	t.Helper()
	observation := ReverseObservation{
		Schema:                 ReverseObservationSchemaV1,
		ExecutionReceiptDigest: "receipt-digest",
		Status:                 status,
		MissingStage:           missingStage,
		ObservedAt:             time.Unix(1_800_000_000, 0).UTC(),
	}
	if status == ProvenanceObserved {
		observation.ObservedOutputDigest = "output-digest"
		observation.VerifierDigest = "verifier-digest"
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

func TestProjectExecutionEnvelopeReverseObservationLSPDiagnostic(t *testing.T) {
	observed := validReverseObservationForDiagnostic(t, ProvenanceObserved, "")
	output := ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(observed, -1, "full-prefix")
	if output.Status != "clear" || output.Publishable ||
		output.Code != "provenance-complete" ||
		output.EvidencePrefixDigest != "full-prefix" {
		t.Fatalf("unexpected observed projection: %#v", output)
	}

	unknown := validReverseObservationForDiagnostic(t, ProvenanceUnknown, "execution-receipt")
	output = ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(unknown, 0, "prefix-unknown")
	if output.Status != "publishable" || !output.Publishable ||
		output.Severity != "error" || output.Code != "provenance-incomplete" ||
		output.MissingStageIndex != 0 {
		t.Fatalf("unexpected unknown projection: %#v", output)
	}

	review := validReverseObservationForDiagnostic(t, ProvenanceReview, "review")
	output = ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(review, 1, "prefix-review")
	if output.Status != "publishable" || !output.Publishable ||
		output.Severity != "warning" || output.Code != "review-required" ||
		output.MissingStageIndex != 1 {
		t.Fatalf("unexpected review projection: %#v", output)
	}

	output = ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(unknown, 0, "")
	if output.Status != "UNKNOWN" || output.Publishable ||
		output.Code != "reverse-observation-prefix" {
		t.Fatalf("unexpected prefix projection: %#v", output)
	}

	tampered := unknown
	tampered.MissingStage = "tampered"
	output = ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(tampered, 0, "prefix")
	if output.Status != "UNKNOWN" || output.Publishable ||
		output.Code != "reverse-observation-integrity" {
		t.Fatalf("unexpected integrity projection: %#v", output)
	}

	output = ProjectExecutionEnvelopeReverseObservationLSPDiagnostic(unknown, -1, "prefix")
	if output.Status != "UNKNOWN" || output.Publishable ||
		output.Code != "reverse-observation-location" {
		t.Fatalf("unexpected location projection: %#v", output)
	}
}