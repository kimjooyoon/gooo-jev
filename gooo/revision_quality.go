package gooo

import "fmt"

type RevisionQuality struct {
	Status               string
	MissingStage         string
	SourceDigest         string
	ProposedSourceDigest string
	ApplicationDigest    string
	MetricsDigest        string
	ChangedByteCount     int
	ChangedLineCount     int
	IRChanged            bool
	ChangeClass          string
	ReviewSignal         string
	QualityDigest        string
	NonExecuting         bool
	NonAuthorizing       bool
}

func EvaluateRevisionMetrics(metrics RevisionMetrics) (RevisionQuality, error) {
	quality := RevisionQuality{
		Status:               "UNKNOWN",
		MissingStage:         "revision-quality",
		SourceDigest:         metrics.SourceDigest,
		ProposedSourceDigest: metrics.ProposedSourceDigest,
		ApplicationDigest:    metrics.ApplicationDigest,
		MetricsDigest:        metrics.MetricsDigest,
		ChangedByteCount:     metrics.ChangedByteCount,
		ChangedLineCount:     metrics.ChangedLineCount,
		IRChanged:            metrics.IRChanged,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if err := metrics.Validate(); err != nil {
		quality.MissingStage = "revision-quality-metrics"
		return quality, fmt.Errorf("gooo revision quality: metrics: %w", err)
	}
	quality.Status = "BOUND"
	quality.MissingStage = ""
	quality.ChangeClass = "localized"
	quality.ReviewSignal = "observe"
	if metrics.IRChanged {
		quality.ChangeClass = "structural"
		quality.ReviewSignal = "inspect"
	} else if metrics.ChangedLineCount > 1 || metrics.ChangedByteCount > 32 {
		quality.ChangeClass = "broad"
		quality.ReviewSignal = "inspect"
	}
	quality.QualityDigest = digestRevisionQuality(quality)
	return quality, nil
}

func (q RevisionQuality) Validate() error {
	if q.Status != "BOUND" {
		return fmt.Errorf("revision quality status must be BOUND")
	}
	if q.MissingStage != "" {
		return fmt.Errorf("revision quality missing stage must be empty")
	}
	if !q.NonExecuting || !q.NonAuthorizing {
		return fmt.Errorf("revision quality must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": q.SourceDigest, "proposed source": q.ProposedSourceDigest,
		"application": q.ApplicationDigest, "metrics": q.MetricsDigest,
		"quality": q.QualityDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision quality %s digest is invalid", name)
		}
	}
	if q.ChangedByteCount < 0 || q.ChangedLineCount < 0 {
		return fmt.Errorf("revision quality change counts must be non-negative")
	}
	if q.ChangeClass != "localized" && q.ChangeClass != "structural" && q.ChangeClass != "broad" {
		return fmt.Errorf("revision quality change class is invalid")
	}
	if q.ReviewSignal != "observe" && q.ReviewSignal != "inspect" {
		return fmt.Errorf("revision quality review signal is invalid")
	}
	if q.IRChanged && q.ChangeClass != "structural" {
		return fmt.Errorf("IR change must be structural")
	}
	if !q.IRChanged && q.ChangeClass == "structural" {
		return fmt.Errorf("structural class requires an IR change")
	}
	if digestRevisionQuality(q) != q.QualityDigest {
		return fmt.Errorf("revision quality digest does not match its fields")
	}
	return nil
}

func digestRevisionQuality(quality RevisionQuality) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d|%d|%t|%s|%s|%t|%t",
		quality.Status,
		quality.MissingStage,
		quality.SourceDigest,
		quality.ProposedSourceDigest,
		quality.ApplicationDigest,
		quality.MetricsDigest,
		quality.ChangedByteCount,
		quality.ChangedLineCount,
		quality.IRChanged,
		quality.ChangeClass,
		quality.ReviewSignal,
		quality.NonExecuting,
		quality.NonAuthorizing,
	))
}
