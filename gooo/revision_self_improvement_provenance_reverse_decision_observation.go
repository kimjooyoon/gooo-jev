package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceReverseDecisionObservation binds observed
// reverse evidence to the next bounded decision without executing a change.
type RevisionSelfImprovementProvenanceReverseDecisionObservation struct {
	Status                    string
	MissingStage              string
	ApplicationFeedbackDigest string
	ReverseObservationDigest  string
	DecisionDigest            string
	PlanDispositionDigest     string
	IterationProvenanceDigest string
	PlanDigest                string
	ApplicationDigest         string
	ReceiptDigest             string
	SourceDigest              string
	ProposedSourceDigest      string
	DecisionProposedSourceDigest string
	GeneratedIRDigest         string
	DecisionGeneratedIRDigest string
	FeedbackSignal            string
	ApplicationSignal         string
	FeedbackApplicationSignal string
	ReverseSignal             string
	ProvenanceReverseSignal   string
	DecisionSignal            string
	DecisionFeedbackSignal    string
	DecisionReason            string
	SignalsAligned            bool
	RequiresObservation       bool
	RequiresReview            bool
	RequiresInspection        bool
	RequiresMeasurement       bool
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ObserveRevisionSelfImprovementProvenanceReverseDecision carries exact
// feedback and reverse evidence into the next decision observation.
func ObserveRevisionSelfImprovementProvenanceReverseDecision(
	feedback RevisionSelfImprovementProvenanceApplicationFeedbackObservation,
	reverse RevisionSelfImprovementProvenanceReverseObservation,
	decision RevisionSelfImprovementDecisionObservation,
) (RevisionSelfImprovementProvenanceReverseDecisionObservation, error) {
	result := RevisionSelfImprovementProvenanceReverseDecisionObservation{
		Status:                       "UNKNOWN",
		MissingStage:                 "revision-self-improvement-provenance-reverse-decision",
		ApplicationFeedbackDigest:    feedback.ObservationDigest,
		ReverseObservationDigest:     reverse.ObservationDigest,
		DecisionDigest:               decision.DecisionDigest,
		PlanDispositionDigest:        feedback.PlanDispositionDigest,
		IterationProvenanceDigest:    feedback.IterationProvenanceDigest,
		PlanDigest:                   feedback.PlanDigest,
		ApplicationDigest:            feedback.ApplicationDigest,
		ReceiptDigest:                feedback.ReceiptDigest,
		SourceDigest:                 feedback.ApplicationSourceDigest,
		ProposedSourceDigest:         feedback.ApplicationProposedSourceDigest,
		DecisionProposedSourceDigest: decision.ProposedSourceDigest,
		GeneratedIRDigest:            reverse.GeneratedIRDigest,
		DecisionGeneratedIRDigest:    decision.GeneratedIRDigest,
		FeedbackSignal:               feedback.FeedbackSignal,
		ApplicationSignal:            feedback.ApplicationSignal,
		FeedbackApplicationSignal:    feedback.FeedbackApplicationSignal,
		ReverseSignal:                reverse.ReverseSignal,
		ProvenanceReverseSignal:      reverse.ProvenanceReverseSignal,
		DecisionSignal:               decision.DecisionSignal,
		DecisionFeedbackSignal:       "provenance-decision-mismatch",
		DecisionReason:               decision.DecisionReason,
		RequiresObservation:          decision.RequiresObservation,
		RequiresReview:               decision.RequiresReview,
		RequiresInspection:           decision.RequiresInspection,
		RequiresMeasurement:          decision.RequiresMeasurement,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setDigest := func() { result.ObservationDigest = digestRevisionSelfImprovementProvenanceReverseDecision(result) }
	setDigest()

	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-feedback"
		setDigest()
		return result, fmt.Errorf("provenance application feedback is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-reverse"
		setDigest()
		return result, fmt.Errorf("provenance reverse observation is not valid: %w", err)
	}
	if err := decision.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-decision"
		setDigest()
		return result, fmt.Errorf("self-improvement decision is not valid: %w", err)
	}
	if reverse.ApplicationFeedbackDigest != feedback.ObservationDigest {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-feedback-link"
		setDigest()
		return result, fmt.Errorf("reverse evidence is not linked to application feedback")
	}
	if decision.SourceDigest != feedback.ApplicationSourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-source-link"
		setDigest()
		return result, fmt.Errorf("decision source is not linked to application feedback source")
	}
	if decision.ProposedSourceDigest != feedback.ApplicationProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-proposed-source-link"
		setDigest()
		return result, fmt.Errorf("decision proposed source is not linked to application feedback")
	}
	if decision.GeneratedIRDigest != reverse.GeneratedIRDigest {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-generated-ir-link"
		setDigest()
		return result, fmt.Errorf("decision generated IR is not linked to reverse evidence")
	}
	if decision.FeedbackSignal != feedback.FeedbackSignal {
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision-feedback-signal-link"
		setDigest()
		return result, fmt.Errorf("decision feedback signal is not linked to application feedback")
	}
	result.SignalsAligned = feedback.SignalsAligned && reverse.SignalsAligned && decision.DecisionSignal == feedback.FeedbackSignal
	if result.SignalsAligned {
		result.DecisionFeedbackSignal = "provenance-decision-aligned"
	} else {
		result.DecisionFeedbackSignal = "provenance-decision-mismatch"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-reverse-decision"
		result.DecisionFeedbackSignal = "provenance-decision-mismatch"
		setDigest()
		return result, fmt.Errorf("provenance reverse decision is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenanceReverseDecisionObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("provenance reverse decision status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound provenance reverse decision has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown provenance reverse decision has no missing stage")
	}
	for name, digest := range map[string]string{
		"application feedback": o.ApplicationFeedbackDigest, "reverse observation": o.ReverseObservationDigest,
		"decision": o.DecisionDigest, "plan disposition": o.PlanDispositionDigest,
		"iteration provenance": o.IterationProvenanceDigest, "plan": o.PlanDigest,
		"application": o.ApplicationDigest, "receipt": o.ReceiptDigest, "source": o.SourceDigest,
		"proposed source": o.ProposedSourceDigest, "decision proposed source": o.DecisionProposedSourceDigest,
		"generated IR": o.GeneratedIRDigest, "decision generated IR": o.DecisionGeneratedIRDigest,
		"observation": o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("provenance reverse decision %s digest is invalid", name)
		}
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "remeasure" && o.FeedbackSignal != "review" && o.FeedbackSignal != "inspect" {
		return fmt.Errorf("provenance reverse decision feedback signal is invalid")
	}
	if o.DecisionSignal != o.FeedbackSignal && o.Status == "BOUND" {
		return fmt.Errorf("provenance reverse decision signal is not linked to feedback")
	}
	if o.ApplicationSignal == "" || o.ReverseSignal == "" || o.ProvenanceReverseSignal == "" || o.DecisionFeedbackSignal == "" {
		return fmt.Errorf("provenance reverse decision signals are incomplete")
	}
	if o.DecisionReason == "" && o.Status == "BOUND" {
		return fmt.Errorf("bound provenance reverse decision has no reason")
	}
	if !o.RequiresObservation {
		return fmt.Errorf("provenance reverse decision must require observation")
	}
	if o.DecisionFeedbackSignal != "provenance-decision-aligned" && o.DecisionFeedbackSignal != "provenance-decision-mismatch" {
		return fmt.Errorf("provenance reverse decision relationship signal is invalid")
	}
	if o.SignalsAligned != (o.DecisionFeedbackSignal == "provenance-decision-aligned") {
		return fmt.Errorf("provenance reverse decision alignment is inconsistent")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("provenance reverse decision must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenanceReverseDecision(o) != o.ObservationDigest {
		return fmt.Errorf("provenance reverse decision digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenanceReverseDecision(o RevisionSelfImprovementProvenanceReverseDecisionObservation) string {
	parts := []string{o.Status, o.MissingStage, o.ApplicationFeedbackDigest, o.ReverseObservationDigest, o.DecisionDigest, o.PlanDispositionDigest, o.IterationProvenanceDigest, o.PlanDigest, o.ApplicationDigest, o.ReceiptDigest, o.SourceDigest, o.ProposedSourceDigest, o.DecisionProposedSourceDigest, o.GeneratedIRDigest, o.DecisionGeneratedIRDigest, o.FeedbackSignal, o.ApplicationSignal, o.FeedbackApplicationSignal, o.ReverseSignal, o.ProvenanceReverseSignal, o.DecisionSignal, o.DecisionFeedbackSignal, o.DecisionReason, strconv.FormatBool(o.SignalsAligned), strconv.FormatBool(o.RequiresObservation), strconv.FormatBool(o.RequiresReview), strconv.FormatBool(o.RequiresInspection), strconv.FormatBool(o.RequiresMeasurement), strconv.FormatBool(o.NonExecuting), strconv.FormatBool(o.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
