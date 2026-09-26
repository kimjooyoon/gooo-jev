package decision

import (
	"fmt"
	"strings"
)

// DecisionConfidenceChangePlanReplayFeedbackBindingInput connects a bound
// self-improvement plan to independently reconciled replay feedback.
type DecisionConfidenceChangePlanReplayFeedbackBindingInput struct {
	PlanProvenance     DecisionConfidenceChangePlanProvenanceBinding
	FeedbackReconciliation ImprovementReplayFeedbackReconciliation
	NonAuthorizing     bool
}

// DecisionConfidenceChangePlanReplayFeedbackBinding records confirmed or
// refuted feedback for a plan without claiming product improvement.
type DecisionConfidenceChangePlanReplayFeedbackBinding struct {
	Status                    string
	FeedbackStatus            string
	MissingStage              string
	ProposalDigest            string
	ChangePlanDigest          string
	SourceReference           string
	PlanEvidenceDigest        string
	MetricName                string
	HistoryDigest             string
	ReplayEvidenceDigest      string
	FeedbackEvidenceDigest    string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

// BindDecisionConfidenceChangePlanReplayFeedback requires both provenance and
// feedback reconciliation before recording a bound feedback disposition.
func BindDecisionConfidenceChangePlanReplayFeedback(input DecisionConfidenceChangePlanReplayFeedbackBindingInput) DecisionConfidenceChangePlanReplayFeedbackBinding {
	output := DecisionConfidenceChangePlanReplayFeedbackBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.PlanProvenance.NonAuthorizing || !input.FeedbackReconciliation.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.PlanProvenance.NonExecuting || !input.FeedbackReconciliation.NonExecuting {
		output.NonExecuting = false
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.PlanProvenance.Validate(); err != nil {
		output.MissingStage = "change-plan-provenance"
		return output
	}
	switch input.FeedbackReconciliation.Status {
	case "review":
		output.Status = "review"
		output.MissingStage = "feedback-reconciliation"
		return output
	case "confirmed", "refuted":
		// Continue only for a reconciled disposition.
	default:
		output.MissingStage = "feedback-reconciliation"
		return output
	}
	feedback := input.FeedbackReconciliation
	if strings.TrimSpace(feedback.MetricName) == "" ||
		strings.TrimSpace(feedback.HistoryDigest) == "" ||
		strings.TrimSpace(feedback.LatestSummaryDigest) == "" ||
		strings.TrimSpace(feedback.ReplayEvidenceDigest) == "" ||
		strings.TrimSpace(feedback.EvidenceDigest) == "" {
		output.MissingStage = "feedback-evidence"
		return output
	}
	evidenceDigest, err := Digest(struct {
		PlanEvidenceDigest     string
		ChangePlanDigest       string
		FeedbackStatus         string
		MetricName             string
		HistoryDigest          string
		ReplayEvidenceDigest   string
		FeedbackEvidenceDigest string
	}{
		PlanEvidenceDigest:     input.PlanProvenance.EvidenceDigest,
		ChangePlanDigest:       input.PlanProvenance.ChangePlanDigest,
		FeedbackStatus:         feedback.Status,
		MetricName:             feedback.MetricName,
		HistoryDigest:          feedback.HistoryDigest,
		ReplayEvidenceDigest:   feedback.ReplayEvidenceDigest,
		FeedbackEvidenceDigest: feedback.EvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "plan-feedback-evidence"
		return output
	}
	output.Status = "bound-" + feedback.Status
	output.FeedbackStatus = feedback.Status
	output.ProposalDigest = input.PlanProvenance.ProposalDigest
	output.ChangePlanDigest = input.PlanProvenance.ChangePlanDigest
	output.SourceReference = input.PlanProvenance.SourceReference
	output.PlanEvidenceDigest = input.PlanProvenance.EvidenceDigest
	output.MetricName = feedback.MetricName
	output.HistoryDigest = feedback.HistoryDigest
	output.ReplayEvidenceDigest = feedback.ReplayEvidenceDigest
	output.FeedbackEvidenceDigest = feedback.EvidenceDigest
	output.EvidenceDigest = evidenceDigest
	return output
}

func (binding DecisionConfidenceChangePlanReplayFeedbackBinding) Validate() error {
	if (binding.Status != "bound-confirmed" && binding.Status != "bound-refuted") ||
		(binding.FeedbackStatus != "confirmed" && binding.FeedbackStatus != "refuted") ||
		!binding.NonExecuting || !binding.NonAuthorizing ||
		strings.TrimSpace(binding.ProposalDigest) == "" ||
		strings.TrimSpace(binding.ChangePlanDigest) == "" ||
		strings.TrimSpace(binding.SourceReference) == "" ||
		strings.TrimSpace(binding.PlanEvidenceDigest) == "" ||
		strings.TrimSpace(binding.MetricName) == "" ||
		strings.TrimSpace(binding.HistoryDigest) == "" ||
		strings.TrimSpace(binding.ReplayEvidenceDigest) == "" ||
		strings.TrimSpace(binding.FeedbackEvidenceDigest) == "" ||
		strings.TrimSpace(binding.EvidenceDigest) == "" ||
		strings.TrimSpace(binding.MissingStage) != "" {
		return fmt.Errorf("change plan replay feedback binding is incomplete")
	}
	expected, err := Digest(struct {
		PlanEvidenceDigest     string
		ChangePlanDigest       string
		FeedbackStatus         string
		MetricName             string
		HistoryDigest          string
		ReplayEvidenceDigest   string
		FeedbackEvidenceDigest string
	}{
		PlanEvidenceDigest:     binding.PlanEvidenceDigest,
		ChangePlanDigest:       binding.ChangePlanDigest,
		FeedbackStatus:         binding.FeedbackStatus,
		MetricName:             binding.MetricName,
		HistoryDigest:          binding.HistoryDigest,
		ReplayEvidenceDigest:   binding.ReplayEvidenceDigest,
		FeedbackEvidenceDigest: binding.FeedbackEvidenceDigest,
	})
	if err != nil {
		return err
	}
	if expected != binding.EvidenceDigest {
		return fmt.Errorf("change plan replay feedback binding evidence digest mismatch")
	}
	return nil
}
