package decision

// ExecutionEnvelopeAuditDispositionInput combines already-observed evidence
// states without granting execution authority.
type ExecutionEnvelopeAuditDispositionInput struct {
	LedgerStatus        string
	BoundaryStatus      string
	CounterexampleStatus string
	NonAuthorizing      bool
}

// ExecutionEnvelopeAuditDisposition classifies a safe observation mode.
type ExecutionEnvelopeAuditDisposition struct {
	Status         string
	Mode           string
	FirstMissing   string
	NonAuthorizing bool
}

// ClassifyExecutionEnvelopeAuditDisposition never returns an execution grant;
// it only selects audit, review, or blocked handling.
func ClassifyExecutionEnvelopeAuditDisposition(input ExecutionEnvelopeAuditDispositionInput) ExecutionEnvelopeAuditDisposition {
	output := ExecutionEnvelopeAuditDisposition{Status: "UNKNOWN", Mode: "blocked", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.FirstMissing = "authorization-boundary"
		return output
	}
	if input.LedgerStatus == "" {
		output.FirstMissing = "evidence-ledger-status"
		return output
	}
	if input.BoundaryStatus == "" {
		output.FirstMissing = "capability-boundary-status"
		return output
	}
	if input.CounterexampleStatus == "" {
		output.FirstMissing = "counterexample-disposition"
		return output
	}
	if input.LedgerStatus != "sealed" {
		output.FirstMissing = "evidence-ledger-status"
		return output
	}
	if input.BoundaryStatus != "bound" {
		output.FirstMissing = "capability-boundary-status"
		return output
	}
	switch input.CounterexampleStatus {
	case "review-required":
		output.Status = "review-required"
		output.Mode = "review-only"
	case "observation-only":
		output.Status = "observation-only"
		output.Mode = "audit-only"
	case "hold", "UNKNOWN":
		output.Status = "hold"
		output.Mode = "blocked"
		output.FirstMissing = "counterexample-evidence"
	default:
		output.FirstMissing = "counterexample-disposition"
	}
	return output
}