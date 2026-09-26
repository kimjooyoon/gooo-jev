package gooo

import "fmt"

type RevisionFeedback struct {
	Status               string
	MissingStage         string
	SourceDigest         string
	ProposedSourceDigest string
	ApplicationDigest    string
	MetricsDigest        string
	QualityDigest        string
	ChangeClass          string
	ActionHint           string
	FeedbackDigest       string
	NonExecuting         bool
	NonAuthorizing       bool
}

func DeriveRevisionFeedback(quality RevisionQuality) (RevisionFeedback, error) {
	feedback := RevisionFeedback{
		Status:               "UNKNOWN",
		MissingStage:         "revision-feedback",
		SourceDigest:         quality.SourceDigest,
		ProposedSourceDigest: quality.ProposedSourceDigest,
		ApplicationDigest:    quality.ApplicationDigest,
		MetricsDigest:        quality.MetricsDigest,
		QualityDigest:        quality.QualityDigest,
		ChangeClass:          quality.ChangeClass,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if err := quality.Validate(); err != nil {
		feedback.MissingStage = "revision-feedback-quality"
		return feedback, fmt.Errorf("gooo revision feedback: quality: %w", err)
	}
	feedback.Status = "BOUND"
	feedback.MissingStage = ""
	switch quality.ChangeClass {
	case "localized":
		feedback.ActionHint = "observe-more"
	case "structural":
		feedback.ActionHint = "inspect-ir"
	case "broad":
		feedback.ActionHint = "bound-scope"
	}
	feedback.FeedbackDigest = digestRevisionFeedback(feedback)
	return feedback, nil
}

func (f RevisionFeedback) Validate() error {
	if f.Status != "BOUND" {
		return fmt.Errorf("revision feedback status must be BOUND")
	}
	if f.MissingStage != "" {
		return fmt.Errorf("revision feedback missing stage must be empty")
	}
	if !f.NonExecuting || !f.NonAuthorizing {
		return fmt.Errorf("revision feedback must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": f.SourceDigest, "proposed source": f.ProposedSourceDigest,
		"application": f.ApplicationDigest, "metrics": f.MetricsDigest,
		"quality": f.QualityDigest, "feedback": f.FeedbackDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision feedback %s digest is invalid", name)
		}
	}
	if f.ChangeClass != "localized" && f.ChangeClass != "structural" && f.ChangeClass != "broad" {
		return fmt.Errorf("revision feedback change class is invalid")
	}
	if f.ActionHint != "observe-more" && f.ActionHint != "inspect-ir" && f.ActionHint != "bound-scope" {
		return fmt.Errorf("revision feedback action hint is invalid")
	}
	if digestRevisionFeedback(f) != f.FeedbackDigest {
		return fmt.Errorf("revision feedback digest does not match its fields")
	}
	return nil
}

func digestRevisionFeedback(feedback RevisionFeedback) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		feedback.Status,
		feedback.MissingStage,
		feedback.SourceDigest,
		feedback.ProposedSourceDigest,
		feedback.ApplicationDigest,
		feedback.MetricsDigest,
		feedback.QualityDigest,
		feedback.ChangeClass,
		feedback.ActionHint,
		feedback.NonExecuting,
		feedback.NonAuthorizing,
	))
}
