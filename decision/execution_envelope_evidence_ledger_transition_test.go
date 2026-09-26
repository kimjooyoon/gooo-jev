package decision

import "testing"

func TestCompareExecutionEnvelopeEvidenceLedgerTransition(t *testing.T) {
	tests := []struct {
		name          string
		input         ExecutionEnvelopeEvidenceLedgerTransitionInput
		status        string
		changed       bool
		direction     string
		firstMismatch string
	}{
		{
			name: "unchanged",
			input: ExecutionEnvelopeEvidenceLedgerTransitionInput{
				PreviousStatus: "sealed", CurrentStatus: "sealed",
				PreviousLedgerDigest: "digest-a", CurrentLedgerDigest: "digest-a",
				NonAuthorizing: true,
			},
			status: "unchanged", changed: false, direction: "stable",
		},
		{
			name: "content changed",
			input: ExecutionEnvelopeEvidenceLedgerTransitionInput{
				PreviousStatus: "sealed", CurrentStatus: "sealed",
				PreviousLedgerDigest: "digest-a", CurrentLedgerDigest: "digest-b",
				NonAuthorizing: true,
			},
			status: "changed", changed: true, direction: "content-transition",
			firstMismatch: "ledger-content",
		},
		{
			name: "status changed",
			input: ExecutionEnvelopeEvidenceLedgerTransitionInput{
				PreviousStatus: "UNKNOWN", CurrentStatus: "sealed",
				PreviousLedgerDigest: "digest-a", CurrentLedgerDigest: "digest-b",
				NonAuthorizing: true,
			},
			status: "changed", changed: true, direction: "status-transition",
			firstMismatch: "ledger-status",
		},
		{
			name: "missing evidence",
			input: ExecutionEnvelopeEvidenceLedgerTransitionInput{
				PreviousStatus: "sealed", CurrentStatus: "sealed",
				PreviousLedgerDigest: "", CurrentLedgerDigest: "digest-b",
				NonAuthorizing: true,
			},
			status: "UNKNOWN", changed: false, firstMismatch: "ledger-transition-evidence",
		},
		{
			name: "authorization claim",
			input: ExecutionEnvelopeEvidenceLedgerTransitionInput{
				PreviousStatus: "sealed", CurrentStatus: "sealed",
				PreviousLedgerDigest: "digest-a", CurrentLedgerDigest: "digest-b",
				NonAuthorizing: false,
			},
			status: "UNKNOWN", changed: false, firstMismatch: "authorization-boundary",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := CompareExecutionEnvelopeEvidenceLedgerTransition(test.input)
			if output.Status != test.status || output.Changed != test.changed ||
				output.Direction != test.direction || output.FirstMismatch != test.firstMismatch {
				t.Fatalf("unexpected transition: %#v", output)
			}
			if output.NonAuthorizing != test.input.NonAuthorizing {
				t.Fatalf("unexpected authorization state: %#v", output)
			}
			if test.status != "UNKNOWN" && output.LedgerDigest != test.input.CurrentLedgerDigest {
				t.Fatalf("unexpected ledger digest: %#v", output)
			}
		})
	}
}