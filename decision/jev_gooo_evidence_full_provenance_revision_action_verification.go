package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationInput
// joins a guarded revision action with an independent reverse observation.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationInput struct {
	Guard            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing   bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding
// preserves revision, guard, and reverse-observation evidence without executing.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding struct {
	Status                           string
	MissingStage                     string
	CandidateID                      string
	RevisionCandidateDigest          string
	RevisionCandidateEvidenceDigest  string
	GuardEvidenceDigest              string
	VerificationStatus               string
	ReverseStatus                    string
	ReverseEvidenceDigest            string
	FirstMismatch                    string
	EvidenceDigest                   string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding) Validate() error {
	if b.Status != "verified" && b.Status != "review" {
		return fmt.Errorf("incomplete Gooo revision action verification binding")
	}
	if b.Status == "verified" && b.MissingStage != "" {
		return fmt.Errorf("verified Gooo revision action must not have a missing stage")
	}
	if b.Status == "review" && b.MissingStage == "" {
		return fmt.Errorf("review Gooo revision action must retain a missing stage")
	}
	if b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationStatus != b.Status ||
		b.ReverseStatus == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo revision action verification evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo revision action verification must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionVerification(
		b.Status,
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.GuardEvidenceDigest,
		b.VerificationStatus,
		b.ReverseStatus,
		b.ReverseEvidenceDigest,
		b.FirstMismatch,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo revision action verification digest mismatch")
	}
	return nil
}

// VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction
// requires an admitted guard and delegates outcome classification to the
// existing independent reverse-observation verifier.
func VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-verification"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding{
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
		return unknown("revision-action-guard-validation")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding{
		Status:                          verification.Status,
		MissingStage:                    verification.MissingStage,
		CandidateID:                     verification.CandidateID,
		RevisionCandidateDigest:         input.Guard.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Guard.RevisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             input.Guard.GuardEvidenceDigest,
		VerificationStatus:              verification.Status,
		ReverseStatus:                   verification.ReverseStatus,
		ReverseEvidenceDigest:           verification.ReverseEvidenceDigest,
		FirstMismatch:                   verification.FirstMismatch,
		EvidenceDigest:                  verification.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionVerification(
		output.Status,
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.GuardEvidenceDigest,
		output.VerificationStatus,
		output.ReverseStatus,
		output.ReverseEvidenceDigest,
		output.FirstMismatch,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-verification-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionVerification(
	status,
	candidateID,
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
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
		VerificationStatus:              verificationStatus,
		ReverseStatus:                   reverseStatus,
		ReverseEvidenceDigest:            reverseEvidenceDigest,
		FirstMismatch:                    firstMismatch,
	})
	if err != nil {
		return ""
	}
	return digest
}