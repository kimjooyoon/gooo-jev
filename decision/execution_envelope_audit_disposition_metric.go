package decision

// ExecutionEnvelopeAuditDispositionMetricInput contains counted observations
// for each safe operating mode.
type ExecutionEnvelopeAuditDispositionMetricInput struct {
	AuditOnlyCount  int
	ReviewOnlyCount int
	BlockedCount    int
	NonAuthorizing  bool
}

// ExecutionEnvelopeAuditDispositionMetric records counts without claiming
// that a dominant mode is an improvement.
type ExecutionEnvelopeAuditDispositionMetric struct {
	Status         string
	Total           int
	DominantMode   string
	NonAuthorizing bool
}

// MeasureExecutionEnvelopeAuditDispositionMetric measures only non-negative,
// non-authorizing observations.
func MeasureExecutionEnvelopeAuditDispositionMetric(input ExecutionEnvelopeAuditDispositionMetricInput) ExecutionEnvelopeAuditDispositionMetric {
	output := ExecutionEnvelopeAuditDispositionMetric{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		return output
	}
	if input.AuditOnlyCount < 0 || input.ReviewOnlyCount < 0 || input.BlockedCount < 0 {
		return output
	}
	output.Total = input.AuditOnlyCount + input.ReviewOnlyCount + input.BlockedCount
	if output.Total == 0 {
		return output
	}
	maximum := input.AuditOnlyCount
	output.DominantMode = "audit-only"
	if input.ReviewOnlyCount > maximum {
		maximum = input.ReviewOnlyCount
		output.DominantMode = "review-only"
	} else if input.ReviewOnlyCount == maximum {
		output.DominantMode = "inconclusive"
	}
	if input.BlockedCount > maximum {
		maximum = input.BlockedCount
		output.DominantMode = "blocked"
	} else if input.BlockedCount == maximum {
		output.DominantMode = "inconclusive"
	}
	if output.DominantMode == "inconclusive" {
		output.Status = "measured"
		return output
	}
	output.Status = "measured"
	return output
}