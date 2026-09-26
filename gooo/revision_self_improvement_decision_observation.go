package gooo

import "fmt"

// RevisionSelfImprovementDecisionObservation records the next bounded
// observation action implied by validated outcome and feedback evidence.
type RevisionSelfImprovementDecisionObservation struct {
	Status                string
	MissingStage          string
	IterationDigest       string
	OutcomeDigest         string
	FeedbackDigest        string
	HistoryDigest         string
	SourceDigest          string
	ProposedSourceDigest  string
	GeneratedIRDigest     string
	ComparisonSignal      string
	FeedbackSignal        string
	DecisionSignal        string
	DecisionReason        string
	RequiresObservation   bool
	RequiresReview        bool
	RequiresInspection    bool
	RequiresMeasurement   bool
	DecisionDigest        string
	NonExecuting          bool
	NonAuthorizing        bool
}

// ObserveRevisionSelfImprovementDecision derives a bounded next-observation
// action without applying, executing, or authorizing a change.
func ObserveRevisionSelfImprovementDecision(iteration RevisionSelfImprovementIteration, outcome RevisionSelfImprovementOutcomeObservation, feedback RevisionSelfImprovementFeedback) (RevisionSelfImprovementDecisionObservation, error) {
	result := RevisionSelfImprovementDecisionObservation{
		Status:               "UNKNOWN",
		MissingStage:         "revision-self-improvement-decision",
		IterationDigest:      iteration.IterationDigest,
		OutcomeDigest:        outcome.OutcomeDigest,
		FeedbackDigest:       feedback.FeedbackDigest,
		HistoryDigest:        iteration.HistoryDigest,
		SourceDigest:         outcome.SourceDigest,
		ProposedSourceDigest: outcome.ProposedSourceDigest,
		GeneratedIRDigest:    outcome.GeneratedIRDigest,
		ComparisonSignal:     feedback.ComparisonSignal,
		FeedbackSignal:       feedback.FeedbackSignal,
		DecisionSignal:       feedback.FeedbackSignal,
		DecisionReason:       feedback.FeedbackReason,
		RequiresObservation: feedback.RequiresObservation,
		RequiresReview:      feedback.RequiresReview,
		RequiresInspection:  feedback.RequiresInspection,
		RequiresMeasurement: feedback.RequiresMeasurement,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setDecisionDigest := func() {
		result.DecisionDigest = digestRevisionSelfImprovementDecision(result)
	}
	setDecisionDigest()

	if err := iteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-decision-iteration"
		setDecisionDigest()
		return result, fmt.Errorf("self-improvement iteration is not valid: %w", err)
	}
	if err := outcome.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-decision-outcome"
		setDecisionDigest()
		return result, fmt.Errorf("self-improvement outcome is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-decision-feedback"
		setDecisionDigest()
		return result, fmt.Errorf("self-improvement feedback is not valid: %w", err)
	}
	if outcome.IterationDigest != iteration.IterationDigest ||
		outcome.SourceDigest != iteration.CandidateSourceDigest ||
		outcome.ProposedSourceDigest != iteration.CandidateProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-decision-iteration-outcome-link"
		setDecisionDigest()
		return result, fmt.Errorf("outcome is not linked to the current self-improvement iteration")
	}
	if feedback.FeedbackDigest != iteration.FeedbackDigest ||
		feedback.WindowDigest != iteration.WindowDigest ||
		feedback.CandidateSourceDigest != iteration.CandidateSourceDigest ||
		feedback.CandidateProposedSourceDigest != iteration.CandidateProposedSourceDigest ||
		feedback.CandidateGeneratedIRDigest != iteration.CandidateGeneratedIRDigest {
		result.MissingStage = "revision-self-improvement-decision-iteration-feedback-link"
		setDecisionDigest()
		return result, fmt.Errorf("feedback is not linked to the current self-improvement iteration")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDecisionDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-decision"
		setDecisionDigest()
		return result, fmt.Errorf("revision self-improvement decision is not valid: %w", err)
	}
	return result, nil
}

func (d RevisionSelfImprovementDecisionObservation) Validate() error {
	if d.Status == "" {
		return fmt.Errorf("revision self-improvement decision status is empty")
	}
	if d.Status == "BOUND" && d.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement decision has a missing stage")
	}
	if d.Status == "UNKNOWN" && d.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement decision has no missing stage")
	}
	for name, digest := range map[string]string{
		"iteration":       d.IterationDigest,
		"outcome":         d.OutcomeDigest,
		"feedback":        d.FeedbackDigest,
		"history":         d.HistoryDigest,
		"source":          d.SourceDigest,
		"proposed source": d.ProposedSourceDigest,
		"generated IR":    d.GeneratedIRDigest,
		"decision":        d.DecisionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision self-improvement decision %s digest is invalid", name)
		}
	}
	if d.ComparisonSignal != "stable" && d.ComparisonSignal != "narrower" &&
		d.ComparisonSignal != "wider" && d.ComparisonSignal != "mixed" {
		return fmt.Errorf("revision self-improvement decision comparison signal is invalid")
	}
	if d.FeedbackSignal != "observe" && d.FeedbackSignal != "remeasure" &&
		d.FeedbackSignal != "review" && d.FeedbackSignal != "inspect" {
		return fmt.Errorf("revision self-improvement decision feedback signal is invalid")
	}
	if d.DecisionSignal != d.FeedbackSignal {
		return fmt.Errorf("revision self-improvement decision signal is not linked to feedback")
	}
	if d.Status == "BOUND" && d.DecisionReason == "" {
		return fmt.Errorf("bound revision self-improvement decision has no reason")
	}
	if !d.RequiresObservation {
		return fmt.Errorf("revision self-improvement decision must require observation")
	}
	switch d.ComparisonSignal {
	case "stable":
		if d.DecisionSignal != "observe" || d.RequiresReview || d.RequiresInspection || d.RequiresMeasurement {
			return fmt.Errorf("stable decision does not match its signal")
		}
	case "narrower":
		if d.DecisionSignal != "remeasure" || d.RequiresReview || d.RequiresInspection || !d.RequiresMeasurement {
			return fmt.Errorf("narrower decision does not match its signal")
		}
	case "wider":
		if d.DecisionSignal != "review" || d.RequiresMeasurement || d.RequiresInspection || !d.RequiresReview {
			return fmt.Errorf("wider decision does not match its signal")
		}
	case "mixed":
		if d.DecisionSignal != "inspect" || d.RequiresReview || !d.RequiresMeasurement || !d.RequiresInspection {
			return fmt.Errorf("mixed decision does not match its signal")
		}
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("revision self-improvement decision must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementDecision(d) != d.DecisionDigest {
		return fmt.Errorf("revision self-improvement decision digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementDecision(decision RevisionSelfImprovementDecisionObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t|%t",
		decision.Status,
		decision.MissingStage,
		decision.IterationDigest,
		decision.OutcomeDigest,
		decision.FeedbackDigest,
		decision.HistoryDigest,
		decision.SourceDigest,
		decision.ProposedSourceDigest,
		decision.GeneratedIRDigest,
		decision.ComparisonSignal,
		decision.FeedbackSignal,
		decision.DecisionSignal,
		decision.DecisionReason,
		decision.RequiresObservation,
		decision.RequiresReview,
		decision.RequiresInspection,
		decision.RequiresMeasurement,
		decision.NonExecuting,
		decision.NonAuthorizing,
	))
}
