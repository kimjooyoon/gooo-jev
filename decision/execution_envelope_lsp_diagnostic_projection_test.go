package decision

import "testing"

func TestProjectExecutionEnvelopeLSPDiagnostic(t *testing.T) {
	output := ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
		Status: "diagnostic", Severity: "error", Code: "provenance-incomplete",
		MissingStage: "generation", MissingStageIndex: 2,
		EvidencePrefixDigest: "prefix-digest", NonAuthorizing: true,
	})
	if output.Status != "publishable" || output.Publishable != true ||
		output.Severity != "error" || output.Code != "provenance-incomplete" ||
		output.MissingStageIndex != 2 ||
		output.EvidencePrefixDigest != "prefix-digest" {
		t.Fatalf("unexpected publishable diagnostic: %#v", output)
	}

	output = ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
		Status: "clear", Severity: "info", Code: "provenance-complete",
		EvidencePrefixDigest: "full-digest", MissingStageIndex: -1,
		NonAuthorizing: true,
	})
	if output.Status != "clear" || output.Publishable ||
		output.Code != "provenance-complete" ||
		output.EvidencePrefixDigest != "full-digest" {
		t.Fatalf("unexpected clear projection: %#v", output)
	}

	output = ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
		Status: "diagnostic", Severity: "error", Code: "provenance-incomplete",
		MissingStage: "generation", MissingStageIndex: 2,
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.Publishable ||
		output.Code != "lsp-diagnostic-evidence" {
		t.Fatalf("unexpected incomplete projection: %#v", output)
	}

	output = ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
		Status: "diagnostic", Severity: "error", Code: "provenance-incomplete",
		MissingStage: "generation", MissingStageIndex: 2,
		EvidencePrefixDigest: "prefix-digest", NonAuthorizing: false,
	})
	if output.Status != "UNKNOWN" || output.Publishable ||
		output.Code != "authorization-boundary" || output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization projection: %#v", output)
	}
}