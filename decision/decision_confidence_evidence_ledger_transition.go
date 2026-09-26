package decision

// DecisionConfidenceEvidenceLedgerTransition is a reverse observation of two
// adjacent ledger states. It describes continuity only and never claims improvement.
type DecisionConfidenceEvidenceLedgerTransition struct {
	PreviousEntryDigest string `json:"previous_entry_digest"`
	CurrentEntryDigest  string `json:"current_entry_digest"`
	SequenceDelta      uint64 `json:"sequence_delta"`
	Status              string `json:"status"`
	NonAuthorizing     bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerTransition(previous, current DecisionConfidenceEvidenceLedgerEntry) DecisionConfidenceEvidenceLedgerTransition {
	transition := DecisionConfidenceEvidenceLedgerTransition{
		PreviousEntryDigest: previous.EntryDigest,
		CurrentEntryDigest:  current.EntryDigest,
		NonAuthorizing:      true,
	}
	if current.Sequence >= previous.Sequence {
		transition.SequenceDelta = current.Sequence - previous.Sequence
	}
	if VerifyDecisionConfidenceEvidenceLedgerEntry(previous) != nil || VerifyDecisionConfidenceEvidenceLedgerEntry(current) != nil {
		transition.Status = "unknown"
		return transition
	}
	if current.Sequence != previous.Sequence+1 || current.PreviousEntryDigest != previous.EntryDigest {
		transition.Status = "unknown"
		return transition
	}
	transition.Status = "extended"
	return transition
}
