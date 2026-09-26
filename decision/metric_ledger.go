package decision

import (
    "errors"
    "strings"
)

type MetricLedgerEntry struct {
    Index          int
    PreviousDigest string
    Delta          MetricDelta
    EntryDigest    string
}

type MetricLedger struct {
    Entries     []MetricLedgerEntry
    LedgerDigest string
}

func NewMetricLedger() (MetricLedger, error) {
    return rebuildMetricLedger(nil)
}

func (ledger MetricLedger) Append(delta MetricDelta) (MetricLedger, error) {
    if err := delta.Validate(); err != nil {
        return MetricLedger{}, err
    }
    if len(ledger.Entries) == 0 && strings.TrimSpace(ledger.LedgerDigest) == "" {
        var err error
        ledger, err = rebuildMetricLedger(nil)
        if err != nil {
            return MetricLedger{}, err
        }
    } else if err := ledger.Validate(); err != nil {
        return MetricLedger{}, err
    }

    entries := append([]MetricLedgerEntry(nil), ledger.Entries...)
    entry := MetricLedgerEntry{
        Index:          len(entries),
        PreviousDigest: ledger.LedgerDigest,
        Delta:          delta,
    }
    entryDigest, err := digestMetricLedgerEntry(entry)
    if err != nil {
        return MetricLedger{}, err
    }
    entry.EntryDigest = entryDigest
    entries = append(entries, entry)
    return rebuildMetricLedger(entries)
}

func (ledger MetricLedger) Validate() error {
    if strings.TrimSpace(ledger.LedgerDigest) == "" {
        return errors.New("metric ledger digest is missing")
    }
    previous := ""
    for index, entry := range ledger.Entries {
        if entry.Index != index {
            return errors.New("metric ledger entry index is not contiguous")
        }
        if entry.PreviousDigest != previous {
            return errors.New("metric ledger entry chain is broken")
        }
        if err := entry.Delta.Validate(); err != nil {
            return err
        }
        digest, err := digestMetricLedgerEntry(entry)
        if err != nil {
            return err
        }
        if digest != entry.EntryDigest {
            return errors.New("metric ledger entry digest does not match its evidence")
        }
        previous = entry.EntryDigest
    }

    expected, err := rebuildMetricLedgerDigest(ledger.Entries)
    if err != nil {
        return err
    }
    if expected != ledger.LedgerDigest {
        return errors.New("metric ledger digest does not match its entries")
    }
    return nil
}

func digestMetricLedgerEntry(entry MetricLedgerEntry) (string, error) {
    copy := entry
    copy.EntryDigest = ""
    return Digest(copy)
}

func rebuildMetricLedger(entries []MetricLedgerEntry) (MetricLedger, error) {
    digest, err := rebuildMetricLedgerDigest(entries)
    if err != nil {
        return MetricLedger{}, err
    }
    return MetricLedger{
        Entries:     append([]MetricLedgerEntry(nil), entries...),
        LedgerDigest: digest,
    }, nil
}

func rebuildMetricLedgerDigest(entries []MetricLedgerEntry) (string, error) {
    return Digest(struct {
        Entries []MetricLedgerEntry
    }{Entries: entries})
}
