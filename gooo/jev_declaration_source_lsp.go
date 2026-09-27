package gooo

import (
	"fmt"
	"strings"
)

const (
	GoooDeclarationSourceDerived = "derived"
	GoooDeclarationSourceUnknown = "UNKNOWN"
	GoooDeclarationSourceError   = "ERROR"

	goooDeclarationSourceDerivedCode  = "JEV_DECLARATION_SOURCE_DERIVED"
	goooDeclarationSourceUnknownCode   = "JEV_DECLARATION_SOURCE_UNKNOWN"
	goooDeclarationSourceBoundaryCode  = "JEV_DECLARATION_SOURCE_BOUNDARY"
	goooDeclarationSourceIntegrityCode = "JEV_DECLARATION_SOURCE_INTEGRITY"
)

type ExecutionEnvelopeDeclarationSourceDigestInput struct {
	DeclarationID  string
	ContractID     string
	SourceText     string
	NonAuthorizing bool
}

type ExecutionEnvelopeDeclarationSourceDigest struct {
	Status            string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

// ComputeExecutionEnvelopeDeclarationSourceDigest binds exact .gooo source
// text to declaration and contract identity without executing the declaration.
func ComputeExecutionEnvelopeDeclarationSourceDigest(
	input ExecutionEnvelopeDeclarationSourceDigestInput,
) ExecutionEnvelopeDeclarationSourceDigest {
	output := ExecutionEnvelopeDeclarationSourceDigest{
		Status:         GoooDeclarationSourceUnknown,
		NonExecuting:   true,
		NonAuthorizing: input.NonAuthorizing,
	}
	if !input.NonAuthorizing {
		output.MissingStage = "authorization-boundary"
		return output
	}
	if strings.TrimSpace(input.DeclarationID) == "" || strings.TrimSpace(input.ContractID) == "" {
		output.MissingStage = "declaration-identity"
		return output
	}
	if strings.TrimSpace(input.SourceText) == "" {
		output.MissingStage = "declaration-source"
		return output
	}
	output.Status = GoooDeclarationSourceDerived
	output.DeclarationID = input.DeclarationID
	output.ContractID = input.ContractID
	output.DeclarationDigest = digestString(strings.Join([]string{
		"gooo-declaration-source",
		input.DeclarationID,
		input.ContractID,
		input.SourceText,
	}, "|"))
	return output
}

type ExecutionEnvelopeDeclarationSourceLSPInput struct {
	Source         ExecutionEnvelopeDeclarationSourceDigest
	NonAuthorizing bool
}

type ExecutionEnvelopeDeclarationSourceLSPProjection struct {
	Status            string
	Code              string
	Severity          string
	Message           string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	EvidenceDigest    string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

// ProjectExecutionEnvelopeDeclarationSourceLSP exposes declaration origin to
// an editor while keeping source evidence non-executing and non-authorizing.
func ProjectExecutionEnvelopeDeclarationSourceLSP(
	input ExecutionEnvelopeDeclarationSourceLSPInput,
) ExecutionEnvelopeDeclarationSourceLSPProjection {
	output := ExecutionEnvelopeDeclarationSourceLSPProjection{
		Status:         GoooDeclarationSourceUnknown,
		Code:           goooDeclarationSourceUnknownCode,
		Severity:       "warning",
		Message:        "JEV declaration source is UNKNOWN: provenance is incomplete",
		MissingStage:   input.Source.MissingStage,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Source.NonAuthorizing {
		output.Status = GoooDeclarationSourceError
		output.Code = goooDeclarationSourceBoundaryCode
		output.Severity = "error"
		output.Message = "JEV declaration source crossed an authorization boundary"
		output.MissingStage = "authorization-boundary"
		output.EvidenceDigest = executionEnvelopeDeclarationSourceLSPDigest(output)
		return output
	}
	if !input.Source.NonExecuting {
		output.Status = GoooDeclarationSourceError
		output.Code = goooDeclarationSourceBoundaryCode
		output.Severity = "error"
		output.Message = "JEV declaration source crossed an execution boundary"
		output.MissingStage = "execution-boundary"
		output.EvidenceDigest = executionEnvelopeDeclarationSourceLSPDigest(output)
		return output
	}
	if input.Source.Status != GoooDeclarationSourceDerived || !validDigest(input.Source.DeclarationDigest) {
		if output.MissingStage == "" {
			output.MissingStage = "declaration-source"
		}
		output.EvidenceDigest = executionEnvelopeDeclarationSourceLSPDigest(output)
		return output
	}
	output.Status = GoooDeclarationSourceDerived
	output.Code = goooDeclarationSourceDerivedCode
	output.Severity = "info"
	output.Message = "JEV declaration source identity is derived from exact source text"
	output.MissingStage = ""
	output.DeclarationID = input.Source.DeclarationID
	output.ContractID = input.Source.ContractID
	output.DeclarationDigest = input.Source.DeclarationDigest
	output.EvidenceDigest = executionEnvelopeDeclarationSourceLSPDigest(output)
	return output
}

func (projection ExecutionEnvelopeDeclarationSourceLSPProjection) Validate() error {
	switch projection.Status {
	case GoooDeclarationSourceDerived, GoooDeclarationSourceUnknown, GoooDeclarationSourceError:
	default:
		return fmt.Errorf("invalid declaration source LSP status %q", projection.Status)
	}
	if !projection.NonExecuting || !projection.NonAuthorizing {
		return fmt.Errorf("declaration source LSP must remain non-executing and non-authorizing")
	}
	switch projection.Status {
	case GoooDeclarationSourceDerived:
		if projection.Code != goooDeclarationSourceDerivedCode || projection.Severity != "info" ||
			projection.MissingStage != "" || projection.DeclarationID == "" || projection.ContractID == "" ||
			!validDigest(projection.DeclarationDigest) {
			return fmt.Errorf("derived declaration source LSP is incomplete")
		}
	case GoooDeclarationSourceUnknown:
		if projection.Code != goooDeclarationSourceUnknownCode || projection.MissingStage == "" {
			return fmt.Errorf("unknown declaration source LSP must preserve its missing stage")
		}
	case GoooDeclarationSourceError:
		if projection.Code != goooDeclarationSourceBoundaryCode || projection.MissingStage == "" {
			return fmt.Errorf("error declaration source LSP must preserve its boundary stage")
		}
	}
	if projection.EvidenceDigest != executionEnvelopeDeclarationSourceLSPDigest(projection) {
		return fmt.Errorf("declaration source LSP evidence digest mismatch")
	}
	return nil
}

func executionEnvelopeDeclarationSourceLSPDigest(
	projection ExecutionEnvelopeDeclarationSourceLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"gooo-declaration-source-lsp|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		projection.Status,
		projection.Code,
		projection.Severity,
		projection.Message,
		projection.DeclarationID,
		projection.ContractID,
		projection.DeclarationDigest,
		projection.MissingStage,
		projection.NonExecuting,
		projection.NonAuthorizing,
	))
}