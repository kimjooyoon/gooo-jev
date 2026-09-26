package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput
// joins a generated candidate guard with reverse observation.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput struct {
	Guard              ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationGuardBinding
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing     bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding
// preserves generated candidate, guard, and reverse verification evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding struct {
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

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding) Validate() error {
	if b.Status != "verified" && b.Status != "review" {
		return fmt.Errorf("incomplete Gooo generated candidate verification binding")
	}
	if b.Status == "verified" && b.MissingStage != "" {
		return fmt.Errorf("verified generated candidate must not have a missing stage")
	}
	if b.Status == "review" && b.MissingStage == "" {
		return fmt.Errorf("review generated candidate must retain a missing stage")
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
		return fmt.Errorf("incomplete Gooo generated candidate verification evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo generated candidate verification must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerification(
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
		return fmt.Errorf("Gooo generated candidate verification digest mismatch")
	}
	return nil
}

// VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration
// requires an admitted guard and delegates classification to reverse observation.
func VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGeneration(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-generation-verification"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding{
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
		return unknown("action-candidate-generation-guard-validation")
	}
	if input.Guard.GuardStatus != "admitted" {
		stage := input.Guard.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "candidate-generation-guard"
		}
		return unknown(stage)
	}
	verification := VerifyJEVActionCandidate(JEVActionCandidateVerificationInput{
		Guard: JEVActionCandidateGuard{
			Status:            input.Guard.GuardStatus,
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
			stage = "candidate-generation-verification"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding{
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
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerification(
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
		return unknown("action-candidate-generation-verification-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerification(
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