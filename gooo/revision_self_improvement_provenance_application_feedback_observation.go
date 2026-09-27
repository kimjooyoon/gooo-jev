package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceApplicationFeedbackObservation preserves the
// handoff from provenance plan disposition to observed application feedback.
type RevisionSelfImprovementProvenanceApplicationFeedbackObservation struct {
	Status                          string
	MissingStage                    string
	PlanDispositionDigest           string
	IterationProvenanceDigest       string
	PlanObservationDigest           string
	ApplicationObservationDigest    string
	PlanDigest                      string
	ApplicationDigest               string
	ReceiptDigest                   string
	CandidateSourceDigest           string
	ApplicationSourceDigest         string
	CandidateProposedSourceDigest   string
	ApplicationProposedSourceDigest string
	CandidateGeneratedIRDigest      string
	ApplicationProposedIRDigest     string
	FeedbackSignal                  string
	ApplicationSignal               string
	DispositionSignal               string
	FeedbackApplicationSignal       string
	SignalsAligned                  bool
	ApplicationObserved             bool
	ObservationDigest               string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

// ObserveRevisionSelfImprovementProvenanceApplicationFeedback binds observed
// application feedback without executing, authorizing, or inferring improvement.
func ObserveRevisionSelfImprovementProvenanceApplicationFeedback(
	planDisposition RevisionSelfImprovementProvenancePlanDispositionObservation,
	application RevisionSelfImprovementApplicationObservation,
) (RevisionSelfImprovementProvenanceApplicationFeedbackObservation, error) {
	result := RevisionSelfImprovementProvenanceApplicationFeedbackObservation{
		Status:                          "UNKNOWN",
		MissingStage:                    "revision-self-improvement-provenance-application-feedback",
		PlanDispositionDigest:           planDisposition.ObservationDigest,
		IterationProvenanceDigest:       planDisposition.IterationProvenanceDigest,
		PlanObservationDigest:           application.PlanObservationDigest,
		ApplicationObservationDigest:    application.ObservationDigest,
		PlanDigest:                      application.PlanDigest,
		ApplicationDigest:               application.ApplicationDigest,
		ReceiptDigest:                  application.ReceiptDigest,
		CandidateSourceDigest:           planDisposition.CandidateSourceDigest,
		ApplicationSourceDigest:         application.SourceDigest,
		CandidateProposedSourceDigest:   planDisposition.CandidateProposedSourceDigest,
		ApplicationProposedSourceDigest: application.ProposedSourceDigest,
		CandidateGeneratedIRDigest:      planDisposition.CandidateGeneratedIRDigest,
		ApplicationProposedIRDigest:     application.ProposedIRDigest,
		FeedbackSignal:                  application.FeedbackSignal,
		ApplicationSignal:               application.ApplicationSignal,
		DispositionSignal:               planDisposition.DispositionSignal,
		ApplicationObserved:             application.ApplicationObserved,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	setDigest := func() { result.ObservationDigest = digestRevisionSelfImprovementProvenanceApplicationFeedback(result) }
	setDigest()

	if err := planDisposition.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-plan-disposition"
		setDigest()
		return result, fmt.Errorf("provenance plan disposition is not valid: %w", err)
	}
	if err := application.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-application"
		setDigest()
		return result, fmt.Errorf("application observation is not valid: %w", err)
	}
	if planDisposition.PlanObservationDigest != application.PlanObservationDigest {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-plan-observation-link"
		setDigest()
		return result, fmt.Errorf("provenance plan observation and application observation are not linked")
	}
	if planDisposition.PlanDigest != application.PlanDigest {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-plan-link"
		setDigest()
		return result, fmt.Errorf("provenance plan and application plan are not linked")
	}
	if planDisposition.CandidateSourceDigest != application.SourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-source-link"
		setDigest()
		return result, fmt.Errorf("provenance candidate source and application source are not linked")
	}
	if planDisposition.CandidateProposedSourceDigest != application.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-proposed-source-link"
		setDigest()
		return result, fmt.Errorf("provenance candidate proposed source and application proposed source are not linked")
	}
	if !application.ApplicationObserved {
		result.MissingStage = "revision-self-improvement-provenance-application-feedback-observed"
		setDigest()
		return result, fmt.Errorf("application feedback was not observed")
	}

	expected := "provenance-plan-unknown"
	switch planDisposition.DispositionSignal {
	case "provenance-observe-plan":
		expected = "observe"
	case "provenance-inspect-plan":
		expected = "inspect"
	case "provenance-plan-mismatch":
		expected = "provenance-plan-mismatch"
	}
	result.SignalsAligned = planDisposition.SignalsAligned && expected == application.FeedbackSignal
	if result.SignalsAligned {
		result.FeedbackApplicationSignal = "provenance-application-aligned"
	} else {
		result.FeedbackApplicationSignal = "provenance-application-mismatch"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-application-feedback"
		result.FeedbackApplicationSignal = "provenance-application-mismatch"
		setDigest()
		return result, fmt.Errorf("provenance application feedback is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenanceApplicationFeedbackObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("provenance application feedback status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound provenance application feedback has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown provenance application feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"plan disposition": o.PlanDispositionDigest, "iteration provenance": o.IterationProvenanceDigest,
		"plan observation": o.PlanObservationDigest, "application observation": o.ApplicationObservationDigest,
		"plan": o.PlanDigest, "application": o.ApplicationDigest, "receipt": o.ReceiptDigest,
		"candidate source": o.CandidateSourceDigest, "application source": o.ApplicationSourceDigest,
		"candidate proposed source": o.CandidateProposedSourceDigest,
		"application proposed source": o.ApplicationProposedSourceDigest,
		"candidate generated IR": o.CandidateGeneratedIRDigest, "application proposed IR": o.ApplicationProposedIRDigest,
		"observation": o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("provenance application feedback %s digest is invalid", name)
		}
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "remeasure" && o.FeedbackSignal != "review" && o.FeedbackSignal != "inspect" {
		return fmt.Errorf("provenance application feedback signal is invalid")
	}
	if o.ApplicationSignal == "" {
		return fmt.Errorf("provenance application feedback application signal is empty")
	}
	if o.DispositionSignal != "provenance-observe-plan" && o.DispositionSignal != "provenance-inspect-plan" && o.DispositionSignal != "provenance-plan-mismatch" && o.DispositionSignal != "provenance-plan-unknown" {
		return fmt.Errorf("provenance application feedback disposition signal is invalid")
	}
	if o.FeedbackApplicationSignal != "provenance-application-aligned" && o.FeedbackApplicationSignal != "provenance-application-mismatch" {
		return fmt.Errorf("provenance application feedback relationship signal is invalid")
	}
	if o.SignalsAligned != (o.FeedbackApplicationSignal == "provenance-application-aligned") {
		return fmt.Errorf("provenance application feedback alignment is inconsistent")
	}
	if !o.ApplicationObserved {
		return fmt.Errorf("provenance application feedback has not been observed")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("provenance application feedback must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceApplicationFeedback(o) != o.ObservationDigest {
		return fmt.Errorf("provenance application feedback digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceApplicationFeedback(o RevisionSelfImprovementProvenanceApplicationFeedbackObservation) string {
	parts := []string{o.Status, o.MissingStage, o.PlanDispositionDigest, o.IterationProvenanceDigest, o.PlanObservationDigest, o.ApplicationObservationDigest, o.PlanDigest, o.ApplicationDigest, o.ReceiptDigest, o.CandidateSourceDigest, o.ApplicationSourceDigest, o.CandidateProposedSourceDigest, o.ApplicationProposedSourceDigest, o.CandidateGeneratedIRDigest, o.ApplicationProposedIRDigest, o.FeedbackSignal, o.ApplicationSignal, o.DispositionSignal, o.FeedbackApplicationSignal, strconv.FormatBool(o.SignalsAligned), strconv.FormatBool(o.ApplicationObserved), strconv.FormatBool(o.NonExecuting), strconv.FormatBool(o.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
