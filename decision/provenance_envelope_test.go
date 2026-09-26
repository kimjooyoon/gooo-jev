package decision

import "testing"

func TestProvenanceEnvelopeBindsEvidenceAndRejectsTampering(t *testing.T) {
    envelope, err := NewProvenanceEnvelope(
        "receipt-digest-1",
        "review-digest-1",
        "capability-digest-1",
        "gooo://jev/decision/receipt/1",
        ProvenanceReview,
    )
    if err != nil {
        t.Fatalf("NewProvenanceEnvelope() error = %v", err)
    }
    if err := envelope.Validate(); err != nil {
        t.Fatalf("Validate() error = %v", err)
    }

    tampered := envelope
    tampered.SourceReference = "gooo://jev/decision/receipt/changed"
    if err := tampered.Validate(); err == nil {
        t.Fatal("tampered source reference unexpectedly validated")
    }

    unknown, err := NewProvenanceEnvelope(
        "receipt-digest-2",
        "review-digest-2",
        "capability-digest-2",
        "gooo://jev/decision/receipt/2",
        ProvenanceUnknown,
    )
    if err != nil {
        t.Fatalf("unknown envelope construction error = %v", err)
    }
    if err := unknown.Validate(); err != nil {
        t.Fatalf("unknown envelope Validate() error = %v", err)
    }

    if _, err := NewProvenanceEnvelope(
        "receipt-digest-3",
        "review-digest-3",
        "capability-digest-3",
        "gooo://jev/decision/receipt/3",
        ProvenanceObservationStatus("accepted"),
    ); err == nil {
        t.Fatal("unsupported status unexpectedly accepted")
    }
}
