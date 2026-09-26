package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput
// models an AX-like plan boundary without executing or authorizing the candidate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput struct {
	Guardian                  ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding
	WorkspaceSource           string
	Gateway                   string
	Model                     string
	NetworkAllowlistDigest    string
	SuspendTokenDigest        string
	NonAuthorizing            bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding
// preserves task, workspace, gateway, model, network, and suspend/resume provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding struct {
	Status                       string
	MissingStage                 string
	PlanStatus                   string
	CandidateStatus              string
	CandidateDigest              string
	RequestedCapability          string
	GuardianDecision             string
	PermissionState              string
	WorkspaceSource              string
	Gateway                      string
	Model                        string
	NetworkAllowlistDigest      string
	SuspendTokenDigest           string
	ResumeTokenDigest            string
	GuardianPolicyEvidenceDigest string
	EvidencePrefixDigest         string
	PlanEvidenceDigest           string
	EvidenceDigest               string
	NonExecuting                 bool
	NonAuthorizing               bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.PlanStatus != "planned" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.RequestedCapability == "" ||
		b.GuardianDecision == "" ||
		b.PermissionState != "not-authorized" ||
		b.WorkspaceSource == "" ||
		b.Gateway == "" ||
		b.Model == "" ||
		b.NetworkAllowlistDigest == "" ||
		b.SuspendTokenDigest == "" ||
		b.ResumeTokenDigest == "" ||
		b.GuardianPolicyEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.PlanEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision execution plan binding")
	}
	expectedResumeDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanResume(
		b.CandidateDigest,
		b.WorkspaceSource,
		b.Gateway,
		b.NetworkAllowlistDigest,
		b.SuspendTokenDigest,
	)
	if b.ResumeTokenDigest != expectedResumeDigest {
		return fmt.Errorf("Gooo extended lineage candidate revision resume token digest mismatch")
	}
	expectedPlanDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanPlan(
		b.CandidateDigest,
		b.RequestedCapability,
		b.GuardianDecision,
		b.PermissionState,
		b.WorkspaceSource,
		b.Gateway,
		b.Model,
		b.NetworkAllowlistDigest,
		b.SuspendTokenDigest,
		b.ResumeTokenDigest,
		b.GuardianPolicyEvidenceDigest,
		b.EvidencePrefixDigest,
	)
	if b.PlanEvidenceDigest != expectedPlanDigest {
		return fmt.Errorf("Gooo extended lineage candidate revision plan evidence digest mismatch")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision execution plan must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision execution plan digest mismatch")
	}
	return nil
}

// PlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan
// creates an auditable plan only; a separate deterministic policy must authorize any action.
func PlanExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-execution-plan"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding{
			Status:          "UNKNOWN",
			MissingStage:    stage,
			PermissionState: "not-authorized",
			NonExecuting:    true,
			NonAuthorizing:  true,
		}
	}
	if !input.NonAuthorizing || !input.Guardian.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Guardian.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Guardian.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-guardian-validation")
	}
	if strings.TrimSpace(input.WorkspaceSource) == "" {
		return unknown("workspace-source")
	}
	if strings.TrimSpace(input.Gateway) == "" {
		return unknown("gateway")
	}
	if strings.TrimSpace(input.Model) == "" {
		return unknown("model")
	}
	if strings.TrimSpace(input.NetworkAllowlistDigest) == "" {
		return unknown("network-allowlist")
	}
	if strings.TrimSpace(input.SuspendTokenDigest) == "" {
		return unknown("suspend-token")
	}
	resumeDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanResume(
		input.Guardian.CandidateDigest,
		input.WorkspaceSource,
		input.Gateway,
		input.NetworkAllowlistDigest,
		input.SuspendTokenDigest,
	)
	planDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanPlan(
		input.Guardian.CandidateDigest,
		input.Guardian.RequestedCapability,
		input.Guardian.GuardianDecision,
		input.Guardian.PermissionState,
		input.WorkspaceSource,
		input.Gateway,
		input.Model,
		input.NetworkAllowlistDigest,
		input.SuspendTokenDigest,
		resumeDigest,
		input.Guardian.GuardianPolicyEvidenceDigest,
		input.Guardian.EvidencePrefixDigest,
	)
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding{
		Status:                       "bound",
		PlanStatus:                   "planned",
		CandidateStatus:              input.Guardian.CandidateStatus,
		CandidateDigest:              input.Guardian.CandidateDigest,
		RequestedCapability:          input.Guardian.RequestedCapability,
		GuardianDecision:             input.Guardian.GuardianDecision,
		PermissionState:              "not-authorized",
		WorkspaceSource:              input.WorkspaceSource,
		Gateway:                      input.Gateway,
		Model:                        input.Model,
		NetworkAllowlistDigest:       input.NetworkAllowlistDigest,
		SuspendTokenDigest:           input.SuspendTokenDigest,
		ResumeTokenDigest:            resumeDigest,
		GuardianPolicyEvidenceDigest: input.Guardian.GuardianPolicyEvidenceDigest,
		EvidencePrefixDigest:         input.Guardian.EvidencePrefixDigest,
		PlanEvidenceDigest:           planDigest,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-execution-plan-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanResume(candidateDigest, workspaceSource, gateway, networkAllowlistDigest, suspendTokenDigest string) string {
	digest, err := Digest(struct {
		CandidateDigest        string
		WorkspaceSource        string
		Gateway                string
		NetworkAllowlistDigest string
		SuspendTokenDigest     string
	}{
		CandidateDigest:        candidateDigest,
		WorkspaceSource:        workspaceSource,
		Gateway:                gateway,
		NetworkAllowlistDigest: networkAllowlistDigest,
		SuspendTokenDigest:     suspendTokenDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanPlan(candidateDigest, requestedCapability, guardianDecision, permissionState, workspaceSource, gateway, model, networkAllowlistDigest, suspendTokenDigest, resumeTokenDigest, guardianPolicyEvidenceDigest, evidencePrefixDigest string) string {
	digest, err := Digest(struct {
		CandidateDigest              string
		RequestedCapability          string
		GuardianDecision             string
		PermissionState              string
		WorkspaceSource              string
		Gateway                      string
		Model                        string
		NetworkAllowlistDigest       string
		SuspendTokenDigest           string
		ResumeTokenDigest            string
		GuardianPolicyEvidenceDigest string
		EvidencePrefixDigest         string
	}{
		CandidateDigest:              candidateDigest,
		RequestedCapability:          requestedCapability,
		GuardianDecision:             guardianDecision,
		PermissionState:              permissionState,
		WorkspaceSource:              workspaceSource,
		Gateway:                      gateway,
		Model:                        model,
		NetworkAllowlistDigest:       networkAllowlistDigest,
		SuspendTokenDigest:           suspendTokenDigest,
		ResumeTokenDigest:            resumeTokenDigest,
		GuardianPolicyEvidenceDigest: guardianPolicyEvidenceDigest,
		EvidencePrefixDigest:         evidencePrefixDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlan(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanBinding) string {
	digest, err := Digest(struct {
		Status                       string
		PlanStatus                   string
		CandidateStatus              string
		CandidateDigest              string
		RequestedCapability          string
		GuardianDecision             string
		PermissionState              string
		WorkspaceSource              string
		Gateway                      string
		Model                        string
		NetworkAllowlistDigest       string
		SuspendTokenDigest           string
		ResumeTokenDigest            string
		GuardianPolicyEvidenceDigest string
		EvidencePrefixDigest         string
		PlanEvidenceDigest           string
		NonExecuting                 bool
		NonAuthorizing               bool
	}{
		Status:                       b.Status,
		PlanStatus:                   b.PlanStatus,
		CandidateStatus:              b.CandidateStatus,
		CandidateDigest:              b.CandidateDigest,
		RequestedCapability:          b.RequestedCapability,
		GuardianDecision:             b.GuardianDecision,
		PermissionState:              b.PermissionState,
		WorkspaceSource:              b.WorkspaceSource,
		Gateway:                      b.Gateway,
		Model:                        b.Model,
		NetworkAllowlistDigest:       b.NetworkAllowlistDigest,
		SuspendTokenDigest:           b.SuspendTokenDigest,
		ResumeTokenDigest:            b.ResumeTokenDigest,
		GuardianPolicyEvidenceDigest: b.GuardianPolicyEvidenceDigest,
		EvidencePrefixDigest:         b.EvidencePrefixDigest,
		PlanEvidenceDigest:           b.PlanEvidenceDigest,
		NonExecuting:                 b.NonExecuting,
		NonAuthorizing:               b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}