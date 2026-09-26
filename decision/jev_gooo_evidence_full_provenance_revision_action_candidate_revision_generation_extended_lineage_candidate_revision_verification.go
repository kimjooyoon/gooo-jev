package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput
// joins a revision candidate guard with reverse observation.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput struct {
	Guard              ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionGuardBinding
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing     bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding
// preserves revision candidate, guard, and reverse verification evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	ParentCandidateDigest           string
	ParentCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	CandidateStatus                 string
	CandidateDigest                 string
	CandidateEvidenceDigest         string
	RevisionSource                  string
	BoundRevisionChangeDigest       string
	GuardEvidenceDigest             string
	VerificationStatus              string
	ReverseStatus                   string
	ReverseEvidenceDigest           string
	FirstMismatch                   string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding) Validate() error {
	if b.Status != "verified" && b.Status != "review" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision verification binding")
	}
	if b.Status == "verified" && b.MissingStage != "" {
		return fmt.Errorf("verified extended lineage candidate revision must not have a missing stage")
	}
	if b.Status == "review" && b.MissingStage == "" {
		return fmt.Errorf("review extended lineage candidate revision must retain a missing stage")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.ParentCandidateDigest == "" ||
		b.ParentCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.CandidateStatus == "" ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationStatus != b.Status ||
		b.ReverseStatus == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision verification evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision verification must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerification(
		b.Status,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.ParentCandidateDigest,
		b.ParentCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.CandidateStatus,
		b.CandidateDigest,
		b.CandidateEvidenceDigest,
		b.RevisionSource,
		b.BoundRevisionChangeDigest,
		b.GuardEvidenceDigest,
		b.VerificationStatus,
		b.ReverseStatus,
		b.ReverseEvidenceDigest,
		b.FirstMismatch,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision verification digest mismatch")
	}
	return nil
}

// VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision
// requires an admitted revision guard and delegates classification to reverse observation.
func VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-generation-extended-lineage-candidate-revision-verification"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding{
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
		return unknown("action-candidate-generation-extended-lineage-candidate-revision-guard-validation")
	}
	if input.Guard.Status != "admitted" {
		stage := input.Guard.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "candidate-generation-extended-lineage-candidate-revision-guard"
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
			stage = "candidate-generation-extended-lineage-candidate-revision-verification"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationBinding{
		Status:                          verification.Status,
		MissingStage:                    verification.MissingStage,
		CandidateID:                     verification.CandidateID,
		SourceCandidateID:               input.Guard.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Guard.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Guard.SourceCandidateEvidenceDigest,
		ParentCandidateDigest:           input.Guard.ParentCandidateDigest,
		ParentCandidateEvidenceDigest:   input.Guard.ParentCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Guard.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Guard.RevisionCandidateEvidenceDigest,
		CandidateStatus:                 input.Guard.CandidateStatus,
		CandidateDigest:                 input.Guard.CandidateDigest,
		CandidateEvidenceDigest:         input.Guard.CandidateEvidenceDigest,
		RevisionSource:                  input.Guard.RevisionSource,
		BoundRevisionChangeDigest:       input.Guard.BoundRevisionChangeDigest,
		GuardEvidenceDigest:             input.Guard.GuardEvidenceDigest,
		VerificationStatus:              verification.Status,
		ReverseStatus:                   verification.ReverseStatus,
		ReverseEvidenceDigest:           verification.ReverseEvidenceDigest,
		FirstMismatch:                   verification.FirstMismatch,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerification(
		output.Status,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.ParentCandidateDigest,
		output.ParentCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.CandidateStatus,
		output.CandidateDigest,
		output.CandidateEvidenceDigest,
		output.RevisionSource,
		output.BoundRevisionChangeDigest,
		output.GuardEvidenceDigest,
		output.VerificationStatus,
		output.ReverseStatus,
		output.ReverseEvidenceDigest,
		output.FirstMismatch,
	)
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-generation-extended-lineage-candidate-revision-verification-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerification(
	status,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	parentCandidateDigest,
	parentCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	candidateStatus,
	candidateDigest,
	candidateEvidenceDigest,
	revisionSource,
	boundRevisionChangeDigest,
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
		ParentCandidateDigest           string
		ParentCandidateEvidenceDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		CandidateStatus                 string
		CandidateDigest                 string
		CandidateEvidenceDigest         string
		RevisionSource                  string
		BoundRevisionChangeDigest       string
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
		ParentCandidateDigest:           parentCandidateDigest,
		ParentCandidateEvidenceDigest:   parentCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		CandidateStatus:                 candidateStatus,
		CandidateDigest:                 candidateDigest,
		CandidateEvidenceDigest:         candidateEvidenceDigest,
		RevisionSource:                  revisionSource,
		BoundRevisionChangeDigest:       boundRevisionChangeDigest,
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