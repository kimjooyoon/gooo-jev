package decision

// ExecutionEnvelopeProvenanceDiagnosticMetricInput contains counts from
// structured editor or LSP diagnostics.
type ExecutionEnvelopeProvenanceDiagnosticMetricInput struct {
	ClearCount       int
	DiagnosticCount  int
	UnknownCount     int
	NonAuthorizing   bool
}

// ExecutionEnvelopeProvenanceDiagnosticMetric records diagnostic counts
// without treating clarity as semantic improvement.
type ExecutionEnvelopeProvenanceDiagnosticMetric struct {
	Status          string
	Total           int
	ClearCount      int
	DiagnosticCount int
	UnknownCount    int
	NonAuthorizing  bool
}

// MeasureExecutionEnvelopeProvenanceDiagnosticMetric preserves unknown counts
// and rejects negative or empty observations.
func MeasureExecutionEnvelopeProvenanceDiagnosticMetric(input ExecutionEnvelopeProvenanceDiagnosticMetricInput) ExecutionEnvelopeProvenanceDiagnosticMetric {
	output := ExecutionEnvelopeProvenanceDiagnosticMetric{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		return output
	}
	if input.ClearCount < 0 || input.DiagnosticCount < 0 || input.UnknownCount < 0 {
		return output
	}
	output.Total = input.ClearCount + input.DiagnosticCount + input.UnknownCount
	if output.Total == 0 {
		return output
	}
	output.ClearCount = input.ClearCount
	output.DiagnosticCount = input.DiagnosticCount
	output.UnknownCount = input.UnknownCount
	output.Status = "measured"
	if output.UnknownCount > 0 {
		output.Status = "measured-with-unknown"
	}
	return output
}