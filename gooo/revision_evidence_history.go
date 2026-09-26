package gooo

import (
	"fmt"
	"strings"
)

type RevisionEvidenceHistory struct {
	Status                    string
	MissingStage              string
	ObservationCount          int
	ChainDigests              []string
	SourceDigests             []string
	ProposedSourceDigests     []string
	InputIRDigests            []string
	ProposedIRDigests         []string
	CandidateDigests          []string
	FirstChainDigest          string
	LastChainDigest           string
	SourceChangeCount         int
	ProposedSourceChangeCount int
	InputIRChangeCount        int
	ProposedIRChangeCount     int
	CandidateChangeCount      int
	ExactSourceMatchCount     int
	StructureMatchCount       int
	ReverseObservedCount      int
	HistoryDigest             string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionEvidenceHistory(chains []RevisionEvidenceChain) (RevisionEvidenceHistory, error) {
	history := RevisionEvidenceHistory{
		Status:           "UNKNOWN",
		MissingStage:     "revision-evidence-history-chains",
		ObservationCount: len(chains),
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if len(chains) == 0 {
		return history, fmt.Errorf("gooo revision evidence history: at least one chain is required")
	}
	history.ChainDigests = make([]string, 0, len(chains))
	history.SourceDigests = make([]string, 0, len(chains))
	history.ProposedSourceDigests = make([]string, 0, len(chains))
	history.InputIRDigests = make([]string, 0, len(chains))
	history.ProposedIRDigests = make([]string, 0, len(chains))
	history.CandidateDigests = make([]string, 0, len(chains))
	for index, chain := range chains {
		if err := chain.Validate(); err != nil {
			history.MissingStage = fmt.Sprintf("revision-evidence-history-chain-%d", index)
			return history, fmt.Errorf("gooo revision evidence history: chain %d: %w", index, err)
		}
		history.ChainDigests = append(history.ChainDigests, chain.ChainDigest)
		history.SourceDigests = append(history.SourceDigests, chain.SourceDigest)
		history.ProposedSourceDigests = append(history.ProposedSourceDigests, chain.ProposedSourceDigest)
		history.InputIRDigests = append(history.InputIRDigests, chain.InputIRDigest)
		history.ProposedIRDigests = append(history.ProposedIRDigests, chain.ProposedIRDigest)
		history.CandidateDigests = append(history.CandidateDigests, chain.CandidateDigest)
		if chain.ExactSourceMatch {
			history.ExactSourceMatchCount++
		}
		if chain.StructureMatch {
			history.StructureMatchCount++
		}
		if chain.ReverseObserved {
			history.ReverseObservedCount++
		}
		if index == 0 {
			continue
		}
		if chain.SourceDigest != chains[index-1].SourceDigest {
			history.SourceChangeCount++
		}
		if chain.ProposedSourceDigest != chains[index-1].ProposedSourceDigest {
			history.ProposedSourceChangeCount++
		}
		if chain.InputIRDigest != chains[index-1].InputIRDigest {
			history.InputIRChangeCount++
		}
		if chain.ProposedIRDigest != chains[index-1].ProposedIRDigest {
			history.ProposedIRChangeCount++
		}
		if chain.CandidateDigest != chains[index-1].CandidateDigest {
			history.CandidateChangeCount++
		}
	}
	history.FirstChainDigest = history.ChainDigests[0]
	history.LastChainDigest = history.ChainDigests[len(history.ChainDigests)-1]
	history.Status = "BOUND"
	history.MissingStage = ""
	history.HistoryDigest = digestRevisionEvidenceHistory(history)
	return history, nil
}

func (h RevisionEvidenceHistory) Validate() error {
	if h.Status != "BOUND" {
		return fmt.Errorf("revision evidence history status must be BOUND")
	}
	if h.MissingStage != "" {
		return fmt.Errorf("revision evidence history missing stage must be empty")
	}
	if h.ObservationCount < 1 {
		return fmt.Errorf("revision evidence history observation count must be positive")
	}
	if len(h.ChainDigests) != h.ObservationCount ||
		len(h.SourceDigests) != h.ObservationCount ||
		len(h.ProposedSourceDigests) != h.ObservationCount ||
		len(h.InputIRDigests) != h.ObservationCount ||
		len(h.ProposedIRDigests) != h.ObservationCount ||
		len(h.CandidateDigests) != h.ObservationCount {
		return fmt.Errorf("revision evidence history arrays do not match observation count")
	}
	if !h.NonExecuting || !h.NonAuthorizing {
		return fmt.Errorf("revision evidence history must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"first chain": h.FirstChainDigest, "last chain": h.LastChainDigest,
		"history": h.HistoryDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision evidence history %s digest is invalid", name)
		}
	}
	for index, digest := range h.ChainDigests {
		if !validDigest(digest) || !validDigest(h.SourceDigests[index]) ||
			!validDigest(h.ProposedSourceDigests[index]) ||
			!validDigest(h.InputIRDigests[index]) ||
			!validDigest(h.ProposedIRDigests[index]) ||
			!validDigest(h.CandidateDigests[index]) {
			return fmt.Errorf("revision evidence history item %d digest is invalid", index)
		}
	}
	if h.ChainDigests[0] != h.FirstChainDigest ||
		h.ChainDigests[len(h.ChainDigests)-1] != h.LastChainDigest {
		return fmt.Errorf("revision evidence history boundary is not linked")
	}
	maxChanges := h.ObservationCount - 1
	if h.SourceChangeCount < 0 || h.SourceChangeCount > maxChanges ||
		h.ProposedSourceChangeCount < 0 || h.ProposedSourceChangeCount > maxChanges ||
		h.InputIRChangeCount < 0 || h.InputIRChangeCount > maxChanges ||
		h.ProposedIRChangeCount < 0 || h.ProposedIRChangeCount > maxChanges ||
		h.CandidateChangeCount < 0 || h.CandidateChangeCount > maxChanges {
		return fmt.Errorf("revision evidence history change counts are invalid")
	}
	if h.ExactSourceMatchCount < 0 || h.ExactSourceMatchCount > h.ObservationCount ||
		h.StructureMatchCount < 0 || h.StructureMatchCount > h.ObservationCount ||
		h.ReverseObservedCount < 0 || h.ReverseObservedCount > h.ObservationCount {
		return fmt.Errorf("revision evidence history observation counts are invalid")
	}
	if digestRevisionEvidenceHistory(h) != h.HistoryDigest {
		return fmt.Errorf("revision evidence history digest does not match its fields")
	}
	return nil
}

func digestRevisionEvidenceHistory(history RevisionEvidenceHistory) string {
	return digestString(fmt.Sprintf("%s|%s|%d|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%d|%d|%d|%d|%d|%d|%s|%t|%t",
		history.Status,
		history.MissingStage,
		history.ObservationCount,
		strings.Join(history.ChainDigests, ","),
		strings.Join(history.SourceDigests, ","),
		strings.Join(history.ProposedSourceDigests, ","),
		strings.Join(history.InputIRDigests, ","),
		strings.Join(history.ProposedIRDigests, ","),
		strings.Join(history.CandidateDigests, ","),
		history.FirstChainDigest,
		history.LastChainDigest,
		history.SourceChangeCount,
		history.ProposedSourceChangeCount,
		history.InputIRChangeCount,
		history.ProposedIRChangeCount,
		history.CandidateChangeCount,
		history.ExactSourceMatchCount,
		history.StructureMatchCount,
		history.ReverseObservedCount,
		"",
		history.NonExecuting,
		history.NonAuthorizing,
	))
}
