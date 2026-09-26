package gooo

import "fmt"

type RevisionCandidateBinding struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	SelectionDigest       string
	MaterializationDigest string
	CandidateDigest       string
	SelectionStatus       string
	Candidate             RevisionCandidate
	BindingDigest         string
	NonExecuting          bool
	NonAuthorizing        bool
}

func BindRevisionCandidateSelection(selection RevisionCandidateSelectionObservation, materialization RevisionCandidateMaterialization) (RevisionCandidateBinding, error) {
	binding := RevisionCandidateBinding{
		Status:                "UNKNOWN",
		MissingStage:          "revision-candidate-binding",
		SourceDigest:          selection.SourceDigest,
		SelectionDigest:       selection.SelectionDigest,
		MaterializationDigest: materialization.MaterializationDigest,
		CandidateDigest:       materialization.CandidateDigest,
		SelectionStatus:       selection.SelectionStatus,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	if err := selection.Validate(); err != nil {
		binding.MissingStage = "revision-candidate-binding-selection"
		return binding, fmt.Errorf("gooo revision candidate binding: selection: %w", err)
	}
	if selection.SelectionStatus != "single" {
		binding.MissingStage = "revision-candidate-binding-selection"
		return binding, fmt.Errorf("gooo revision candidate binding: selection is not single")
	}
	if err := materialization.Validate(); err != nil {
		binding.MissingStage = "revision-candidate-binding-materialization"
		return binding, fmt.Errorf("gooo revision candidate binding: materialization: %w", err)
	}
	if selection.SourceDigest != materialization.SourceDigest {
		binding.MissingStage = "revision-candidate-binding-source"
		return binding, fmt.Errorf("gooo revision candidate binding: source digests do not match")
	}
	if selection.SelectedCandidateDigest != materialization.CandidateDigest {
		binding.MissingStage = "revision-candidate-binding-candidate"
		return binding, fmt.Errorf("gooo revision candidate binding: selected candidate is not materialized")
	}
	binding.Status = "BOUND"
	binding.MissingStage = ""
	binding.SourceDigest = materialization.SourceDigest
	binding.Candidate = materialization.Candidate
	binding.BindingDigest = digestRevisionCandidateBinding(binding)
	return binding, nil
}

func (b RevisionCandidateBinding) Validate() error {
	if b.Status != "BOUND" {
		return fmt.Errorf("revision candidate binding status must be BOUND")
	}
	if b.MissingStage != "" {
		return fmt.Errorf("revision candidate binding missing stage must be empty")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("revision candidate binding must remain non-executing and non-authorizing")
	}
	if b.SelectionStatus != "single" {
		return fmt.Errorf("revision candidate binding selection status must be single")
	}
	for name, digest := range map[string]string{
		"source": b.SourceDigest, "selection": b.SelectionDigest,
		"materialization": b.MaterializationDigest, "candidate": b.CandidateDigest,
		"binding": b.BindingDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision candidate binding %s digest is invalid", name)
		}
	}
	if err := b.Candidate.Validate(); err != nil {
		return fmt.Errorf("revision candidate binding candidate: %w", err)
	}
	if b.Candidate.CandidateDigest != b.CandidateDigest {
		return fmt.Errorf("revision candidate binding candidate digest is not linked")
	}
	if digestRevisionCandidateBinding(b) != b.BindingDigest {
		return fmt.Errorf("revision candidate binding digest does not match its fields")
	}
	return nil
}

func digestRevisionCandidateBinding(binding RevisionCandidateBinding) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%t|%t",
		binding.Status,
		binding.MissingStage,
		binding.SourceDigest,
		binding.SelectionDigest,
		binding.MaterializationDigest,
		binding.CandidateDigest,
		binding.SelectionStatus,
		binding.NonExecuting,
		binding.NonAuthorizing,
	))
}
