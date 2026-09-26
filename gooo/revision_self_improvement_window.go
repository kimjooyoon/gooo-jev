package gooo

import "fmt"

// RevisionSelfImprovementWindow compares two bound self-improvement receipts
// without claiming that either observation is an improvement.
type RevisionSelfImprovementWindow struct {
	Status                 string
	MissingStage           string
	BaselineReceiptDigest  string
	CandidateReceiptDigest string
	BaselineSourceDigest   string
	CandidateSourceDigest  string
	BaselineProposedSourceDigest  string
	CandidateProposedSourceDigest string
	BaselineApplicationDigest      string
	CandidateApplicationDigest     string
	BaselineMetricsDigest          string
	CandidateMetricsDigest         string
	BaselineGeneratedIRDigest      string
	CandidateGeneratedIRDigest     string
	BaselineStructureDigest        string
	CandidateStructureDigest       string
	BaselineChangedByteCount       int
	CandidateChangedByteCount      int
	BaselineChangedLineCount       int
	CandidateChangedLineCount      int
	ByteDelta              int
	LineDelta              int
	SourceStable           bool
	CandidateStable        bool
	IRChangeStable         bool
	StructureStable        bool
	ExactSourceStable      bool
	GeneratedIRStable      bool
	ComparisonSignal       string
	WindowDigest           string
	NonExecuting           bool
	NonAuthorizing         bool
}

// ObserveRevisionSelfImprovementWindow compares bounded evidence from two
// observations. It never executes, authorizes, persists, or labels improvement.
func ObserveRevisionSelfImprovementWindow(baseline RevisionSelfImprovementReceipt, candidate RevisionSelfImprovementReceipt) (RevisionSelfImprovementWindow, error) {
	window := RevisionSelfImprovementWindow{
		Status:                         "UNKNOWN",
		MissingStage:                   "revision-self-improvement-window",
		BaselineReceiptDigest:          baseline.ReceiptDigest,
		CandidateReceiptDigest:         candidate.ReceiptDigest,
		BaselineSourceDigest:           baseline.SourceDigest,
		CandidateSourceDigest:          candidate.SourceDigest,
		BaselineProposedSourceDigest:   baseline.ProposedSourceDigest,
		CandidateProposedSourceDigest:  candidate.ProposedSourceDigest,
		BaselineApplicationDigest:      baseline.ApplicationDigest,
		CandidateApplicationDigest:     candidate.ApplicationDigest,
		BaselineMetricsDigest:          baseline.MetricsDigest,
		CandidateMetricsDigest:          candidate.MetricsDigest,
		BaselineGeneratedIRDigest:      baseline.GeneratedIRDigest,
		CandidateGeneratedIRDigest:     candidate.GeneratedIRDigest,
		BaselineStructureDigest:        baseline.StructureDigest,
		CandidateStructureDigest:        candidate.StructureDigest,
		BaselineChangedByteCount:       baseline.ChangedByteCount,
		CandidateChangedByteCount:      candidate.ChangedByteCount,
		BaselineChangedLineCount:       baseline.ChangedLineCount,
		CandidateChangedLineCount:      candidate.ChangedLineCount,
		NonExecuting:                   true,
		NonAuthorizing:                 true,
	}
	setWindowDigest := func() {
		window.WindowDigest = digestRevisionSelfImprovementWindow(window)
	}
	setWindowDigest()

	if err := validateBoundSelfImprovementReceipt(baseline); err != nil {
		window.MissingStage = "revision-self-improvement-window-baseline"
		setWindowDigest()
		return window, fmt.Errorf("baseline self-improvement receipt is not valid: %w", err)
	}
	if err := validateBoundSelfImprovementReceipt(candidate); err != nil {
		window.MissingStage = "revision-self-improvement-window-candidate"
		setWindowDigest()
		return window, fmt.Errorf("candidate self-improvement receipt is not valid: %w", err)
	}

	window.ByteDelta = candidate.ChangedByteCount - baseline.ChangedByteCount
	window.LineDelta = candidate.ChangedLineCount - baseline.ChangedLineCount
	window.SourceStable = baseline.SourceDigest == candidate.SourceDigest &&
		baseline.ProposedSourceDigest == candidate.ProposedSourceDigest
	window.CandidateStable = baseline.CandidateDigest == candidate.CandidateDigest &&
		baseline.EditDigest == candidate.EditDigest
	window.IRChangeStable = baseline.IRChanged == candidate.IRChanged
	window.StructureStable = baseline.StructureDigest == candidate.StructureDigest &&
		baseline.StructureMatch == candidate.StructureMatch
	window.ExactSourceStable = baseline.ExactSourceMatch == candidate.ExactSourceMatch
	window.GeneratedIRStable = baseline.GeneratedIRDigest == candidate.GeneratedIRDigest
	window.ComparisonSignal = revisionSelfImprovementWindowSignal(window)
	window.Status = "BOUND"
	window.MissingStage = ""
	setWindowDigest()
	if err := window.Validate(); err != nil {
		window.Status = "UNKNOWN"
		window.MissingStage = "revision-self-improvement-window"
		setWindowDigest()
		return window, fmt.Errorf("revision self-improvement window is not valid: %w", err)
	}
	return window, nil
}

