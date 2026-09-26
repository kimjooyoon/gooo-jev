package decision

import (
	"fmt"
	"strings"
	"time"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardInput
// carries an extended lineage candidate into the typed action guard.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardInput struct {
	Candidate             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding
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

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding
// preserves extended lineage and guard evidence without authorizing execution.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	GuardEvidenceDigest             string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding) Validate() error {
	if b.Status != "admitted" && b.Status != "review" && b.Status != "rejected" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision guard binding")
	}
	if b.MissingStage == "" && b.GuardEvidenceDigest == "" {
		return fmt.Errorf("extended lineage candidate revision guard evidence is missing")
	}
	if b.MissingStage == "" && b.Status != "admitted" && b.Status != "review" && b.Status != "rejected" {
		return fmt.Errorf("extended lineage candidate revision guard status is invalid")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision guard evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision guard must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuard(
		b.Status,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.GuardEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision guard digest mismatch")
	}
	return nil
}

// GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage
// reuses the existing typed guard without executing or authorizing a candidate.
func GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-revision-generation-extended-lineage-guard"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding{
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
		return unknown("revision-action-candidate-generation-extended-lineage-validation")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding{
		Status:                          guard.Status,
		MissingStage:                    guard.MissingStage,
		CandidateID:                     guard.CandidateID,
		SourceCandidateID:               input.Candidate.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Candidate.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Candidate.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Candidate.CandidateDigest,
		RevisionCandidateEvidenceDigest: input.Candidate.CandidateEvidenceDigest,
		GuardEvidenceDigest:             guard.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuard(
		output.Status,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.GuardEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-revision-generation-extended-lineage-guard-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuard(
	status,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
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
		GuardEvidenceDigest             string
	}{
		Status:                          status,
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}
