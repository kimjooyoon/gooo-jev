package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput
// adapts revision verification into the existing editor-facing projection.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput struct {
	Verification       ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionVerificationBinding
	MissingStageIndex  int
	EvidencePrefixDigest string
	NonAuthorizing     bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding
// preserves revision provenance while exposing only evidence-backed LSP state.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding struct {
	Status                           string
	MissingStage                     string
	CandidateID                      string
	RevisionCandidateDigest          string
	RevisionCandidateEvidenceDigest  string
	VerificationEvidenceDigest       string
	ProjectionStatus                 string
	Publishable                      bool
	Severity                         string
	Code                             string
	MissingStageIndex                int
	EvidencePrefixDigest             string
	EvidenceDigest                   string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding) Validate() error {
	if b.Status != "clear" && b.Status != "publishable" {
		return fmt.Errorf("incomplete Gooo revision action LSP binding")
	}
	if b.Status == "clear" && b.MissingStage != "" {
		return fmt.Errorf("clear Gooo revision action LSP state must not have a missing stage")
	}
	if b.Status == "publishable" && (b.MissingStage == "" || b.MissingStageIndex < 0) {
		return fmt.Errorf("publishable Gooo revision action LSP state must retain its missing stage")
	}
	if b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ProjectionStatus != b.Status ||
		b.Code == "" ||
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo revision action LSP evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo revision action LSP must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionLSP(
		b.Status,
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.Publishable,
		b.Severity,
		b.Code,
		b.MissingStageIndex,
		b.EvidencePrefixDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo revision action LSP digest mismatch")
	}
	return nil
}

// ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP
// keeps incomplete or missing-prefix evidence UNKNOWN instead of publishable.
func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-lsp"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding{
			Status:          "UNKNOWN",
			MissingStage:    stage,
			MissingStageIndex: -1,
			NonExecuting:    true,
			NonAuthorizing:  true,
		}
	}
	if !input.NonAuthorizing || !input.Verification.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Verification.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Verification.Validate(); err != nil {
		return unknown("revision-action-verification")
	}
	projection := ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification: JEVActionCandidateVerification{
			Status:                input.Verification.Status,
			CandidateID:          input.Verification.CandidateID,
			GuardEvidenceDigest:  input.Verification.GuardEvidenceDigest,
			ReverseStatus:        input.Verification.ReverseStatus,
			ReverseEvidenceDigest: input.Verification.ReverseEvidenceDigest,
			FirstMismatch:        input.Verification.FirstMismatch,
			EvidenceDigest:       input.Verification.EvidenceDigest,
			MissingStage:         input.Verification.MissingStage,
			NonExecuting:         true,
			NonAuthorizing:       true,
		},
		MissingStageIndex:   input.MissingStageIndex,
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		NonAuthorizing:      true,
	})
	if projection.Status != "clear" && projection.Status != "publishable" {
		return unknown(projection.Code)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPBinding{
		Status:                          projection.Status,
		MissingStage:                    input.Verification.MissingStage,
		CandidateID:                     input.Verification.CandidateID,
		RevisionCandidateDigest:         input.Verification.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Verification.RevisionCandidateEvidenceDigest,
		VerificationEvidenceDigest:      input.Verification.EvidenceDigest,
		ProjectionStatus:                projection.Status,
		Publishable:                     projection.Publishable,
		Severity:                        projection.Severity,
		Code:                            projection.Code,
		MissingStageIndex:               projection.MissingStageIndex,
		EvidencePrefixDigest:            projection.EvidencePrefixDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionLSP(
		output.Status,
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.Publishable,
		output.Severity,
		output.Code,
		output.MissingStageIndex,
		output.EvidencePrefixDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-lsp-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionLSP(
	status,
	candidateID,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	verificationEvidenceDigest string,
	publishable bool,
	severity,
	code string,
	missingStageIndex int,
	evidencePrefixDigest string,
) string {
	digest, err := Digest(struct {
		Status                          string
		CandidateID                     string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		VerificationEvidenceDigest      string
		Publishable                     bool
		Severity                        string
		Code                            string
		MissingStageIndex               int
		EvidencePrefixDigest            string
	}{
		Status:                          status,
		CandidateID:                     candidateID,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		Publishable:                     publishable,
		Severity:                        severity,
		Code:                            code,
		MissingStageIndex:               missingStageIndex,
		EvidencePrefixDigest:            evidencePrefixDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}