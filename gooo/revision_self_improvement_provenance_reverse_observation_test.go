package gooo

import "testing"

func provenanceReverseObservationInputs(t *testing.T) (RevisionSelfImprovementProvenanceApplicationFeedbackObservation, RevisionSelfImprovementReverseObservation) {
	t.Helper()
	disposition, application := provenanceApplicationFeedbackInputs(t)
	feedback, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	lifecycle, generation := selfImprovementReverseObservationInputs(t)
	reverse, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementReverseGeneration() error = %v", err) }
	return feedback, reverse
}

func TestObserveRevisionSelfImprovementProvenanceReverseObservationBindsEvidence(t *testing.T) {
	feedback, reverse := provenanceReverseObservationInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverse)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReverseObservation() error = %v", err) }
	if result.Status != "BOUND" || result.ReverseSignal != "reverse-observed" || result.ProvenanceReverseSignal != "provenance-reverse-mismatch" || !result.ExactIRMatch || !result.ExactStructureMatch { t.Fatalf("unexpected provenance reverse observation: %#v", result) }
	if result.ApplicationFeedbackDigest != feedback.ObservationDigest || result.ReverseObservationDigest != reverse.ObservationDigest { t.Fatalf("reverse observation lost feedback links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestObserveRevisionSelfImprovementProvenanceReverseObservationRetainsSourceFailure(t *testing.T) {
	feedback, reverse := provenanceReverseObservationInputs(t)
	reverse.ProposedSourceDigest = digestString("other-proposed-source")
	reverse.ObservationDigest = digestRevisionSelfImprovementReverseObservation(reverse)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverse)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceReverseObservation() error = nil, want proposed source link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-reverse-proposed-source-link" { t.Fatalf("unexpected unknown proposed-source observation: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceReverseObservationRetainsReverseFailure(t *testing.T) {
	feedback, reverse := provenanceReverseObservationInputs(t)
	reverse.Status = "UNKNOWN"
	reverse.MissingStage = "revision-self-improvement-reverse-structure"
	reverse.ReverseSignal = "reverse-observation-unknown"
	reverse.ObservationDigest = digestRevisionSelfImprovementReverseObservation(reverse)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverse)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceReverseObservation() error = nil, want reverse failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-reverse-observed" { t.Fatalf("unexpected unknown reverse observation: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceReverseObservationIsDeterministic(t *testing.T) {
	feedback, reverse := provenanceReverseObservationInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverse)
	if err != nil { t.Fatalf("first ObserveRevisionSelfImprovementProvenanceReverseObservation() error = %v", err) }
	second, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverse)
	if err != nil { t.Fatalf("second ObserveRevisionSelfImprovementProvenanceReverseObservation() error = %v", err) }
	if first.ObservationDigest != second.ObservationDigest { t.Fatal("same feedback and reverse evidence produced different observation digest") }
}
