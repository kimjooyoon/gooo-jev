package decision

import "fmt"

// DecisionConfidenceEvidenceLedgerEntry is one append-only provenance link.
// It records identity and ordering only; it never authorizes execution or change.
type DecisionConfidenceEvidenceLedgerEntry struct {
	Sequence            uint64 `json:"sequence"`
	Kind                string `json:"kind"`
	SourceDigest        string `json:"source_digest"`
	PreviousEntryDigest string `json:"previous_entry_digest,omitempty"`
	EntryDigest         string `json:"entry_digest"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

type decisionConfidenceEvidenceLedgerIdentity struct {
	Sequence            uint64 `json:"sequence"`
	Kind                string `json:"kind"`
	SourceDigest        string `json:"source_digest"`
	PreviousEntryDigest string `json:"previous_entry_digest,omitempty"`
}

func NewDecisionConfidenceEvidenceLedgerEntry(sequence uint64, kind, sourceDigest, previousEntryDigest string) (DecisionConfidenceEvidenceLedgerEntry, error) {
	if sequence == 0 {
		return DecisionConfidenceEvidenceLedgerEntry{}, fmt.Errorf("sequence must be positive")
	}
	if kind == "" {
		return DecisionConfidenceEvidenceLedgerEntry{}, fmt.Errorf("kind must not be empty")
	}
	if sourceDigest == "" {
		return DecisionConfidenceEvidenceLedgerEntry{}, fmt.Errorf("source digest must not be empty")
	}
	if sequence > 1 && previousEntryDigest == "" {
		return DecisionConfidenceEvidenceLedgerEntry{}, fmt.Errorf("previous entry digest is required after the first entry")
	}
	identity := decisionConfidenceEvidenceLedgerIdentity{
		Sequence:            sequence,
		Kind:                kind,
		SourceDigest:        sourceDigest,
		PreviousEntryDigest: previousEntryDigest,
	}
	digest, err := Digest(identity)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerEntry{}, err
	}
	return DecisionConfidenceEvidenceLedgerEntry{
		Sequence:            sequence,
		Kind:                kind,
		SourceDigest:        sourceDigest,
		PreviousEntryDigest: previousEntryDigest,
		EntryDigest:         digest,
		NonAuthorizing:      true,
	}, nil
}

func VerifyDecisionConfidenceEvidenceLedgerEntry(entry DecisionConfidenceEvidenceLedgerEntry) error {
	if !entry.NonAuthorizing {
		return fmt.Errorf("evidence ledger entry must remain non-authorizing")
	}
	verified, err := NewDecisionConfidenceEvidenceLedgerEntry(
		entry.Sequence,
		entry.Kind,
		entry.SourceDigest,
		entry.PreviousEntryDigest,
	)
	if err != nil {
		return err
	}
	if verified.EntryDigest != entry.EntryDigest {
		return fmt.Errorf("entry digest mismatch")
	}
	return nil
}

func AppendDecisionConfidenceEvidenceLedgerEntry(ledger []DecisionConfidenceEvidenceLedgerEntry, entry DecisionConfidenceEvidenceLedgerEntry) ([]DecisionConfidenceEvidenceLedgerEntry, error) {
	if err := VerifyDecisionConfidenceEvidenceLedgerEntry(entry); err != nil {
		return nil, err
	}
	if len(ledger) == 0 {
		if entry.Sequence != 1 || entry.PreviousEntryDigest != "" {
			return nil, fmt.Errorf("first entry must start at sequence one without a predecessor")
		}
	} else {
		previous := ledger[len(ledger)-1]
		if err := VerifyDecisionConfidenceEvidenceLedgerEntry(previous); err != nil {
			return nil, err
		}
		if entry.Sequence != previous.Sequence+1 || entry.PreviousEntryDigest != previous.EntryDigest {
			return nil, fmt.Errorf("entry does not continue the ledger")
		}
	}
	return append(append([]DecisionConfidenceEvidenceLedgerEntry(nil), ledger...), entry), nil
}
