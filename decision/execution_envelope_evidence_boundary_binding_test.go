package decision

import "testing"

func TestBindExecutionEnvelopeEvidenceBoundary(t *testing.T) {
	input := ExecutionEnvelopeEvidenceBoundaryBindingInput{
		LedgerDigest:   "ledger-digest",
		BoundaryDigest: "boundary-digest",
		NonAuthorizing: true,
	}
	output := BindExecutionEnvelopeEvidenceBoundary(input)
	if output.Status != "bound" || output.BindingDigest == "" ||
		output.MissingStage != "" || output.NonAuthorizing != true {
		t.Fatalf("unexpected binding: %#v", output)
	}

	input.BoundaryDigest = ""
	output = BindExecutionEnvelopeEvidenceBoundary(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "capability-boundary" ||
		output.BindingDigest != "" {
		t.Fatalf("unexpected missing boundary: %#v", output)
	}

	input.LedgerDigest = ""
	input.BoundaryDigest = "boundary-digest"
	output = BindExecutionEnvelopeEvidenceBoundary(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "evidence-ledger" {
		t.Fatalf("unexpected missing ledger: %#v", output)
	}

	input.LedgerDigest = "ledger-digest"
	input.NonAuthorizing = false
	output = BindExecutionEnvelopeEvidenceBoundary(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "authorization-boundary" ||
		output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization boundary: %#v", output)
	}
}