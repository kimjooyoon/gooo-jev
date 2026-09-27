package gooo

import "testing"

func provenanceReverseDecisionInputs(t *testing.T) (RevisionSelfImprovementProvenanceApplicationFeedbackObservation, RevisionSelfImprovementProvenanceReverseObservation, RevisionSelfImprovementDecisionObservation) {
	t.Helper()
	feedback, reverseObservation := provenanceReverseObservationInputs(t)
	reverse, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverseObservation)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReverseObservation() error = %v", err) }
	iteration, outcome, feedbackWindow := selfImprovementDecisionInputs(t)
	decision, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedbackWindow)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementDecision() error = %v", err) }
	return feedback, reverse, decision
}

func TestObserveRevisionSelfImprovementProvenanceReverseDecisionBindsNextAction(t *testing.T) {
	feedback, reverse, decision := provenanceReverseDecisionInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decision)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReverseDecision() error = %v", err) }
	if result.Status != "BOUND" || result.DecisionSignal != decision.DecisionSignal || result.DecisionFeedbackSignal != "provenance-decision-mismatch" || result.ProvenanceReverseSignal != reverse.ProvenanceReverseSignal { t.Fatalf("unexpected provenance reverse decision: %#v", result) }
	if result.ApplicationFeedbackDigest != feedback.ObservationDigest || result.ReverseObservationDigest != reverse.ObservationDigest || result.DecisionDigest != decision.DecisionDigest { t.Fatalf("provenance reverse decision lost links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestObserveRevisionSelfImprovementProvenanceReverseDecisionRetainsSourceFailure(t *testing.T) {
	feedback, reverse, decision := provenanceReverseDecisionInputs(t)
	decision.SourceDigest = digestString("other-source")
	decision.DecisionDigest = digestRevisionSelfImprovementDecision(decision)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decision)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceReverseDecision() error = nil, want source link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-reverse-decision-source-link" { t.Fatalf("unexpected unknown source decision: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceReverseDecisionRetainsIRFailure(t *testing.T) {
	feedback, reverse, decision := provenanceReverseDecisionInputs(t)
	decision.GeneratedIRDigest = digestString("other-generated-ir")
	decision.DecisionDigest = digestRevisionSelfImprovementDecision(decision)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decision)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceReverseDecision() error = nil, want generated IR link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-reverse-decision-generated-ir-link" { t.Fatalf("unexpected unknown IR decision: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceReverseDecisionIsDeterministic(t *testing.T) {
	feedback, reverse, decision := provenanceReverseDecisionInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decision)
	if err != nil { t.Fatalf("first ObserveRevisionSelfImprovementProvenanceReverseDecision() error = %v", err) }
	second, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decision)
	if err != nil { t.Fatalf("second ObserveRevisionSelfImprovementProvenanceReverseDecision() error = %v", err) }
	if first.ObservationDigest != second.ObservationDigest { t.Fatal("same reverse evidence and decision produced different observation digest") }
}
