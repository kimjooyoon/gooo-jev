package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput
// adapts candidate revision verification into the editor-facing projection.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput struct {
	Verification         ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationBinding
	MissingStageIndex    int
	EvidencePrefixDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding
// preserves source and generated candidate provenance in LSP state.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding struct {
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

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding) Validate() error {
	if b.Status != "clear" && b.Status != "publishable" {
		return fmt.Errorf("incomplete Gooo action-derived candidate revision LSP binding")
	}
	if b.Status == "clear" && (b.MissingStage != "" || b.MissingStageIndex != -1) {
		return fmt.Errorf("clear Gooo action-derived candidate revision LSP state must not have a missing stage")
	}
	if b.Status == "publishable" && (b.MissingStage == "" || b.MissingStageIndex < 0) {
		return fmt.Errorf("publishable Gooo action-derived candidate revision LSP state must retain its missing stage")
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
		return fmt.Errorf("incomplete Gooo action-derived candidate revision LSP evidence")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate revision LSP must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(
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
		return fmt.Errorf("Gooo action-derived candidate revision LSP digest mismatch")
	}
	return nil
}

// ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP
// keeps incomplete or missing-prefix evidence UNKNOWN instead of publishable.
func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-revision-lsp"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding{
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
		return unknown("action-candidate-revision-verification")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding{
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
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(
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
		return unknown("revision-action-candidate-revision-lsp-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(
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