package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput
// joins an extended lineage candidate guard with reverse observation.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput struct {
	Guard              ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageGuardBinding
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing     bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding
// preserves extended lineage candidate, guard, and reverse verification evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	GuardEvidenceDigest             string
	VerificationStatus              string
	ReverseStatus                   string
	ReverseEvidenceDigest           string
	FirstMismatch                   string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding) Validate() error {
	if b.Status != "verified" && b.Status != "review" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate verification binding")
	}
	if b.Status == "verified" && b.MissingStage != "" {
		return fmt.Errorf("verified extended lineage candidate must not have a missing stage")
	}
	if b.Status == "review" && b.MissingStage == "" {
		return fmt.Errorf("review extended lineage candidate must retain a missing stage")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationStatus != b.Status ||
		b.ReverseStatus == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate verification evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate verification must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerification(
		b.Status,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.GuardEvidenceDigest,
		b.VerificationStatus,
		b.ReverseStatus,
		b.ReverseEvidenceDigest,
		b.FirstMismatch,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate verification digest mismatch")
	}
	return nil
}

// VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage
// requires an admitted lineage guard and delegates classification to reverse observation.
func VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-generation-extended-lineage-verification"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Guard.NonAuthorizing || !input.ReverseObservation.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Guard.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Guard.Validate(); err != nil {
		return unknown("action-candidate-generation-extended-lineage-guard-validation")
	}
	if input.Guard.Status != "admitted" {
		stage := input.Guard.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "candidate-generation-extended-lineage-guard"
		}
		return unknown(stage)
	}
	verification := VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard: JEVActionCandidateGuard{
			Status:            input.Guard.Status,
			CandidateID:       input.Guard.CandidateID,
			EvidenceDigest:    input.Guard.GuardEvidenceDigest,
			NonExecuting:      true,
			NonAuthorizing:    true,
		},
		ReverseObservation: input.ReverseObservation,
		NonAuthorizing:     true,
	})
	if verification.EvidenceDigest == "" {
		stage := verification.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "candidate-generation-extended-lineage-verification"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerificationBinding{
		Status:                          verification.Status,
		MissingStage:                    verification.MissingStage,
		CandidateID:                     verification.CandidateID,
		SourceCandidateID:               input.Guard.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Guard.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Guard.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Guard.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Guard.RevisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             input.Guard.GuardEvidenceDigest,
		VerificationStatus:              verification.Status,
		ReverseStatus:                   verification.ReverseStatus,
		ReverseEvidenceDigest:           verification.ReverseEvidenceDigest,
		FirstMismatch:                   verification.FirstMismatch,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerification(
		output.Status,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.GuardEvidenceDigest,
		output.VerificationStatus,
		output.ReverseStatus,
		output.ReverseEvidenceDigest,
		output.FirstMismatch,
	)
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-generation-extended-lineage-verification-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageVerification(
	status,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	guardEvidenceDigest,
	verificationStatus,
	reverseStatus,
	reverseEvidenceDigest,
	firstMismatch string,
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
		VerificationStatus              string
		ReverseStatus                   string
		ReverseEvidenceDigest           string
		FirstMismatch                   string
	}{
		Status:                          status,
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
		VerificationStatus:              verificationStatus,
		ReverseStatus:                   reverseStatus,
		ReverseEvidenceDigest:           reverseEvidenceDigest,
		FirstMismatch:                   firstMismatch,
	})
	if err != nil {
		return ""
	}
	return digest
}
