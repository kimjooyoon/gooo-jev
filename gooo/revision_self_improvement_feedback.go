package gooo

import "fmt"

// RevisionSelfImprovementFeedback is a bounded next-observation signal derived
// from a self-improvement window. It never applies or authorizes a change.
type RevisionSelfImprovementFeedback struct {
	Status                    string
	MissingStage              string
	WindowDigest              string
	BaselineReceiptDigest     string
	CandidateReceiptDigest    string
	ComparisonSignal          string
	FeedbackSignal            string
	FeedbackReason            string
	CandidateSourceDigest     string
	CandidateProposedSourceDigest string
	CandidateGeneratedIRDigest    string
	RequiresObservation       bool
	RequiresReview            bool
	RequiresInspection        bool
	RequiresMeasurement       bool
	FeedbackDigest             string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ObserveRevisionSelfImprovementFeedback maps observed window evidence to a
// bounded next-observation signal without executing or authorizing a change.
func ObserveRevisionSelfImprovementFeedback(window RevisionSelfImprovementWindow) (RevisionSelfImprovementFeedback, error) {
	feedback := RevisionSelfImprovementFeedback{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-feedback",
		WindowDigest:              window.WindowDigest,
		BaselineReceiptDigest:     window.BaselineReceiptDigest,
		CandidateReceiptDigest:    window.CandidateReceiptDigest,
		ComparisonSignal:          window.ComparisonSignal,
		CandidateSourceDigest:     window.CandidateSourceDigest,
		CandidateProposedSourceDigest: window.CandidateProposedSourceDigest,
		CandidateGeneratedIRDigest:    window.CandidateGeneratedIRDigest,
		RequiresObservation:       true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setFeedbackDigest := func() {
		feedback.FeedbackDigest = digestRevisionSelfImprovementFeedback(feedback)
	}
	setFeedbackDigest()

	if err := validateBoundSelfImprovementWindow(window); err != nil {
		feedback.MissingStage = "revision-self-improvement-feedback-window"
		setFeedbackDigest()
		return feedback, fmt.Errorf("self-improvement window is not valid: %w", err)
	}

	switch window.ComparisonSignal {
	case "stable":
		feedback.FeedbackSignal = "observe"
		feedback.FeedbackReason = "stable-observation"
	case "narrower":
		feedback.FeedbackSignal = "remeasure"
		feedback.FeedbackReason = "narrower-surface-not-improvement"
		feedback.RequiresMeasurement = true
	case "wider":
		feedback.FeedbackSignal = "review"
		feedback.FeedbackReason = "wider-surface-requires-review"
		feedback.RequiresReview = true
	case "mixed":
		feedback.FeedbackSignal = "inspect"
		feedback.FeedbackReason = "mixed-evidence-requires-inspection"
		feedback.RequiresInspection = true
	default:
		feedback.MissingStage = "revision-self-improvement-feedback-signal"
		setFeedbackDigest()
		return feedback, fmt.Errorf("self-improvement window signal is not recognized")
	}

	feedback.Status = "BOUND"
	feedback.MissingStage = ""
	setFeedbackDigest()
	if err := feedback.Validate(); err != nil {
		feedback.Status = "UNKNOWN"
		feedback.MissingStage = "revision-self-improvement-feedback"
		setFeedbackDigest()
		return feedback, fmt.Errorf("revision self-improvement feedback is not valid: %w", err)
	}
	return feedback, nil
}

func validateBoundSelfImprovementWindow(window RevisionSelfImprovementWindow) error {
	if err := window.Validate(); err != nil {
		return err
	}
	if window.Status != "BOUND" || window.MissingStage != "" {
		return fmt.Errorf("window is not BOUND")
	}
	return nil
}

func (f RevisionSelfImprovementFeedback) Validate() error {
	if f.Status == "" {
		return fmt.Errorf("revision self-improvement feedback status is empty")
	}
	if f.Status == "BOUND" && f.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement feedback has a missing stage")
	}
	if f.Status == "UNKNOWN" && f.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement feedback has no missing stage")
	}
	if f.ComparisonSignal != "narrower" && f.ComparisonSignal != "wider" &&
		f.ComparisonSignal != "stable" && f.ComparisonSignal != "mixed" {
		return fmt.Errorf("revision self-improvement feedback comparison signal is invalid")
	}
	if f.FeedbackSignal != "observe" && f.FeedbackSignal != "remeasure" &&
		f.FeedbackSignal != "review" && f.FeedbackSignal != "inspect" {
		return fmt.Errorf("revision self-improvement feedback signal is invalid")
	}
	if f.Status == "BOUND" && f.FeedbackReason == "" {
		return fmt.Errorf("bound revision self-improvement feedback has no reason")
	}
	if !f.RequiresObservation {
		return fmt.Errorf("revision self-improvement feedback must require observation")
	}
	if !f.NonExecuting || !f.NonAuthorizing {
		return fmt.Errorf("revision self-improvement feedback must remain non-executing and non-authorizing")
	}
	if f.Status == "BOUND" {
		switch f.ComparisonSignal {
		case "stable":
			if f.FeedbackSignal != "observe" || f.RequiresReview || f.RequiresInspection || f.RequiresMeasurement {
				return fmt.Errorf("stable feedback does not match its signal")
			}
		case "narrower":
			if f.FeedbackSignal != "remeasure" || f.RequiresReview || f.RequiresInspection || !f.RequiresMeasurement {
				return fmt.Errorf("narrower feedback does not match its signal")
			}
		case "wider":
			if f.FeedbackSignal != "review" || !f.RequiresMeasurement || f.RequiresInspection || f.RequiresReview == false {
				return fmt.Errorf("wider feedback does not match its signal")
			}
		case "mixed":
			if f.FeedbackSignal != "inspect" || f.RequiresReview || !f.RequiresMeasurement || f.RequiresInspection == false {
				return fmt.Errorf("mixed feedback does not match its signal")
			}
		}
	}
	if digestRevisionSelfImprovementFeedback(f) != f.FeedbackDigest {
		return fmt.Errorf("revision self-improvement feedback digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementFeedback(feedback RevisionSelfImprovementFeedback) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t|%t",
		feedback.Status,
		feedback.MissingStage,
		feedback.WindowDigest,
		feedback.BaselineReceiptDigest,
		feedback.CandidateReceiptDigest,
		feedback.ComparisonSignal,
		feedback.FeedbackSignal,
		feedback.FeedbackReason,
		feedback.CandidateSourceDigest,
		feedback.CandidateProposedSourceDigest,
		feedback.CandidateGeneratedIRDigest,
		feedback.RequiresObservation,
		feedback.RequiresReview,
		feedback.RequiresInspection,
		feedback.RequiresMeasurement,
		feedback.NonExecuting,
		feedback.NonAuthorizing,
	))
}