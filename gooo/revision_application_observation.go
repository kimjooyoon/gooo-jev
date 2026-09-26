package gooo

import "fmt"

// RevisionApplicationObservation links a plan observation to an in-memory
// application receipt without executing, persisting, or authorizing it.
type RevisionApplicationObservation struct {
	Status                    string
	MissingStage              string
	PlanObservationDigest     string
	PlanDigest                string
	ApplicationDigest         string
	ReceiptDigest             string
	SourceDigest              string
	ProposedSourceDigest      string
	InputIRDigest             string
	ProposedIRDigest          string
	CandidateDigest           string
	EditDigest                string
	ApplicationSignal         string
	ApplicationObserved       bool
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ObserveRevisionApplication binds an application receipt to its observed
// candidate plan while keeping the result informational only.
func ObserveRevisionApplication(planObservation RevisionCandidatePlanObservation, receipt RevisionApplicationReceipt) (RevisionApplicationObservation, error) {
	result := RevisionApplicationObservation{
		Status:                "UNKNOWN",
		MissingStage:          "revision-application-observation",
		PlanObservationDigest: planObservation.ObservationDigest,
		PlanDigest:            receipt.PlanDigest,
		ApplicationDigest:     receipt.ApplicationDigest,
		ReceiptDigest:         receipt.ReceiptDigest,
		SourceDigest:          receipt.SourceDigest,
		ProposedSourceDigest:  receipt.ProposedSourceDigest,
		InputIRDigest:         receipt.InputIRDigest,
		ProposedIRDigest:      receipt.ProposedIRDigest,
		CandidateDigest:       receipt.CandidateDigest,
		EditDigest:            receipt.EditDigest,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionApplicationObservation(result)
	}
	setObservationDigest()

	if err := planObservation.Validate(); err != nil {
		result.MissingStage = "revision-application-observation-plan"
		setObservationDigest()
		return result, fmt.Errorf("revision candidate plan observation is not valid: %w", err)
	}
	if err := receipt.Validate(); err != nil {
		result.MissingStage = "revision-application-observation-receipt"
		setObservationDigest()
		return result, fmt.Errorf("revision application receipt is not valid: %w", err)
	}
	if planObservation.PlanDigest != receipt.PlanDigest ||
		planObservation.SourceDigest != receipt.SourceDigest ||
		planObservation.CandidateDigest != receipt.CandidateDigest ||
		planObservation.EditDigest != receipt.EditDigest {
		result.MissingStage = "revision-application-observation-link"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan observation and application receipt are not linked")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ApplicationObserved = true
	result.ApplicationSignal = "application-observed"
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-application-observation"
		setObservationDigest()
		return result, fmt.Errorf("revision application observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionApplicationObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision application observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision application observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision application observation has no missing stage")
	}
	if o.Status == "BOUND" && (!o.ApplicationObserved || o.ApplicationSignal == "") {
		return fmt.Errorf("bound revision application observation is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision application observation must remain non-executing and non-authorizing")
	}
	if digestRevisionApplicationObservation(o) != o.ObservationDigest {
		return fmt.Errorf("revision application observation digest does not match its fields")
	}
	return nil
}

func digestRevisionApplicationObservation(observation RevisionApplicationObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.PlanObservationDigest,
		observation.PlanDigest,
		observation.ApplicationDigest,
		observation.ReceiptDigest,
		observation.SourceDigest,
		observation.ProposedSourceDigest,
		observation.InputIRDigest,
		observation.ProposedIRDigest,
		observation.CandidateDigest,
		observation.EditDigest,
		observation.ApplicationSignal,
		observation.ApplicationObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}