package gooo

import "fmt"

// RevisionImprovementObservation compares two bound evidence summaries without
// claiming that either revision is beneficial or authorizing a change.
type RevisionImprovementObservation struct {
	Status                         string
	MissingStage                   string
	BaselineSummaryDigest          string
	CandidateSummaryDigest         string
	BaselineSourceDigest           string
	CandidateSourceDigest          string
	BaselineProposedSourceDigest   string
	CandidateProposedSourceDigest  string
	BaselineProposedIRDigest       string
	CandidateProposedIRDigest      string
	BaselineGeneratedSourceDigest  string
	CandidateGeneratedSourceDigest string
	BaselineGeneratedIRDigest      string
	CandidateGeneratedIRDigest     string
	BaselineStructureDigest        string
	CandidateStructureDigest       string
	BaselineStageCount             int
	CandidateStageCount            int
	StageCountDelta                int
	SourceStable                   bool
	ProposedIRChanged              bool
	GeneratedSourceChanged         bool
	GeneratedIRChanged             bool
	StructureStable                bool
	ObservationClass               string
	ReviewSignal                   string
	ObservationDigest              string
	NonExecuting                   bool
	NonAuthorizing                 bool
}

// ObserveRevisionImprovement compares exact provenance fields from two bound
// summaries and reports only observed differences, never an improvement claim.
func ObserveRevisionImprovement(baseline, candidate RevisionEvidenceSummary) (RevisionImprovementObservation, error) {
	observation := RevisionImprovementObservation{
		Status:                         "UNKNOWN",
		MissingStage:                   "revision-improvement-observation",
		BaselineSummaryDigest:          baseline.SummaryDigest,
		CandidateSummaryDigest:         candidate.SummaryDigest,
		BaselineSourceDigest:            baseline.SourceDigest,
		CandidateSourceDigest:           candidate.SourceDigest,
		BaselineProposedSourceDigest:    baseline.ProposedSourceDigest,
		CandidateProposedSourceDigest:   candidate.ProposedSourceDigest,
		BaselineProposedIRDigest:        baseline.ProposedIRDigest,
		CandidateProposedIRDigest:       candidate.ProposedIRDigest,
		BaselineGeneratedSourceDigest:   baseline.GeneratedSourceDigest,
		CandidateGeneratedSourceDigest:  candidate.GeneratedSourceDigest,
		BaselineGeneratedIRDigest:       baseline.GeneratedIRDigest,
		CandidateGeneratedIRDigest:      candidate.GeneratedIRDigest,
		BaselineStructureDigest:         baseline.StructureDigest,
		CandidateStructureDigest:        candidate.StructureDigest,
		BaselineStageCount:              baseline.EvidenceStageCount,
		CandidateStageCount:             candidate.EvidenceStageCount,
		StageCountDelta:                 candidate.EvidenceStageCount - baseline.EvidenceStageCount,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	setObservationDigest := func() {
		observation.ObservationDigest = digestRevisionImprovementObservation(observation)
	}
	setObservationDigest()

	if err := baseline.Validate(); err != nil {
		observation.MissingStage = "revision-improvement-observation-baseline"
		setObservationDigest()
		return observation, fmt.Errorf("baseline revision evidence summary is not valid: %w", err)
	}
	if err := candidate.Validate(); err != nil {
		observation.MissingStage = "revision-improvement-observation-candidate"
		setObservationDigest()
		return observation, fmt.Errorf("candidate revision evidence summary is not valid: %w", err)
	}

	observation.Status = "BOUND"
	observation.MissingStage = ""
	observation.SourceStable = baseline.SourceDigest == candidate.SourceDigest
	observation.ProposedIRChanged = baseline.ProposedIRDigest != candidate.ProposedIRDigest
	observation.GeneratedSourceChanged = baseline.GeneratedSourceDigest != candidate.GeneratedSourceDigest
	observation.GeneratedIRChanged = baseline.GeneratedIRDigest != candidate.GeneratedIRDigest
	observation.StructureStable = baseline.StructureDigest == candidate.StructureDigest
	observation.ObservationClass = revisionObservationClass(observation)
	observation.ReviewSignal = revisionObservationReviewSignal(observation)
	setObservationDigest()
	if err := observation.Validate(); err != nil {
		observation.Status = "UNKNOWN"
		observation.MissingStage = "revision-improvement-observation"
		setObservationDigest()
		return observation, fmt.Errorf("revision improvement observation is not valid: %w", err)
	}
	return observation, nil
}

func (o RevisionImprovementObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision improvement observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision improvement observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision improvement observation has no missing stage")
	}
	if o.BaselineStageCount < 0 || o.CandidateStageCount < 0 || o.StageCountDelta != o.CandidateStageCount-o.BaselineStageCount {
		return fmt.Errorf("revision improvement observation stage counts are inconsistent")
	}
	if o.Status == "BOUND" && o.ObservationClass == "" {
		return fmt.Errorf("bound revision improvement observation has no class")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision improvement observation must remain non-executing and non-authorizing")
	}
	if digestRevisionImprovementObservation(o) != o.ObservationDigest {
		return fmt.Errorf("revision improvement observation digest does not match its fields")
	}
	return nil
}

func revisionObservationClass(observation RevisionImprovementObservation) string {
	if !observation.SourceStable {
		return "source-change"
	}
	if observation.ProposedIRChanged || observation.GeneratedIRChanged {
		return "semantic-change"
	}
	if observation.GeneratedSourceChanged {
		return "generated-source-change"
	}
	if !observation.StructureStable {
		return "structure-change"
	}
	return "stable"
}

func revisionObservationReviewSignal(observation RevisionImprovementObservation) string {
	if observation.ProposedIRChanged || observation.GeneratedIRChanged || !observation.StructureStable {
		return "inspect"
	}
	if !observation.SourceStable || observation.GeneratedSourceChanged {
		return "review"
	}
	return "observe"
}

func digestRevisionImprovementObservation(observation RevisionImprovementObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%d|%t|%t|%t|%t|%t|%s|%s|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.BaselineSummaryDigest,
		observation.CandidateSummaryDigest,
		observation.BaselineSourceDigest,
		observation.CandidateSourceDigest,
		observation.BaselineProposedSourceDigest,
		observation.CandidateProposedSourceDigest,
		observation.BaselineProposedIRDigest,
		observation.CandidateProposedIRDigest,
		observation.BaselineGeneratedSourceDigest,
		observation.CandidateGeneratedSourceDigest,
		observation.BaselineGeneratedIRDigest,
		observation.CandidateGeneratedIRDigest,
		observation.BaselineStructureDigest,
		observation.CandidateStructureDigest,
		observation.BaselineStageCount,
		observation.CandidateStageCount,
		observation.StageCountDelta,
		observation.SourceStable,
		observation.ProposedIRChanged,
		observation.GeneratedSourceChanged,
		observation.GeneratedIRChanged,
		observation.StructureStable,
		observation.ObservationClass,
		observation.ReviewSignal,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}