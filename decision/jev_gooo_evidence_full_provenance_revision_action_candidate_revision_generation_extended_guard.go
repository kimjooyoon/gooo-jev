package decision

import (
	"fmt"
	"strings"
	"time"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardInput
// carries an extended generated candidate into the typed action guard.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardInput struct {
	Candidate             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedBinding
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

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding
// preserves generated candidate lineage and guard evidence without authorizing execution.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	VerificationEvidenceDigest      string
	EvidencePrefixDigest            string
	GuardStatus                     string
	GuardEvidenceDigest             string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding) Validate() error {
	if b.Status != "admitted" && b.Status != "review" && b.Status != "rejected" {
		return fmt.Errorf("incomplete Gooo extended generated candidate revision guard binding")
	}
	if b.MissingStage == "" && b.GuardStatus != b.Status {
		return fmt.Errorf("extended generated candidate revision guard status mismatch")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.GuardStatus == "" ||
		b.GuardEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended generated candidate revision guard evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended generated candidate revision guard must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuard(
		b.Status,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.EvidencePrefixDigest,
		b.GuardStatus,
		b.GuardEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended generated candidate revision guard digest mismatch")
	}
	return nil
}

// GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended
// reuses the existing typed guard without executing or authorizing a candidate.
func GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-revision-generation-extended-guard"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding{
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
		return unknown("revision-action-candidate-generation-extended-validation")
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
		WriteAllowed:          input.WriteAllowed,
		ConfirmationRequired: input.ConfirmationRequired,
		ConfirmationPresent:  input.ConfirmationPresent,
		NonAuthorizing:       true,
	})
	if guard.EvidenceDigest == "" {
		return unknown(guard.MissingStage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding{
		Status:                          guard.Status,
		MissingStage:                    guard.MissingStage,
		CandidateID:                     guard.CandidateID,
		SourceCandidateID:               input.Candidate.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Candidate.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Candidate.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Candidate.CandidateDigest,
		RevisionCandidateEvidenceDigest: input.Candidate.CandidateEvidenceDigest,
		VerificationEvidenceDigest:      input.Candidate.VerificationEvidenceDigest,
		EvidencePrefixDigest:            input.Candidate.EvidencePrefixDigest,
		GuardStatus:                     guard.Status,
		GuardEvidenceDigest:             guard.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuard(
		output.Status,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.EvidencePrefixDigest,
		output.GuardStatus,
		output.GuardEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-revision-generation-extended-guard-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuard(
	status,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	verificationEvidenceDigest,
	evidencePrefixDigest,
	guardStatus,
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
		VerificationEvidenceDigest      string
		EvidencePrefixDigest            string
		GuardStatus                     string
		GuardEvidenceDigest             string
	}{
		Status:                          status,
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		EvidencePrefixDigest:            evidencePrefixDigest,
		GuardStatus:                     guardStatus,
		GuardEvidenceDigest:             guardEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}
