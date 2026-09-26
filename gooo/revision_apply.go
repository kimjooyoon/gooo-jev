package gooo

import (
	"fmt"
	"strings"
)

// SourceEdit is a single-line, source-bound proposed edit with an exclusive end.
type SourceEdit struct {
	Start       Position
	End         Position
	Replacement string
	Digest      string
}

// RevisionApplication is a non-executing result of applying and reverse-observing one edit.
type RevisionApplication struct {
	Status               string
	MissingStage         string
	SourceDigest         string
	ProposedSourceDigest string
	InputIRDigest        string
	ProposedIRDigest     string
	CandidateDigest      string
	EditDigest           string
	ApplicationDigest    string
	ProposedSource       string
	NonExecuting         bool
	NonAuthorizing       bool
}

// ApplyRevision applies one bounded edit only when the expected source digest matches.
func ApplyRevision(source, expectedSourceDigest string, candidate RevisionCandidate, edit SourceEdit) (RevisionApplication, error) {
	application := RevisionApplication{
		Status:          "UNKNOWN",
		MissingStage:    "revision-application",
		SourceDigest:    expectedSourceDigest,
		CandidateDigest: candidate.CandidateDigest,
		EditDigest:      edit.Digest,
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	if !validDigest(expectedSourceDigest) || digestString(source) != expectedSourceDigest {
		application.MissingStage = "revision-source-precondition"
		return application, fmt.Errorf("gooo revision application: source digest precondition failed")
	}
	if err := candidate.Validate(); err != nil {
		application.MissingStage = "revision-candidate"
		return application, fmt.Errorf("gooo revision application: candidate: %w", err)
	}
	if edit.Digest == "" || digestSourceEdit(edit) != edit.Digest {
		application.MissingStage = "revision-edit-evidence"
		return application, fmt.Errorf("gooo revision application: edit digest does not match its fields")
	}
	if edit.Start.Line != edit.End.Line || strings.Contains(edit.Replacement, "\n") {
		application.MissingStage = "revision-range"
		return application, fmt.Errorf("gooo revision application: edit must stay on one line")
	}
	start, err := positionOffset(source, edit.Start)
	if err != nil {
		application.MissingStage = "revision-range"
		return application, fmt.Errorf("gooo revision application: start: %w", err)
	}
	end, err := positionOffset(source, edit.End)
	if err != nil || end < start {
		application.MissingStage = "revision-range"
		if err == nil {
			err = fmt.Errorf("end precedes start")
		}
		return application, fmt.Errorf("gooo revision application: range: %w", err)
	}

	inputIRDigest, err := parseRevisionSource(source)
	if err != nil {
		application.MissingStage = "revision-source-parse"
		return application, fmt.Errorf("gooo revision application: source: %w", err)
	}
	proposed := source[:start] + edit.Replacement + source[end:]
	application.ProposedSource = proposed
	application.ProposedSourceDigest = digestString(proposed)
	proposedIRDigest, err := parseRevisionSource(proposed)
	if err != nil {
		application.MissingStage = "revision-reverse-observation"
		return application, fmt.Errorf("gooo revision application: proposed source: %w", err)
	}

	application.Status = "BOUND"
	application.MissingStage = ""
	application.InputIRDigest = inputIRDigest
	application.ProposedIRDigest = proposedIRDigest
	application.ApplicationDigest = digestRevisionApplication(application)
	return application, nil
}

// Validate checks the source, edit, candidate, and reverse-observed IR evidence.
func (a RevisionApplication) Validate() error {
	if a.Status != "BOUND" {
		return fmt.Errorf("application status must be BOUND")
	}
	if a.MissingStage != "" {
		return fmt.Errorf("application missing stage must be empty")
	}
	if !a.NonExecuting || !a.NonAuthorizing {
		return fmt.Errorf("application must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": a.SourceDigest, "proposed source": a.ProposedSourceDigest,
		"input IR": a.InputIRDigest, "proposed IR": a.ProposedIRDigest,
		"candidate": a.CandidateDigest, "edit": a.EditDigest, "application": a.ApplicationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("application %s digest is invalid", name)
		}
	}
	if a.ProposedSource == "" {
		return fmt.Errorf("application proposed source is required")
	}
	if expected := digestRevisionApplication(a); expected != a.ApplicationDigest {
		return fmt.Errorf("application digest does not match its fields")
	}
	return nil
}

func parseRevisionSource(source string) (string, error) {
	for _, line := range strings.Split(source, "\n") {
		text := strings.TrimSpace(line)
		if strings.HasPrefix(text, "decision ") {
			document, err := ParseDecision(source)
			if err != nil {
				return "", err
			}
			return document.IRDigest, nil
		}
	}
	document, err := Parse(source)
	if err != nil {
		return "", err
	}
	return document.IRDigest, nil
}

func digestSourceEdit(edit SourceEdit) string {
	return digestString(fmt.Sprintf("%d:%d-%d:%d|%s", edit.Start.Line, edit.Start.Column, edit.End.Line, edit.End.Column, edit.Replacement))
}

func digestRevisionApplication(application RevisionApplication) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
		application.SourceDigest,
		application.ProposedSourceDigest,
		application.InputIRDigest,
		application.ProposedIRDigest,
		application.CandidateDigest,
		application.EditDigest,
		application.ProposedSource,
	))
}

func positionOffset(source string, position Position) (int, error) {
	if position.Line < 1 || position.Column < 1 {
		return 0, fmt.Errorf("position must be positive")
	}
	lines := strings.SplitAfter(source, "\n")
	if position.Line > len(lines) {
		return 0, fmt.Errorf("line %d is outside source", position.Line)
	}
	offset := 0
	for _, line := range lines[:position.Line-1] {
		offset += len(line)
	}
	line := lines[position.Line-1]
	contentLength := len(line)
	if strings.HasSuffix(line, "\n") {
		contentLength--
	}
	if position.Column > contentLength+1 {
		return 0, fmt.Errorf("column %d is outside line %d", position.Column, position.Line)
	}
	return offset + position.Column - 1, nil
}
