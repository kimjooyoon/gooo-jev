package gooo

import (
	"fmt"
	"strings"
)

const (
	lspSemanticTokenTypeGooProvenance  = "gooo_provenance"
	lspSemanticTokenProvenanceObserved = "observed"
	lspSemanticTokenProvenanceUnknown  = "unknown"
	lspSemanticTokenModifiersObserved  = "source|ir|generated|reverse_observed"
	lspSemanticTokenModifiersUnknown   = "unknown"
)

type LSPSemanticTokenProvenance struct {
	Status                   string
	MissingStage             string
	ProvenanceStatus         string
	TokenType                string
	TokenModifiers           string
	SourceDigest             string
	IRDigest                 string
	GenerationDigest         string
	ReverseObservationDigest string
	InputObservationDigest   string
	ObservationDigest        string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveLSPSemanticTokenProvenance(
	input RevisionSelfImprovementCycleJEVSubsetConvergenceObservation,
) LSPSemanticTokenProvenance {
	result := LSPSemanticTokenProvenance{
		Status:                   "UNKNOWN",
		MissingStage:             "lsp-semantic-token-provenance",
		ProvenanceStatus:         lspSemanticTokenProvenanceUnknown,
		TokenType:                lspSemanticTokenTypeGooProvenance,
		TokenModifiers:           lspSemanticTokenModifiersUnknown,
		SourceDigest:              input.DeclarationDigest,
		IRDigest:                 input.IRDigest,
		GenerationDigest:         input.GenerationDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		InputObservationDigest:   input.ObservationDigest,
		ReadOnly:                 true,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ObservationDigest = digestLSPSemanticTokenProvenance(result)
	}
	setDigest()

	if err := input.Validate(); err != nil {
		result.MissingStage = "lsp-semantic-token-provenance-input"
		setDigest()
		return result
	}

	result.Status = "BOUND"
	result.MissingStage = input.MissingStage
	if input.Status == "BOUND" {
		result.ProvenanceStatus = lspSemanticTokenProvenanceObserved
		result.TokenModifiers = lspSemanticTokenModifiersObserved
	} else {
		result.ProvenanceStatus = lspSemanticTokenProvenanceUnknown
		result.TokenModifiers = lspSemanticTokenModifiersUnknown
	}
	setDigest()
	return result
}

func (value LSPSemanticTokenProvenance) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("LSP semantic token provenance status is invalid")
	}
	if value.TokenType != lspSemanticTokenTypeGooProvenance {
		return fmt.Errorf("LSP semantic token provenance type is invalid")
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP semantic token provenance must remain read-only, non-executing, and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.ProvenanceStatus != lspSemanticTokenProvenanceUnknown ||
			value.TokenModifiers != lspSemanticTokenModifiersUnknown {
			return fmt.Errorf("unknown LSP semantic token provenance must preserve its missing stage")
		}
	} else {
		if value.InputObservationDigest == "" {
			return fmt.Errorf("bound LSP semantic token provenance is missing its input observation")
		}
		switch value.ProvenanceStatus {
		case lspSemanticTokenProvenanceObserved:
			if value.TokenModifiers != lspSemanticTokenModifiersObserved {
				return fmt.Errorf("observed LSP semantic token provenance modifiers are invalid")
			}
		case lspSemanticTokenProvenanceUnknown:
			if value.TokenModifiers != lspSemanticTokenModifiersUnknown {
				return fmt.Errorf("unknown LSP semantic token provenance modifiers are invalid")
			}
		default:
			return fmt.Errorf("LSP semantic token provenance state is invalid")
		}
	}
	if value.ObservationDigest != digestLSPSemanticTokenProvenance(value) {
		return fmt.Errorf("LSP semantic token provenance digest does not match its fields")
	}
	return nil
}

func digestLSPSemanticTokenProvenance(value LSPSemanticTokenProvenance) string {
	return digestString(strings.Join([]string{
		value.Status,
		value.MissingStage,
		value.ProvenanceStatus,
		value.TokenType,
		value.TokenModifiers,
		value.SourceDigest,
		value.IRDigest,
		value.GenerationDigest,
		value.ReverseObservationDigest,
		value.InputObservationDigest,
		fmt.Sprintf("%t", value.ReadOnly),
		fmt.Sprintf("%t", value.NonExecuting),
		fmt.Sprintf("%t", value.NonAuthorizing),
	}, "|"))
}
