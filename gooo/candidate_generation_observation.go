package gooo

import "fmt"

type CandidateGenerationObservation struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	GeneratedSourceDigest string
	GeneratedIRDigest     string
	StructureDigest       string
	CandidateDigest       string
	MaterializationDigest string
	ExactSourceMatch      bool
	StructureMatch        bool
	ReverseObserved       bool
	ObservationDigest     string
	NonExecuting          bool
	NonAuthorizing        bool
}

func ObserveCandidateGeneration(document DocumentIR, materialization RevisionCandidateMaterialization) (CandidateGenerationObservation, error) {
	observation := CandidateGenerationObservation{
		Status:                "UNKNOWN",
		MissingStage:          "candidate-generation",
		SourceDigest:          document.SourceDigest,
		CandidateDigest:       materialization.CandidateDigest,
		MaterializationDigest: materialization.MaterializationDigest,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	if err := materialization.Validate(); err != nil {
		observation.MissingStage = "candidate-generation-materialization"
		return observation, fmt.Errorf("gooo candidate generation: materialization: %w", err)
	}
	if document.SourceDigest != materialization.SourceDigest {
		observation.MissingStage = "candidate-generation-source"
		return observation, fmt.Errorf("gooo candidate generation: document and materialization source digests do not match")
	}
	generation, err := Generate(document)
	if err != nil {
		observation.MissingStage = "candidate-generation-" + generation.MissingStage
		return observation, fmt.Errorf("gooo candidate generation: %w", err)
	}
	observation.Status = "BOUND"
	observation.MissingStage = ""
	observation.SourceDigest = generation.SourceDigest
	observation.GeneratedSourceDigest = generation.GeneratedSourceDigest
	observation.GeneratedIRDigest = generation.GeneratedIRDigest
	observation.StructureDigest = generation.StructureDigest
	observation.ExactSourceMatch = generation.ExactSourceMatch
	observation.StructureMatch = generation.StructureMatch
	observation.ReverseObserved = true
	observation.ObservationDigest = digestCandidateGenerationObservation(observation)
	return observation, nil
}

func (o CandidateGenerationObservation) Validate() error {
	if o.Status != "BOUND" {
		return fmt.Errorf("candidate generation observation status must be BOUND")
	}
	if o.MissingStage != "" {
		return fmt.Errorf("candidate generation observation missing stage must be empty")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("candidate generation observation must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"source": o.SourceDigest, "generated source": o.GeneratedSourceDigest,
		"generated IR": o.GeneratedIRDigest, "structure": o.StructureDigest,
		"candidate": o.CandidateDigest, "materialization": o.MaterializationDigest,
		"observation": o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("candidate generation observation %s digest is invalid", name)
		}
	}
	if !o.StructureMatch || !o.ReverseObserved {
		return fmt.Errorf("candidate generation observation lacks reverse structure evidence")
	}
	if digestCandidateGenerationObservation(o) != o.ObservationDigest {
		return fmt.Errorf("candidate generation observation digest does not match its fields")
	}
	return nil
}

func digestCandidateGenerationObservation(observation CandidateGenerationObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.SourceDigest,
		observation.GeneratedSourceDigest,
		observation.GeneratedIRDigest,
		observation.StructureDigest,
		observation.CandidateDigest,
		observation.ExactSourceMatch,
		observation.StructureMatch,
		observation.ReverseObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}
