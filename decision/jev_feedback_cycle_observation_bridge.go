package decision

import "strings"

// ExecutionEnvelopeJEVImprovementFeedbackCycleInput feeds validated feedback
// aggregation back into the existing improvement-cycle observation contract.
type ExecutionEnvelopeJEVImprovementFeedbackCycleInput struct {
	DeclarationDigest        string
	IRDigest                 string
	GenerationDigest         string
	ReverseObservationDigest string
	MetricDigest             string
	ChangePlanDigest         string
	Feedback                 JEVImprovementFeedbackAggregation
	NonAuthorizing           bool
}

// ObserveJEVImprovementFeedbackCycle returns an existing cycle observation so
// reverse feedback can become the next cycle input without execution or
// authorization.
func ObserveJEVImprovementFeedbackCycle(input ExecutionEnvelopeJEVImprovementFeedbackCycleInput) JEVImprovementCycleObservation {
	output := JEVImprovementCycleObservation{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	stages := []struct {
		name  string
		value string
	}{
		{"declaration", input.DeclarationDigest},
		{"ir", input.IRDigest},
		{"generation", input.GenerationDigest},
		{"reverse-observation", input.ReverseObservationDigest},
		{"metric", input.MetricDigest},
		{"change-plan", input.ChangePlanDigest},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	if err := input.Feedback.Validate(); err != nil {
		output.MissingStage = input.Feedback.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "feedback-aggregation"
		}
		return output
	}
	boundChangePlanDigest, err := Digest(struct {
		ChangePlanDigest      string
		FeedbackEvidenceDigest string
	}{
		ChangePlanDigest:       input.ChangePlanDigest,
		FeedbackEvidenceDigest: input.Feedback.EvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "change-plan-feedback-binding"
		return output
	}
	var lspCode string
	switch input.Feedback.Status {
	case "stable-for-review":
		lspCode = "jev.change-plan.ready-for-external-apply-review"
	case "needs-revision":
		lspCode = "jev.change-plan.abort"
	case "hold":
		lspCode = "jev.change-plan.hold"
	default:
		output.MissingStage = "feedback-disposition"
		return output
	}
	output.DeclarationDigest = input.DeclarationDigest
	output.IRDigest = input.IRDigest
	output.GenerationDigest = input.GenerationDigest
	output.ReverseObservationDigest = input.ReverseObservationDigest
	output.MetricDigest = input.MetricDigest
	output.ChangePlanDigest = boundChangePlanDigest
	output.DispositionStatus = input.Feedback.Status
	output.LSPCode = lspCode
	output.EvidenceDigest = digestJEVImprovementCycleObservation(output.DeclarationDigest, output.IRDigest, output.GenerationDigest, output.ReverseObservationDigest, output.MetricDigest, output.ChangePlanDigest, output.DispositionStatus, output.LSPCode)
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "feedback-cycle-observation-evidence"
		output.EvidenceDigest = ""
		return output
	}
	output.Status = input.Feedback.Status
	return output
}
