package gooo

import (
	"encoding/hex"
	"fmt"
)

// DecisionObservation is an externally supplied, provenance-bound decision result.
type DecisionObservation struct {
	DecisionName  string
	ObservedValue string
	Score         *float64
	Boolean       *bool
	EvidenceDigest string
}

// DecisionAssessment is a validated receipt, not an execution or authorization.
type DecisionAssessment struct {
	Status           string
	MissingStage     string
	DecisionName     string
	DecisionKind     DecisionKind
	DecisionID       string
	ObservedValue    string
	Score            float64
	HasScore         bool
	Boolean          bool
	HasBoolean       bool
	EvidenceDigest   string
	DecisionDigest   string
	AssessmentDigest string
	NonExecuting     bool
	NonAuthorizing   bool
}

// AssessDecision validates an external observation against a typed decision declaration.
func AssessDecision(document DecisionDocumentIR, observation DecisionObservation) (DecisionAssessment, error) {
	assessment := DecisionAssessment{
		Status:         "UNKNOWN",
		MissingStage:   "decision-assessment",
		DecisionName:   observation.DecisionName,
		ObservedValue:  observation.ObservedValue,
		EvidenceDigest: observation.EvidenceDigest,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if err := document.Validate(); err != nil {
		assessment.MissingStage = "decision-validation"
		return assessment, fmt.Errorf("gooo decision assessment: %w", err)
	}

	var declaration *DecisionDecl
	for index := range document.Decisions {
		if document.Decisions[index].Name == observation.DecisionName {
			declaration = &document.Decisions[index]
			break
		}
	}
	if declaration == nil {
		assessment.MissingStage = "decision-resolution"
		return assessment, fmt.Errorf("gooo decision assessment: decision %q is undefined", observation.DecisionName)
	}
	assessment.DecisionKind = declaration.Kind
	assessment.DecisionID = declaration.ID
	assessment.DecisionDigest = declaration.Digest

	if !validDigest(observation.EvidenceDigest) {
		assessment.MissingStage = "decision-evidence"
		return assessment, fmt.Errorf("gooo decision assessment: evidence digest is required")
	}

	switch declaration.Kind {
	case ChoiceDecision:
		if observation.ObservedValue == "" || observation.Score != nil || observation.Boolean != nil {
			assessment.MissingStage = "decision-shape"
			return assessment, fmt.Errorf("gooo decision assessment: choice requires a value only")
		}
	case ScoreDecision:
		if observation.Score == nil || observation.ObservedValue != "" || observation.Boolean != nil {
			assessment.MissingStage = "decision-shape"
			return assessment, fmt.Errorf("gooo decision assessment: score requires a score only")
		}
		if *observation.Score < 0 || *observation.Score > 1 {
			assessment.MissingStage = "decision-score"
			return assessment, fmt.Errorf("gooo decision assessment: score must be between 0 and 1")
		}
		assessment.Score = *observation.Score
		assessment.HasScore = true
	case BooleanDecision:
		if observation.Boolean == nil || observation.ObservedValue != "" || observation.Score != nil {
			assessment.MissingStage = "decision-shape"
			return assessment, fmt.Errorf("gooo decision assessment: boolean requires a boolean only")
		}
		assessment.Boolean = *observation.Boolean
		assessment.HasBoolean = true
	default:
		assessment.MissingStage = "decision-shape"
		return assessment, fmt.Errorf("gooo decision assessment: unsupported decision kind %q", declaration.Kind)
	}

	assessment.Status = "BOUND"
	assessment.MissingStage = ""
	assessment.AssessmentDigest = digestAssessment(assessment)
	return assessment, nil
}

// Validate checks that an assessment remains source and evidence bound.
func (a DecisionAssessment) Validate() error {
	if a.Status != "BOUND" {
		return fmt.Errorf("assessment status must be BOUND")
	}
	if a.MissingStage != "" {
		return fmt.Errorf("assessment missing stage must be empty")
	}
	if a.DecisionName == "" || a.DecisionID == "" || a.DecisionDigest == "" {
		return fmt.Errorf("assessment decision identity is incomplete")
	}
	if a.DecisionKind != ChoiceDecision && a.DecisionKind != ScoreDecision && a.DecisionKind != BooleanDecision {
		return fmt.Errorf("assessment decision kind is invalid")
	}
	if !validDigest(a.EvidenceDigest) || !validDigest(a.AssessmentDigest) {
		return fmt.Errorf("assessment evidence digests are invalid")
	}
	if !a.NonExecuting || !a.NonAuthorizing {
		return fmt.Errorf("assessment must remain non-executing and non-authorizing")
	}
	switch a.DecisionKind {
	case ChoiceDecision:
		if a.ObservedValue == "" || a.HasScore || a.HasBoolean {
			return fmt.Errorf("choice assessment shape is invalid")
		}
	case ScoreDecision:
		if !a.HasScore || a.ObservedValue != "" || a.HasBoolean || a.Score < 0 || a.Score > 1 {
			return fmt.Errorf("score assessment shape is invalid")
		}
	case BooleanDecision:
		if !a.HasBoolean || a.ObservedValue != "" || a.HasScore {
			return fmt.Errorf("boolean assessment shape is invalid")
		}
	}
	if expected := digestAssessment(a); expected != a.AssessmentDigest {
		return fmt.Errorf("assessment digest does not match its fields")
	}
	return nil
}

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func digestAssessment(assessment DecisionAssessment) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%t|%t|%s",
		assessment.DecisionName,
		assessment.DecisionKind,
		assessment.DecisionID,
		assessment.ObservedValue,
		formatScore(assessment.Score, assessment.HasScore),
		assessment.Boolean,
		assessment.HasBoolean,
		assessment.EvidenceDigest,
	))
}

func formatScore(score float64, present bool) string {
	if !present {
		return ""
	}
	return fmt.Sprintf("%.17g", score)
}
