package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput
// joins an extended generated candidate guard with reverse observation.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput struct {
	Guard              ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedGuardBinding
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing     bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding
// preserves extended generated candidate, guard, and reverse verification evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding struct {
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
	GuardEvidenceDigest             string
	VerificationStatus              string
	ReverseStatus                   string
	ReverseEvidenceDigest           string
	FirstMismatch                   string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding) Validate() error {
	if b.Status != "verified" && b.Status != "review" {
		return fmt.Errorf("incomplete Gooo extended generated candidate verification binding")
	}
	if b.Status == "verified" && b.MissingStage != "" {
		return fmt.Errorf("verified extended generated candidate must not have a missing stage")
	}
	if b.Status == "review" && b.MissingStage == "" {
		return fmt.Errorf("review extended generated candidate must retain a missing stage")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationStatus != b.Status ||
		b.ReverseStatus == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended generated candidate verification evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended generated candidate verification must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerification(
		b.Status,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.EvidencePrefixDigest,
		b.GuardEvidenceDigest,
		b.VerificationStatus,
		b.ReverseStatus,
		b.ReverseEvidenceDigest,
		b.FirstMismatch,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended generated candidate verification digest mismatch")
	}
	return nil
}

// VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended
// requires an admitted extended guard and delegates classification to reverse observation.
func VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-generation-extended-verification"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding{
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
		return unknown("action-candidate-generation-extended-guard-validation")
	}
	if input.Guard.GuardStatus != "admitted" {
		stage := input.Guard.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "candidate-generation-extended-guard"
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
			stage = "candidate-generation-extended-verification"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding{
		Status:                          verification.Status,
		MissingStage:                    verification.MissingStage,
		CandidateID:                     verification.CandidateID,
		SourceCandidateID:               input.Guard.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Guard.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Guard.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Guard.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Guard.RevisionCandidateEvidenceDigest,
		VerificationEvidenceDigest:      input.Guard.VerificationEvidenceDigest,
		EvidencePrefixDigest:            input.Guard.EvidencePrefixDigest,
		GuardEvidenceDigest:             input.Guard.GuardEvidenceDigest,
		VerificationStatus:              verification.Status,
		ReverseStatus:                   verification.ReverseStatus,
		ReverseEvidenceDigest:           verification.ReverseEvidenceDigest,
		FirstMismatch:                   verification.FirstMismatch,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerification(
		output.Status,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.EvidencePrefixDigest,
		output.GuardEvidenceDigest,
		output.VerificationStatus,
		output.ReverseStatus,
		output.ReverseEvidenceDigest,
		output.FirstMismatch,
	)
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-generation-extended-verification-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerification(
	status,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	verificationEvidenceDigest,
	evidencePrefixDigest,
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
		VerificationEvidenceDigest      string
		EvidencePrefixDigest            string
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
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		EvidencePrefixDigest:            evidencePrefixDigest,
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
