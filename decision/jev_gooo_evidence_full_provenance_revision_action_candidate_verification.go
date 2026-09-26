package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationInput
// joins an action-derived candidate guard with reverse observation.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationInput struct {
	Guard            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateGuardBinding
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing   bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding
// preserves candidate, guard, and reverse verification evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	RevisionCandidateDigest         string
	CandidateEvidenceDigest         string
	GuardEvidenceDigest             string
	VerificationStatus              string
	ReverseStatus                   string
	ReverseEvidenceDigest           string
	FirstMismatch                   string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding) Validate() error {
	if b.Status != "verified" && b.Status != "review" {
		return fmt.Errorf("incomplete Gooo action-derived candidate verification binding")
	}
	if b.Status == "verified" && b.MissingStage != "" {
		return fmt.Errorf("verified action-derived candidate must not have a missing stage")
	}
	if b.Status == "review" && b.MissingStage == "" {
		return fmt.Errorf("review action-derived candidate must retain a missing stage")
	}
	if b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationStatus != b.Status ||
		b.ReverseStatus == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate verification evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate verification must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		Status                  string
		CandidateID             string
		RevisionCandidateDigest string
		CandidateEvidenceDigest string
		GuardEvidenceDigest     string
		VerificationStatus      string
		ReverseStatus           string
		ReverseEvidenceDigest   string
		FirstMismatch           string
	}{
		Status:                  b.Status,
		CandidateID:             b.CandidateID,
		RevisionCandidateDigest: b.RevisionCandidateDigest,
		CandidateEvidenceDigest: b.CandidateEvidenceDigest,
		GuardEvidenceDigest:     b.GuardEvidenceDigest,
		VerificationStatus:      b.VerificationStatus,
		ReverseStatus:           b.ReverseStatus,
		ReverseEvidenceDigest:   b.ReverseEvidenceDigest,
		FirstMismatch:            b.FirstMismatch,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate verification digest mismatch")
	}
	return nil
}

// VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate
// requires an admitted guard and delegates classification to the existing
// reverse-observation verifier.
func VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "action-candidate-verification"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding{
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
		return unknown("action-candidate-guard-validation")
	}
	if input.Guard.GuardStatus != "admitted" {
		stage := input.Guard.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "candidate-guard"
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
			stage = "candidate-verification"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding{
		Status:                  verification.Status,
		MissingStage:            verification.MissingStage,
		CandidateID:             verification.CandidateID,
		RevisionCandidateDigest: input.Guard.RevisionCandidateDigest,
		CandidateEvidenceDigest: input.Guard.CandidateEvidenceDigest,
		GuardEvidenceDigest:     input.Guard.GuardEvidenceDigest,
		VerificationStatus:      verification.Status,
		ReverseStatus:           verification.ReverseStatus,
		ReverseEvidenceDigest:   verification.ReverseEvidenceDigest,
		FirstMismatch:            verification.FirstMismatch,
		NonExecuting:            true,
		NonAuthorizing:          true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		Status                  string
		CandidateID             string
		RevisionCandidateDigest string
		CandidateEvidenceDigest string
		GuardEvidenceDigest     string
		VerificationStatus      string
		ReverseStatus           string
		ReverseEvidenceDigest   string
		FirstMismatch            string
	}{
		Status:                  output.Status,
		CandidateID:             output.CandidateID,
		RevisionCandidateDigest: output.RevisionCandidateDigest,
		CandidateEvidenceDigest: output.CandidateEvidenceDigest,
		GuardEvidenceDigest:     output.GuardEvidenceDigest,
		VerificationStatus:      output.VerificationStatus,
		ReverseStatus:           output.ReverseStatus,
		ReverseEvidenceDigest:   output.ReverseEvidenceDigest,
		FirstMismatch:            output.FirstMismatch,
	})
	if err := output.Validate(); err != nil {
		return unknown("action-candidate-verification-evidence")
	}
	return output
}