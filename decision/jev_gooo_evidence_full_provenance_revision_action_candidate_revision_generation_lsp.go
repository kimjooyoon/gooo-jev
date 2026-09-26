package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPInput
// adapts generated candidate verification into the editor-facing projection.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPInput struct {
	Verification         ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationVerificationBinding
	MissingStageIndex    int
	EvidencePrefixDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding
// preserves generated candidate provenance in LSP state.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	VerificationEvidenceDigest      string
	ProjectionStatus                string
	Publishable                     bool
	Severity                        string
	Code                            string
	MissingStageIndex               int
	EvidencePrefixDigest            string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding) Validate() error {
	if b.Status != "clear" && b.Status != "publishable" {
		return fmt.Errorf("incomplete Gooo generated candidate LSP binding")
	}
	if b.Status == "clear" && (b.MissingStage != "" || b.MissingStageIndex != -1) {
		return fmt.Errorf("clear Gooo generated candidate LSP state must not have a missing stage")
	}
	if b.Status == "publishable" && (b.MissingStage == "" || b.MissingStageIndex < 0) {
		return fmt.Errorf("publishable Gooo generated candidate LSP state must retain its missing stage")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ProjectionStatus != b.Status ||
		b.Code == "" ||
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo generated candidate LSP evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo generated candidate LSP must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(
		b.Status,
		b.MissingStage,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
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
		return fmt.Errorf("Gooo generated candidate LSP digest mismatch")
	}
	return nil
}

// ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP
// keeps incomplete or missing-prefix evidence UNKNOWN instead of publishable.
func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-lsp"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding{
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
		return unknown("action-candidate-generation-verification")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSPBinding{
		Status:                          projection.Status,
		MissingStage:                    input.Verification.MissingStage,
		CandidateID:                     input.Verification.CandidateID,
		SourceCandidateID:               input.Verification.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Verification.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Verification.SourceCandidateEvidenceDigest,
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
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(
		output.Status,
		output.MissingStage,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
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
		return unknown("revision-action-candidate-generation-lsp-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationLSP(
	status,
	missingStage,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
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
		MissingStage                    string
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		SourceCandidateEvidenceDigest   string
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
		MissingStage:                    missingStage,
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
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