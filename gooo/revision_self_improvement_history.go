package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementHistory preserves an ordered series of self-improvement
// windows and their next-observation feedback without claiming improvement.
type RevisionSelfImprovementHistory struct {
	Status                string
	MissingStage          string
	ObservationCount      int
	WindowDigests         []string
	FeedbackDigests       []string
	FirstWindowDigest     string
	LastWindowDigest      string
	FirstFeedbackDigest   string
	LastFeedbackDigest    string
	StableCount           int
	NarrowerCount         int
	WiderCount             int
	MixedCount             int
	ObserveCount          int
	RemeasureCount        int
	ReviewCount           int
	InspectCount          int
	HistorySignal         string
	HistoryDigest         string
	NonExecuting          bool
	NonAuthorizing        bool
}

// ObserveRevisionSelfImprovementHistory binds ordered windows to their
// corresponding feedback and preserves the first unresolved stage.
func ObserveRevisionSelfImprovementHistory(windows []RevisionSelfImprovementWindow, feedback []RevisionSelfImprovementFeedback) (RevisionSelfImprovementHistory, error) {
	history := RevisionSelfImprovementHistory{
		Status:           "UNKNOWN",
		MissingStage:     "revision-self-improvement-history-windows",
		ObservationCount: len(windows),
		WindowDigests:    make([]string, 0, len(windows)),
		FeedbackDigests:  make([]string, 0, len(feedback)),
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setHistoryDigest := func() {
		history.HistoryDigest = digestRevisionSelfImprovementHistory(history)
	}
	setHistoryDigest()

	if len(windows) == 0 {
		return history, fmt.Errorf("revision self-improvement history requires at least one window")
	}
	if len(windows) != len(feedback) {
		history.MissingStage = "revision-self-improvement-history-feedback-count"
		setHistoryDigest()
		return history, fmt.Errorf("revision self-improvement history windows and feedback counts differ")
	}

	for index, window := range windows {
		if err := validateBoundSelfImprovementWindow(window); err != nil {
			history.MissingStage = fmt.Sprintf("revision-self-improvement-history-window-%d", index)
			setHistoryDigest()
			return history, fmt.Errorf("window %d is not valid: %w", index, err)
		}
		history.WindowDigests = append(history.WindowDigests, window.WindowDigest)
		switch window.ComparisonSignal {
		case "stable":
			history.StableCount++
		case "narrower":
			history.NarrowerCount++
		case "wider":
			history.WiderCount++
		case "mixed":
			history.MixedCount++
		}

		currentFeedback := feedback[index]
		if err := validateBoundSelfImprovementFeedback(currentFeedback); err != nil {
			history.MissingStage = fmt.Sprintf("revision-self-improvement-history-feedback-%d", index)
			setHistoryDigest()
			return history, fmt.Errorf("feedback %d is not valid: %w", index, err)
		}
		if currentFeedback.WindowDigest != window.WindowDigest {
			history.MissingStage = fmt.Sprintf("revision-self-improvement-history-link-%d", index)
			setHistoryDigest()
			return history, fmt.Errorf("feedback %d is not linked to window %d", index, index)
		}
		history.FeedbackDigests = append(history.FeedbackDigests, currentFeedback.FeedbackDigest)
		switch currentFeedback.FeedbackSignal {
		case "observe":
			history.ObserveCount++
		case "remeasure":
			history.RemeasureCount++
		case "review":
			history.ReviewCount++
		case "inspect":
			history.InspectCount++
		}
	}

	history.FirstWindowDigest = history.WindowDigests[0]
	history.LastWindowDigest = history.WindowDigests[len(history.WindowDigests)-1]
	history.FirstFeedbackDigest = history.FeedbackDigests[0]
	history.LastFeedbackDigest = history.FeedbackDigests[len(history.FeedbackDigests)-1]
	history.HistorySignal = revisionSelfImprovementHistorySignal(history)
	history.Status = "BOUND"
	history.MissingStage = ""
	setHistoryDigest()
	if err := history.Validate(); err != nil {
		history.Status = "UNKNOWN"
		history.MissingStage = "revision-self-improvement-history"
		setHistoryDigest()
		return history, fmt.Errorf("revision self-improvement history is not valid: %w", err)
	}
	return history, nil
}

func validateBoundSelfImprovementFeedback(feedback RevisionSelfImprovementFeedback) error {
	if err := feedback.Validate(); err != nil {
		return err
	}
	if feedback.Status != "BOUND" || feedback.MissingStage != "" {
		return fmt.Errorf("feedback is not BOUND")
	}
	return nil
}

func revisionSelfImprovementHistorySignal(history RevisionSelfImprovementHistory) string {
	switch {
	case history.StableCount == history.ObservationCount:
		return "stable"
	case history.NarrowerCount == history.ObservationCount:
		return "narrower"
	case history.WiderCount == history.ObservationCount:
		return "wider"
	case history.MixedCount == history.ObservationCount:
		return "mixed"
	default:
		return "mixed"
	}
}

func (h RevisionSelfImprovementHistory) Validate() error {
	if h.Status == "" {
		return fmt.Errorf("revision self-improvement history status is empty")
	}
	if h.Status == "BOUND" && h.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement history has a missing stage")
	}
	if h.Status == "UNKNOWN" && h.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement history has no missing stage")
	}
	if h.ObservationCount < 1 {
		return fmt.Errorf("revision self-improvement history observation count must be positive")
	}
	if len(h.WindowDigests) != h.ObservationCount || len(h.FeedbackDigests) != h.ObservationCount {
		return fmt.Errorf("revision self-improvement history digest counts do not match observations")
	}
	if h.WindowDigests[0] != h.FirstWindowDigest ||
		h.WindowDigests[len(h.WindowDigests)-1] != h.LastWindowDigest ||
		h.FeedbackDigests[0] != h.FirstFeedbackDigest ||
		h.FeedbackDigests[len(h.FeedbackDigests)-1] != h.LastFeedbackDigest {
		return fmt.Errorf("revision self-improvement history boundary digests are not linked")
	}
	for index, digest := range h.WindowDigests {
		if !validDigest(digest) {
			return fmt.Errorf("revision self-improvement history window digest %d is invalid", index)
		}
		if !validDigest(h.FeedbackDigests[index]) {
			return fmt.Errorf("revision self-improvement history feedback digest %d is invalid", index)
		}
	}
	if h.StableCount < 0 || h.NarrowerCount < 0 || h.WiderCount < 0 || h.MixedCount < 0 ||
		h.StableCount+h.NarrowerCount+h.WiderCount+h.MixedCount != h.ObservationCount {
		return fmt.Errorf("revision self-improvement history window signal counts are invalid")
	}
	if h.ObserveCount < 0 || h.RemeasureCount < 0 || h.ReviewCount < 0 || h.InspectCount < 0 ||
		h.ObserveCount+h.RemeasureCount+h.ReviewCount+h.InspectCount != h.ObservationCount {
		return fmt.Errorf("revision self-improvement history feedback signal counts are invalid")
	}
	if h.HistorySignal != "stable" && h.HistorySignal != "narrower" &&
		h.HistorySignal != "wider" && h.HistorySignal != "mixed" {
		return fmt.Errorf("revision self-improvement history signal is invalid")
	}
	if revisionSelfImprovementHistorySignal(h) != h.HistorySignal {
		return fmt.Errorf("revision self-improvement history signal does not match counts")
	}
	if !h.NonExecuting || !h.NonAuthorizing {
		return fmt.Errorf("revision self-improvement history must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementHistory(h) != h.HistoryDigest {
		return fmt.Errorf("revision self-improvement history digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementHistory(history RevisionSelfImprovementHistory) string {
	parts := []string{
		history.Status,
		history.MissingStage,
		strconv.Itoa(history.ObservationCount),
		strings.Join(history.WindowDigests, ","),
		strings.Join(history.FeedbackDigests, ","),
		history.FirstWindowDigest,
		history.LastWindowDigest,
		history.FirstFeedbackDigest,
		history.LastFeedbackDigest,
		strconv.Itoa(history.StableCount),
		strconv.Itoa(history.NarrowerCount),
		strconv.Itoa(history.WiderCount),
		strconv.Itoa(history.MixedCount),
		strconv.Itoa(history.ObserveCount),
		strconv.Itoa(history.RemeasureCount),
		strconv.Itoa(history.ReviewCount),
		strconv.Itoa(history.InspectCount),
		history.HistorySignal,
		strconv.FormatBool(history.NonExecuting),
		strconv.FormatBool(history.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}