package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceReverseObservation preserves the bridge
// from application feedback to exact reverse observation evidence.
type RevisionSelfImprovementProvenanceReverseObservation struct {
	Status                        string
	MissingStage                  string
	ApplicationFeedbackDigest     string
	ReverseObservationDigest      string
	PlanDispositionDigest         string
	IterationProvenanceDigest     string
	PlanDigest                    string
	ApplicationDigest             string
	ReceiptDigest                 string
	ProposedSourceDigest          string
	ReverseProposedSourceDigest   string
	ProposedIRDigest              string
	ReverseProposedIRDigest       string
	GeneratedIRDigest             string
	ReverseGeneratedIRDigest      string
	StructureDigest               string
	ReverseStructureDigest        string
	FeedbackSignal                string
	ApplicationSignal             string
	FeedbackApplicationSignal     string
	ReverseSignal                 string
	ProvenanceReverseSignal       string
	SignalsAligned                bool
	ExactIRMatch                  bool
	ExactStructureMatch           bool
	ObservationDigest             string
	NonExecuting                  bool
	NonAuthorizing                bool
}

// ObserveRevisionSelfImprovementProvenanceReverseObservation binds feedback
// to exact reverse evidence without executing or authorizing a revision.
func ObserveRevisionSelfImprovementProvenanceReverseObservation(
	feedback RevisionSelfImprovementProvenanceApplicationFeedbackObservation,
	reverse RevisionSelfImprovementReverseObservation,
) (RevisionSelfImprovementProvenanceReverseObservation, error) {
	result := RevisionSelfImprovementProvenanceReverseObservation{
		Status:                      "UNKNOWN",
		MissingStage:                "revision-self-improvement-provenance-reverse-observation",
		ApplicationFeedbackDigest:   feedback.ObservationDigest,
		ReverseObservationDigest:    reverse.ObservationDigest,
		PlanDispositionDigest:       feedback.PlanDispositionDigest,
		IterationProvenanceDigest:   feedback.IterationProvenanceDigest,
		PlanDigest:                  feedback.PlanDigest,
		ApplicationDigest:           feedback.ApplicationDigest,
		ReceiptDigest:               feedback.ReceiptDigest,
		ProposedSourceDigest:        feedback.ApplicationProposedSourceDigest,
		ReverseProposedSourceDigest: reverse.ProposedSourceDigest,
		ProposedIRDigest:            feedback.ApplicationProposedIRDigest,
		ReverseProposedIRDigest:     reverse.ProposedIRDigest,
		GeneratedIRDigest:           reverse.GeneratedIRDigest,
		ReverseGeneratedIRDigest:    reverse.ObservedGeneratedIRDigest,
		StructureDigest:             reverse.StructureDigest,
		ReverseStructureDigest:      reverse.ObservedStructureDigest,
		FeedbackSignal:              feedback.FeedbackSignal,
		ApplicationSignal:           feedback.ApplicationSignal,
		FeedbackApplicationSignal:   feedback.FeedbackApplicationSignal,
		ReverseSignal:               reverse.ReverseSignal,
		ExactIRMatch:                reverse.ExactIRMatch,
		ExactStructureMatch:         reverse.ExactStructureMatch,
		NonExecuting:                true,
		NonAuthorizing:              true,
	}
	setDigest := func() { result.ObservationDigest = digestRevisionSelfImprovementProvenanceReverseObservation(result) }
	setDigest()

	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-reverse-feedback"
		setDigest()
		return result, fmt.Errorf("provenance application feedback is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-reverse-observation-input"
		setDigest()
		return result, fmt.Errorf("reverse observation is not valid: %w", err)
	}
	if reverse.Status != "BOUND" || reverse.ReverseSignal != "reverse-observed" {
		result.MissingStage = "revision-self-improvement-provenance-reverse-observed"
		setDigest()
		return result, fmt.Errorf("reverse observation is not a bounded exact observation")
	}
	if feedback.ApplicationProposedSourceDigest != reverse.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-reverse-proposed-source-link"
		setDigest()
		return result, fmt.Errorf("application proposed source is not linked to reverse observation")
	}
	if feedback.ApplicationProposedIRDigest != reverse.ProposedIRDigest {
		result.MissingStage = "revision-self-improvement-provenance-reverse-proposed-ir-link"
		setDigest()
		return result, fmt.Errorf("application proposed IR is not linked to reverse observation")
	}
	result.SignalsAligned = feedback.SignalsAligned && reverse.ExactIRMatch && reverse.ExactStructureMatch
	if result.SignalsAligned {
		result.ProvenanceReverseSignal = "provenance-reverse-aligned"
	} else {
		result.ProvenanceReverseSignal = "provenance-reverse-mismatch"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-reverse-observation"
		result.ProvenanceReverseSignal = "provenance-reverse-mismatch"
		setDigest()
		return result, fmt.Errorf("provenance reverse observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenanceReverseObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("provenance reverse observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound provenance reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown provenance reverse observation has no missing stage")
	}
	for name, digest := range map[string]string{
		"application feedback": o.ApplicationFeedbackDigest,
		"reverse observation": o.ReverseObservationDigest,
		"plan disposition": o.PlanDispositionDigest,
		"iteration provenance": o.IterationProvenanceDigest,
		"plan": o.PlanDigest,
		"application": o.ApplicationDigest,
		"receipt": o.ReceiptDigest,
		"proposed source": o.ProposedSourceDigest,
		"reverse proposed source": o.ReverseProposedSourceDigest,
		"proposed IR": o.ProposedIRDigest,
		"reverse proposed IR": o.ReverseProposedIRDigest,
		"generated IR": o.GeneratedIRDigest,
		"reverse generated IR": o.ReverseGeneratedIRDigest,
		"structure": o.StructureDigest,
		"reverse structure": o.ReverseStructureDigest,
		"observation": o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("provenance reverse observation %s digest is invalid", name)
		}
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "remeasure" && o.FeedbackSignal != "review" && o.FeedbackSignal != "inspect" {
		return fmt.Errorf("provenance reverse observation feedback signal is invalid")
	}
	if o.ApplicationSignal == "" {
		return fmt.Errorf("provenance reverse observation application signal is empty")
	}
	if o.FeedbackApplicationSignal != "provenance-application-aligned" && o.FeedbackApplicationSignal != "provenance-application-mismatch" {
		return fmt.Errorf("provenance reverse observation feedback relationship is invalid")
	}
	if o.ReverseSignal != "reverse-observed" && o.ReverseSignal != "reverse-observation-unknown" && o.ReverseSignal != "reverse-source-mismatch" && o.ReverseSignal != "reverse-ir-mismatch" && o.ReverseSignal != "reverse-structure-mismatch" {
		return fmt.Errorf("provenance reverse observation signal is invalid")
	}
	if o.ProvenanceReverseSignal != "provenance-reverse-aligned" && o.ProvenanceReverseSignal != "provenance-reverse-mismatch" {
		return fmt.Errorf("provenance reverse observation relationship signal is invalid")
	}
	if o.SignalsAligned != (o.ProvenanceReverseSignal == "provenance-reverse-aligned") {
		return fmt.Errorf("provenance reverse observation alignment is inconsistent")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("provenance reverse observation must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceReverseObservation(o) != o.ObservationDigest {
		return fmt.Errorf("provenance reverse observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceReverseObservation(o RevisionSelfImprovementProvenanceReverseObservation) string {
	parts := []string{o.Status, o.MissingStage, o.ApplicationFeedbackDigest, o.ReverseObservationDigest, o.PlanDispositionDigest, o.IterationProvenanceDigest, o.PlanDigest, o.ApplicationDigest, o.ReceiptDigest, o.ProposedSourceDigest, o.ReverseProposedSourceDigest, o.ProposedIRDigest, o.ReverseProposedIRDigest, o.GeneratedIRDigest, o.ReverseGeneratedIRDigest, o.StructureDigest, o.ReverseStructureDigest, o.FeedbackSignal, o.ApplicationSignal, o.FeedbackApplicationSignal, o.ReverseSignal, o.ProvenanceReverseSignal, strconv.FormatBool(o.SignalsAligned), strconv.FormatBool(o.ExactIRMatch), strconv.FormatBool(o.ExactStructureMatch), strconv.FormatBool(o.NonExecuting), strconv.FormatBool(o.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
