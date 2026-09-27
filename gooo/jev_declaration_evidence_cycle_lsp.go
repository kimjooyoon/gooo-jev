package gooo

import "fmt"

const (
	ExecutionEnvelopeDeclarationEvidenceCycleBound    = "BOUND"
	ExecutionEnvelopeDeclarationEvidenceCycleDeferred = "DEFERRED"
	ExecutionEnvelopeDeclarationEvidenceCycleUnknown  = "UNKNOWN"
	ExecutionEnvelopeDeclarationEvidenceCycleError    = "ERROR"

	executionEnvelopeDeclarationEvidenceCycleBoundCode     = "gooo.declaration_evidence_cycle.bound"
	executionEnvelopeDeclarationEvidenceCycleDeferredCode  = "gooo.declaration_evidence_cycle.deferred"
	executionEnvelopeDeclarationEvidenceCycleUnknownCode   = "gooo.declaration_evidence_cycle.unknown"
	executionEnvelopeDeclarationEvidenceCycleIntegrityCode = "gooo.declaration_evidence_cycle.integrity"
	executionEnvelopeDeclarationEvidenceCycleBoundaryCode  = "gooo.declaration_evidence_cycle.capability_boundary"
)

type ExecutionEnvelopeDeclarationEvidenceCycleInput struct {
	SourceStatus              string
	DeclarationID             string
	ContractID                string
	DeclarationDigest         string
	IRStatus                  string
	IRDigest                  string
	GenerationDigest          string
	BindingDigest             string
	ReverseStatus             string
	ReverseObservationDigest  string
	ReverseFirstMismatch      string
	NonExecuting              bool
	NonAuthorizing            bool
}

