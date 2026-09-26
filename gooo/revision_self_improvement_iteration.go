package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementIteration rebinds the current language snapshot to
// the last ordered self-improvement observation without executing or
// authorizing the next change.
type RevisionSelfImprovementIteration struct {
	Status                            string
	MissingStage                      string
	StageCount                        int
	SourceDigest                      string
	IRDigest                          string
	HistoryDigest                     string
	WindowDigest                      string
	FeedbackDigest                    string
	CandidateSourceDigest             string
	CandidateProposedSourceDigest     string
	CandidateGeneratedIRDigest        string
	HistorySignal                     string
	FeedbackSignal                    string
	ObservationCount                  int
	IterationDigest                   string
	NonExecuting                      bool
	NonAuthorizing                    bool
}

// ObserveRevisionSelfImprovementIteration binds the current snapshot, ordered
// history, and last feedback into one evidence-linked iteration boundary.
func ObserveRevisionSelfImprovementIteration(source string, history RevisionSelfImprovementHistory, feedback RevisionSelfImprovementFeedback) (RevisionSelfImprovementIteration, error) {
	iteration := RevisionSelfImprovementIteration{
		Status:                        "UNKNOWN",
		MissingStage:                  "revision-self-improvement-iteration",
		StageCount:                    3,
		SourceDigest:                  digestString(source),
		HistoryDigest:                 history.HistoryDigest,
		WindowDigest:                  history.LastWindowDigest,
		FeedbackDigest:                feedback.FeedbackDigest,
		CandidateSourceDigest:         history.LastCandidateSourceDigest,
		CandidateProposedSourceDigest: history.LastCandidateProposedSourceDigest,
		CandidateGeneratedIRDigest:    history.LastCandidateGeneratedIRDigest,
		HistorySignal:                 history.HistorySignal,
		FeedbackSignal:                feedback.FeedbackSignal,
		ObservationCount:              history.ObservationCount,
		NonExecuting:                  true,
		NonAuthorizing:                true,
	}
	setIterationDigest := func() {
		iteration.IterationDigest = digestRevisionSelfImprovementIteration(iteration)
	}
	setIterationDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		iteration.MissingStage = "revision-self-improvement-iteration-snapshot"
		setIterationDigest()
		return iteration, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := validateBoundSelfImprovementHistory(history); err != nil {
		iteration.MissingStage = "revision-self-improvement-iteration-history"
		setIterationDigest()
		return iteration, fmt.Errorf("self-improvement history is not valid: %w", err)
	}
	if err := validateBoundSelfImprovementFeedback(feedback); err != nil {
		iteration.MissingStage = "revision-self-improvement-iteration-feedback"
		setIterationDigest()
		return iteration, fmt.Errorf("self-improvement feedback is not valid: %w", err)
	}
	if snapshot.SourceDigest != history.LastCandidateSourceDigest {
		iteration.MissingStage = "revision-self-improvement-iteration-source-link"
		setIterationDigest()
		return iteration, fmt.Errorf("current source does not match the last history candidate source")
	}
	if history.LastFeedbackDigest != feedback.FeedbackDigest ||
		history.LastWindowDigest != feedback.WindowDigest {
		iteration.MissingStage = "revision-self-improvement-iteration-feedback-link"
		setIterationDigest()
		return iteration, fmt.Errorf("feedback is not linked to the last history observation")
	}

	iteration.Status = "BOUND"
	iteration.MissingStage = ""
	iteration.IRDigest = snapshot.IRDigest
	iteration.SourceDigest = snapshot.SourceDigest
	setIterationDigest()
	if err := iteration.Validate(); err != nil {
		iteration.Status = "UNKNOWN"
		iteration.MissingStage = "revision-self-improvement-iteration"
		setIterationDigest()
		return iteration, fmt.Errorf("revision self-improvement iteration is not valid: %w", err)
	}
	return iteration, nil
}

func validateBoundSelfImprovementHistory(history RevisionSelfImprovementHistory) error {
	if err := history.Validate(); err != nil {
		return err
	}
	if history.Status != "BOUND" || history.MissingStage != "" {
		return fmt.Errorf("history is not BOUND")
	}
	return nil
}

func (i RevisionSelfImprovementIteration) Validate() error {
	if i.Status == "" {
		return fmt.Errorf("revision self-improvement iteration status is empty")
	}
	if i.Status == "BOUND" && i.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement iteration has a missing stage")
	}
	if i.Status == "UNKNOWN" && i.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement iteration has no missing stage")
	}
	if i.StageCount != 3 {
		return fmt.Errorf("revision self-improvement iteration stage count must be 3")
	}
	if i.ObservationCount < 1 {
		return fmt.Errorf("revision self-improvement iteration observation count must be positive")
	}
	if !validDigest(i.SourceDigest) || !validDigest(i.HistoryDigest) ||
		!validDigest(i.WindowDigest) || !validDigest(i.FeedbackDigest) ||
		!validDigest(i.CandidateSourceDigest) ||
		!validDigest(i.CandidateProposedSourceDigest) ||
		!validDigest(i.CandidateGeneratedIRDigest) {
		return fmt.Errorf("revision self-improvement iteration provenance digest is invalid")
	}
	if i.HistorySignal != "stable" && i.HistorySignal != "narrower" &&
		i.HistorySignal != "wider" && i.HistorySignal != "mixed" {
		return fmt.Errorf("revision self-improvement iteration history signal is invalid")
	}
	if i.FeedbackSignal != "observe" && i.FeedbackSignal != "remeasure" &&
		i.FeedbackSignal != "review" && i.FeedbackSignal != "inspect" {
		return fmt.Errorf("revision self-improvement iteration feedback signal is invalid")
	}
	if !i.NonExecuting || !i.NonAuthorizing {
		return fmt.Errorf("revision self-improvement iteration must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementIteration(i) != i.IterationDigest {
		return fmt.Errorf("revision self-improvement iteration digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementIteration(iteration RevisionSelfImprovementIteration) string {
	parts := []string{
		iteration.Status,
		iteration.MissingStage,
		strconv.Itoa(iteration.StageCount),
		iteration.SourceDigest,
		iteration.IRDigest,
		iteration.HistoryDigest,
		iteration.WindowDigest,
		iteration.FeedbackDigest,
		iteration.CandidateSourceDigest,
		iteration.CandidateProposedSourceDigest,
		iteration.CandidateGeneratedIRDigest,
		iteration.HistorySignal,
		iteration.FeedbackSignal,
		strconv.Itoa(iteration.ObservationCount),
		strconv.FormatBool(iteration.NonExecuting),
		strconv.FormatBool(iteration.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}