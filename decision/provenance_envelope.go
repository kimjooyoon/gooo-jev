package decision

import (
    "errors"
    "strings"
)

type ProvenanceObservationStatus string

const (
    ProvenanceObserved ProvenanceObservationStatus = "observed"
    ProvenanceUnknown  ProvenanceObservationStatus = "unknown"
    ProvenanceReview   ProvenanceObservationStatus = "review"
)

type ProvenanceEnvelope struct {
    ReceiptDigest           string
    ReviewDigest            string
    CapabilityBoundaryDigest string
    SourceReference         string
    ObservationStatus       ProvenanceObservationStatus
    EnvelopeDigest          string
}

func NewProvenanceEnvelope(
    receiptDigest,
    reviewDigest,
    capabilityBoundaryDigest,
    sourceReference string,
    status ProvenanceObservationStatus,
) (ProvenanceEnvelope, error) {
    envelope := ProvenanceEnvelope{
        ReceiptDigest:            receiptDigest,
        ReviewDigest:             reviewDigest,
        CapabilityBoundaryDigest: capabilityBoundaryDigest,
        SourceReference:          sourceReference,
        ObservationStatus:        status,
    }
    if err := envelope.validateShape(); err != nil {
        return ProvenanceEnvelope{}, err
    }
    digest, err := Digest(envelope)
    if err != nil {
        return ProvenanceEnvelope{}, err
    }
    envelope.EnvelopeDigest = digest
    return envelope, nil
}

func (envelope ProvenanceEnvelope) Validate() error {
    if err := envelope.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(envelope.EnvelopeDigest) == "" {
        return errors.New("provenance envelope digest is missing")
    }
    copy := envelope
    copy.EnvelopeDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != envelope.EnvelopeDigest {
        return errors.New("provenance envelope digest does not match its evidence")
    }
    return nil
}

func (envelope ProvenanceEnvelope) validateShape() error {
    if strings.TrimSpace(envelope.ReceiptDigest) == "" ||
        strings.TrimSpace(envelope.ReviewDigest) == "" ||
        strings.TrimSpace(envelope.CapabilityBoundaryDigest) == "" ||
        strings.TrimSpace(envelope.SourceReference) == "" {
        return errors.New("provenance envelope is incomplete")
    }
    switch envelope.ObservationStatus {
    case ProvenanceObserved, ProvenanceUnknown, ProvenanceReview:
        return nil
    default:
        return errors.New("unsupported provenance observation status")
    }
}
