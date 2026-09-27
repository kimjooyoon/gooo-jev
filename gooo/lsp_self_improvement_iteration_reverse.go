package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementIterationReverseSymbolResult is a read-only LSP projection
// joining iteration provenance with reverse generation evidence.
type SelfImprovementIterationReverseSymbolResult struct {
	Status                            string
	MissingStage                      string
	SourceDigest                      string
	IRDigest                          string
	SymbolName                        string
	SymbolKind                        SymbolKind
	SymbolDigest                      string
	CurrentLifecycleObservationDigest string
	IterationObservationDigest        string
	ReverseObservationDigest          string
	CandidateSourceDigest             string
	CandidateProposedSourceDigest     string
	CandidateGeneratedIRDigest       string
	ReverseSourceDigest               string
	ReverseProposedSourceDigest       string
	ReverseGeneratedIRDigest          string
	ReverseGeneratedIRMatchesCandidate bool
	MetricsBindingDigest              string
	GenerationAssessmentDigest        string
	ExactSourceMatch                  bool
	ExactIRMatch                      bool
	ExactStructureMatch               bool
	ReverseSignal                     string
	RelationSignal                    string
	ResultDigest                      string
	NonExecuting                      bool
	NonAuthorizing                    bool
}

// ExplainSelfImprovementIterationReverseObservation joins the current
// declaration, iteration provenance, and reverse generation evidence without
// executing or authorizing a candidate change.
func ExplainSelfImprovementIterationReverseObservation(
	source, symbolName string,
	current RevisionSelfImprovementLifecycleObservation,
	iteration RevisionSelfImprovementIterationProvenanceObservation,
	reverse RevisionSelfImprovementReverseObservation,
) (SelfImprovementIterationReverseSymbolResult, error) {
	result := SelfImprovementIterationReverseSymbolResult{
		Status:                            "UNKNOWN",
		MissingStage:                      "lsp-self-improvement-iteration-reverse-provenance",
		SourceDigest:                      digestString(source),
		CurrentLifecycleObservationDigest: current.ObservationDigest,
		IterationObservationDigest:        iteration.ObservationDigest,
		ReverseObservationDigest:          reverse.ObservationDigest,
		CandidateSourceDigest:             iteration.CandidateSourceDigest,
		CandidateProposedSourceDigest:     iteration.CandidateProposedSourceDigest,
		CandidateGeneratedIRDigest:        iteration.CandidateGeneratedIRDigest,
		ReverseSourceDigest:               reverse.SourceDigest,
		ReverseProposedSourceDigest:       reverse.ProposedSourceDigest,
		ReverseGeneratedIRDigest:          reverse.GeneratedIRDigest,
		ReverseGeneratedIRMatchesCandidate: reverse.GeneratedIRDigest == iteration.CandidateGeneratedIRDigest,
		MetricsBindingDigest:              reverse.MetricsBindingDigest,
		GenerationAssessmentDigest:        reverse.GenerationAssessmentDigest,
		ExactSourceMatch:                  reverse.ExactSourceMatch,
		ExactIRMatch:                      reverse.ExactIRMatch,
		ExactStructureMatch:               reverse.ExactStructureMatch,
		ReverseSignal:                     reverse.ReverseSignal,
		RelationSignal:                    "iteration-reverse-unknown",
		NonExecuting:                      true,
		NonAuthorizing:                    true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementIterationReverse(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-reverse-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := current.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-reverse-current"
		setResultDigest()
		return result, fmt.Errorf("current self-improvement lifecycle is not valid: %w", err)
	}
	if err := iteration.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-reverse-iteration"
		setResultDigest()
		return result, fmt.Errorf("iteration provenance observation is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-reverse-reverse"
		setResultDigest()
		return result, fmt.Errorf("reverse observation is not valid: %w", err)
	}
	if snapshot.SourceDigest != current.SourceDigest ||
		iteration.SourceDigest != snapshot.SourceDigest ||
		reverse.SourceDigest != current.SourceDigest ||
		reverse.LifecycleObservationDigest != current.ObservationDigest {
		result.MissingStage = "lsp-self-improvement-iteration-reverse-source-link"
		setResultDigest()
		return result, fmt.Errorf("iteration and reverse evidence are not linked to current source")
	}
	if iteration.IRDigest != snapshot.IRDigest ||
		reverse.InputIRDigest != snapshot.IRDigest {
		result.MissingStage = "lsp-self-improvement-iteration-reverse-ir-link"
		setResultDigest()
		return result, fmt.Errorf("iteration and reverse evidence are not linked to current IR")
	}

	for _, symbol := range snapshot.Symbols {
		if symbol.Name == symbolName {
			result.Status = "BOUND"
			result.MissingStage = ""
			result.IRDigest = snapshot.IRDigest
			result.SymbolName = symbol.Name
			result.SymbolKind = symbol.Kind
			result.SymbolDigest = symbol.Digest
			if result.ReverseGeneratedIRMatchesCandidate {
				result.RelationSignal = "iteration-reverse-aligned"
			} else {
				result.RelationSignal = "iteration-reverse-distinct"
			}
			setResultDigest()
			if err := result.Validate(); err != nil {
				result.Status = "UNKNOWN"
				result.MissingStage = "lsp-self-improvement-iteration-reverse-provenance"
				result.RelationSignal = "iteration-reverse-unknown"
				setResultDigest()
				return result, fmt.Errorf("LSP iteration reverse result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-iteration-reverse-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementIterationReverseSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP iteration reverse status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP iteration reverse has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP iteration reverse has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP iteration reverse must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP iteration reverse is incomplete")
	}
	for name, digest := range map[string]string{
		"source":              r.SourceDigest,
		"ir":                  r.IRDigest,
		"symbol":              r.SymbolDigest,
		"lifecycle":           r.CurrentLifecycleObservationDigest,
		"iteration":           r.IterationObservationDigest,
		"reverse":             r.ReverseObservationDigest,
		"candidate-source":    r.CandidateSourceDigest,
		"candidate-proposed":  r.CandidateProposedSourceDigest,
		"candidate-generated": r.CandidateGeneratedIRDigest,
		"reverse-source":      r.ReverseSourceDigest,
		"reverse-proposed":    r.ReverseProposedSourceDigest,
		"reverse-generated":   r.ReverseGeneratedIRDigest,
		"metrics-binding":     r.MetricsBindingDigest,
		"generation-assessment": r.GenerationAssessmentDigest,
		"result":              r.ResultDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("LSP iteration reverse %s digest is invalid", name)
		}
	}
	if r.ReverseSignal == "" {
		return fmt.Errorf("LSP iteration reverse signal is empty")
	}
	if r.RelationSignal != "iteration-reverse-aligned" &&
		r.RelationSignal != "iteration-reverse-distinct" &&
		r.RelationSignal != "iteration-reverse-unknown" {
		return fmt.Errorf("LSP iteration reverse relation signal is invalid")
	}
	if digestLSPSelfImprovementIterationReverse(r) != r.ResultDigest {
		return fmt.Errorf("LSP iteration reverse digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementIterationReverse(result SelfImprovementIterationReverseSymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		string(result.SymbolKind),
		result.SymbolDigest,
		result.CurrentLifecycleObservationDigest,
		result.IterationObservationDigest,
		result.ReverseObservationDigest,
		result.CandidateSourceDigest,
		result.CandidateProposedSourceDigest,
		result.CandidateGeneratedIRDigest,
		result.ReverseSourceDigest,
		result.ReverseProposedSourceDigest,
		result.ReverseGeneratedIRDigest,
		strconv.FormatBool(result.ReverseGeneratedIRMatchesCandidate),
		result.MetricsBindingDigest,
		result.GenerationAssessmentDigest,
		strconv.FormatBool(result.ExactSourceMatch),
		strconv.FormatBool(result.ExactIRMatch),
		strconv.FormatBool(result.ExactStructureMatch),
		result.ReverseSignal,
		result.RelationSignal,
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
