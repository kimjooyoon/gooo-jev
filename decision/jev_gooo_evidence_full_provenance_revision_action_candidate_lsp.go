package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput
// adapts action-derived candidate verification into the editor-facing projection.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput struct {
	Verification         ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateVerificationBinding
	MissingStageIndex    int
	EvidencePrefixDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding
// preserves candidate provenance while exposing only evidence-backed LSP state.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding struct {
	Status                     string
	MissingStage               string
	CandidateID                string
	RevisionCandidateDigest    string
	CandidateEvidenceDigest    string
	GuardEvidenceDigest        string
	VerificationEvidenceDigest string
	ProjectionStatus           string
	Publishable                bool
	Severity                   string
	Code                       string
	MissingStageIndex          int
	EvidencePrefixDigest       string
	EvidenceDigest             string
	NonExecuting               bool
	NonAuthorizing             bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding) Validate() error {
	if b.Status != "clear" && b.Status != "publishable" {
		return fmt.Errorf("incomplete Gooo action-derived candidate LSP binding")
	}
	if b.Status == "clear" && (b.MissingStage != "" || b.MissingStageIndex != -1) {
		return fmt.Errorf("clear Gooo action-derived candidate LSP state must not have a missing stage")
	}
	if b.Status == "publishable" && (b.MissingStage == "" || b.MissingStageIndex < 0) {
		return fmt.Errorf("publishable Gooo action-derived candidate LSP state must retain its missing stage")
	}
	if b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ProjectionStatus != b.Status ||
		b.Code == "" ||
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate LSP evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate LSP must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateLSP(
		b.Status,
		b.MissingStage,
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.CandidateEvidenceDigest,
		b.GuardEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.Publishable,
		b.Severity,
		b.Code,
		b.MissingStageIndex,
		b.EvidencePrefixDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate LSP digest mismatch")
	}
	return nil
}

// ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP
// keeps incomplete or missing-prefix evidence UNKNOWN instead of publishable.
func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-lsp"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding{
			Status:            "UNKNOWN",
			MissingStage:      stage,
			MissingStageIndex: -1,
			NonExecuting:      true,
			NonAuthorizing:    true,
		}
	}
	if !input.NonAuthorizing || !input.Verification.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Verification.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Verification.Validate(); err != nil {
		return unknown("revision-action-candidate-verification")
	}
	projection := ProjectJEVActionCandidateLSP(JEVActionCandidateLSPProjectionInput{
		Verification: JEVActionCandidateVerification{
			Status:                input.Verification.Status,
			CandidateID:           input.Verification.CandidateID,
			GuardEvidenceDigest:   input.Verification.GuardEvidenceDigest,
			ReverseStatus:         input.Verification.ReverseStatus,
			ReverseEvidenceDigest: input.Verification.ReverseEvidenceDigest,
			FirstMismatch:         input.Verification.FirstMismatch,
			EvidenceDigest:        input.Verification.EvidenceDigest,
			MissingStage:          input.Verification.MissingStage,
			NonExecuting:          true,
			NonAuthorizing:        true,
		},
		MissingStageIndex:    input.MissingStageIndex,
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		NonAuthorizing:       true,
	})
	if projection.Status != "clear" && projection.Status != "publishable" {
		return unknown(projection.Code)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPBinding{
		Status:                     projection.Status,
		MissingStage:               input.Verification.MissingStage,
		CandidateID:                input.Verification.CandidateID,
		RevisionCandidateDigest:    input.Verification.RevisionCandidateDigest,
		CandidateEvidenceDigest:    input.Verification.CandidateEvidenceDigest,
		GuardEvidenceDigest:        input.Verification.GuardEvidenceDigest,
		VerificationEvidenceDigest: input.Verification.EvidenceDigest,
		ProjectionStatus:           projection.Status,
		Publishable:                projection.Publishable,
		Severity:                   projection.Severity,
		Code:                       projection.Code,
		MissingStageIndex:          projection.MissingStageIndex,
		EvidencePrefixDigest:       projection.EvidencePrefixDigest,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateLSP(
		output.Status,
		output.MissingStage,
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.CandidateEvidenceDigest,
		output.GuardEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.Publishable,
		output.Severity,
		output.Code,
		output.MissingStageIndex,
		output.EvidencePrefixDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-lsp-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateLSP(
	status,
	missingStage,
	candidateID,
	revisionCandidateDigest,
	candidateEvidenceDigest,
	guardEvidenceDigest,
	verificationEvidenceDigest string,
	publishable bool,
	severity,
	code string,
	missingStageIndex int,
	evidencePrefixDigest string,
) string {
	digest, err := Digest(struct {
		Status                     string
		MissingStage               string
		CandidateID                string
		RevisionCandidateDigest    string
		CandidateEvidenceDigest    string
		GuardEvidenceDigest        string
		VerificationEvidenceDigest string
		Publishable                bool
		Severity                   string
		Code                       string
		MissingStageIndex          int
		EvidencePrefixDigest       string
	}{
		Status:                     status,
		MissingStage:               missingStage,
		CandidateID:                candidateID,
		RevisionCandidateDigest:    revisionCandidateDigest,
		CandidateEvidenceDigest:    candidateEvidenceDigest,
		GuardEvidenceDigest:        guardEvidenceDigest,
		VerificationEvidenceDigest: verificationEvidenceDigest,
		Publishable:                publishable,
		Severity:                   severity,
		Code:                       code,
		MissingStageIndex:          missingStageIndex,
		EvidencePrefixDigest:       evidencePrefixDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}