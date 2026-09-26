package decision

import "fmt"

type LedgerEntry struct {
	Index          int
	PreviousDigest string
	Receipt        Receipt
	EntryDigest    string
}

type Ledger struct {
	Entries []LedgerEntry
}

func (ledger Ledger) Append(receipt Receipt) (Ledger, error) {
	if err := receipt.Validate(); err != nil {
		return Ledger{}, err
	}
	if err := ledger.Validate(); err != nil {
		return Ledger{}, err
	}
	entries := append([]LedgerEntry(nil), ledger.Entries...)
	previousDigest := ""
	if len(entries) > 0 {
		previousDigest = entries[len(entries)-1].EntryDigest
	}
	entry := LedgerEntry{
		Index:          len(entries),
		PreviousDigest: previousDigest,
		Receipt:        receipt,
	}
	entry.EntryDigest = digestLedgerEntry(entry)
	entries = append(entries, entry)
	return Ledger{Entries: entries}, nil
}

func (ledger Ledger) Validate() error {
	previousDigest := ""
	for index, entry := range ledger.Entries {
		if entry.Index != index {
			return fmt.Errorf("ledger entry index %d does not match position", entry.Index)
		}
		if entry.PreviousDigest != previousDigest {
			return fmt.Errorf("ledger entry %d previous digest does not match", index)
		}
		if err := entry.Receipt.Validate(); err != nil {
			return fmt.Errorf("ledger entry %d receipt: %w", index, err)
		}
		if entry.EntryDigest == "" || entry.EntryDigest != digestLedgerEntry(entry) {
			return fmt.Errorf("ledger entry %d digest does not match", index)
		}
		previousDigest = entry.EntryDigest
	}
	return nil
}

func (ledger Ledger) Digest() string {
	if len(ledger.Entries) == 0 {
		return ""
	}
	return ledger.Entries[len(ledger.Entries)-1].EntryDigest
}

func digestLedgerEntry(entry LedgerEntry) string {
	entry.EntryDigest = ""
	digest, _ := Digest(entry)
	return digest
}
