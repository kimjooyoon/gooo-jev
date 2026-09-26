package gooo

import "fmt"

// SymbolKind identifies a declaration surfaced to language tooling.
type SymbolKind string

const (
	EntitySymbol   SymbolKind = "entity"
	ActivitySymbol SymbolKind = "activity"
)

// DiagnosticSeverity follows the LSP severity ordering.
type DiagnosticSeverity int

const (
	SeverityError DiagnosticSeverity = 1
)

// Symbol is a provenance-linked declaration projection.
type Symbol struct {
	Name     string
	Kind     SymbolKind
	Position Position
	Digest   string
}

// Diagnostic is a source-bound parser or IR diagnostic.
type Diagnostic struct {
	Stage        string
	Message      string
	Position     Position
	Severity     DiagnosticSeverity
	SourceDigest string
}

// LanguageSnapshot is a deterministic, non-executing language-tooling view.
type LanguageSnapshot struct {
	Status         string
	MissingStage   string
	Package        string
	Namespace      string
	SourceDigest   string
	IRDigest       string
	Symbols        []Symbol
	Diagnostics    []Diagnostic
	NonExecuting   bool
	NonAuthorizing bool
}

// Analyze projects parser and IR provenance into a tooling snapshot.
func Analyze(source string) LanguageSnapshot {
	snapshot := LanguageSnapshot{
		Status:         "UNKNOWN",
		MissingStage:   "parse",
		SourceDigest:   digestString(source),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	document, err := Parse(source)
	if err != nil {
		parseError := ParseError{Stage: "unknown", Line: 1, Column: 1, Cause: err.Error()}
		if typed, ok := err.(ParseError); ok {
			parseError = typed
		}
		snapshot.Status = document.Status
		snapshot.MissingStage = parseError.Stage
		snapshot.Package = document.Package
		snapshot.Namespace = document.Namespace
		snapshot.SourceDigest = document.SourceDigest
		snapshot.IRDigest = document.IRDigest
		snapshot.Diagnostics = []Diagnostic{{
			Stage:        parseError.Stage,
			Message:      parseError.Cause,
			Position:     Position{Line: parseError.Line, Column: parseError.Column},
			Severity:     SeverityError,
			SourceDigest: document.SourceDigest,
		}}
		return snapshot
	}

	snapshot.Status = document.Status
	snapshot.MissingStage = document.MissingStage
	snapshot.Package = document.Package
	snapshot.Namespace = document.Namespace
	snapshot.SourceDigest = document.SourceDigest
	snapshot.IRDigest = document.IRDigest
	for _, entity := range document.Entities {
		snapshot.Symbols = append(snapshot.Symbols, Symbol{
			Name: entity.Name, Kind: EntitySymbol, Position: entity.Position, Digest: entity.Digest,
		})
	}
	for _, activity := range document.Activities {
		snapshot.Symbols = append(snapshot.Symbols, Symbol{
			Name: activity.Name, Kind: ActivitySymbol, Position: activity.Position, Digest: activity.Digest,
		})
	}
	return snapshot
}

// Validate ensures the tooling projection remains deterministic and non-authorizing.
func (s LanguageSnapshot) Validate() error {
	if s.Status != "BOUND" && s.Status != "UNKNOWN" {
		return fmt.Errorf("snapshot status %q is invalid", s.Status)
	}
	if s.SourceDigest == "" {
		return fmt.Errorf("snapshot source digest is required")
	}
	if !s.NonExecuting || !s.NonAuthorizing {
		return fmt.Errorf("snapshot must remain non-executing and non-authorizing")
	}
	seen := map[string]struct{}{}
	for _, symbol := range s.Symbols {
		if symbol.Name == "" || symbol.Digest == "" {
			return fmt.Errorf("snapshot symbol is incomplete")
		}
		if _, exists := seen[symbol.Name]; exists {
			return fmt.Errorf("snapshot symbol %q is duplicated", symbol.Name)
		}
		seen[symbol.Name] = struct{}{}
	}
	for _, diagnostic := range s.Diagnostics {
		if diagnostic.Stage == "" || diagnostic.Message == "" || diagnostic.SourceDigest != s.SourceDigest {
			return fmt.Errorf("snapshot diagnostic is not source-bound")
		}
		if diagnostic.Severity != SeverityError {
			return fmt.Errorf("snapshot diagnostic severity is invalid")
		}
	}
	if s.Status == "BOUND" && len(s.Diagnostics) != 0 {
		return fmt.Errorf("bound snapshot cannot contain diagnostics")
	}
	if s.Status == "UNKNOWN" && len(s.Diagnostics) == 0 {
		return fmt.Errorf("unknown snapshot must retain a diagnostic")
	}
	return nil
}
