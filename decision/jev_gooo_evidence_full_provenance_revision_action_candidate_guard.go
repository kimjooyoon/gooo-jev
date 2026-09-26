package decision

import (
	"fmt"
	"strings"
	"time"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardInput
// carries an action-derived revision candidate into the typed action guard.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardInput struct {
	Candidate             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding
	CandidateID           string
	ObservationDigest     string
	ObservationAt         time.Time
	Now                   time.Time
	MaxAge                time.Duration
	RequestedPath         string
	AllowedPathPrefix     string
	WriteRequested        bool
	WriteAllowed          bool
	ConfirmationRequired  bool
	ConfirmationPresent   bool
	NonAuthorizing        bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding
// keeps action feedback and guard evidence in one non-executing binding.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	RevisionCandidateDigest         string
	CandidateEvidenceDigest         string
	GuardStatus                     string
	GuardEvidenceDigest             string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding) Validate() error {
	if b.Status != "admitted" && b.Status != "review" && b.Status != "rejected" {
		return fmt.Errorf("incomplete Gooo action-derived candidate guard binding")
	}
	if b.MissingStage == "" && b.GuardStatus != b.Status {
		return fmt.Errorf("action-derived candidate guard status mismatch")
	}
	if b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.GuardStatus == "" ||
		b.GuardEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate guard evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate guard must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		Status                  string
		CandidateID             string
		RevisionCandidateDigest string
		CandidateEvidenceDigest string
		GuardStatus             string
		GuardEvidenceDigest     string
	}{
		Status:                  b.Status,
		CandidateID:             b.CandidateID,
		RevisionCandidateDigest: b.RevisionCandidateDigest,
		CandidateEvidenceDigest: b.CandidateEvidenceDigest,
		GuardStatus:             b.GuardStatus,
		GuardEvidenceDigest:     b.GuardEvidenceDigest,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate guard digest mismatch")
	}
	return nil
}

// GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate
// reuses the existing typed guard without executing or authorizing a candidate.
func GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-guard"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Candidate.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Candidate.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Candidate.Validate(); err != nil {
		return unknown("revision-action-candidate-validation")
	}
	guard := GuardJEVActionCandidate(JEVActionCandidateGuardInput{
		CandidateID:          input.CandidateID,
		ObservationDigest:    input.ObservationDigest,
		ObservationAt:        input.ObservationAt,
		Now:                  input.Now,
		MaxAge:               input.MaxAge,
		RequestedPath:        input.RequestedPath,
		AllowedPathPrefix:    input.AllowedPathPrefix,
		WriteRequested:       input.WriteRequested,
		WriteAllowed:         input.WriteAllowed,
		ConfirmationRequired: input.ConfirmationRequired,
		ConfirmationPresent:  input.ConfirmationPresent,
		NonAuthorizing:       true,
	})
	if guard.EvidenceDigest == "" {
		return unknown(guard.MissingStage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding{
		Status:                  guard.Status,
		MissingStage:            guard.MissingStage,
		CandidateID:             guard.CandidateID,
		RevisionCandidateDigest: input.Candidate.CandidateDigest,
		CandidateEvidenceDigest: input.Candidate.CandidateEvidenceDigest,
		GuardStatus:             guard.Status,
		GuardEvidenceDigest:     guard.EvidenceDigest,
		NonExecuting:            true,
		NonAuthorizing:          true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		Status                  string
		CandidateID             string
		RevisionCandidateDigest string
		CandidateEvidenceDigest string
		GuardStatus             string
		GuardEvidenceDigest     string
	}{
		Status:                  output.Status,
		CandidateID:             output.CandidateID,
		RevisionCandidateDigest: output.RevisionCandidateDigest,
		CandidateEvidenceDigest: output.CandidateEvidenceDigest,
		GuardStatus:             output.GuardStatus,
		GuardEvidenceDigest:     output.GuardEvidenceDigest,
	})
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-guard-evidence")
	}
	return output
}