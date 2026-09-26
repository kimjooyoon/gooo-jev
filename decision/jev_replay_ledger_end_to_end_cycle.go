package decision

// ExecutionEnvelopeJEVReplayLedgerEndToEndCycleInput is the complete
// non-executing input from provenance checkpoint to cycle observation.
type ExecutionEnvelopeJEVReplayLedgerEndToEndCycleInput struct {
	Checkpoint       JEVFullProvenanceReplayLedgerCheckpoint
	Ledger           JEVImprovementReplayFeedbackLedger
	ChangePlanDigest string
	NonAuthorizing   bool
}

// ObserveJEVReplayLedgerCycleFromLedger derives aggregation from the ledger
// before reusing the audited cycle bridge.
func ObserveJEVReplayLedgerCycleFromLedger(input ExecutionEnvelopeJEVReplayLedgerEndToEndCycleInput) JEVReplayLedgerCycleObservation {
	output := JEVReplayLedgerCycleObservation{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Checkpoint.NonAuthorizing || !input.Ledger.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Checkpoint.NonExecuting || !input.Ledger.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	aggregation := AggregateJEVReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerAggregationInput{
		Ledger:         input.Ledger,
		NonAuthorizing: true,
	})
	if aggregation.Status == jevImprovementFeedbackAggregateUnknown {
		output.MissingStage = aggregation.MissingStage
		if output.MissingStage == "" {
			output.MissingStage = "feedback-aggregation"
		}
		return output
	}
	return ObserveJEVReplayLedgerCycle(ExecutionEnvelopeJEVReplayLedgerCycleInput{
		Checkpoint:       input.Checkpoint,
		ChangePlanDigest: input.ChangePlanDigest,
		Feedback:         aggregation,
		NonAuthorizing:   true,
	})
}