package decision

import "fmt"

// ExecutionEnvelopeJEVReplayLedgerMetricInput measures one validated replay
// ledger without executing or authorizing any candidate.
type ExecutionEnvelopeJEVReplayLedgerMetricInput struct {
	Ledger         JEVImprovementReplayFeedbackLedger
	NonAuthorizing bool
}

// JEVReplayLedgerMetric records deterministic counts and permille ratios.
type JEVReplayLedgerMetric struct {
	Status             string
	MissingStage       string
	Total              int
	Confirmed          int
	Refuted            int
	Unknown            int
	ConfirmedPermille  int
	RefutedPermille    int
	UnknownPermille    int
	InputLedgerDigest  string
	MetricDigest       string
	EvidenceDigest     string
	NonExecuting       bool
	NonAuthorizing     bool
}

func (m JEVReplayLedgerMetric) Validate() error {
	if m.Status != "ready" || m.Total <= 0 ||
		m.InputLedgerDigest == "" || m.MetricDigest == "" || m.EvidenceDigest == "" {
		return fmt.Errorf("incomplete JEV replay ledger metric")
	}
	if m.Confirmed < 0 || m.Refuted < 0 || m.Unknown < 0 ||
		m.Confirmed+m.Refuted+m.Unknown != m.Total {
		return fmt.Errorf("invalid JEV replay ledger metric counts")
	}
	if m.ConfirmedPermille != replayLedgerPermille(m.Confirmed, m.Total) ||
		m.RefutedPermille != replayLedgerPermille(m.Refuted, m.Total) ||
		m.UnknownPermille != replayLedgerPermille(m.Unknown, m.Total) {
		return fmt.Errorf("invalid JEV replay ledger metric ratios")
	}
	if !m.NonExecuting {
		return fmt.Errorf("JEV replay ledger metric must be non-executing")
	}
	if !m.NonAuthorizing {
		return fmt.Errorf("JEV replay ledger metric must be non-authorizing")
	}
	expectedMetric, err := Digest(struct {
		Total            int
		Confirmed        int
		Refuted          int
		Unknown          int
		ConfirmedPermille int
		RefutedPermille   int
		UnknownPermille   int
		InputLedgerDigest string
	}{
		Total:             m.Total,
		Confirmed:         m.Confirmed,
		Refuted:           m.Refuted,
		Unknown:            m.Unknown,
		ConfirmedPermille: m.ConfirmedPermille,
		RefutedPermille:   m.RefutedPermille,
		UnknownPermille:   m.UnknownPermille,
		InputLedgerDigest: m.InputLedgerDigest,
	})
	if err != nil || m.MetricDigest != expectedMetric {
		return fmt.Errorf("JEV replay ledger metric digest mismatch")
	}
	expectedEvidence, err := Digest(struct {
		MetricDigest      string
		InputLedgerDigest string
	}{
		MetricDigest:      m.MetricDigest,
		InputLedgerDigest: m.InputLedgerDigest,
	})
	if err != nil || m.EvidenceDigest != expectedEvidence {
		return fmt.Errorf("JEV replay ledger metric evidence mismatch")
	}
	return nil
}

// MeasureJEVReplayLedgerMetric computes ratios from the validated ledger.
// It never treats cache presence, execution, or authorization as evidence.
func MeasureJEVReplayLedgerMetric(input ExecutionEnvelopeJEVReplayLedgerMetricInput) JEVReplayLedgerMetric {
	output := JEVReplayLedgerMetric{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Ledger.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Ledger.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Ledger.Validate(); err != nil {
		output.MissingStage = "replay-feedback-ledger"
		return output
	}
	output.Status = "ready"
	output.Total = input.Ledger.ConfirmedCount + input.Ledger.RefutedCount + input.Ledger.UnknownCount
	output.Confirmed = input.Ledger.ConfirmedCount
	output.Refuted = input.Ledger.RefutedCount
	output.Unknown = input.Ledger.UnknownCount
	if output.Total <= 0 {
		output.Status = "UNKNOWN"
		output.MissingStage = "feedback-history"
		return output
	}
	output.ConfirmedPermille = replayLedgerPermille(output.Confirmed, output.Total)
	output.RefutedPermille = replayLedgerPermille(output.Refuted, output.Total)
	output.UnknownPermille = replayLedgerPermille(output.Unknown, output.Total)
	output.InputLedgerDigest = input.Ledger.EvidenceDigest
	var err error
	output.MetricDigest, err = Digest(struct {
		Total             int
		Confirmed         int
		Refuted           int
		Unknown           int
		ConfirmedPermille int
		RefutedPermille   int
		UnknownPermille   int
		InputLedgerDigest string
	}{
		Total:             output.Total,
		Confirmed:         output.Confirmed,
		Refuted:           output.Refuted,
		Unknown:            output.Unknown,
		ConfirmedPermille: output.ConfirmedPermille,
		RefutedPermille:   output.RefutedPermille,
		UnknownPermille:   output.UnknownPermille,
		InputLedgerDigest: output.InputLedgerDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "replay-ledger-metric"
		return output
	}
	output.EvidenceDigest, err = Digest(struct {
		MetricDigest      string
		InputLedgerDigest string
	}{
		MetricDigest:      output.MetricDigest,
		InputLedgerDigest: output.InputLedgerDigest,
	})
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "replay-ledger-metric-evidence"
		output.MetricDigest = ""
		return output
	}
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "replay-ledger-metric-evidence"
		output.MetricDigest = ""
		output.EvidenceDigest = ""
	}
	return output
}

func replayLedgerPermille(part, total int) int {
	return int((int64(part) * 1000) / int64(total))
}