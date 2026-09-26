package gooo

import "fmt"

// RevisionGenerationAssessment binds a revision assessment to generation and
// reverse observation of the proposed source. It never executes or authorizes.
type RevisionGenerationAssessment struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	ProposedSourceDigest  string
	ApplicationDigest     string
	AssessmentDigest      string
	GenerationSourceDigest string
	GeneratedSourceDigest string
	GeneratedIRDigest     string
	StructureDigest       string
	StructureMatch        bool
	ExactSourceMatch      bool
	SourceMatchClass      string
	GenerationDigest      string
	NonExecuting          bool
	NonAuthorizing        bool
}

// AssessRevisionGeneration links the proposed source digest to a successful
// generation receipt and preserves reverse-observation failures as UNKNOWN.
func AssessRevisionGeneration(assessment RevisionApplicationAssessment, generation GenerationReceipt) (RevisionGenerationAssessment, error) {
	result := RevisionGenerationAssessment{
		Status:                 "UNKNOWN",
		MissingStage:           "revision-generation-assessment",
		SourceDigest:           assessment.SourceDigest,
		ProposedSourceDigest:   assessment.ProposedSourceDigest,
		ApplicationDigest:      assessment.ApplicationDigest,
		AssessmentDigest:       assessment.AssessmentDigest,
		GenerationSourceDigest: generation.SourceDigest,
		GeneratedSourceDigest:  generation.GeneratedSourceDigest,
		GeneratedIRDigest:      generation.GeneratedIRDigest,
		StructureDigest:        generation.StructureDigest,
		StructureMatch:         generation.StructureMatch,
		ExactSourceMatch:       generation.ExactSourceMatch,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	setGenerationDigest := func() {
		result.GenerationDigest = digestRevisionGenerationAssessment(result)
	}
	setGenerationDigest()

	if err := assessment.Validate(); err != nil {
		result.MissingStage = "revision-generation-assessment-application"
		setGenerationDigest()
		return result, fmt.Errorf("revision application assessment is not valid: %w", err)
	}
	if err := validateGenerationReceipt(generation); err != nil {
		result.MissingStage = "revision-generation-assessment-generation"
		setGenerationDigest()
		return result, fmt.Errorf("generation receipt is not valid: %w", err)
	}
	if assessment.ProposedSourceDigest != generation.SourceDigest {
		result.MissingStage = "revision-generation-assessment-link"
		setGenerationDigest()
		return result, fmt.Errorf("proposed source digest does not match generation source digest")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	if generation.ExactSourceMatch {
		result.SourceMatchClass = "exact"
	} else {
		result.SourceMatchClass = "canonical"
	}
	setGenerationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-generation-assessment"
		setGenerationDigest()
		return result, fmt.Errorf("revision generation assessment is not valid: %w", err)
	}
	return result, nil
}

func validateGenerationReceipt(generation GenerationReceipt) error {
	if generation.Status != "BOUND" || generation.MissingStage != "" {
		return fmt.Errorf("generation receipt is not bound")
	}
	if generation.SourceDigest == "" || generation.GeneratedSourceDigest == "" || generation.GeneratedIRDigest == "" || generation.StructureDigest == "" {
		return fmt.Errorf("generation receipt is missing a digest")
	}
	if !generation.StructureMatch {
		return fmt.Errorf("generation structure does not match reverse observation")
	}
	if !generation.NonExecuting || !generation.NonAuthorizing {
		return fmt.Errorf("generation receipt must remain non-executing and non-authorizing")
	}
	return nil
}

func (a RevisionGenerationAssessment) Validate() error {
	if a.Status == "" {
		return fmt.Errorf("revision generation assessment status is empty")
	}
	if a.Status == "BOUND" && a.MissingStage != "" {
		return fmt.Errorf("bound revision generation assessment has a missing stage")
	}
	if a.Status == "UNKNOWN" && a.MissingStage == "" {
		return fmt.Errorf("unknown revision generation assessment has no missing stage")
	}
	if a.Status == "BOUND" && a.SourceMatchClass == "" {
		return fmt.Errorf("bound revision generation assessment has no source match class")
	}
	if !a.NonExecuting || !a.NonAuthorizing {
		return fmt.Errorf("revision generation assessment must remain non-executing and non-authorizing")
	}
	if digestRevisionGenerationAssessment(a) != a.GenerationDigest {
		return fmt.Errorf("revision generation assessment digest does not match its fields")
	}
	return nil
}

func digestRevisionGenerationAssessment(assessment RevisionGenerationAssessment) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t",
		assessment.Status,
		assessment.MissingStage,
		assessment.SourceDigest,
		assessment.ProposedSourceDigest,
		assessment.ApplicationDigest,
		assessment.AssessmentDigest,
		assessment.GenerationSourceDigest,
		assessment.GeneratedSourceDigest,
		assessment.GeneratedIRDigest,
		assessment.StructureDigest,
		assessment.SourceMatchClass,
		assessment.NonExecuting,
		assessment.NonAuthorizing,
	))
}