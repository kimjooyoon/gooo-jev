package gooo

import (
	"fmt"
	"strings"
)

type RevisionTrendWindow struct {
	Status                  string
	MissingStage            string
	ObservationCount        int
	TrendDigests            []string
	FirstTrendDigest        string
	LastTrendDigest         string
	NarrowerCount           int
	WiderCount              int
	StableCount             int
	MixedCount              int
	CandidateStableCount    int
	IRChangeStableCount     int
	SourceStableCount       int
	WindowSignal            string
	WindowDigest            string
	NonExecuting            bool
	NonAuthorizing          bool
}

func ObserveRevisionTrendWindow(observations []RevisionTrendObservation) (RevisionTrendWindow, error) {
	window := RevisionTrendWindow{
		Status:           "UNKNOWN",
		MissingStage:     "revision-trend-window-observations",
		ObservationCount: len(observations),
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if len(observations) == 0 {
		return window, fmt.Errorf("gooo revision trend window: at least one observation is required")
	}
	window.TrendDigests = make([]string, 0, len(observations))
	for index, observation := range observations {
		if err := observation.Validate(); err != nil {
			window.MissingStage = fmt.Sprintf("revision-trend-window-observation-%d", index)
			return window, fmt.Errorf("gooo revision trend window: observation %d: %w", index, err)
		}
		window.TrendDigests = append(window.TrendDigests, observation.TrendDigest)
		switch observation.ChangeSignal {
		case "narrower":
			window.NarrowerCount++
		case "wider":
			window.WiderCount++
		case "stable":
			window.StableCount++
		case "mixed":
			window.MixedCount++
		}
		if observation.CandidateStable {
			window.CandidateStableCount++
		}
		if observation.IRChangeStable {
			window.IRChangeStableCount++
		}
		if observation.SourceStable {
			window.SourceStableCount++
		}
	}
	window.FirstTrendDigest = window.TrendDigests[0]
	window.LastTrendDigest = window.TrendDigests[len(window.TrendDigests)-1]
	window.Status = "BOUND"
	window.MissingStage = ""
	window.WindowSignal = revisionTrendWindowSignal(window)
	window.WindowDigest = digestRevisionTrendWindow(window)
	return window, nil
}

func revisionTrendWindowSignal(window RevisionTrendWindow) string {
	switch {
	case window.NarrowerCount == window.ObservationCount:
		return "narrower"
	case window.WiderCount == window.ObservationCount:
		return "wider"
	case window.StableCount == window.ObservationCount:
		return "stable"
	case window.MixedCount == window.ObservationCount:
		return "mixed"
	default:
		return "mixed"
	}
}

func (w RevisionTrendWindow) Validate() error {
	if w.Status != "BOUND" {
		return fmt.Errorf("revision trend window status must be BOUND")
	}
	if w.MissingStage != "" {
		return fmt.Errorf("revision trend window missing stage must be empty")
	}
	if w.ObservationCount < 1 {
		return fmt.Errorf("revision trend window observation count must be positive")
	}
	if len(w.TrendDigests) != w.ObservationCount {
		return fmt.Errorf("revision trend window observation count does not match digests")
	}
	if !w.NonExecuting || !w.NonAuthorizing {
		return fmt.Errorf("revision trend window must remain non-executing and non-authorizing")
	}
	if !validDigest(w.FirstTrendDigest) || !validDigest(w.LastTrendDigest) || !validDigest(w.WindowDigest) {
		return fmt.Errorf("revision trend window boundary or window digest is invalid")
	}
	if w.TrendDigests[0] != w.FirstTrendDigest || w.TrendDigests[len(w.TrendDigests)-1] != w.LastTrendDigest {
		return fmt.Errorf("revision trend window boundary digest is not linked")
	}
	for index, digest := range w.TrendDigests {
		if !validDigest(digest) {
			return fmt.Errorf("revision trend window observation %d digest is invalid", index)
		}
	}
	if w.NarrowerCount < 0 || w.WiderCount < 0 || w.StableCount < 0 || w.MixedCount < 0 {
		return fmt.Errorf("revision trend window signal counts must be non-negative")
	}
	if w.NarrowerCount+w.WiderCount+w.StableCount+w.MixedCount != w.ObservationCount {
		return fmt.Errorf("revision trend window signal counts do not cover observations")
	}
	if w.CandidateStableCount < 0 || w.CandidateStableCount > w.ObservationCount ||
		w.IRChangeStableCount < 0 || w.IRChangeStableCount > w.ObservationCount ||
		w.SourceStableCount < 0 || w.SourceStableCount > w.ObservationCount {
		return fmt.Errorf("revision trend window stability counts are invalid")
	}
	if w.WindowSignal != "narrower" && w.WindowSignal != "wider" &&
		w.WindowSignal != "stable" && w.WindowSignal != "mixed" {
		return fmt.Errorf("revision trend window signal is invalid")
	}
	if revisionTrendWindowSignal(w) != w.WindowSignal {
		return fmt.Errorf("revision trend window signal does not match counts")
	}
	if digestRevisionTrendWindow(w) != w.WindowDigest {
		return fmt.Errorf("revision trend window digest does not match its fields")
	}
	return nil
}

func digestRevisionTrendWindow(window RevisionTrendWindow) string {
	return digestString(fmt.Sprintf("%s|%s|%d|%s|%s|%s|%d|%d|%d|%d|%d|%d|%d|%s|%t|%t",
		window.Status,
		window.MissingStage,
		window.ObservationCount,
		strings.Join(window.TrendDigests, ","),
		window.FirstTrendDigest,
		window.LastTrendDigest,
		window.NarrowerCount,
		window.WiderCount,
		window.StableCount,
		window.MixedCount,
		window.CandidateStableCount,
		window.IRChangeStableCount,
		window.SourceStableCount,
		window.WindowSignal,
		window.NonExecuting,
		window.NonAuthorizing,
	))
}
