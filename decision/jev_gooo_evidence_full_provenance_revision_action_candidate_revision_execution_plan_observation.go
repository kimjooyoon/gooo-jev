package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput
// reverse-observes an execution plan lifecycle without executing the plan.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput struct {
	Plan                     ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding
	LifecycleEvent           string
	ObservationEvidenceDigest string
	ReverseObservationDigest string
	NonAuthorizing           bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding
// preserves lifecycle, plan, suspend/resume, and reverse-observation provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding struct {
	Status                    string
	MissingStage              string
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
	EvidencePrefixDigest      string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
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
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage execution plan observation binding")
	}
	expectedDisposition := lifecycleDispositionForGoooExtendedLineageExecutionPlanObservation(b.ObservationStatus)
	if expectedDisposition == "" || b.LifecycleDisposition != expectedDisposition {
		return fmt.Errorf("invalid Gooo extended lineage execution plan lifecycle disposition")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage execution plan observation must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage execution plan observation digest mismatch")
	}
	return nil
}

// ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation
// records lifecycle facts and leaves execution and authorization to later deterministic stages.
func ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			PermissionState: "not-authorized",
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Plan.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Plan.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Plan.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-execution-plan-validation")
	}
	event := strings.TrimSpace(input.LifecycleEvent)
	disposition := lifecycleDispositionForGoooExtendedLineageExecutionPlanObservation(event)
	if disposition == "" {
		return unknown("lifecycle-event")
	}
	if strings.TrimSpace(input.ObservationEvidenceDigest) == "" {
		return unknown("observation-evidence")
	}
	if strings.TrimSpace(input.ReverseObservationDigest) == "" {
		return unknown("reverse-observation")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding{
		Status:                    "bound",
		ObservationStatus:         event,
		LifecycleDisposition:      disposition,
		PlanStatus:                input.Plan.PlanStatus,
		CandidateStatus:           input.Plan.CandidateStatus,
		CandidateDigest:           input.Plan.CandidateDigest,
		PermissionState:           "not-authorized",
		WorkspaceSource:           input.Plan.WorkspaceSource,
		Gateway:                   input.Plan.Gateway,
		Model:                     input.Plan.Model,
		NetworkAllowlistDigest:    input.Plan.NetworkAllowlistDigest,
		SuspendTokenDigest:        input.Plan.SuspendTokenDigest,
		ResumeTokenDigest:         input.Plan.ResumeTokenDigest,
		PlanEvidenceDigest:        input.Plan.PlanEvidenceDigest,
		ObservationEvidenceDigest: input.ObservationEvidenceDigest,
		ReverseObservationDigest:  input.ReverseObservationDigest,
		EvidencePrefixDigest:      input.Plan.EvidencePrefixDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-evidence")
	}
	return output
}

func lifecycleDispositionForGoooExtendedLineageExecutionPlanObservation(event string) string {
	switch strings.TrimSpace(event) {
	case "not-started":
		return "awaiting-policy"
	case "suspended":
		return "resume-required"
	case "resumed":
		return "resume-observed"
	default:
		return ""
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding) string {
	digest, err := Digest(struct {
		Status                    string
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
		EvidencePrefixDigest      string
		NonExecuting              bool
		NonAuthorizing            bool
	}{
		Status:                    b.Status,
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
		EvidencePrefixDigest:      b.EvidencePrefixDigest,
		NonExecuting:              b.NonExecuting,
		NonAuthorizing:            b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}