package decision

import "testing"

func TestBindExecutionEnvelopeReverseObservationFromSource(t *testing.T) {
	input := ExecutionEnvelopeReverseObservationSourceInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "evidence-1",
		ExpectedEvidenceDigest: "evidence-1",
		SourceText:             "reverse_observation: ready evidence-1",
		NonAuthorizing:         true,
	}
	reproduced := BindExecutionEnvelopeReverseObservationFromSource(input)
	if reproduced.Status != "reproduced" {
		t.Fatalf("status = %q, want reproduced", reproduced.Status)
	}
	if reproduced.SourceDigest == "" || reproduced.ObservationDigest == "" {
		t.Fatal("reproduced observation must retain both source and observation digests")
	}
	changed := input
	changed.SourceText += "\n"
	changedSource := BindExecutionEnvelopeReverseObservationFromSource(changed)
	if changedSource.SourceDigest == reproduced.SourceDigest {
		t.Fatal("changed reverse-observation source must change its source digest")
	}
	counterexample := input
	counterexample.ExpectedEvidenceDigest = "different-evidence"
	counterexampleBinding := BindExecutionEnvelopeReverseObservationFromSource(counterexample)
	if counterexampleBinding.Status != "counterexample" {
		t.Fatalf("counterexample status = %q, want counterexample", counterexampleBinding.Status)
	}
	if counterexampleBinding.FirstMismatch == "" {
		t.Fatal("counterexample must preserve its first mismatch")
	}
	missingSource := input
	missingSource.SourceText = ""
	missingBinding := BindExecutionEnvelopeReverseObservationFromSource(missingSource)
	if missingBinding.Status != "UNKNOWN" || missingBinding.FirstMismatch != "reverse-observation-source" {
		t.Fatalf("missing source = %#v, want UNKNOWN at reverse-observation-source", missingBinding)
	}
	unauthorized := input
	unauthorized.NonAuthorizing = false
	unauthorizedBinding := BindExecutionEnvelopeReverseObservationFromSource(unauthorized)
	if unauthorizedBinding.Status != "UNKNOWN" || unauthorizedBinding.NonAuthorizing || unauthorizedBinding.FirstMismatch != "authorization-boundary" {
		t.Fatalf("unauthorized source = %#v, want non-authorizing UNKNOWN", unauthorizedBinding)
	}
}