func validateBoundSelfImprovementReceipt(receipt RevisionSelfImprovementReceipt) error {
	if err := receipt.Validate(); err != nil {
		return err
	}
	if receipt.Status != "BOUND" || receipt.MissingStage != "" {
		return fmt.Errorf("receipt is not BOUND")
	}
	return nil
}

func revisionSelfImprovementWindowSignal(window RevisionSelfImprovementWindow) string {
	if window.ByteDelta == 0 && window.LineDelta == 0 && window.IRChangeStable &&
		window.SourceStable && window.CandidateStable && window.StructureStable &&
		window.ExactSourceStable && window.GeneratedIRStable {
		return "stable"
	}
	if window.IRChangeStable && window.ByteDelta < 0 && window.LineDelta <= 0 {
		return "narrower"
	}
	if window.IRChangeStable && window.ByteDelta > 0 && window.LineDelta >= 0 {
		return "wider"
	}
	return "mixed"
}

func (w RevisionSelfImprovementWindow) Validate() error {
	if w.Status == "" {
		return fmt.Errorf("revision self-improvement window status is empty")
	}
	if w.Status == "BOUND" && w.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement window has a missing stage")
	}
	if w.Status == "UNKNOWN" && w.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement window has no missing stage")
	}
	if w.BaselineChangedByteCount < 0 || w.CandidateChangedByteCount < 0 ||
		w.BaselineChangedLineCount < 0 || w.CandidateChangedLineCount < 0 {
		return fmt.Errorf("revision self-improvement window change counts must be non-negative")
	}
	if w.ByteDelta != w.CandidateChangedByteCount-w.BaselineChangedByteCount ||
		w.LineDelta != w.CandidateChangedLineCount-w.BaselineChangedLineCount {
		return fmt.Errorf("revision self-improvement window deltas are not linked")
	}
	if w.ComparisonSignal != "narrower" && w.ComparisonSignal != "wider" &&
		w.ComparisonSignal != "stable" && w.ComparisonSignal != "mixed" {
		return fmt.Errorf("revision self-improvement window signal is invalid")
	}
	if !w.NonExecuting || !w.NonAuthorizing {
		return fmt.Errorf("revision self-improvement window must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementWindow(w) != w.WindowDigest {
		return fmt.Errorf("revision self-improvement window digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementWindow(window RevisionSelfImprovementWindow) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%d|%d|%d|%d|%t|%t|%t|%t|%t|%t|%s|%t|%t",
		window.Status,
		window.MissingStage,
		window.BaselineReceiptDigest,
		window.CandidateReceiptDigest,
		window.BaselineSourceDigest,
		window.CandidateSourceDigest,
		window.BaselineProposedSourceDigest,
		window.CandidateProposedSourceDigest,
		window.BaselineApplicationDigest,
		window.CandidateApplicationDigest,
		window.BaselineMetricsDigest,
		window.CandidateMetricsDigest,
		window.BaselineGeneratedIRDigest,
		window.CandidateGeneratedIRDigest,
		window.BaselineStructureDigest,
		window.CandidateStructureDigest,
		window.BaselineChangedByteCount,
		window.CandidateChangedByteCount,
		window.BaselineChangedLineCount,
		window.CandidateChangedLineCount,
		window.ByteDelta,
		window.LineDelta,
		window.SourceStable,
		window.CandidateStable,
		window.IRChangeStable,
		window.StructureStable,
		window.ExactSourceStable,
		window.GeneratedIRStable,
		window.ComparisonSignal,
		window.NonExecuting,
		window.NonAuthorizing,
	))
}