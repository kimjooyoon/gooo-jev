package decision

import "testing"

func TestClassifyExecutionEnvelopeAuditDisposition(t *testing.T) {
	base := ExecutionEnvelopeAuditDispositionInput{
		LedgerStatus: "sealed", BoundaryStatus: "bound",
		CounterexampleStatus: "observation-only", NonAuthorizing: true,
	}
	output := ClassifyExecutionEnvelopeAuditDisposition(base)
	if output.Status != "observation-only" || output.Mode != "audit-only" ||
		output.FirstMissing != "" || output.NonAuthorizing != true {
		t.Fatalf("unexpected audit disposition: %#v", output)
	}

	base.CounterexampleStatus = "review-required"
	output = ClassifyExecutionEnvelopeAuditDisposition(base)
	if output.Status != "review-required" || output.Mode != "review-only" {
		t.Fatalf("unexpected review disposition: %#v", output)
	}

	base.CounterexampleStatus = "hold"
	output = ClassifyExecutionEnvelopeAuditDisposition(base)
	if output.Status != "hold" || output.Mode != "blocked" ||
		output.FirstMissing != "counterexample-evidence" {
		t.Fatalf("unexpected hold disposition: %#v", output)
	}

	base.CounterexampleStatus = "observation-only"
	base.BoundaryStatus = "UNKNOWN"
	output = ClassifyExecutionEnvelopeAuditDisposition(base)
	if output.Status != "UNKNOWN" || output.Mode != "blocked" ||
		output.FirstMissing != "capability-boundary-status" {
		t.Fatalf("unexpected boundary disposition: %#v", output)
	}

	base.BoundaryStatus = "bound"
	base.NonAuthorizing = false
	output = ClassifyExecutionEnvelopeAuditDisposition(base)
	if output.Status != "UNKNOWN" || output.Mode != "blocked" ||
		output.FirstMissing != "authorization-boundary" ||
		output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization disposition: %#v", output)
	}
}