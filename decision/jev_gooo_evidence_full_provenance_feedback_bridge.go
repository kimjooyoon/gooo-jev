package decision

import "strings"

// ExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeInput connects a
// candidate admission decision to the existing append-only replay ledger.
type ExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeInput struct {
	Previous             JEVImprovementReplayFeedbackLedger
	CandidateGate        ExecutionEnvelopeGoooEvidenceFullProvenanceCandidateGate
	ReplayObservationDigest string
	NonAuthorizing       bool
}

// AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge converts
// ADMIT/HOLD into existing CONFIRMED/REFUTED feedback and appends it without
// applying the candidate or granting capability.
func AppendExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridge(input ExecutionEnvelopeGoooEvidenceFullProvenanceFeedbackBridgeInput) JEVImprovementReplayFeedbackLedger {
	unknown := JEVImprovementReplayFeedbackLedger{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.CandidateGate.NonAuthorizing {
		if !input.NonAuthorizing {
			unknown.NonAuthorizing = false
		}
		unknown.MissingStage = "authorization-boundary"
		return unknown
	}
	if err := input.CandidateGate.Validate(); err != nil {
		unknown.MissingStage = "candidate-gate-validation"
		return unknown
	}
	if strings.TrimSpace(input.ReplayObservationDigest) == "" {
		unknown.MissingStage = "replay-observation"
		return unknown
	}
	feedbackStatus := jevReplayFeedbackRefuted
	if input.CandidateGate.AdmissionStatus == "ADMIT" {
		feedbackStatus = jevReplayFeedbackConfirmed
	}
	feedback := JEVImprovementReplayFeedback{
		Status:                  feedbackStatus,
		FeedbackKind:            feedbackStatus,
		CandidateDigest:         input.CandidateGate.CandidateSourceDigest,
		ReplayObservationDigest: input.ReplayObservationDigest,
		MetricDigest:            input.CandidateGate.MetricEvidenceDigest,
		NonExecuting:            true,
		NonAuthorizing:          true,
	}
	feedback.FeedbackEvidenceDigest = digestJEVImprovementReplayFeedbackEvidence(
		feedback.ReplayObservationDigest,
		feedback.MetricDigest,
		feedback.FeedbackKind,
	)
	feedback.EvidenceDigest = digestJEVImprovementReplayFeedback(
		feedback.Status,
		feedback.FeedbackKind,
		feedback.CandidateDigest,
		feedback.ReplayObservationDigest,
		feedback.MetricDigest,
		feedback.FeedbackEvidenceDigest,
	)
	return AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       input.Previous,
		Feedback:       feedback,
		MetricDigest:   feedback.MetricDigest,
		NonAuthorizing: true,
	})
}