type ExecutionEnvelopeDeclarationEvidenceCycleProjection struct {
	Status                    string
	Code                      string
	Severity                  string
	Message                   string
	DeclarationID             string
	ContractID                string
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	BindingDigest             string
	ReverseObservationDigest  string
	FirstMismatch             string
	MissingStage              string
	EvidenceDigest             string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP exposes the first
// unresolved stage across source, IR/generation, and reverse observation.
func ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(
	input ExecutionEnvelopeDeclarationEvidenceCycleInput,
) ExecutionEnvelopeDeclarationEvidenceCycleProjection {
	output := ExecutionEnvelopeDeclarationEvidenceCycleProjection{
		Status:                   ExecutionEnvelopeDeclarationEvidenceCycleUnknown,
		Code:                     executionEnvelopeDeclarationEvidenceCycleUnknownCode,
		Severity:                 "warning",
		Message:                  "declaration evidence cycle is unresolved",
		DeclarationID:            input.DeclarationID,
		ContractID:               input.ContractID,
		DeclarationDigest:         input.DeclarationDigest,
		IRDigest:                 input.IRDigest,
		GenerationDigest:         input.GenerationDigest,
		BindingDigest:            input.BindingDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		FirstMismatch:            input.ReverseFirstMismatch,
		MissingStage:             "declaration_source",
		NonExecuting:             true,
		NonAuthorizing:           true,
	}

	switch {
	case !input.NonExecuting || !input.NonAuthorizing:
		output.Status = ExecutionEnvelopeDeclarationEvidenceCycleError
		output.Code = executionEnvelopeDeclarationEvidenceCycleBoundaryCode
		output.Severity = "error"
		output.Message = "declaration evidence cycle crossed a capability boundary"
		output.MissingStage = "capability_boundary"
	case input.SourceStatus != GoooDeclarationSourceDerived:
		output.MissingStage = "declaration_source"
		if input.DeclarationID == "" || input.ContractID == "" {
			output.FirstMismatch = "declaration-identity"
		}
	case input.IRStatus != ExecutionEnvelopeDeclarationIRGenerationLSPBound:
		output.MissingStage = "declaration_ir_generation"
	case input.ReverseStatus == "DEFERRED":
		output.Status = ExecutionEnvelopeDeclarationEvidenceCycleDeferred
		output.Code = executionEnvelopeDeclarationEvidenceCycleDeferredCode
		output.Severity = "warning"
		output.Message = "declaration evidence cycle is deferred at reverse observation"
		output.MissingStage = "reverse_observation"
	case input.ReverseStatus != "BOUND":
		output.MissingStage = "reverse_observation"
	case !validDigest(input.DeclarationDigest) ||
		!validDigest(input.IRDigest) ||
		!validDigest(input.GenerationDigest) ||
		!validDigest(input.BindingDigest) ||
		!validDigest(input.ReverseObservationDigest):
		output.Status = ExecutionEnvelopeDeclarationEvidenceCycleError
		output.Code = executionEnvelopeDeclarationEvidenceCycleIntegrityCode
		output.Severity = "error"
		output.Message = "declaration evidence cycle failed digest validation"
		output.MissingStage = "evidence_integrity"
	default:
		output.Status = ExecutionEnvelopeDeclarationEvidenceCycleBound
		output.Code = executionEnvelopeDeclarationEvidenceCycleBoundCode
		output.Severity = "info"
		output.Message = "declaration source, IR generation, and reverse observation are bound"
		output.MissingStage = ""
		output.FirstMismatch = ""
	}

	output.EvidenceDigest = executionEnvelopeDeclarationEvidenceCycleDigest(output)
	return output
}

func (projection ExecutionEnvelopeDeclarationEvidenceCycleProjection) Validate() error {
	switch projection.Status {
	case ExecutionEnvelopeDeclarationEvidenceCycleBound,
		ExecutionEnvelopeDeclarationEvidenceCycleDeferred,
		ExecutionEnvelopeDeclarationEvidenceCycleUnknown,
		ExecutionEnvelopeDeclarationEvidenceCycleError:
	default:
		return fmt.Errorf("invalid declaration evidence cycle status %q", projection.Status)
	}
	if !projection.NonExecuting || !projection.NonAuthorizing {
		return fmt.Errorf("declaration evidence cycle must remain non-executing and non-authorizing")
	}
	switch projection.Status {
	case ExecutionEnvelopeDeclarationEvidenceCycleBound:
		if projection.Code != executionEnvelopeDeclarationEvidenceCycleBoundCode ||
			projection.MissingStage != "" ||
			projection.FirstMismatch != "" ||
			!validDigest(projection.DeclarationDigest) ||
			!validDigest(projection.IRDigest) ||
			!validDigest(projection.GenerationDigest) ||
			!validDigest(projection.BindingDigest) ||
			!validDigest(projection.ReverseObservationDigest) {
			return fmt.Errorf("bound declaration evidence cycle is incomplete")
		}
	case ExecutionEnvelopeDeclarationEvidenceCycleDeferred:
		if projection.Code != executionEnvelopeDeclarationEvidenceCycleDeferredCode ||
			projection.MissingStage != "reverse_observation" ||
			projection.FirstMismatch == "" {
			return fmt.Errorf("deferred declaration evidence cycle lost its reverse stage")
		}
	case ExecutionEnvelopeDeclarationEvidenceCycleUnknown:
		if projection.Code != executionEnvelopeDeclarationEvidenceCycleUnknownCode ||
			projection.MissingStage == "" {
			return fmt.Errorf("unknown declaration evidence cycle lost its missing stage")
		}
	case ExecutionEnvelopeDeclarationEvidenceCycleError:
		if projection.MissingStage == "" ||
			(projection.Code != executionEnvelopeDeclarationEvidenceCycleIntegrityCode &&
				projection.Code != executionEnvelopeDeclarationEvidenceCycleBoundaryCode) {
			return fmt.Errorf("error declaration evidence cycle lost its failure stage")
		}
	}
	if projection.EvidenceDigest != executionEnvelopeDeclarationEvidenceCycleDigest(projection) {
		return fmt.Errorf("declaration evidence cycle digest mismatch")
	}
	return nil
}

func executionEnvelopeDeclarationEvidenceCycleDigest(
	projection ExecutionEnvelopeDeclarationEvidenceCycleProjection,
) string {
	return digestString(fmt.Sprintf(
		"gooo-declaration-evidence-cycle|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		projection.Status,
		projection.Code,
		projection.Severity,
		projection.Message,
		projection.DeclarationID,
		projection.ContractID,
		projection.DeclarationDigest,
		projection.IRDigest,
		projection.GenerationDigest,
		projection.BindingDigest,
		projection.ReverseObservationDigest,
		projection.FirstMismatch,
		projection.MissingStage,
		projection.NonExecuting,
		projection.NonAuthorizing,
	))
}

