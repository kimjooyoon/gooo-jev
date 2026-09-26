package decision

import "strings"

// ExecutionEnvelopeJEVImprovementReplayFeedbackLSPInput adapts replay feedback
// evidence to an editor-facing projection without authorizing execution.
type ExecutionEnvelopeJEVImprovementReplayFeedbackLSPInput struct {
	Feedback       JEVImprovementReplayFeedback
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection is a deterministic
// diagnostic projection that preserves UNKNOWN provenance and feedback evidence.
type ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection struct {
	Status                  string
	Code                    string
	Severity                string
	Message                 string
	CandidateDigest         string
	ReplayObservationDigest string
	MetricDigest            string
	FeedbackEvidenceDigest  string
	EvidenceDigest          string
	MissingStage             string
	NonExecuting             bool
	NonAuthorizing           bool
}

func digestExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection(projection ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection) (string, error) {
	return Digest(struct {
		Status                  string
		Code                    string
		Severity                string
		Message                 string
		CandidateDigest         string
		ReplayObservationDigest string
		MetricDigest            string
		FeedbackEvidenceDigest  string
		MissingStage             string
	}{
		Status:                  projection.Status,
		Code:                    projection.Code,
		Severity:                 projection.Severity,
		Message:                 projection.Message,
		CandidateDigest:          projection.CandidateDigest,
		ReplayObservationDigest: projection.ReplayObservationDigest,
		MetricDigest:             projection.MetricDigest,
		FeedbackEvidenceDigest:  projection.FeedbackEvidenceDigest,
		MissingStage:             projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSP projects confirmed
// or refuted replay feedback while keeping the result non-executing and
// non-authorizing for editor consumers.
func ProjectExecutionEnvelopeJEVImprovementReplayFeedbackLSP(input ExecutionEnvelopeJEVImprovementReplayFeedbackLSPInput) ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection {
	output := ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection{
		Status: "UNKNOWN", Code: "JEV_REPLAY_FEEDBACK_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV replay feedback is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeJEVImprovementReplayFeedbackLSP(output)
	}
	if !input.Feedback.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV replay feedback is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeJEVImprovementReplayFeedbackLSP(output)
	}
	if strings.TrimSpace(input.Feedback.CandidateDigest) == "" || strings.TrimSpace(input.Feedback.ReplayObservationDigest) == "" || strings.TrimSpace(input.Feedback.MetricDigest) == "" || strings.TrimSpace(input.Feedback.FeedbackEvidenceDigest) == "" || strings.TrimSpace(input.Feedback.EvidenceDigest) == "" {
		output.MissingStage = input.Feedback.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "feedback-evidence"
		}
		output.Message = "JEV replay feedback is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeJEVImprovementReplayFeedbackLSP(output)
	}
	output.CandidateDigest = input.Feedback.CandidateDigest
	output.ReplayObservationDigest = input.Feedback.ReplayObservationDigest
	output.MetricDigest = input.Feedback.MetricDigest
	output.FeedbackEvidenceDigest = input.Feedback.FeedbackEvidenceDigest
	switch input.Feedback.Status {
	case jevReplayFeedbackConfirmed:
		output.Status = "confirmed"
		output.Code = "JEV_REPLAY_FEEDBACK_CONFIRMED"
		output.Severity = "info"
		output.Message = "JEV replay feedback confirmed the candidate observation"
	case jevReplayFeedbackRefuted:
		output.Status = "refuted"
		output.Code = "JEV_REPLAY_FEEDBACK_REFUTED"
		output.Severity = "error"
		output.Message = "JEV replay feedback refuted the candidate observation"
	default:
		output.MissingStage = input.Feedback.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "feedback-status"
		}
		output.Message = "JEV replay feedback is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeJEVImprovementReplayFeedbackLSP(output)
	}
	output.EvidenceDigest = input.Feedback.EvidenceDigest
	return finalizeExecutionEnvelopeJEVImprovementReplayFeedbackLSP(output)
}

func finalizeExecutionEnvelopeJEVImprovementReplayFeedbackLSP(output ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection) ExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection {
	digest, err := digestExecutionEnvelopeJEVImprovementReplayFeedbackLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_REPLAY_FEEDBACK_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV replay feedback is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
