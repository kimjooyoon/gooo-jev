package gooo

import "fmt"

const (
	ProvenanceLineageLSPBound  = "BOUND"
	ProvenanceLineageLSPUnknown = "UNKNOWN"

	provenanceLineageLSPComplete = "provenance-lineage-complete"
	provenanceLineageLSPUnknown  = "provenance-lineage"
	provenanceLineageLSPIntegrity = "lineage-integrity"
	provenanceLineageLSPBoundary = "capability-boundary"
)

// ProvenanceLineageLSPProjection exposes a lineage receipt to an editor
// without treating evidence availability as correctness or improvement.
type ProvenanceLineageLSPProjection struct {
	Status             string
	Publishable        bool
	Severity           string
	Code               string
	MissingStage       string
	EdgeCount          int
	ChainDigest        string
	EvidenceDigest     string
	ReadOnly           bool
	ClaimsImprovement  bool
	CanExecute         bool
	CanAuthorize       bool
	NonAuthorizing     bool
	ObservationDigest  string
}

// ProjectProvenanceLineageLSP preserves unknown stages and publishes only a
// complete, validated lineage receipt as a clear read-only projection.
func ProjectProvenanceLineageLSP(receipt LineageReceipt) ProvenanceLineageLSPProjection {
	output := ProvenanceLineageLSPProjection{
		Status:            ProvenanceLineageLSPUnknown,
		Publishable:       false,
		Severity:          "error",
		Code:              provenanceLineageLSPUnknown,
		MissingStage:      receipt.MissingStage,
		ReadOnly:          true,
		ClaimsImprovement: false,
		CanExecute:        false,
		CanAuthorize:      false,
		NonAuthorizing:    true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "lineage-observation"
	}
	if !receipt.NonExecuting || !receipt.NonAuthorizing {
		output.NonAuthorizing = false
		output.Code = provenanceLineageLSPBoundary
		output.MissingStage = "capability-boundary"
		output.ObservationDigest = provenanceLineageLSPDigest(output)
		return output
	}
	if receipt.Status != "BOUND" {
		output.ObservationDigest = provenanceLineageLSPDigest(output)
		return output
	}
	if err := receipt.Validate(); err != nil {
		output.Code = provenanceLineageLSPIntegrity
		output.MissingStage = "lineage-validation"
		output.ObservationDigest = provenanceLineageLSPDigest(output)
		return output
	}

	output.Status = ProvenanceLineageLSPBound
	output.Publishable = false
	output.Severity = "info"
	output.Code = provenanceLineageLSPComplete
	output.MissingStage = ""
	output.EdgeCount = len(receipt.Edges)
	output.ChainDigest = receipt.ChainDigest
	output.EvidenceDigest = receipt.EvidenceDigest
	output.ObservationDigest = provenanceLineageLSPDigest(output)
	return output
}

func (projection ProvenanceLineageLSPProjection) Validate() error {
	if projection.Status != ProvenanceLineageLSPBound && projection.Status != ProvenanceLineageLSPUnknown {
		return fmt.Errorf("invalid provenance lineage LSP status %q", projection.Status)
	}
	if !projection.ReadOnly || projection.ClaimsImprovement || projection.CanExecute || projection.CanAuthorize {
		return fmt.Errorf("provenance lineage LSP crossed a forbidden boundary")
	}
	if projection.Status == ProvenanceLineageLSPUnknown {
		if projection.MissingStage == "" || projection.Code == provenanceLineageLSPComplete {
			return fmt.Errorf("unknown provenance lineage LSP must preserve its unresolved stage")
		}
	} else {
		if projection.MissingStage != "" ||
			projection.Code != provenanceLineageLSPComplete ||
			projection.EdgeCount != 4 ||
			!validDigest(projection.ChainDigest) ||
			!validDigest(projection.EvidenceDigest) {
			return fmt.Errorf("bound provenance lineage LSP is incomplete")
		}
	}
	if projection.ObservationDigest != provenanceLineageLSPDigest(projection) {
		return fmt.Errorf("provenance lineage LSP observation digest mismatch")
	}
	return nil
}

func provenanceLineageLSPDigest(projection ProvenanceLineageLSPProjection) string {
	return digestString(fmt.Sprintf(
		"provenance-lineage-lsp|%s|%t|%s|%s|%s|%d|%s|%s|%t|%t|%t|%t|%t",
		projection.Status,
		projection.Publishable,
		projection.Severity,
		projection.Code,
		projection.MissingStage,
		projection.EdgeCount,
		projection.ChainDigest,
		projection.EvidenceDigest,
		projection.ReadOnly,
		projection.ClaimsImprovement,
		projection.CanExecute,
		projection.CanAuthorize,
		projection.NonAuthorizing,
	))
}