package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceHistoryObservation preserves an ordered
// series of lifecycle provenance transitions without claiming improvement.
type RevisionSelfImprovementProvenanceHistoryObservation struct {
	Status                            string
	MissingStage                      string
	ObservationCount                  int
	TransitionDigests                 []string
	PreviousLifecycleObservationDigests []string
	CurrentLifecycleObservationDigests  []string
	FirstTransitionDigest             string
	LastTransitionDigest              string
	FirstPreviousLifecycleDigest      string
	LastCurrentLifecycleDigest        string
	StableCount                       int
	TransitionedCount                 int
	MetricsChangeCount               int
	GenerationChangeCount            int
	HistorySignal                    string
	HistoryDigest                    string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

// ObserveRevisionSelfImprovementProvenanceHistory binds ordered transition
// evidence and preserves the first unresolved transition stage.
func ObserveRevisionSelfImprovementProvenanceHistory(
	transitions []RevisionSelfImprovementProvenanceTransitionObservation,
) (RevisionSelfImprovementProvenanceHistoryObservation, error) {
	history := RevisionSelfImprovementProvenanceHistoryObservation{
		Status:                              "UNKNOWN",
		MissingStage:                        "revision-self-improvement-provenance-history-transitions",
		ObservationCount:                    len(transitions),
		TransitionDigests:                   make([]string, 0, len(transitions)),
		PreviousLifecycleObservationDigests: make([]string, 0, len(transitions)),
		CurrentLifecycleObservationDigests:  make([]string, 0, len(transitions)),
		NonExecuting:                        true,
		NonAuthorizing:                      true,
	}
	setHistoryDigest := func() {
		history.HistoryDigest = digestRevisionSelfImprovementProvenanceHistory(history)
	}
	setHistoryDigest()

	if len(transitions) == 0 {
		return history, fmt.Errorf("self-improvement provenance history requires at least one transition")
	}

	for index, transition := range transitions {
		if err := validateBoundSelfImprovementProvenanceTransition(transition); err != nil {
			history.MissingStage = fmt.Sprintf("revision-self-improvement-provenance-history-transition-%d", index)
			setHistoryDigest()
			return history, fmt.Errorf("transition %d is not valid: %w", index, err)
		}
		history.TransitionDigests = append(history.TransitionDigests, transition.ObservationDigest)
		history.PreviousLifecycleObservationDigests = append(
			history.PreviousLifecycleObservationDigests,
			transition.PreviousLifecycleObservationDigest,
		)
		history.CurrentLifecycleObservationDigests = append(
			history.CurrentLifecycleObservationDigests,
			transition.CurrentLifecycleObservationDigest,
		)
		switch transition.TransitionSignal {
		case "provenance-stable":
			history.StableCount++
		case "provenance-transition-bound":
			history.TransitionedCount++
		}
		if transition.MetricsChanged {
			history.MetricsChangeCount++
		}
		if transition.GenerationChanged {
			history.GenerationChangeCount++
		}
	}

	history.FirstTransitionDigest = history.TransitionDigests[0]
	history.LastTransitionDigest = history.TransitionDigests[len(history.TransitionDigests)-1]
	history.FirstPreviousLifecycleDigest = history.PreviousLifecycleObservationDigests[0]
	history.LastCurrentLifecycleDigest = history.CurrentLifecycleObservationDigests[len(history.CurrentLifecycleObservationDigests)-1]
	history.HistorySignal = revisionSelfImprovementProvenanceHistorySignal(history)
	history.Status = "BOUND"
	history.MissingStage = ""
	setHistoryDigest()
	if err := history.Validate(); err != nil {
		history.Status = "UNKNOWN"
		history.MissingStage = "revision-self-improvement-provenance-history"
		setHistoryDigest()
		return history, fmt.Errorf("self-improvement provenance history is not valid: %w", err)
	}
	return history, nil
}

func validateBoundSelfImprovementProvenanceTransition(
	transition RevisionSelfImprovementProvenanceTransitionObservation,
) error {
	if err := transition.Validate(); err != nil {
		return err
	}
	if transition.Status != "BOUND" || transition.MissingStage != "" {
		return fmt.Errorf("provenance transition is not BOUND")
	}
	if transition.TransitionSignal != "provenance-stable" &&
		transition.TransitionSignal != "provenance-transition-bound" {
		return fmt.Errorf("provenance transition signal is not actionable")
	}
	return nil
}

func revisionSelfImprovementProvenanceHistorySignal(
	history RevisionSelfImprovementProvenanceHistoryObservation,
) string {
	switch {
	case history.StableCount == history.ObservationCount:
		return "stable"
	case history.TransitionedCount == history.ObservationCount:
		return "transitioned"
	default:
		return "mixed"
	}
}

func (h RevisionSelfImprovementProvenanceHistoryObservation) Validate() error {
	if h.Status == "" {
		return fmt.Errorf("self-improvement provenance history status is empty")
	}
	if h.Status == "BOUND" && h.MissingStage != "" {
		return fmt.Errorf("bound self-improvement provenance history has a missing stage")
	}
	if h.Status == "UNKNOWN" && h.MissingStage == "" {
		return fmt.Errorf("unknown self-improvement provenance history has no missing stage")
	}
	if h.ObservationCount < 1 {
		return fmt.Errorf("self-improvement provenance history observation count must be positive")
	}
	if len(h.TransitionDigests) != h.ObservationCount ||
		len(h.PreviousLifecycleObservationDigests) != h.ObservationCount ||
		len(h.CurrentLifecycleObservationDigests) != h.ObservationCount {
		return fmt.Errorf("self-improvement provenance history digest counts do not match observations")
	}
	if h.TransitionDigests[0] != h.FirstTransitionDigest ||
		h.TransitionDigests[len(h.TransitionDigests)-1] != h.LastTransitionDigest ||
		h.PreviousLifecycleObservationDigests[0] != h.FirstPreviousLifecycleDigest ||
		h.CurrentLifecycleObservationDigests[len(h.CurrentLifecycleObservationDigests)-1] != h.LastCurrentLifecycleDigest {
		return fmt.Errorf("self-improvement provenance history boundary digests are not linked")
	}
	for index, digest := range h.TransitionDigests {
		if !validDigest(digest) ||
			!validDigest(h.PreviousLifecycleObservationDigests[index]) ||
			!validDigest(h.CurrentLifecycleObservationDigests[index]) {
			return fmt.Errorf("self-improvement provenance history digest %d is invalid", index)
		}
	}
	if h.StableCount < 0 || h.TransitionedCount < 0 ||
		h.StableCount+h.TransitionedCount != h.ObservationCount {
		return fmt.Errorf("self-improvement provenance history transition counts are invalid")
	}
	if h.MetricsChangeCount < 0 || h.MetricsChangeCount > h.ObservationCount ||
		h.GenerationChangeCount < 0 || h.GenerationChangeCount > h.ObservationCount {
		return fmt.Errorf("self-improvement provenance history change counts are invalid")
	}
	if h.HistorySignal != "stable" && h.HistorySignal != "transitioned" &&
		h.HistorySignal != "mixed" {
		return fmt.Errorf("self-improvement provenance history signal is invalid")
	}
	if revisionSelfImprovementProvenanceHistorySignal(h) != h.HistorySignal {
		return fmt.Errorf("self-improvement provenance history signal does not match counts")
	}
	if !h.NonExecuting || !h.NonAuthorizing {
		return fmt.Errorf("self-improvement provenance history must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceHistory(h) != h.HistoryDigest {
		return fmt.Errorf("self-improvement provenance history digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceHistory(
	history RevisionSelfImprovementProvenanceHistoryObservation,
) string {
	parts := []string{
		history.Status,
		history.MissingStage,
		strconv.Itoa(history.ObservationCount),
		strings.Join(history.TransitionDigests, ","),
		strings.Join(history.PreviousLifecycleObservationDigests, ","),
		strings.Join(history.CurrentLifecycleObservationDigests, ","),
		history.FirstTransitionDigest,
		history.LastTransitionDigest,
		history.FirstPreviousLifecycleDigest,
		history.LastCurrentLifecycleDigest,
		strconv.Itoa(history.StableCount),
		strconv.Itoa(history.TransitionedCount),
		strconv.Itoa(history.MetricsChangeCount),
		strconv.Itoa(history.GenerationChangeCount),
		history.HistorySignal,
		strconv.FormatBool(history.NonExecuting),
		strconv.FormatBool(history.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
