package decision

import (
	"strings"
	"testing"
)

func TestProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPBound(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput{
		Binding: ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding{
			Status: "bound", BindingDigest: "binding-digest", NonExecuting: true, NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if got.Status != "bound" || got.Code != "JEV_REVISION_FEEDBACK_BOUND" || got.Severity != "info" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing || got.BindingDigest != "binding-digest" {
		t.Fatalf("missing boundary %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPPreservesUnknownStage(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput{
		Binding: ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding{
			Status: "UNKNOWN", MissingStage: "feedback-aggregation", NonExecuting: true, NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.Code != "JEV_REVISION_FEEDBACK_UNKNOWN" || got.MissingStage != "feedback-aggregation" {
		t.Fatalf("got %+v", got)
	}
	if !strings.Contains(got.Message, "feedback-aggregation") || got.EvidenceDigest == "" {
		t.Fatalf("missing diagnostic evidence %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPRejectsExecutingBinding(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput{
		Binding: ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding{
			Status: "bound", BindingDigest: "binding-digest", NonExecuting: false, NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "execution-boundary" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput{
		NonAuthorizing: false,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
