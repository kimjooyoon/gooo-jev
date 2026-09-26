package decision

import "testing"

func TestMetricLedgerPreservesImprovementRegressionAndUnknown(t *testing.T) {
    improved, err := NewMetricDelta("quality", 0.5, 0.75, "evidence-improved", MetricObserved)
    if err != nil {
        t.Fatalf("improved delta error = %v", err)
    }
    regressed, err := NewMetricDelta("quality", 0.75, 0.6, "evidence-regressed", MetricObserved)
    if err != nil {
        t.Fatalf("regressed delta error = %v", err)
    }
    unknown, err := NewMetricDelta("quality", 0, 0, "evidence-unknown", MetricUnknown)
    if err != nil {
        t.Fatalf("unknown delta error = %v", err)
    }

    ledger, err := NewMetricLedger()
    if err != nil {
        t.Fatalf("NewMetricLedger() error = %v", err)
    }
    ledger, err = ledger.Append(improved)
    if err != nil {
        t.Fatalf("first Append() error = %v", err)
    }
    ledger, err = ledger.Append(regressed)
    if err != nil {
        t.Fatalf("second Append() error = %v", err)
    }
    ledger, err = ledger.Append(unknown)
    if err != nil {
        t.Fatalf("third Append() error = %v", err)
    }
    if err := ledger.Validate(); err != nil {
        t.Fatalf("Validate() error = %v", err)
    }
    if len(ledger.Entries) != 3 {
        t.Fatalf("entry count = %d, want 3", len(ledger.Entries))
    }
    if ledger.Entries[0].Delta.Direction != MetricImproved ||
        ledger.Entries[1].Delta.Direction != MetricRegressed ||
        ledger.Entries[2].Delta.Direction != MetricDirectionUnknown {
        t.Fatal("metric direction history was not preserved")
    }
    if ledger.Entries[0].PreviousDigest != "" {
        t.Fatal("first entry unexpectedly has a predecessor")
    }
    if ledger.Entries[1].PreviousDigest != ledger.Entries[0].EntryDigest ||
        ledger.Entries[2].PreviousDigest != ledger.Entries[1].EntryDigest {
        t.Fatal("metric entry predecessor chain is broken")
    }

    tampered := ledger
    tampered.Entries = append([]MetricLedgerEntry(nil), ledger.Entries...)
    tampered.Entries[1].Delta.After = 0.9
    if err := tampered.Validate(); err == nil {
        t.Fatal("tampered ledger unexpectedly validated")
    }
}
