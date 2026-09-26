package gooo

import "fmt"

// RevisionCandidateObservation links a materialized candidate to an evidence
// comparison without authorizing its application or execution.
type RevisionCandidateObservation struct {
	Status                string
	MissingStage          string
	ObservationDigest     string
	MaterializationDigest string
	SourceDigest          string
	CandidateDigest       string
	FeedbackDigest        string
	DirectionDigest       string
	AssessmentDigest      string
	EvidenceDigest        string
	ObservationClass      string
	ReviewSignal          string
	MaterializationSignal string
	CandidateAvailable    bool
	ResultDigest          string
	NonExecuting          bool
	NonAuthorizing        bool
}

// ObserveRevisionCandidateForImprovement connects a bounded comparison to a
// materialized candidate and returns only an evidence signal.
func ObserveRevisionCandidateForImprovement(observation RevisionImprovementObservation, materialization RevisionCandidateMaterialization) (RevisionCandidateObservation, error) {
	result := RevisionCandidateObservation{
		Status:                "UNKNOWN",
		MissingStage:          "revision-candidate-observation",
		ObservationDigest:     observation.ObservationDigest,
		MaterializationDigest: materialization.MaterializationDigest,
		SourceDigest:          materialization.SourceDigest,
		CandidateDigest:       materialization.CandidateDigest,
		FeedbackDigest:        materialization.FeedbackDigest,
		DirectionDigest:       materialization.DirectionDigest,
		AssessmentDigest:      materialization.AssessmentDigest,
		EvidenceDigest:        materialization.EvidenceDigest,
		ObservationClass:      observation.ObservationClass,
		ReviewSignal:          observation.ReviewSignal,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestRevisionCandidateObservation(result)
	}
	setResultDigest()

	if err := observation.Validate(); err != nil {
		result.MissingStage = "revision-candidate-observation-comparison"
		setResultDigest()
		return result, fmt.Errorf("revision improvement observation is not valid: %w", err)
	}
	if err := materialization.Validate(); err != nil {
		result.MissingStage = "revision-candidate-observation-materialization"
		setResultDigest()
		return result, fmt.Errorf("revision candidate materialization is not valid: %w", err)
	}
	if observation.BaselineSourceDigest != materialization.SourceDigest {
		result.MissingStage = "revision-candidate-observation-link"
		setResultDigest()
		return result, fmt.Errorf("candidate materialization source is not linked to the observation baseline")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.CandidateAvailable = true
	if observation.ReviewSignal == "inspect" {
		result.MaterializationSignal = "inspect-candidate"
	} else {
		result.MaterializationSignal = "candidate-observed"
	}
	setResultDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-candidate-observation"
		setResultDigest()
		return result, fmt.Errorf("revision candidate observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionCandidateObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision candidate observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision candidate observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision candidate observation has no missing stage")
	}
	if o.Status == "BOUND" && (!o.CandidateAvailable || o.MaterializationSignal == "") {
		return fmt.Errorf("bound revision candidate observation is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision candidate observation must remain non-executing and non-authorizing")
	}
	if digestRevisionCandidateObservation(o) != o.ResultDigest {
		return fmt.Errorf("revision candidate observation digest does not match its fields")
	}
	return nil
}

func digestRevisionCandidateObservation(observation RevisionCandidateObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.ObservationDigest,
		observation.MaterializationDigest,
		observation.SourceDigest,
		observation.CandidateDigest,
		observation.FeedbackDigest,
		observation.DirectionDigest,
		observation.AssessmentDigest,
		observation.EvidenceDigest,
		observation.ObservationClass,
		observation.ReviewSignal,
		observation.MaterializationSignal,
		observation.CandidateAvailable,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}