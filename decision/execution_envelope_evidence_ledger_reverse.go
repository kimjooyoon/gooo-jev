package decision

// ExecutionEnvelopeEvidenceLedgerReverseInput compares an observed sealed
// ledger with an expected sealed ledger.
type ExecutionEnvelopeEvidenceLedgerReverseInput struct {
	ObservedStatus         string `json:"observed_status"`
	ExpectedStatus         string `json:"expected_status"`
	ObservedLedgerDigest   string `json:"observed_ledger_digest"`
	ExpectedLedgerDigest   string `json:"expected_ledger_digest"`
	ObservedMissingStage   string `json:"observed_missing_stage"`
	ExpectedMissingStage   string `json:"expected_missing_stage"`
	NonAuthorizing         bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeEvidenceLedgerReverse records exact reproduction or the
// first ledger divergence without treating reproduction as authorization.
type ExecutionEnvelopeEvidenceLedgerReverse struct {
	Status         string `json:"status"`
	FirstMismatch  string `json:"first_mismatch"`
	LedgerDigest   string `json:"ledger_digest"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ObserveExecutionEnvelopeEvidenceLedgerReverse is fail-closed for incomplete
// expected/observed ledger evidence.
func ObserveExecutionEnvelopeEvidenceLedgerReverse(input ExecutionEnvelopeEvidenceLedgerReverseInput) ExecutionEnvelopeEvidenceLedgerReverse {
	output := ExecutionEnvelopeEvidenceLedgerReverse{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.FirstMismatch = "authorization-boundary"
		return output
	}
	checks := []struct {
		name     string
		observed string
		expected string
	}{
		{name: "status", observed: input.ObservedStatus, expected: input.ExpectedStatus},
		{name: "missing_stage", observed: input.ObservedMissingStage, expected: input.ExpectedMissingStage},
		{name: "ledger_digest", observed: input.ObservedLedgerDigest, expected: input.ExpectedLedgerDigest},
	}
	for _, check := range checks {
		if check.name == "missing_stage" && check.observed == "" && check.expected == "" && input.ObservedStatus == "sealed" && input.ExpectedStatus == "sealed" {
			continue
		}
		if check.observed == "" || check.expected == "" {
			output.FirstMismatch = check.name
			return output
		}
		if check.observed != check.expected {
			output.Status = "counterexample"
			output.FirstMismatch = check.name
			output.LedgerDigest = input.ObservedLedgerDigest
			return output
		}
	}
	output.Status = "reproduced"
	output.LedgerDigest = input.ObservedLedgerDigest
	return output
}
