package decision

import (
	"errors"
	"math"
	"strings"
	"time"
)

type FeedbackKind string

const (
	FeedbackConfirmed FeedbackKind = "confirmed"
	FeedbackRefuted   FeedbackKind = "refuted"
	FeedbackUnknown   FeedbackKind = "unknown"
)

type FeedbackObservation struct {
	ReceiptDigest  string
	LedgerDigest   string
	Kind           FeedbackKind
	MetricName     string
	MetricValue    float64
	EvidenceDigest string
	RecordedAt     time.Time
	NonAuthorizing bool
	FeedbackDigest string
}

func ObserveFeedback(
	ledger Ledger,
	receipt Receipt,
	kind FeedbackKind,
	metricName string,
	metricValue float64,
	evidenceDigest string,
	recordedAt time.Time,
) (FeedbackObservation, error) {
	if err := ledger.Validate(); err != nil {
		return FeedbackObservation{}, err
	}
	if err := receipt.Validate(); err != nil {
		return FeedbackObservation{}, err
	}
	if strings.TrimSpace(metricName) == "" {
		return FeedbackObservation{}, errors.New("feedback metric name is required")
	}
	if math.IsNaN(metricValue) || math.IsInf(metricValue, 0) {
		return FeedbackObservation{}, errors.New("feedback metric value must be finite")
	}
	if strings.TrimSpace(evidenceDigest) == "" {
		return FeedbackObservation{}, errors.New("feedback evidence digest is required")
	}
	if recordedAt.IsZero() {
		return FeedbackObservation{}, errors.New("feedback recorded time is required")
	}
	switch kind {
	case FeedbackConfirmed, FeedbackRefuted, FeedbackUnknown:
	default:
		return FeedbackObservation{}, errors.New("unsupported feedback kind")
	}
	receiptDigest, err := Digest(receipt)
	if err != nil {
		return FeedbackObservation{}, err
	}
	found := false
	for _, entry := range ledger.Entries {
		if entry.Receipt == receipt {
			found = true
			break
		}
	}
	if !found {
		return FeedbackObservation{}, errors.New("feedback receipt is not present in ledger")
	}
	feedback := FeedbackObservation{
		ReceiptDigest:  receiptDigest,
		LedgerDigest:   ledger.Digest(),
		Kind:           kind,
		MetricName:     metricName,
		MetricValue:    metricValue,
		EvidenceDigest: evidenceDigest,
		RecordedAt:     recordedAt.UTC(),
		NonAuthorizing: true,
	}
	feedback.FeedbackDigest, err = Digest(feedback)
	if err != nil {
		return FeedbackObservation{}, err
	}
	return feedback, nil
}

func (feedback FeedbackObservation) Validate() error {
	if strings.TrimSpace(feedback.ReceiptDigest) == "" ||
		strings.TrimSpace(feedback.LedgerDigest) == "" ||
		strings.TrimSpace(feedback.MetricName) == "" ||
		strings.TrimSpace(feedback.EvidenceDigest) == "" ||
		strings.TrimSpace(feedback.FeedbackDigest) == "" {
		return errors.New("feedback observation is incomplete")
	}
	if math.IsNaN(feedback.MetricValue) || math.IsInf(feedback.MetricValue, 0) {
		return errors.New("feedback metric value must be finite")
	}
	if feedback.RecordedAt.IsZero() {
		return errors.New("feedback recorded time is required")
	}
	if !feedback.NonAuthorizing {
		return errors.New("feedback observation must remain non-authorizing")
	}
	switch feedback.Kind {
	case FeedbackConfirmed, FeedbackRefuted, FeedbackUnknown:
	default:
		return errors.New("unsupported feedback kind")
	}
	copy := feedback
	copy.FeedbackDigest = ""
	digest, err := Digest(copy)
	if err != nil {
		return err
	}
	if digest != feedback.FeedbackDigest {
		return errors.New("feedback digest does not match its evidence")
	}
	return nil
}
