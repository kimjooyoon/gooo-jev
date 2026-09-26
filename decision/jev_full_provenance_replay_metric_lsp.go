package decision

// ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPInput projects the dual
// metric binding into an editor diagnostic without authorizing execution.
type ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPInput struct {
	Binding       JEVFullProvenanceReplayMetricBinding
	ReplayMetric  JEVReplayLedgerMetric
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPBinding exposes metric
// quality and evidence boundaries to the editor.
type ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPBinding struct {
	Status                string
	Publishable           bool
	Severity              string
	Code                  string
	MissingStage          string
	MissingStageIndex     int
	EvidenceDigest        string
	EvidencePrefixDigest  string
	ReplayMetricDigest    string
	RefutedPermille       int
	UnknownPermille       int
	NonExecuting          bool
	NonAuthorizing        bool
}

// ProjectJEVFullProvenanceReplayMetricToLSP preserves metric regressions and
// UNKNOWN ratios as diagnostics rather than flattening them into success.
func ProjectJEVFullProvenanceReplayMetricToLSP(input ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPInput) ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPBinding {
	output := ExecutionEnvelopeJEVFullProvenanceReplayMetricLSPBinding{
		Status:         "UNKNOWN",
		MissingStageIndex: -1,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing || !input.ReplayMetric.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Code = "authorization-boundary"
		return output
	}
	if !input.Binding.NonExecuting || !input.ReplayMetric.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Code = "execution-boundary"
		return output
	}
	if err := input.Binding.Validate(); err != nil {
		output.MissingStage = "replay-metric-binding"
		output.Code = "jev-replay-metric-binding-integrity"
		return output
	}
	if err := input.ReplayMetric.Validate(); err != nil {
		output.MissingStage = "replay-metric"
		output.Code = "jev-replay-metric-integrity"
		return output
	}
	if input.Binding.ReplayMetricDigest != input.ReplayMetric.MetricDigest {
		output.MissingStage = "replay-metric-binding"
		output.Code = "jev-replay-metric-binding"
		return output
	}
	prefix, err := Digest(struct {
		DeclarationDigest      string
		IRDigest               string
		GenerationDigest       string
		ReverseObservationDigest string
		ProvenanceMetricDigest string
		ReplayMetricDigest     string
		RefutedPermille        int
		UnknownPermille        int
	}{
		DeclarationDigest:        input.Binding.DeclarationDigest,
		IRDigest:                 input.Binding.IRDigest,
		GenerationDigest:         input.Binding.GenerationDigest,
		ReverseObservationDigest: input.Binding.ReverseObservationDigest,
		ProvenanceMetricDigest:   input.Binding.ProvenanceMetricDigest,
		ReplayMetricDigest:       input.ReplayMetric.MetricDigest,
		RefutedPermille:          input.ReplayMetric.RefutedPermille,
		UnknownPermille:          input.ReplayMetric.UnknownPermille,
	})
	if err != nil {
		output.MissingStage = "replay-metric-prefix"
		output.Code = "jev-replay-metric-prefix"
		return output
	}
	output.Publishable = true
	output.EvidenceDigest = input.Binding.EvidenceDigest
	output.EvidencePrefixDigest = prefix
	output.ReplayMetricDigest = input.ReplayMetric.MetricDigest
	output.RefutedPermille = input.ReplayMetric.RefutedPermille
	output.UnknownPermille = input.ReplayMetric.UnknownPermille
	switch {
	case input.ReplayMetric.UnknownPermille > 0:
		output.Status = "diagnostic"
		output.Severity = "warning"
		output.Code = "jev.replay-metric.unknown"
	case input.ReplayMetric.RefutedPermille > 0:
		output.Status = "diagnostic"
		output.Severity = "warning"
		output.Code = "jev.replay-metric.refuted"
	default:
		output.Status = "ready"
		output.Severity = "info"
		output.Code = "jev.replay-metric.stable"
	}
	return output
}