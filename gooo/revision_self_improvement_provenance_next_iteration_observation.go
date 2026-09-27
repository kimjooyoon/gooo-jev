package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceNextIterationObservation records the
// non-executing boundary from a replan to the next iteration evidence.
type RevisionSelfImprovementProvenanceNextIterationObservation struct {
	Status                    string
	MissingStage              string
	ReplanDigest              string
	NextIterationDigest       string
	SourceDigest              string
	NextSourceDigest          string
	IRDigest                  string
	NextIRDigest              string
	HistoryDigest             string
	WindowDigest              string
	FeedbackDigest            string
	CandidateSourceDigest     string
	CandidateProposedSourceDigest string
	CandidateGeneratedIRDigest string
	DecisionSignal            string
	FeedbackSignal            string
	ReplanSignal              string
	IterationHistorySignal    string
	TransitionSignal          string
	SignalsAligned            bool
	ObservationCount          int
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ObserveRevisionSelfImprovementProvenanceNextIteration binds a replan to the
// next iteration boundary without executing or authorizing a revision.
func ObserveRevisionSelfImprovementProvenanceNextIteration(
	replan RevisionSelfImprovementProvenanceReplanObservation,
	nextIteration RevisionSelfImprovementIteration,
) (RevisionSelfImprovementProvenanceNextIterationObservation, error) {
	result := RevisionSelfImprovementProvenanceNextIterationObservation{
		Status:                      "UNKNOWN",
		MissingStage:                "revision-self-improvement-provenance-next-iteration",
		ReplanDigest:                replan.ObservationDigest,
		NextIterationDigest:         nextIteration.IterationDigest,
		SourceDigest:                replan.SourceDigest,
		NextSourceDigest:            nextIteration.SourceDigest,
		IRDigest:                    replan.InputIRDigest,
		NextIRDigest:                nextIteration.IRDigest,
		HistoryDigest:               nextIteration.HistoryDigest,
		WindowDigest:                nextIteration.WindowDigest,
		FeedbackDigest:              nextIteration.FeedbackDigest,
		CandidateSourceDigest:       nextIteration.CandidateSourceDigest,
		CandidateProposedSourceDigest: nextIteration.CandidateProposedSourceDigest,
		CandidateGeneratedIRDigest:  nextIteration.CandidateGeneratedIRDigest,
		DecisionSignal:              replan.DecisionSignal,
		FeedbackSignal:              nextIteration.FeedbackSignal,
		ReplanSignal:                replan.ReplanSignal,
		IterationHistorySignal:      nextIteration.HistorySignal,
		TransitionSignal:            "provenance-next-iteration-mismatch",
		ObservationCount:            nextIteration.ObservationCount,
		NonExecuting:                true,
		NonAuthorizing:              true,
	}
	setDigest := func() { result.ObservationDigest = digestRevisionSelfImprovementProvenanceNextIteration(result) }
	setDigest()

	if err := replan.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-next-iteration-replan"
		setDigest()
		return result, fmt.Errorf("provenance replan is not valid: %w", err)
	}
	if err := nextIteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-next-iteration-iteration"
		setDigest()
		return result, fmt.Errorf("next self-improvement iteration is not valid: %w", err)
	}
	if replan.SourceDigest != nextIteration.SourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-next-iteration-source-link"
		setDigest()
		return result, fmt.Errorf("next iteration source is not linked to the replan source")
	}
	if nextIteration.SourceDigest != nextIteration.CandidateSourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-next-iteration-candidate-link"
		setDigest()
		return result, fmt.Errorf("next iteration source is not linked to its candidate source")
	}
	result.SignalsAligned = replan.SignalsAligned && replan.DecisionSignal == nextIteration.FeedbackSignal
	if result.SignalsAligned {
		result.TransitionSignal = "provenance-next-iteration-aligned"
	} else {
		result.TransitionSignal = "provenance-next-iteration-mismatch"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-next-iteration"
		result.TransitionSignal = "provenance-next-iteration-mismatch"
		setDigest()
		return result, fmt.Errorf("provenance next iteration observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenanceNextIterationObservation) Validate() error {
	if o.Status == "" { return fmt.Errorf("provenance next iteration status is empty") }
	if o.Status == "BOUND" && o.MissingStage != "" { return fmt.Errorf("bound provenance next iteration has a missing stage") }
	if o.Status == "UNKNOWN" && o.MissingStage == "" { return fmt.Errorf("unknown provenance next iteration has no missing stage") }
	for name, digest := range map[string]string{
		"replan": o.ReplanDigest, "next iteration": o.NextIterationDigest, "source": o.SourceDigest,
		"next source": o.NextSourceDigest, "ir": o.IRDigest, "next ir": o.NextIRDigest,
		"history": o.HistoryDigest, "window": o.WindowDigest, "feedback": o.FeedbackDigest,
		"candidate source": o.CandidateSourceDigest, "candidate proposed source": o.CandidateProposedSourceDigest,
		"candidate generated ir": o.CandidateGeneratedIRDigest, "observation": o.ObservationDigest,
	} {
		if !validDigest(digest) { return fmt.Errorf("provenance next iteration %s digest is invalid", name) }
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "remeasure" && o.DecisionSignal != "review" && o.DecisionSignal != "inspect" { return fmt.Errorf("provenance next iteration decision signal is invalid") }
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "remeasure" && o.FeedbackSignal != "review" && o.FeedbackSignal != "inspect" { return fmt.Errorf("provenance next iteration feedback signal is invalid") }
	if o.ReplanSignal != "provenance-replan-aligned" && o.ReplanSignal != "provenance-replan-mismatch" { return fmt.Errorf("provenance next iteration replan signal is invalid") }
	if o.IterationHistorySignal != "stable" && o.IterationHistorySignal != "narrower" && o.IterationHistorySignal != "wider" && o.IterationHistorySignal != "mixed" { return fmt.Errorf("provenance next iteration history signal is invalid") }
	if o.TransitionSignal != "provenance-next-iteration-aligned" && o.TransitionSignal != "provenance-next-iteration-mismatch" { return fmt.Errorf("provenance next iteration transition signal is invalid") }
	if o.SignalsAligned != (o.TransitionSignal == "provenance-next-iteration-aligned") { return fmt.Errorf("provenance next iteration alignment is inconsistent") }
	if o.ObservationCount < 1 { return fmt.Errorf("provenance next iteration observation count must be positive") }
	if !o.NonExecuting || !o.NonAuthorizing { return fmt.Errorf("provenance next iteration must remain non-executing and non-authorizing") }
	if digestRevisionSelfImprovementProvenanceNextIteration(o) != o.ObservationDigest { return fmt.Errorf("provenance next iteration digest does not match its fields") }
	return nil
}

func digestRevisionSelfImprovementProvenanceNextIteration(o RevisionSelfImprovementProvenanceNextIterationObservation) string {
	parts := []string{o.Status, o.MissingStage, o.ReplanDigest, o.NextIterationDigest, o.SourceDigest, o.NextSourceDigest, o.IRDigest, o.NextIRDigest, o.HistoryDigest, o.WindowDigest, o.FeedbackDigest, o.CandidateSourceDigest, o.CandidateProposedSourceDigest, o.CandidateGeneratedIRDigest, o.DecisionSignal, o.FeedbackSignal, o.ReplanSignal, o.IterationHistorySignal, o.TransitionSignal, strconv.FormatBool(o.SignalsAligned), strconv.Itoa(o.ObservationCount), strconv.FormatBool(o.NonExecuting), strconv.FormatBool(o.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
