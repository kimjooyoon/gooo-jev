package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	jevReplayFeedbackLedgerReady   = "ready"
	jevReplayFeedbackLedgerUnknown = "UNKNOWN"
)

// ExecutionEnvelopeJEVReplayFeedbackLedgerInput appends one validated replay
// feedback to an immutable, non-executing evidence chain.
type ExecutionEnvelopeJEVReplayFeedbackLedgerInput struct {
	Previous       JEVImprovementReplayFeedbackLedger
	Feedback       JEVImprovementReplayFeedback
	MetricDigest   string
	NonAuthorizing bool
}

// JEVImprovementReplayFeedbackLedger preserves replay feedback history without
// applying a candidate or granting capability.
type JEVImprovementReplayFeedbackLedger struct {
	Status                 string
	MissingStage           string
	FeedbackStatus         string
	PreviousEvidenceDigest string
	FeedbackEvidenceDigest string
	MetricDigest           string
	EntryDigest            string
	ConfirmedCount         int
	RefutedCount           int
	UnknownCount           int
	EvidenceDigest         string
	NonExecuting           bool
	NonAuthorizing         bool
}

func (l JEVImprovementReplayFeedbackLedger) Validate() error {
	if l.Status != jevReplayFeedbackLedgerReady ||
		l.FeedbackStatus == "" ||
		l.PreviousEvidenceDigest == "" ||
		l.FeedbackEvidenceDigest == "" ||
		l.MetricDigest == "" ||
		l.EntryDigest == "" ||
		l.EvidenceDigest == "" {
		return fmt.Errorf("incomplete JEV improvement replay feedback ledger")
	}
	switch l.FeedbackStatus {
	case jevReplayFeedbackConfirmed, jevReplayFeedbackRefuted, jevReplayFeedbackUnknown:
	default:
		return fmt.Errorf("invalid JEV improvement replay feedback ledger status")
	}
	if l.ConfirmedCount < 0 || l.RefutedCount < 0 || l.UnknownCount < 0 {
		return fmt.Errorf("negative JEV improvement replay feedback ledger count")
	}
	if !l.NonExecuting {
		return fmt.Errorf("JEV improvement replay feedback ledger must be non-executing")
	}
	if !l.NonAuthorizing {
		return fmt.Errorf("JEV improvement replay feedback ledger must be non-authorizing")
	}
	expectedEntry := digestJEVImprovementReplayFeedbackLedgerEntry(
		l.PreviousEvidenceDigest,
		l.FeedbackEvidenceDigest,
		l.MetricDigest,
		l.FeedbackStatus,
		l.ConfirmedCount,
		l.RefutedCount,
		l.UnknownCount,
	)
	if l.EntryDigest != expectedEntry {
		return fmt.Errorf("JEV improvement replay feedback ledger entry digest mismatch")
	}
	expected := digestJEVImprovementReplayFeedbackLedger(
		l.Status,
		l.FeedbackStatus,
		l.PreviousEvidenceDigest,
		l.FeedbackEvidenceDigest,
		l.MetricDigest,
		l.EntryDigest,
		l.ConfirmedCount,
		l.RefutedCount,
		l.UnknownCount,
	)
	if l.EvidenceDigest != expected {
		return fmt.Errorf("JEV improvement replay feedback ledger evidence digest mismatch")
	}
	return nil
}

// AppendJEVImprovementReplayFeedbackLedger records one feedback result and
// keeps UNKNOWN causes outside the success path.
func AppendJEVImprovementReplayFeedbackLedger(input ExecutionEnvelopeJEVReplayFeedbackLedgerInput) JEVImprovementReplayFeedbackLedger {
	output := JEVImprovementReplayFeedbackLedger{
		Status:         jevReplayFeedbackLedgerUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing ||
		(input.Previous.Status != "" && !input.Previous.NonAuthorizing) {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Feedback.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Feedback.Validate(); err != nil {
		output.MissingStage = "feedback-validation"
		return output
	}
	if strings.TrimSpace(input.MetricDigest) == "" {
		output.MissingStage = "metric"
		return output
	}
	if input.MetricDigest != input.Feedback.MetricDigest {
		output.MissingStage = "metric-binding"
		return output
	}

	previousDigest := "genesis"
	confirmedCount := 0
	refutedCount := 0
	unknownCount := 0
	if input.Previous.Status != "" {
		if err := input.Previous.Validate(); err != nil {
			output.MissingStage = "ledger-validation"
			return output
		}
		previousDigest = input.Previous.EvidenceDigest
		confirmedCount = input.Previous.ConfirmedCount
		refutedCount = input.Previous.RefutedCount
		unknownCount = input.Previous.UnknownCount
	}
	switch input.Feedback.Status {
	case jevReplayFeedbackConfirmed:
		confirmedCount++
	case jevReplayFeedbackRefuted:
		refutedCount++
	case jevReplayFeedbackUnknown:
		unknownCount++
	}
	output.Status = jevReplayFeedbackLedgerReady
	output.FeedbackStatus = input.Feedback.Status
	output.PreviousEvidenceDigest = previousDigest
	output.FeedbackEvidenceDigest = input.Feedback.EvidenceDigest
	output.MetricDigest = input.MetricDigest
	output.ConfirmedCount = confirmedCount
	output.RefutedCount = refutedCount
	output.UnknownCount = unknownCount
	output.EntryDigest = digestJEVImprovementReplayFeedbackLedgerEntry(
		output.PreviousEvidenceDigest,
		output.FeedbackEvidenceDigest,
		output.MetricDigest,
		output.FeedbackStatus,
		output.ConfirmedCount,
		output.RefutedCount,
		output.UnknownCount,
	)
	output.EvidenceDigest = digestJEVImprovementReplayFeedbackLedger(
		output.Status,
		output.FeedbackStatus,
		output.PreviousEvidenceDigest,
		output.FeedbackEvidenceDigest,
		output.MetricDigest,
		output.EntryDigest,
		output.ConfirmedCount,
		output.RefutedCount,
		output.UnknownCount,
	)
	if err := output.Validate(); err != nil {
		output.Status = jevReplayFeedbackLedgerUnknown
		output.MissingStage = "ledger-evidence"
		output.EvidenceDigest = ""
	}
	return output
}

func digestJEVImprovementReplayFeedbackLedgerEntry(
	previousEvidenceDigest,
	feedbackEvidenceDigest,
	metricDigest,
	feedbackStatus string,
	confirmedCount,
	refutedCount,
	unknownCount int,
) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf(
		"%s|%s|%s|%s|%d|%d|%d",
		previousEvidenceDigest,
		feedbackEvidenceDigest,
		metricDigest,
		feedbackStatus,
		confirmedCount,
		refutedCount,
		unknownCount,
	)))
	return hex.EncodeToString(sum[:])
}

func digestJEVImprovementReplayFeedbackLedger(
	status,
	feedbackStatus,
	previousEvidenceDigest,
	feedbackEvidenceDigest,
	metricDigest,
	entryDigest string,
	confirmedCount,
	refutedCount,
	unknownCount int,
) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%d|%d|%d",
		status,
		feedbackStatus,
		previousEvidenceDigest,
		feedbackEvidenceDigest,
		metricDigest,
		entryDigest,
		confirmedCount,
		refutedCount,
		unknownCount,
	)))
	return hex.EncodeToString(sum[:])
}