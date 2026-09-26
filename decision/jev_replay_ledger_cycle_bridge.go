package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ExecutionEnvelopeJEVReplayLedgerCycleInput feeds a validated checkpoint
// into the existing non-executing improvement-cycle observation.
type ExecutionEnvelopeJEVReplayLedgerCycleInput struct {
	Checkpoint      JEVFullProvenanceReplayLedgerCheckpoint
	ChangePlanDigest string
	Feedback        JEVImprovementFeedbackAggregation
	NonAuthorizing  bool
}

// JEVReplayLedgerCycleObservation preserves both cycle and replay-ledger
// evidence so later stages cannot silently drop the ledger boundary.
type JEVReplayLedgerCycleObservation struct {
	Status                   string
	MissingStage             string
	CycleStatus              string
	CycleEvidenceDigest      string
	CheckpointEvidenceDigest string
	LedgerEvidenceDigest     string
	MetricDigest             string
	EvidenceDigest           string
	NonExecuting             bool
	NonAuthorizing           bool
}

func (o JEVReplayLedgerCycleObservation) Validate() error {
	if o.Status == "" ||
		o.CycleStatus == "" ||
		o.CycleEvidenceDigest == "" ||
		o.CheckpointEvidenceDigest == "" ||
		o.LedgerEvidenceDigest == "" ||
		o.MetricDigest == "" ||
		o.EvidenceDigest == "" {
		return fmt.Errorf("incomplete JEV replay ledger cycle observation")
	}
	switch o.Status {
	case "stable-for-review", "needs-revision", "hold":
	default:
		return fmt.Errorf("invalid JEV replay ledger cycle observation status")
	}
	if o.Status != o.CycleStatus {
		return fmt.Errorf("JEV replay ledger cycle status mismatch")
	}
	if !o.NonExecuting {
		return fmt.Errorf("JEV replay ledger cycle observation must be non-executing")
	}
	if !o.NonAuthorizing {
		return fmt.Errorf("JEV replay ledger cycle observation must be non-authorizing")
	}
	expected := digestJEVReplayLedgerCycleObservation(
		o.Status,
		o.CycleStatus,
		o.CycleEvidenceDigest,
		o.CheckpointEvidenceDigest,
		o.LedgerEvidenceDigest,
		o.MetricDigest,
	)
	if o.EvidenceDigest != expected {
		return fmt.Errorf("JEV replay ledger cycle observation digest mismatch")
	}
	return nil
}

// ObserveJEVReplayLedgerCycle preserves UNKNOWN at the first invalid boundary.
func ObserveJEVReplayLedgerCycle(input ExecutionEnvelopeJEVReplayLedgerCycleInput) JEVReplayLedgerCycleObservation {
	output := JEVReplayLedgerCycleObservation{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Checkpoint.NonAuthorizing || !input.Feedback.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Checkpoint.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Checkpoint.Validate(); err != nil {
		output.MissingStage = "replay-ledger-checkpoint"
		return output
	}
	if strings.TrimSpace(input.ChangePlanDigest) == "" {
		output.MissingStage = "change-plan"
		return output
	}
	cycle := ObserveJEVImprovementCycleFromFullProvenance(ExecutionEnvelopeJEVFullProvenanceCycleInput{
		Provenance: ExecutionEnvelopeFullProvenanceSourceBinding{
			Status:                 "complete",
			DeclarationDigest:      input.Checkpoint.DeclarationDigest,
			IRDigest:               input.Checkpoint.IRDigest,
			GenerationDigest:       input.Checkpoint.GenerationDigest,
			ReverseObservationDigest: input.Checkpoint.ReverseObservationDigest,
			MetricDigest:            input.Checkpoint.MetricDigest,
			EvidenceDigest:          input.Checkpoint.ProvenanceEvidenceDigest,
			CompletenessDigest:      input.Checkpoint.CompletenessDigest,
			NonExecuting:            true,
			NonAuthorizing:          true,
		},
		ChangePlanDigest: input.ChangePlanDigest,
		Feedback:         input.Feedback,
		NonAuthorizing:   true,
	})
	if cycle.Status == "UNKNOWN" {
		output.MissingStage = cycle.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "improvement-cycle-observation"
		}
		return output
	}
	output.Status = cycle.Status
	output.CycleStatus = cycle.Status
	output.CycleEvidenceDigest = cycle.EvidenceDigest
	output.CheckpointEvidenceDigest = input.Checkpoint.CheckpointEvidenceDigest
	output.LedgerEvidenceDigest = input.Checkpoint.LedgerEvidenceDigest
	output.MetricDigest = input.Checkpoint.MetricDigest
	output.EvidenceDigest = digestJEVReplayLedgerCycleObservation(
		output.Status,
		output.CycleStatus,
		output.CycleEvidenceDigest,
		output.CheckpointEvidenceDigest,
		output.LedgerEvidenceDigest,
		output.MetricDigest,
	)
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "cycle-ledger-evidence"
		output.EvidenceDigest = ""
	}
	return output
}

func digestJEVReplayLedgerCycleObservation(
	status,
	cycleStatus,
	cycleEvidenceDigest,
	checkpointEvidenceDigest,
	ledgerEvidenceDigest,
	metricDigest string,
) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s",
		status,
		cycleStatus,
		cycleEvidenceDigest,
		checkpointEvidenceDigest,
		ledgerEvidenceDigest,
		metricDigest,
	)))
	return hex.EncodeToString(sum[:])
}