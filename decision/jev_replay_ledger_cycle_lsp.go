package decision

import "strings"

// ExecutionEnvelopeJEVReplayLedgerCycleLSPInput projects a replay-ledger
// cycle observation without turning editor output into authorization.
type ExecutionEnvelopeJEVReplayLedgerCycleLSPInput struct {
	Observation    JEVReplayLedgerCycleObservation
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVReplayLedgerCycleLSPBinding is the editor-facing
// projection of cycle, checkpoint, and ledger evidence.
type ExecutionEnvelopeJEVReplayLedgerCycleLSPBinding struct {
	Status                   string
	Publishable              bool
	Severity                 string
	Code                     string
	MissingStage             string
	MissingStageIndex        int
	EvidenceDigest            string
	EvidencePrefixDigest      string
	LedgerEvidenceDigest      string
	NonExecuting              bool
	NonAuthorizing            bool
}

type replayLedgerCycleLSPStage struct {
	Name  string
	Value string
}

func replayLedgerCycleLSPStageIndex(stage string) (int, bool) {
	switch stage {
	case "replay-ledger-checkpoint", "declaration":
		return 0, true
	case "ir":
		return 1, true
	case "generation":
		return 2, true
	case "reverse-observation":
		return 3, true
	case "metric":
		return 4, true
	case "change-plan", "feedback-aggregation", "replay-feedback-ledger":
		return 5, true
	case "cycle-ledger-evidence", "improvement-cycle-observation":
		return 6, true
	default:
		return -1, false
	}
}

func deriveReplayLedgerCycleLSPEvidencePrefixDigest(observation JEVReplayLedgerCycleObservation, missingStageIndex int) (string, error) {
	stages := []replayLedgerCycleLSPStage{
		{Name: "checkpoint", Value: observation.CheckpointEvidenceDigest},
		{Name: "ledger", Value: observation.LedgerEvidenceDigest},
		{Name: "metric", Value: observation.MetricDigest},
		{Name: "cycle", Value: observation.CycleEvidenceDigest},
	}
	limit := len(stages)
	if missingStageIndex >= 0 && missingStageIndex < limit {
		limit = missingStageIndex
	}
	return Digest(struct {
		MissingStageIndex int
		Stages            []replayLedgerCycleLSPStage
	}{
		MissingStageIndex: missingStageIndex,
		Stages:            stages[:limit],
	})
}

// ProjectJEVReplayLedgerCycleToLSP preserves cycle status and UNKNOWN prefix
// evidence while keeping the projection non-executing and non-authorizing.
func ProjectJEVReplayLedgerCycleToLSP(input ExecutionEnvelopeJEVReplayLedgerCycleLSPInput) ExecutionEnvelopeJEVReplayLedgerCycleLSPBinding {
	output := ExecutionEnvelopeJEVReplayLedgerCycleLSPBinding{
		Status: "UNKNOWN", MissingStageIndex: -1,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Observation.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.Code = "authorization-boundary"
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Observation.NonExecuting {
		output.Code = "execution-boundary"
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Observation.Status == "UNKNOWN" {
		output.MissingStage = input.Observation.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.Code = "jev-replay-ledger-cycle-evidence"
			return output
		}
		index, ok := replayLedgerCycleLSPStageIndex(output.MissingStage)
		if !ok {
			output.Code = "jev-replay-ledger-cycle-location"
			return output
		}
		prefix, err := deriveReplayLedgerCycleLSPEvidencePrefixDigest(input.Observation, index)
		if err != nil {
			output.Code = "jev-replay-ledger-cycle-prefix"
			return output
		}
		output.Status = "diagnostic"
		output.Publishable = true
		output.Severity = "error"
		output.Code = "jev.replay-ledger-cycle." + strings.ReplaceAll(output.MissingStage, "_", "-")
		output.MissingStageIndex = index
		output.EvidencePrefixDigest = prefix
		return output
	}
	if err := input.Observation.Validate(); err != nil {
		output.Code = "jev-replay-ledger-cycle-integrity"
		return output
	}
	severity := "info"
	switch input.Observation.Status {
	case "stable-for-review":
		severity = "info"
	case "needs-revision", "hold":
		severity = "warning"
	default:
		output.Code = "jev-replay-ledger-cycle-status"
		return output
	}
	prefix, err := deriveReplayLedgerCycleLSPEvidencePrefixDigest(input.Observation, -1)
	if err != nil {
		output.Code = "jev-replay-ledger-cycle-prefix"
		return output
	}
	output.Status = "ready"
	output.Publishable = true
	output.Severity = severity
	output.Code = "jev.replay-ledger-cycle." + strings.ReplaceAll(input.Observation.Status, "_", "-")
	output.EvidenceDigest = input.Observation.EvidenceDigest
	output.EvidencePrefixDigest = prefix
	output.LedgerEvidenceDigest = input.Observation.LedgerEvidenceDigest
	return output
}