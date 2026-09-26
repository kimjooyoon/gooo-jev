package gooo

import (
	"fmt"
	"strings"
)

type RevisionCandidateSelectionObservation struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	ObservationCount      int
	ChainDigests          []string
	CandidateDigests      []string
	UniqueCandidateDigests []string
	SelectedCandidateDigest string
	SelectionStatus       string
	SelectionDigest       string
	NonExecuting          bool
	NonAuthorizing        bool
}

func ObserveRevisionCandidateSelection(chains []RevisionEvidenceChain) (RevisionCandidateSelectionObservation, error) {
	observation := RevisionCandidateSelectionObservation{
		Status:           "UNKNOWN",
		MissingStage:     "revision-candidate-selection-chains",
		ObservationCount: len(chains),
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if len(chains) == 0 {
		return observation, fmt.Errorf("gooo revision candidate selection: at least one chain is required")
	}
	observation.ChainDigests = make([]string, 0, len(chains))
	observation.CandidateDigests = make([]string, 0, len(chains))
	observation.UniqueCandidateDigests = make([]string, 0, len(chains))
	seenCandidates := map[string]struct{}{}
	for index, chain := range chains {
		if err := chain.Validate(); err != nil {
			observation.MissingStage = fmt.Sprintf("revision-candidate-selection-chain-%d", index)
			return observation, fmt.Errorf("gooo revision candidate selection: chain %d: %w", index, err)
		}
		if index == 0 {
			observation.SourceDigest = chain.SourceDigest
		} else if chain.SourceDigest != observation.SourceDigest {
			observation.MissingStage = "revision-candidate-selection-source"
			return observation, fmt.Errorf("gooo revision candidate selection: source digests do not match")
		}
		observation.ChainDigests = append(observation.ChainDigests, chain.ChainDigest)
		observation.CandidateDigests = append(observation.CandidateDigests, chain.CandidateDigest)
		if _, exists := seenCandidates[chain.CandidateDigest]; !exists {
			seenCandidates[chain.CandidateDigest] = struct{}{}
			observation.UniqueCandidateDigests = append(observation.UniqueCandidateDigests, chain.CandidateDigest)
		}
	}
	if len(observation.UniqueCandidateDigests) == 1 {
		observation.SelectionStatus = "single"
		observation.SelectedCandidateDigest = observation.UniqueCandidateDigests[0]
	} else {
		observation.SelectionStatus = "ambiguous"
	}
	observation.Status = "BOUND"
	observation.MissingStage = ""
	observation.SelectionDigest = digestRevisionCandidateSelection(observation)
	return observation, nil
}

func (o RevisionCandidateSelectionObservation) Validate() error {
	if o.Status != "BOUND" {
		return fmt.Errorf("revision candidate selection status must be BOUND")
	}
	if o.MissingStage != "" {
		return fmt.Errorf("revision candidate selection missing stage must be empty")
	}
	if o.ObservationCount < 1 ||
		len(o.ChainDigests) != o.ObservationCount ||
		len(o.CandidateDigests) != o.ObservationCount {
		return fmt.Errorf("revision candidate selection observation count is invalid")
	}
	if len(o.UniqueCandidateDigests) < 1 || len(o.UniqueCandidateDigests) > o.ObservationCount {
		return fmt.Errorf("revision candidate selection unique candidate count is invalid")
	}
	if !validDigest(o.SourceDigest) || !validDigest(o.SelectionDigest) {
		return fmt.Errorf("revision candidate selection digest is invalid")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision candidate selection must remain non-executing and non-authorizing")
	}
	for index, digest := range o.ChainDigests {
		if !validDigest(digest) || !validDigest(o.CandidateDigests[index]) {
			return fmt.Errorf("revision candidate selection item %d digest is invalid", index)
		}
	}
	for _, digest := range o.UniqueCandidateDigests {
		if !validDigest(digest) {
			return fmt.Errorf("revision candidate selection unique candidate digest is invalid")
		}
	}
	expectedUnique := uniqueRevisionCandidateDigests(o.CandidateDigests)
	if len(expectedUnique) != len(o.UniqueCandidateDigests) {
		return fmt.Errorf("revision candidate selection unique candidates are not linked")
	}
	for index, digest := range expectedUnique {
		if digest != o.UniqueCandidateDigests[index] {
			return fmt.Errorf("revision candidate selection unique candidate order is not linked")
		}
	}
	switch o.SelectionStatus {
	case "single":
		if len(o.UniqueCandidateDigests) != 1 ||
			o.SelectedCandidateDigest != o.UniqueCandidateDigests[0] {
			return fmt.Errorf("single candidate selection is incomplete")
		}
	case "ambiguous":
		if len(o.UniqueCandidateDigests) < 2 || o.SelectedCandidateDigest != "" {
			return fmt.Errorf("ambiguous candidate selection is incomplete")
		}
	default:
		return fmt.Errorf("revision candidate selection status is invalid")
	}
	if digestRevisionCandidateSelection(o) != o.SelectionDigest {
		return fmt.Errorf("revision candidate selection digest does not match its fields")
	}
	return nil
}

func uniqueRevisionCandidateDigests(digests []string) []string {
	unique := make([]string, 0, len(digests))
	seen := map[string]struct{}{}
	for _, digest := range digests {
		if _, exists := seen[digest]; exists {
			continue
		}
		seen[digest] = struct{}{}
		unique = append(unique, digest)
	}
	return unique
}

func digestRevisionCandidateSelection(observation RevisionCandidateSelectionObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%d|%s|%s|%s|%s|%s|%s|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.SourceDigest,
		observation.ObservationCount,
		strings.Join(observation.ChainDigests, ","),
		strings.Join(observation.CandidateDigests, ","),
		strings.Join(observation.UniqueCandidateDigests, ","),
		observation.SelectedCandidateDigest,
		observation.SelectionStatus,
		"",
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}
