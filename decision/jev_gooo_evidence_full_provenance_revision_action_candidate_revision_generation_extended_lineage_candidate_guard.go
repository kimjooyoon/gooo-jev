package decision

import (
	"fmt"
	"strings"
	"time"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardInput
// carries a generated extended lineage candidate into the typed action guard.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardInput struct {
	Candidate            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateBinding
	CandidateID          string
	ObservationDigest    string
	ObservationAt        time.Time
	Now                  time.Time
	MaxAge               time.Duration
	RequestedPath        string
	AllowedPathPrefix    string
	WriteRequested       bool
	WriteAllowed         bool
	ConfirmationRequired bool
	ConfirmationPresent  bool
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding
// preserves generated candidate lineage and guard evidence without authorizing execution.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	CandidateStatus                 string
	CandidateDigest                 string
	CandidateEvidenceDigest         string
	RevisionSource                  string
	BoundRevisionChangeDigest       string
	GuardEvidenceDigest             string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding) Validate() error {
	if b.Status != "admitted" && b.Status != "review" && b.Status != "rejected" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate guard binding")
	}
	if b.MissingStage == "" && b.GuardEvidenceDigest == "" {
		return fmt.Errorf("extended lineage candidate guard evidence is missing")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.CandidateStatus == "" ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate guard evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate guard must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuard(
		b.Status,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.CandidateStatus,
		b.CandidateDigest,
		b.CandidateEvidenceDigest,
		b.RevisionSource,
		b.BoundRevisionChangeDigest,
		b.GuardEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate guard digest mismatch")
	}
	return nil
}

// GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate
// reuses the typed action guard without executing or authorizing a candidate.
func GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-revision-generation-extended-lineage-candidate-guard"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding{
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
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-validation")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuardBinding{
		Status:                          guard.Status,
		MissingStage:                    guard.MissingStage,
		CandidateID:                     guard.CandidateID,
		SourceCandidateID:               input.Candidate.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Candidate.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Candidate.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Candidate.CandidateDigest,
		RevisionCandidateEvidenceDigest: input.Candidate.CandidateEvidenceDigest,
		CandidateStatus:                 input.Candidate.CandidateStatus,
		CandidateDigest:                 input.Candidate.CandidateDigest,
		CandidateEvidenceDigest:         input.Candidate.CandidateEvidenceDigest,
		RevisionSource:                  input.Candidate.RevisionSource,
		BoundRevisionChangeDigest:       input.Candidate.BoundRevisionChangeDigest,
		GuardEvidenceDigest:             guard.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuard(
		output.Status,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.CandidateStatus,
		output.CandidateDigest,
		output.CandidateEvidenceDigest,
		output.RevisionSource,
		output.BoundRevisionChangeDigest,
		output.GuardEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-revision-generation-extended-lineage-candidate-guard-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateGuard(
	status,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	candidateStatus,
	candidateDigest,
	candidateEvidenceDigest,
	revisionSource,
	boundRevisionChangeDigest,
	guardEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		Status                          string
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		SourceCandidateEvidenceDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		CandidateStatus                 string
		CandidateDigest                 string
		CandidateEvidenceDigest         string
		RevisionSource                  string
		BoundRevisionChangeDigest       string
		GuardEvidenceDigest             string
	}{
		Status:                          status,
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		CandidateStatus:                 candidateStatus,
		CandidateDigest:                 candidateDigest,
		CandidateEvidenceDigest:         candidateEvidenceDigest,
		RevisionSource:                  revisionSource,
		BoundRevisionChangeDigest:       boundRevisionChangeDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}