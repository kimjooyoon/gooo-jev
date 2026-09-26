package decision

import "testing"

func validChangePlanReplayFeedbackBinding(t *testing.T, kind FeedbackKind) DecisionConfidenceChangePlanReplayFeedbackBinding {
	t.Helper()
	plan := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              validChangePlanForProvenance(t),
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          true,
	})
	feedback := makeImprovementReplayFeedback(t, kind)
	reconciliation := ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, kind),
		FeedbackHistory: makeReplayFeedbackHistory(t, feedback),
		NonAuthorizing:  true,
	})
	return BindDecisionConfidenceChangePlanReplayFeedback(DecisionConfidenceChangePlanReplayFeedbackBindingInput{
		PlanProvenance:         plan,
		FeedbackReconciliation: reconciliation,
		NonAuthorizing:         true,
	})
}

func TestBindDecisionConfidenceChangePlanReplayFeedbackBindsConfirmedAndRefuted(t *testing.T) {
	for _, kind := range []FeedbackKind{FeedbackConfirmed, FeedbackRefuted} {
		binding := validChangePlanReplayFeedbackBinding(t, kind)
		if binding.Status != "bound-"+string(kind) || binding.FeedbackStatus != string(kind) || binding.MissingStage != "" {
			t.Fatalf("unexpected bound feedback for %q: %+v", kind, binding)
		}
		if binding.EvidenceDigest == "" || !binding.NonExecuting || !binding.NonAuthorizing {
			t.Fatalf("bound feedback lost evidence or safety boundary: %+v", binding)
		}
		if err := binding.Validate(); err != nil {
			t.Fatalf("bound feedback should validate: %v", err)
		}
	}
}

func TestBindDecisionConfidenceChangePlanReplayFeedbackKeepsReview(t *testing.T) {
	plan := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              validChangePlanForProvenance(t),
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          true,
	})
	unknown := makeImprovementReplayFeedback(t, FeedbackUnknown)
	reconciliation := ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
		FeedbackHistory: makeReplayFeedbackHistory(t, unknown),
		NonAuthorizing:  true,
	})
	binding := BindDecisionConfidenceChangePlanReplayFeedback(DecisionConfidenceChangePlanReplayFeedbackBindingInput{
		PlanProvenance:         plan,
		FeedbackReconciliation: reconciliation,
		NonAuthorizing:         true,
	})
	if binding.Status != "review" || binding.MissingStage != "feedback-reconciliation" {
		t.Fatalf("unexpected review binding: %+v", binding)
	}
}

func TestBindDecisionConfidenceChangePlanReplayFeedbackFailsClosed(t *testing.T) {
	binding := validChangePlanReplayFeedbackBinding(t, FeedbackConfirmed)
	binding.EvidenceDigest = "tampered"
	if err := binding.Validate(); err == nil {
		t.Fatal("tampered binding unexpectedly validated")
	}

	binding = validChangePlanReplayFeedbackBinding(t, FeedbackConfirmed)
	binding.NonAuthorizing = false
	output := BindDecisionConfidenceChangePlanReplayFeedback(DecisionConfidenceChangePlanReplayFeedbackBindingInput{
		PlanProvenance:         BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
			ChangePlan:              validChangePlanForProvenance(t),
			DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
			NonAuthorizing:          true,
		}),
		FeedbackReconciliation: ImprovementReplayFeedbackReconciliation{
			Status: "confirmed", MetricName: "metric", HistoryDigest: "history", LatestSummaryDigest: "summary",
			ReplayEvidenceDigest: "replay", EvidenceDigest: "feedback", NonExecuting: true, NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if output.Status != "bound-confirmed" {
		t.Fatalf("unexpected valid output after local tamper fixture: %+v", output)
	}
}

func TestBindDecisionConfidenceChangePlanReplayFeedbackRejectsTamperedFeedbackAndAuthorization(t *testing.T) {
	plan := BindDecisionConfidenceChangePlanProvenance(DecisionConfidenceChangePlanProvenanceBindingInput{
		ChangePlan:              validChangePlanForProvenance(t),
		DeclarationIRGeneration: validDeclarationIRGenerationBinding(),
		NonAuthorizing:          true,
	})
	feedback := makeImprovementReplayFeedback(t, FeedbackConfirmed)
	reconciliation := ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
		ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
		FeedbackHistory: makeReplayFeedbackHistory(t, feedback),
		NonAuthorizing:  true,
	})
	reconciliation.EvidenceDigest = "tampered"
	output := BindDecisionConfidenceChangePlanReplayFeedback(DecisionConfidenceChangePlanReplayFeedbackBindingInput{
		PlanProvenance:         plan,
		FeedbackReconciliation: reconciliation,
		NonAuthorizing:         true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "feedback-reconciliation" {
		t.Fatalf("unexpected tampered feedback output: %+v", output)
	}

	output = BindDecisionConfidenceChangePlanReplayFeedback(DecisionConfidenceChangePlanReplayFeedbackBindingInput{
		PlanProvenance:         plan,
		FeedbackReconciliation: ReconcileImprovementReplayFeedback(ImprovementReplayFeedbackReconciliationInput{
			ReplayFeedback:  makeReplayFeedbackBinding(t, FeedbackConfirmed),
			FeedbackHistory: makeReplayFeedbackHistory(t, feedback),
			NonAuthorizing:  true,
		}),
		NonAuthorizing: false,
	})
	if output.Status != "UNKNOWN" || output.NonAuthorizing || output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected authorization output: %+v", output)
	}
}
