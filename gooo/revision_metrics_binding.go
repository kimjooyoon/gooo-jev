package gooo

import "fmt"

type RevisionMetricsBinding struct {
	Status               string
	MissingStage         string
	ReceiptDigest        string
	MetricsDigest        string
	ApplicationDigest    string
	SourceDigest         string
	ProposedSourceDigest string
	InputIRDigest        string
	ProposedIRDigest     string
	CandidateDigest      string
	EditDigest           string
	SourceBytesBefore    int
	SourceBytesAfter     int
	SourceBytesDelta     int
	ChangedByteCount     int
	SourceLinesBefore    int
	SourceLinesAfter     int
	SourceLinesDelta     int
	ChangedLineCount     int
	IRChanged            bool
	BindingDigest        string
	NonExecuting         bool
	NonAuthorizing       bool
}

func ObserveRevisionMetricsBinding(receipt RevisionApplicationReceipt, metrics RevisionMetrics) (RevisionMetricsBinding, error) {
	binding := RevisionMetricsBinding{
		Status:               "UNKNOWN",
		MissingStage:         "revision-metrics-binding",
		ReceiptDigest:        receipt.ReceiptDigest,
		MetricsDigest:        metrics.MetricsDigest,
		ApplicationDigest:    metrics.ApplicationDigest,
		SourceDigest:         metrics.SourceDigest,
		ProposedSourceDigest: metrics.ProposedSourceDigest,
		InputIRDigest:        metrics.InputIRDigest,
		ProposedIRDigest:     metrics.ProposedIRDigest,
		CandidateDigest:      metrics.CandidateDigest,
		EditDigest:           metrics.EditDigest,
		SourceBytesBefore:    metrics.SourceBytesBefore,
		SourceBytesAfter:     metrics.SourceBytesAfter,
		SourceBytesDelta:     metrics.SourceBytesDelta,
		ChangedByteCount:     metrics.ChangedByteCount,
		SourceLinesBefore:    metrics.SourceLinesBefore,
		SourceLinesAfter:     metrics.SourceLinesAfter,
		SourceLinesDelta:     metrics.SourceLinesDelta,
		ChangedLineCount:     metrics.ChangedLineCount,
		IRChanged:            metrics.IRChanged,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if err := receipt.Validate(); err != nil {
		binding.MissingStage = "revision-metrics-binding-receipt"
		return binding, fmt.Errorf("gooo revision metrics binding: receipt: %w", err)
	}
	if err := metrics.Validate(); err != nil {
		binding.MissingStage = "revision-metrics-binding-metrics"
		return binding, fmt.Errorf("gooo revision metrics binding: metrics: %w", err)
	}
	if receipt.ApplicationDigest != metrics.ApplicationDigest ||
		receipt.SourceDigest != metrics.SourceDigest ||
		receipt.ProposedSourceDigest != metrics.ProposedSourceDigest ||
		receipt.InputIRDigest != metrics.InputIRDigest ||
		receipt.ProposedIRDigest != metrics.ProposedIRDigest ||
		receipt.CandidateDigest != metrics.CandidateDigest ||
		receipt.EditDigest != metrics.EditDigest {
		binding.MissingStage = "revision-metrics-binding-link"
		return binding, fmt.Errorf("gooo revision metrics binding: metrics are not linked to application receipt")
	}
	binding.Status = "BOUND"
	binding.MissingStage = ""
	binding.BindingDigest = digestRevisionMetricsBinding(binding)
	return binding, nil
}

func (b RevisionMetricsBinding) Validate() error {
	if b.Status != "BOUND" {
		return fmt.Errorf("revision metrics binding status must be BOUND")
	}
	if b.MissingStage != "" {
		return fmt.Errorf("revision metrics binding missing stage must be empty")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("revision metrics binding must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"receipt": b.ReceiptDigest, "metrics": b.MetricsDigest,
		"application": b.ApplicationDigest, "source": b.SourceDigest,
		"proposed source": b.ProposedSourceDigest, "input IR": b.InputIRDigest,
		"proposed IR": b.ProposedIRDigest, "candidate": b.CandidateDigest,
		"edit": b.EditDigest, "binding": b.BindingDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision metrics binding %s digest is invalid", name)
		}
	}
	if b.SourceBytesDelta != b.SourceBytesAfter-b.SourceBytesBefore ||
		b.SourceLinesDelta != b.SourceLinesAfter-b.SourceLinesBefore {
		return fmt.Errorf("revision metrics binding deltas are inconsistent")
	}
	if b.SourceBytesBefore < 0 || b.SourceBytesAfter < 0 || b.ChangedByteCount < 0 ||
		b.SourceLinesBefore < 0 || b.SourceLinesAfter < 0 || b.ChangedLineCount < 0 {
		return fmt.Errorf("revision metrics binding counts must be non-negative")
	}
	if digestRevisionMetricsBinding(b) != b.BindingDigest {
		return fmt.Errorf("revision metrics binding digest does not match its fields")
	}
	return nil
}

func digestRevisionMetricsBinding(binding RevisionMetricsBinding) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%d|%d|%d|%d|%d|%d|%t|%t",
		binding.Status,
		binding.MissingStage,
		binding.ReceiptDigest,
		binding.MetricsDigest,
		binding.ApplicationDigest,
		binding.SourceDigest,
		binding.ProposedSourceDigest,
		binding.InputIRDigest,
		binding.ProposedIRDigest,
		binding.CandidateDigest,
		binding.EditDigest,
		binding.SourceBytesBefore,
		binding.SourceBytesAfter,
		binding.SourceBytesDelta,
		binding.ChangedByteCount,
		binding.SourceLinesBefore,
		binding.SourceLinesAfter,
		binding.SourceLinesDelta,
		binding.ChangedLineCount,
		binding.IRChanged,
		binding.NonExecuting,
		binding.NonAuthorizing,
	))
}
