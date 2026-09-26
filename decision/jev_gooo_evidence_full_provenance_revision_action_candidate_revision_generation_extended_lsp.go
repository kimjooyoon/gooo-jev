package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput
// adapts extended generated candidate verification into the editor-facing projection.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput struct {
	Verification         ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedVerificationBinding
	MissingStageIndex    int
	EvidencePrefixDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding
// preserves extended generated candidate provenance in LSP state.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	GuardEvidenceDigest             string
	VerificationEvidenceDigest      string
	ReverseEvidenceDigest           string
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

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding) Validate() error {
	if b.Status != "clear" && b.Status != "publishable" {
		return fmt.Errorf("incomplete Gooo extended generated candidate LSP binding")
	}
	if b.Status == "clear" && (b.MissingStage != "" || b.MissingStageIndex != -1) {
		return fmt.Errorf("clear Gooo extended generated candidate LSP state must not have a missing stage")
	}
	if b.Status == "publishable" && (b.MissingStage == "" || b.MissingStageIndex < 0) {
		return fmt.Errorf("publishable Gooo extended generated candidate LSP state must retain its missing stage")
	}
	if b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.ProjectionStatus != b.Status ||
		b.Code == "" ||
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended generated candidate LSP evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended generated candidate LSP must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(
		b.Status,
		b.MissingStage,
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.GuardEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.ReverseEvidenceDigest,
		b.Publishable,
		b.Severity,
		b.Code,
		b.MissingStageIndex,
		b.EvidencePrefixDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended generated candidate LSP digest mismatch")
	}
	return nil
}

// ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP
// keeps incomplete or missing-prefix evidence UNKNOWN instead of publishable.
func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lsp"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding{
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
		return unknown("action-candidate-generation-extended-verification")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSPBinding{
		Status:                          projection.Status,
		MissingStage:                    input.Verification.MissingStage,
		CandidateID:                     input.Verification.CandidateID,
		SourceCandidateID:               input.Verification.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Verification.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Verification.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Verification.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Verification.RevisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             input.Verification.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Verification.EvidenceDigest,
		ReverseEvidenceDigest:           input.Verification.ReverseEvidenceDigest,
		ProjectionStatus:                projection.Status,
		Publishable:                     projection.Publishable,
		Severity:                        projection.Severity,
		Code:                            projection.Code,
		MissingStageIndex:               projection.MissingStageIndex,
		EvidencePrefixDigest:            projection.EvidencePrefixDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(
		output.Status,
		output.MissingStage,
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.GuardEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.ReverseEvidenceDigest,
		output.Publishable,
		output.Severity,
		output.Code,
		output.MissingStageIndex,
		output.EvidencePrefixDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lsp-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLSP(
	status,
	missingStage,
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	guardEvidenceDigest,
	verificationEvidenceDigest,
	reverseEvidenceDigest string,
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
		GuardEvidenceDigest             string
		VerificationEvidenceDigest      string
		ReverseEvidenceDigest           string
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
		GuardEvidenceDigest:             guardEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		ReverseEvidenceDigest:           reverseEvidenceDigest,
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
