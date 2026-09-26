package decision

import "testing"

func TestDiagnoseExecutionEnvelopeProvenance(t *testing.T) {
	output := DiagnoseExecutionEnvelopeProvenance(ExecutionEnvelopeProvenanceDiagnosticInput{
		Status: "sealed", EvidencePrefixDigest: "full-digest",
		MissingStageIndex: -1, NonAuthorizing: true,
	})
	if output.Status != "clear" || output.Severity != "info" ||
		output.Code != "provenance-complete" ||
		output.EvidencePrefixDigest != "full-digest" {
		t.Fatalf("unexpected complete diagnostic: %#v", output)
	}

	output = DiagnoseExecutionEnvelopeProvenance(ExecutionEnvelopeProvenanceDiagnosticInput{
		Status: "UNKNOWN", MissingStage: "generation", MissingStageIndex: 2,
		EvidencePrefixDigest: "prefix-digest", NonAuthorizing: true,
	})
	if output.Status != "diagnostic" || output.Severity != "error" ||
		output.Code != "provenance-incomplete" ||
		output.MissingStage != "generation" || output.MissingStageIndex != 2 {
		t.Fatalf("unexpected incomplete diagnostic: %#v", output)
	}

	output = DiagnoseExecutionEnvelopeProvenance(ExecutionEnvelopeProvenanceDiagnosticInput{
		Status: "review-required", MissingStage: "reverse-observation",
		MissingStageIndex: 3, EvidencePrefixDigest: "prefix-digest",
		NonAuthorizing: true,
	})
	if output.Status != "diagnostic" || output.Severity != "warning" ||
		output.Code != "review-required" {
		t.Fatalf("unexpected review diagnostic: %#v", output)
	}

	output = DiagnoseExecutionEnvelopeProvenance(ExecutionEnvelopeProvenanceDiagnosticInput{
		Status: "UNKNOWN", MissingStage: "generation", MissingStageIndex: 2,
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.Code != "ledger-diagnostic-evidence" ||
		output.MissingStageIndex != -1 {
		t.Fatalf("unexpected missing evidence diagnostic: %#v", output)
	}

	output = DiagnoseExecutionEnvelopeProvenance(ExecutionEnvelopeProvenanceDiagnosticInput{
		Status: "sealed", MissingStageIndex: -1,
		EvidencePrefixDigest: "full-digest", NonAuthorizing: false,
	})
	if output.Status != "UNKNOWN" || output.Code != "authorization-boundary" ||
		output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization diagnostic: %#v", output)
	}
}