package decision

// ExecutionEnvelopeEvidenceLedgerTransitionInput compares two ordered ledger observations.
type ExecutionEnvelopeEvidenceLedgerTransitionInput struct {
	PreviousStatus       string
	CurrentStatus        string
	PreviousLedgerDigest string
	CurrentLedgerDigest  string
	NonAuthorizing       bool
}

// ExecutionEnvelopeEvidenceLedgerTransition records provenance drift without
// treating a changed ledger as an improvement.
type ExecutionEnvelopeEvidenceLedgerTransition struct {
	Status         string
	Changed        bool
	Direction      string
	FirstMismatch  string
	LedgerDigest   string
	NonAuthorizing bool
}

// CompareExecutionEnvelopeEvidenceLedgerTransition compares status before
// content so the first observable mismatch remains explicit.
func CompareExecutionEnvelopeEvidenceLedgerTransition(input ExecutionEnvelopeEvidenceLedgerTransitionInput) ExecutionEnvelopeEvidenceLedgerTransition {
	output := ExecutionEnvelopeEvidenceLedgerTransition{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.FirstMismatch = "authorization-boundary"
		return output
	}
	if input.PreviousStatus == "" || input.CurrentStatus == "" ||
		input.PreviousLedgerDigest == "" || input.CurrentLedgerDigest == "" {
		output.FirstMismatch = "ledger-transition-evidence"
		return output
	}
	if input.PreviousStatus != input.CurrentStatus {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "status-transition"
		output.FirstMismatch = "ledger-status"
		output.LedgerDigest = input.CurrentLedgerDigest
		return output
	}
	if input.PreviousLedgerDigest != input.CurrentLedgerDigest {
		output.Status = "changed"
		output.Changed = true
		output.Direction = "content-transition"
		output.FirstMismatch = "ledger-content"
		output.LedgerDigest = input.CurrentLedgerDigest
		return output
	}
	output.Status = "unchanged"
	output.Direction = "stable"
	output.LedgerDigest = input.CurrentLedgerDigest
	return output
}