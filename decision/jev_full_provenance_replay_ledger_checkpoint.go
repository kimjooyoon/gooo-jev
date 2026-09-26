package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	jevFullProvenanceReplayLedgerCheckpointComplete = "complete"
	jevFullProvenanceReplayLedgerCheckpointUnknown  = "UNKNOWN"
)

// ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput joins the full
// declaration-to-metric provenance with the replay feedback ledger.
type ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput struct {
	Provenance     ExecutionEnvelopeFullProvenanceSourceBinding
	Ledger         JEVImprovementReplayFeedbackLedger
	NonAuthorizing bool
}

// JEVFullProvenanceReplayLedgerCheckpoint is a non-executing checkpoint for
// one complete provenance chain and its accumulated replay evidence.
type JEVFullProvenanceReplayLedgerCheckpoint struct {
	Status                    string
	MissingStage              string
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	ReverseObservationDigest  string
	MetricDigest              string
	ProvenanceEvidenceDigest  string
	CompletenessDigest        string
	LedgerEvidenceDigest      string
	LedgerEntryDigest         string
	CheckpointEvidenceDigest  string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (c JEVFullProvenanceReplayLedgerCheckpoint) Validate() error {
	if c.Status != jevFullProvenanceReplayLedgerCheckpointComplete ||
		c.DeclarationDigest == "" ||
		c.IRDigest == "" ||
		c.GenerationDigest == "" ||
		c.ReverseObservationDigest == "" ||
		c.MetricDigest == "" ||
		c.ProvenanceEvidenceDigest == "" ||
		c.CompletenessDigest == "" ||
		c.LedgerEvidenceDigest == "" ||
		c.LedgerEntryDigest == "" ||
		c.CheckpointEvidenceDigest == "" {
		return fmt.Errorf("incomplete JEV full provenance replay ledger checkpoint")
	}
	if !c.NonExecuting {
		return fmt.Errorf("JEV full provenance replay ledger checkpoint must be non-executing")
	}
	if !c.NonAuthorizing {
		return fmt.Errorf("JEV full provenance replay ledger checkpoint must be non-authorizing")
	}
	expected := digestJEVFullProvenanceReplayLedgerCheckpoint(
		c.DeclarationDigest,
		c.IRDigest,
		c.GenerationDigest,
		c.ReverseObservationDigest,
		c.MetricDigest,
		c.ProvenanceEvidenceDigest,
		c.CompletenessDigest,
		c.LedgerEvidenceDigest,
		c.LedgerEntryDigest,
	)
	if c.CheckpointEvidenceDigest != expected {
		return fmt.Errorf("JEV full provenance replay ledger checkpoint digest mismatch")
	}
	return nil
}

// BindJEVFullProvenanceReplayLedgerCheckpoint creates a checkpoint only when
// both the complete source chain and replay ledger validate.
func BindJEVFullProvenanceReplayLedgerCheckpoint(input ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput) JEVFullProvenanceReplayLedgerCheckpoint {
	output := JEVFullProvenanceReplayLedgerCheckpoint{
		Status:         jevFullProvenanceReplayLedgerCheckpointUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Provenance.NonAuthorizing || !input.Ledger.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Provenance.NonExecuting || !input.Ledger.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Provenance.Status != "complete" {
		output.MissingStage = input.Provenance.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "full-provenance"
		}
		return output
	}
	provenanceStages := []struct {
		name  string
		value string
	}{
		{"declaration", input.Provenance.DeclarationDigest},
		{"ir", input.Provenance.IRDigest},
		{"generation", input.Provenance.GenerationDigest},
		{"reverse-observation", input.Provenance.ReverseObservationDigest},
		{"metric", input.Provenance.MetricDigest},
		{"provenance-evidence", input.Provenance.EvidenceDigest},
		{"provenance-completeness", input.Provenance.CompletenessDigest},
	}
	for _, stage := range provenanceStages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	if err := input.Ledger.Validate(); err != nil {
		output.MissingStage = "replay-feedback-ledger"
		return output
	}
	output.Status = jevFullProvenanceReplayLedgerCheckpointComplete
	output.DeclarationDigest = input.Provenance.DeclarationDigest
	output.IRDigest = input.Provenance.IRDigest
	output.GenerationDigest = input.Provenance.GenerationDigest
	output.ReverseObservationDigest = input.Provenance.ReverseObservationDigest
	output.MetricDigest = input.Provenance.MetricDigest
	output.ProvenanceEvidenceDigest = input.Provenance.EvidenceDigest
	output.CompletenessDigest = input.Provenance.CompletenessDigest
	output.LedgerEvidenceDigest = input.Ledger.EvidenceDigest
	output.LedgerEntryDigest = input.Ledger.EntryDigest
	output.CheckpointEvidenceDigest = digestJEVFullProvenanceReplayLedgerCheckpoint(
		output.DeclarationDigest,
		output.IRDigest,
		output.GenerationDigest,
		output.ReverseObservationDigest,
		output.MetricDigest,
		output.ProvenanceEvidenceDigest,
		output.CompletenessDigest,
		output.LedgerEvidenceDigest,
		output.LedgerEntryDigest,
	)
	if err := output.Validate(); err != nil {
		output.Status = jevFullProvenanceReplayLedgerCheckpointUnknown
		output.MissingStage = "checkpoint-evidence"
		output.CheckpointEvidenceDigest = ""
	}
	return output
}

func digestJEVFullProvenanceReplayLedgerCheckpoint(
	declarationDigest,
	irDigest,
	generationDigest,
	reverseObservationDigest,
	metricDigest,
	provenanceEvidenceDigest,
	completenessDigest,
	ledgerEvidenceDigest,
	ledgerEntryDigest string,
) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%s|%s|%s",
		declarationDigest,
		irDigest,
		generationDigest,
		reverseObservationDigest,
		metricDigest,
		provenanceEvidenceDigest,
		completenessDigest,
		ledgerEvidenceDigest,
		ledgerEntryDigest,
	)))
	return hex.EncodeToString(sum[:])
}