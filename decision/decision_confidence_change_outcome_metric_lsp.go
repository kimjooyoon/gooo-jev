package decision

import "fmt"

// DecisionConfidenceChangeOutcomeMetricLSPDiagnostic is a read-only editor
// projection of one observed external outcome.
type DecisionConfidenceChangeOutcomeMetricLSPDiagnostic struct {
	Status            string
	Outcome           string
	Code              string
	Severity          string
	Message           string
	MissingStage      string
	ObservationDigest string
	MetricDigest      string
	EvidenceDigest    string
	ReadOnly          bool
	NonAuthorizing    bool
}

// ProjectDecisionConfidenceChangeOutcomeMetricLSP preserves outcome evidence
// while refusing to turn UNKNOWN or invalid evidence into a success signal.
func ProjectDecisionConfidenceChangeOutcomeMetricLSP(
	metric DecisionConfidenceChangeOutcomeMetric,
) DecisionConfidenceChangeOutcomeMetricLSPDiagnostic {
	output := DecisionConfidenceChangeOutcomeMetricLSPDiagnostic{
		Status:         "UNKNOWN",
		Outcome:        "unknown",
		Code:           "JEV_DECISION_OUTCOME_UNKNOWN",
		Severity:       "warning",
		Message:        "JEV decision outcome is UNKNOWN; external review remains required",
		MissingStage:   "outcome-metric",
		ReadOnly:       true,
		NonAuthorizing: true,
	}
	if err := metric.Validate(); err != nil {
		if !metric.NonAuthorizing {
			output.NonAuthorizing = false
			output.MissingStage = "authorization-boundary"
		}
		return finalizeDecisionConfidenceChangeOutcomeMetricLSP(output)
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.ObservationDigest = metric.ObservationDigest
	output.MetricDigest = metric.MetricDigest
	switch {
	case metric.AppliedCount == 1:
		output.Outcome = string(DecisionConfidenceChangePlanObservedApplied)
		output.Code = "JEV_DECISION_OUTCOME_APPLIED"
		output.Severity = "info"
		output.Message = "JEV decision outcome was observed as applied; no automatic action was taken"
	case metric.AbortedCount == 1:
		output.Outcome = string(DecisionConfidenceChangePlanObservedAborted)
		output.Code = "JEV_DECISION_OUTCOME_ABORTED"
		output.Severity = "warning"
		output.Message = "JEV decision outcome was observed as aborted; no automatic action was taken"
	case metric.RolledBackCount == 1:
		output.Outcome = string(DecisionConfidenceChangePlanObservedRolledBack)
		output.Code = "JEV_DECISION_OUTCOME_ROLLED_BACK"
		output.Severity = "warning"
		output.Message = "JEV decision outcome was observed as rolled back; no automatic action was taken"
	case metric.NotAppliedCount == 1:
		output.Outcome = string(DecisionConfidenceChangePlanObservedNotApplied)
		output.Code = "JEV_DECISION_OUTCOME_NOT_APPLIED"
		output.Severity = "warning"
		output.Message = "JEV decision outcome was observed as not applied; no automatic action was taken"
	default:
		output.Outcome = string(DecisionConfidenceChangePlanObservedUnknown)
		output.Code = "JEV_DECISION_OUTCOME_UNKNOWN"
		output.Severity = "warning"
		output.Message = "JEV decision outcome evidence is bound but the external outcome remains UNKNOWN"
	}
	return finalizeDecisionConfidenceChangeOutcomeMetricLSP(output)
}

func (d DecisionConfidenceChangeOutcomeMetricLSPDiagnostic) Validate() error {
	if d.Status != "BOUND" && d.Status != "UNKNOWN" {
		return fmt.Errorf("invalid outcome LSP status %q", d.Status)
	}
	if d.Outcome != string(DecisionConfidenceChangePlanObservedApplied) &&
		d.Outcome != string(DecisionConfidenceChangePlanObservedAborted) &&
		d.Outcome != string(DecisionConfidenceChangePlanObservedRolledBack) &&
		d.Outcome != string(DecisionConfidenceChangePlanObservedNotApplied) &&
		d.Outcome != string(DecisionConfidenceChangePlanObservedUnknown) {
		return fmt.Errorf("invalid outcome LSP outcome %q", d.Outcome)
	}
	if d.Status == "BOUND" {
		if d.MissingStage != "" || d.ObservationDigest == "" || d.MetricDigest == "" {
			return fmt.Errorf("bound outcome LSP diagnostic is incomplete")
		}
	} else if d.MissingStage == "" {
		return fmt.Errorf("unknown outcome LSP diagnostic has no missing stage")
	}
	if !d.ReadOnly || !d.NonAuthorizing {
		return fmt.Errorf("outcome LSP diagnostic must remain read-only and non-authorizing")
	}
	if d.EvidenceDigest == "" {
		return fmt.Errorf("outcome LSP evidence digest is required")
	}
	expected := digestDecisionConfidenceChangeOutcomeMetricLSP(d)
	if expected != d.EvidenceDigest {
		return fmt.Errorf("outcome LSP evidence digest mismatch")
	}
	return nil
}

func finalizeDecisionConfidenceChangeOutcomeMetricLSP(
	output DecisionConfidenceChangeOutcomeMetricLSPDiagnostic,
) DecisionConfidenceChangeOutcomeMetricLSPDiagnostic {
	digest, err := digestDecisionConfidenceChangeOutcomeMetricLSP(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Outcome = string(DecisionConfidenceChangePlanObservedUnknown)
		output.Code = "JEV_DECISION_OUTCOME_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV decision outcome is UNKNOWN: missing lsp-evidence-digest"
		output.ObservationDigest = ""
		output.MetricDigest = ""
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}

func digestDecisionConfidenceChangeOutcomeMetricLSP(
	diagnostic DecisionConfidenceChangeOutcomeMetricLSPDiagnostic,
) (string, error) {
	return Digest(struct {
		Status            string
		Outcome           string
		Code              string
		Severity          string
		Message           string
		MissingStage      string
		ObservationDigest string
		MetricDigest      string
		ReadOnly          bool
		NonAuthorizing    bool
	}{
		Status:            diagnostic.Status,
		Outcome:           diagnostic.Outcome,
		Code:              diagnostic.Code,
		Severity:          diagnostic.Severity,
		Message:           diagnostic.Message,
		MissingStage:      diagnostic.MissingStage,
		ObservationDigest: diagnostic.ObservationDigest,
		MetricDigest:      diagnostic.MetricDigest,
		ReadOnly:          diagnostic.ReadOnly,
		NonAuthorizing:    diagnostic.NonAuthorizing,
	})
}
