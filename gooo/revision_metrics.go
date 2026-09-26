package gooo

import (
	"fmt"
	"strings"
)

type RevisionMetrics struct {
	Status               string
	MissingStage         string
	SourceDigest         string
	ProposedSourceDigest string
	InputIRDigest        string
	ProposedIRDigest     string
	CandidateDigest      string
	EditDigest           string
	ApplicationDigest    string
	MetricsDigest        string
	SourceBytesBefore    int
	SourceBytesAfter     int
	SourceBytesDelta     int
	ChangedByteCount     int
	SourceLinesBefore    int
	SourceLinesAfter     int
	SourceLinesDelta     int
	ChangedLineCount     int
	IRChanged            bool
	NonExecuting         bool
	NonAuthorizing       bool
}

func MeasureRevision(source string, application RevisionApplication) (RevisionMetrics, error) {
	metrics := RevisionMetrics{
		Status:               "UNKNOWN",
		MissingStage:         "revision-metrics",
		SourceDigest:         application.SourceDigest,
		ProposedSourceDigest: application.ProposedSourceDigest,
		InputIRDigest:        application.InputIRDigest,
		ProposedIRDigest:     application.ProposedIRDigest,
		CandidateDigest:      application.CandidateDigest,
		EditDigest:           application.EditDigest,
		ApplicationDigest:    application.ApplicationDigest,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if err := application.Validate(); err != nil {
		metrics.MissingStage = "revision-metrics-application"
		return metrics, fmt.Errorf("gooo revision metrics: application: %w", err)
	}
	if !validDigest(application.SourceDigest) || digestString(source) != application.SourceDigest {
		metrics.MissingStage = "revision-metrics-source"
		return metrics, fmt.Errorf("gooo revision metrics: source digest does not match application")
	}

	beforeBytes := []byte(source)
	afterBytes := []byte(application.ProposedSource)
	beforeLines := strings.Split(source, "\n")
	afterLines := strings.Split(application.ProposedSource, "\n")
	metrics.SourceBytesBefore = len(beforeBytes)
	metrics.SourceBytesAfter = len(afterBytes)
	metrics.SourceBytesDelta = len(afterBytes) - len(beforeBytes)
	metrics.ChangedByteCount = changedBytes(beforeBytes, afterBytes)
	metrics.SourceLinesBefore = len(beforeLines)
	metrics.SourceLinesAfter = len(afterLines)
	metrics.SourceLinesDelta = len(afterLines) - len(beforeLines)
	metrics.ChangedLineCount = changedLines(beforeLines, afterLines)
	metrics.IRChanged = application.InputIRDigest != application.ProposedIRDigest
	metrics.Status = "BOUND"
	metrics.MissingStage = ""
	metrics.MetricsDigest = digestRevisionMetrics(metrics)
	return metrics, nil
}

func (m RevisionMetrics) Validate() error {
	if m.Status != "BOUND" {
		return fmt.Errorf("revision metrics status must be BOUND")
	}
	if m.MissingStage != "" {
		return fmt.Errorf("revision metrics missing stage must be empty")
	}
	if !m.NonExecuting || !m.NonAuthorizing {
		return fmt.Errorf("revision metrics must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": m.SourceDigest, "proposed source": m.ProposedSourceDigest,
		"input IR": m.InputIRDigest, "proposed IR": m.ProposedIRDigest,
		"candidate": m.CandidateDigest, "edit": m.EditDigest,
		"application": m.ApplicationDigest, "metrics": m.MetricsDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision metrics %s digest is invalid", name)
		}
	}
	if m.SourceBytesBefore < 0 || m.SourceBytesAfter < 0 || m.ChangedByteCount < 0 {
		return fmt.Errorf("revision metrics byte counts must be non-negative")
	}
	if m.SourceLinesBefore < 0 || m.SourceLinesAfter < 0 || m.ChangedLineCount < 0 {
		return fmt.Errorf("revision metrics line counts must be non-negative")
	}
	if m.SourceBytesDelta != m.SourceBytesAfter-m.SourceBytesBefore {
		return fmt.Errorf("revision metrics byte delta is inconsistent")
	}
	if m.SourceLinesDelta != m.SourceLinesAfter-m.SourceLinesBefore {
		return fmt.Errorf("revision metrics line delta is inconsistent")
	}
	if m.IRChanged != (m.InputIRDigest != m.ProposedIRDigest) {
		return fmt.Errorf("revision metrics IR change flag is inconsistent")
	}
	if digestRevisionMetrics(m) != m.MetricsDigest {
		return fmt.Errorf("revision metrics digest does not match its fields")
	}
	return nil
}

func changedBytes(before, after []byte) int {
	changed := absInt(len(after) - len(before))
	for index := 0; index < len(before) && index < len(after); index++ {
		if before[index] != after[index] {
			changed++
		}
	}
	return changed
}

func changedLines(before, after []string) int {
	changed := absInt(len(after) - len(before))
	for index := 0; index < len(before) && index < len(after); index++ {
		if before[index] != after[index] {
			changed++
		}
	}
	return changed
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func digestRevisionMetrics(metrics RevisionMetrics) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%d|%d|%d|%d|%d|%d|%t|%t|%t",
		metrics.Status,
		metrics.MissingStage,
		metrics.SourceDigest,
		metrics.ProposedSourceDigest,
		metrics.InputIRDigest,
		metrics.ProposedIRDigest,
		metrics.CandidateDigest,
		metrics.EditDigest,
		metrics.ApplicationDigest,
		metrics.SourceBytesBefore,
		metrics.SourceBytesAfter,
		metrics.SourceBytesDelta,
		metrics.ChangedByteCount,
		metrics.SourceLinesBefore,
		metrics.SourceLinesAfter,
		metrics.SourceLinesDelta,
		metrics.ChangedLineCount,
		metrics.IRChanged,
		metrics.NonExecuting,
		metrics.NonAuthorizing,
	))
}
