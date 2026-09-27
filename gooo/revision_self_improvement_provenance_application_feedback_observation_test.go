package gooo

import "testing"

func provenanceApplicationFeedbackInputs(t *testing.T) (RevisionSelfImprovementProvenancePlanDispositionObservation, RevisionSelfImprovementApplicationObservation) {
	t.Helper()
	iteration, plan := provenancePlanDispositionInputs(t)
	disposition, err := ObserveRevisionSelfImprovementProvenancePlanDisposition(iteration, plan)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenancePlanDisposition() error = %v", err) }
	selfImprovementIteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	application, err := ObserveRevisionSelfImprovementApplication(selfImprovementIteration, planObservation, applicationObservation)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementApplication() error = %v", err) }
	return disposition, application
}

func TestObserveRevisionSelfImprovementProvenanceApplicationFeedbackBindsObservation(t *testing.T) {
	disposition, application := provenanceApplicationFeedbackInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	if result.Status != "BOUND" || !result.ApplicationObserved || result.FeedbackApplicationSignal != "provenance-application-mismatch" || result.DispositionSignal != disposition.DispositionSignal { t.Fatalf("unexpected provenance application feedback: %#v", result) }
	if result.PlanDispositionDigest != disposition.ObservationDigest || result.ApplicationObservationDigest != application.ObservationDigest { t.Fatalf("application feedback lost provenance links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestObserveRevisionSelfImprovementProvenanceApplicationFeedbackRetainsPlanFailure(t *testing.T) {
	disposition, application := provenanceApplicationFeedbackInputs(t)
	application.PlanDigest = digestString("other-plan")
	application.ObservationDigest = digestRevisionSelfImprovementApplicationObservation(application)
	result, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = nil, want plan link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-application-feedback-plan-link" { t.Fatalf("unexpected unknown plan-link feedback: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceApplicationFeedbackRetainsSourceFailure(t *testing.T) {
	disposition, application := provenanceApplicationFeedbackInputs(t)
	application.SourceDigest = digestString("other-source")
	application.ObservationDigest = digestRevisionSelfImprovementApplicationObservation(application)
	result, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = nil, want source link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-application-feedback-source-link" { t.Fatalf("unexpected unknown source-link feedback: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceApplicationFeedbackIsDeterministic(t *testing.T) {
	disposition, application := provenanceApplicationFeedbackInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err != nil { t.Fatalf("first ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	second, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err != nil { t.Fatalf("second ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	if first.ObservationDigest != second.ObservationDigest { t.Fatal("same provenance disposition and application produced different observation digest") }
}
