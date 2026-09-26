package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput
// binds an externally supplied metric to a reverse-observed plan lifecycle.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput struct {
	Observation              ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding
	MetricName               string
	MetricValue              int
	MetricUnit               string
	MetricEvidenceDigest     string
	NonAuthorizing           bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding
// preserves measurement value, unit, plan lifecycle, and reverse-observation provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding struct {
	Status                    string
	MissingStage              string
	MetricStatus              string
	ObservationStatus         string
	LifecycleDisposition      string
	PlanStatus                string
	CandidateStatus           string
	CandidateDigest           string
	PermissionState           string
	WorkspaceSource           string
	Gateway                   string
	Model                     string
	NetworkAllowlistDigest    string
	SuspendTokenDigest        string
	ResumeTokenDigest         string
	PlanEvidenceDigest        string
	ObservationEvidenceDigest string
	ReverseObservationDigest  string
	MetricName                string
	MetricValue               int
	MetricUnit                string
	MetricEvidenceDigest      string
	EvidencePrefixDigest      string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.MetricStatus != "observed" ||
		b.ObservationStatus == "" ||
		b.LifecycleDisposition == "" ||
		b.PlanStatus != "planned" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.PermissionState != "not-authorized" ||
		b.WorkspaceSource == "" ||
		b.Gateway == "" ||
		b.Model == "" ||
		b.NetworkAllowlistDigest == "" ||
		b.SuspendTokenDigest == "" ||
		b.ResumeTokenDigest == "" ||
		b.PlanEvidenceDigest == "" ||
		b.ObservationEvidenceDigest == "" ||
		b.ReverseObservationDigest == "" ||
		b.MetricName == "" ||
		b.MetricValue < 0 ||
		b.MetricUnit == "" ||
		b.MetricEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage execution plan observation metric binding")
	}
	expectedDisposition := lifecycleDispositionForGoooExtendedLineageExecutionPlanObservation(b.ObservationStatus)
	if expectedDisposition == "" || b.LifecycleDisposition != expectedDisposition {
		return fmt.Errorf("invalid Gooo extended lineage execution plan metric lifecycle disposition")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage execution plan observation metric must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage execution plan observation metric digest mismatch")
	}
	return nil
}

// MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric
// records an externally measured value and never infers it from plan execution.
func MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-metric"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			PermissionState: "not-authorized",
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Observation.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Observation.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Observation.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-validation")
	}
	if strings.TrimSpace(input.MetricName) == "" {
		return unknown("metric-name")
	}
	if input.MetricValue < 0 {
		return unknown("metric-value")
	}
	if strings.TrimSpace(input.MetricUnit) == "" {
		return unknown("metric-unit")
	}
	if strings.TrimSpace(input.MetricEvidenceDigest) == "" {
		return unknown("metric-evidence")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding{
		Status:                    "bound",
		MetricStatus:              "observed",
		ObservationStatus:         input.Observation.ObservationStatus,
		LifecycleDisposition:      input.Observation.LifecycleDisposition,
		PlanStatus:                input.Observation.PlanStatus,
		CandidateStatus:           input.Observation.CandidateStatus,
		CandidateDigest:           input.Observation.CandidateDigest,
		PermissionState:            "not-authorized",
		WorkspaceSource:            input.Observation.WorkspaceSource,
		Gateway:                   input.Observation.Gateway,
		Model:                     input.Observation.Model,
		NetworkAllowlistDigest:    input.Observation.NetworkAllowlistDigest,
		SuspendTokenDigest:        input.Observation.SuspendTokenDigest,
		ResumeTokenDigest:         input.Observation.ResumeTokenDigest,
		PlanEvidenceDigest:        input.Observation.PlanEvidenceDigest,
		ObservationEvidenceDigest: input.Observation.ObservationEvidenceDigest,
		ReverseObservationDigest:  input.Observation.ReverseObservationDigest,
		MetricName:                strings.TrimSpace(input.MetricName),
		MetricValue:               input.MetricValue,
		MetricUnit:                strings.TrimSpace(input.MetricUnit),
		MetricEvidenceDigest:      input.MetricEvidenceDigest,
		EvidencePrefixDigest:      input.Observation.EvidencePrefixDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-metric-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding) string {
	digest, err := Digest(struct {
		Status                    string
		MetricStatus              string
		ObservationStatus         string
		LifecycleDisposition      string
		PlanStatus                string
		CandidateStatus           string
		CandidateDigest           string
		PermissionState           string
		WorkspaceSource            string
		Gateway                   string
		Model                     string
		NetworkAllowlistDigest    string
		SuspendTokenDigest        string
		ResumeTokenDigest         string
		PlanEvidenceDigest        string
		ObservationEvidenceDigest string
		ReverseObservationDigest  string
		MetricName                string
		MetricValue               int
		MetricUnit                string
		MetricEvidenceDigest      string
		EvidencePrefixDigest      string
		NonExecuting              bool
		NonAuthorizing            bool
	}{
		Status:                    b.Status,
		MetricStatus:              b.MetricStatus,
		ObservationStatus:         b.ObservationStatus,
		LifecycleDisposition:      b.LifecycleDisposition,
		PlanStatus:                b.PlanStatus,
		CandidateStatus:           b.CandidateStatus,
		CandidateDigest:           b.CandidateDigest,
		PermissionState:           b.PermissionState,
		WorkspaceSource:            b.WorkspaceSource,
		Gateway:                   b.Gateway,
		Model:                     b.Model,
		NetworkAllowlistDigest:    b.NetworkAllowlistDigest,
		SuspendTokenDigest:         b.SuspendTokenDigest,
		ResumeTokenDigest:          b.ResumeTokenDigest,
		PlanEvidenceDigest:         b.PlanEvidenceDigest,
		ObservationEvidenceDigest: b.ObservationEvidenceDigest,
		ReverseObservationDigest:  b.ReverseObservationDigest,
		MetricName:                b.MetricName,
		MetricValue:               b.MetricValue,
		MetricUnit:                b.MetricUnit,
		MetricEvidenceDigest:      b.MetricEvidenceDigest,
		EvidencePrefixDigest:      b.EvidencePrefixDigest,
		NonExecuting:              b.NonExecuting,
		NonAuthorizing:            b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}