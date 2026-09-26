package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DecisionConfidenceImprovementReviewSchemaV1 = "gooo/jev-decision-confidence-improvement-review/v1"

type DecisionConfidenceImprovementReviewDecision string

const (
	DecisionConfidenceImprovementApproved  DecisionConfidenceImprovementReviewDecision = "approved"
	DecisionConfidenceImprovementRejected   DecisionConfidenceImprovementReviewDecision = "rejected"
	DecisionConfidenceImprovementNoChange  DecisionConfidenceImprovementReviewDecision = "no-change"
	DecisionConfidenceImprovementReviewUnknown DecisionConfidenceImprovementReviewDecision = "unknown"
)

type DecisionConfidenceImprovementReviewReceipt struct {
	Schema               string
	SignalDigest         string
	ReviewerDigest       string
	ProposedChangeDigest string
	Decision             DecisionConfidenceImprovementReviewDecision
	MissingStage         string
	NonExecution         bool
	ReviewedAt           time.Time
	ReceiptDigest        string
}

func ReviewDecisionConfidenceImprovementSignal(
	signal DecisionConfidenceImprovementSignal,
	reviewerDigest string,
	proposedChangeDigest string,
	approved bool,
	reviewedAt time.Time,
) (DecisionConfidenceImprovementReviewReceipt, error) {
	if reviewedAt.IsZero() {
		return DecisionConfidenceImprovementReviewReceipt{}, errors.New("decision confidence improvement review time is required")
	}
	receipt := DecisionConfidenceImprovementReviewReceipt{
		Schema:               DecisionConfidenceImprovementReviewSchemaV1,
		SignalDigest:         signal.SignalDigest,
		ReviewerDigest:       strings.TrimSpace(reviewerDigest),
		ProposedChangeDigest: strings.TrimSpace(proposedChangeDigest),
		Decision:             DecisionConfidenceImprovementApproved,
		NonExecution:         true,
		ReviewedAt:           reviewedAt.UTC(),
	}
	switch {
	case signal.Validate() != nil:
		receipt.Decision = DecisionConfidenceImprovementReviewUnknown
		receipt.MissingStage = "signal"
	case signal.Status == DecisionConfidenceNoChange:
		receipt.Decision = DecisionConfidenceImprovementNoChange
	case signal.Status == DecisionConfidenceSignalUnknown:
		receipt.Decision = DecisionConfidenceImprovementReviewUnknown
		receipt.MissingStage = "signal-status"
		if strings.TrimSpace(signal.MissingStage) != "" {
			receipt.MissingStage = "signal:" + signal.MissingStage
		}
	case strings.TrimSpace(reviewerDigest) == "":
		receipt.Decision = DecisionConfidenceImprovementReviewUnknown
		receipt.MissingStage = "reviewer"
	case strings.TrimSpace(proposedChangeDigest) == "":
		receipt.Decision = DecisionConfidenceImprovementReviewUnknown
		receipt.MissingStage = "proposed-change"
	case approved:
		receipt.Decision = DecisionConfidenceImprovementApproved
	default:
		receipt.Decision = DecisionConfidenceImprovementRejected
	}
	if err := receipt.assignDigest(); err != nil {
		return DecisionConfidenceImprovementReviewReceipt{}, fmt.Errorf("digest decision confidence improvement review: %w", err)
	}
	if err := receipt.Validate(); err != nil {
		return DecisionConfidenceImprovementReviewReceipt{}, err
	}
	return receipt, nil
}

func (receipt *DecisionConfidenceImprovementReviewReceipt) assignDigest() error {
	digest, err := receipt.computeDigest()
	if err != nil {
		return err
	}
	receipt.ReceiptDigest = digest
	return nil
}

func (receipt DecisionConfidenceImprovementReviewReceipt) Validate() error {
	if err := receipt.validateShape(); err != nil {
		return err
	}
	expected, err := receipt.computeDigest()
	if err != nil {
		return fmt.Errorf("digest decision confidence improvement review: %w", err)
	}
	if receipt.ReceiptDigest != expected {
		return errors.New("decision confidence improvement review digest mismatch")
	}
	return nil
}

func (receipt DecisionConfidenceImprovementReviewReceipt) validateShape() error {
	if receipt.Schema != DecisionConfidenceImprovementReviewSchemaV1 {
		return errors.New("unsupported decision confidence improvement review schema")
	}
	if receipt.ReviewedAt.IsZero() {
		return errors.New("decision confidence improvement review time is required")
	}
	if strings.TrimSpace(receipt.ReceiptDigest) == "" {
		return errors.New("decision confidence improvement review digest is required")
	}
	if !receipt.NonExecution {
		return errors.New("decision confidence improvement review must remain non-executing")
	}
	switch receipt.Decision {
	case DecisionConfidenceImprovementApproved, DecisionConfidenceImprovementRejected:
		if strings.TrimSpace(receipt.SignalDigest) == "" ||
			strings.TrimSpace(receipt.ReviewerDigest) == "" ||
			strings.TrimSpace(receipt.ProposedChangeDigest) == "" ||
			strings.TrimSpace(receipt.MissingStage) != "" {
			return errors.New("reviewed decision confidence improvement is incomplete")
		}
	case DecisionConfidenceImprovementNoChange:
		if strings.TrimSpace(receipt.SignalDigest) == "" ||
			strings.TrimSpace(receipt.MissingStage) != "" {
			return errors.New("no-change decision confidence review is incomplete")
		}
	case DecisionConfidenceImprovementReviewUnknown:
		if strings.TrimSpace(receipt.MissingStage) == "" {
			return errors.New("unknown decision confidence improvement review requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported decision confidence improvement review decision %q", receipt.Decision)
	}
	return nil
}

func (receipt DecisionConfidenceImprovementReviewReceipt) computeDigest() (string, error) {
	return Digest(struct {
		Schema               string
		SignalDigest         string
		ReviewerDigest       string
		ProposedChangeDigest string
		Decision             DecisionConfidenceImprovementReviewDecision
		MissingStage         string
		NonExecution         bool
		ReviewedAt           time.Time
	}{
		Schema:               receipt.Schema,
		SignalDigest:         receipt.SignalDigest,
		ReviewerDigest:       receipt.ReviewerDigest,
		ProposedChangeDigest: receipt.ProposedChangeDigest,
		Decision:             receipt.Decision,
		MissingStage:         receipt.MissingStage,
		NonExecution:         receipt.NonExecution,
		ReviewedAt:           receipt.ReviewedAt.UTC(),
	})
}
