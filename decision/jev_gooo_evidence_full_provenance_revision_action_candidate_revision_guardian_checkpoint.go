package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput
// carries a calibrated review to a deterministic guardian checkpoint without granting capability.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput struct {
	Confidence                    ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding
	RequestedCapability           string
	GuardianPolicyEvidenceDigest  string
	NonAuthorizing                bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding
// records the guardian boundary as a non-authorizing policy checkpoint.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding struct {
	Status                       string
	MissingStage                 string
	GuardianStatus               string
	CandidateStatus              string
	CandidateDigest              string
	CandidateSource              string
	RevisionSource               string
	ReviewDisposition            string
	ConfidenceBand               string
	AutomationMode               string
	RequestedCapability          string
	GuardianDecision             string
	PermissionState              string
	GuardianPolicyEvidenceDigest string
	FallbackStage                string
	EvidencePrefixDigest         string
	EvidenceDigest               string
	NonExecuting                 bool
	NonAuthorizing               bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.GuardianStatus != "policy-checkpoint" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.RevisionSource == "" ||
		b.ReviewDisposition == "" ||
		b.ConfidenceBand == "" ||
		b.AutomationMode == "" ||
		b.RequestedCapability == "" ||
		b.GuardianDecision == "" ||
		b.PermissionState != "not-authorized" ||
		b.GuardianPolicyEvidenceDigest == "" ||
		b.FallbackStage == "" ||
		b.EvidencePrefixDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision guardian checkpoint binding")
	}
	expectedDecision := guardianDecisionForGoooExtendedLineageCandidateRevision(b.ReviewDisposition, b.AutomationMode)
	if expectedDecision == "" || b.GuardianDecision != expectedDecision {
		return fmt.Errorf("invalid Gooo extended lineage candidate revision guardian decision")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision guardian checkpoint must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision guardian checkpoint digest mismatch")
	}
	return nil
}

// EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint
// keeps permission denied until deterministic policy authorizes independently.
func EvaluateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-guardian-checkpoint"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			PermissionState: "not-authorized",
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Confidence.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Confidence.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Confidence.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-confidence-validation")
	}
	if strings.TrimSpace(input.RequestedCapability) == "" {
		return unknown("requested-capability")
	}
	if strings.TrimSpace(input.GuardianPolicyEvidenceDigest) == "" {
		return unknown("guardian-policy-evidence")
	}
	decision := guardianDecisionForGoooExtendedLineageCandidateRevision(input.Confidence.ReviewDisposition, input.Confidence.AutomationMode)
	if decision == "" {
		return unknown("guardian-decision")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding{
		Status:                       "bound",
		GuardianStatus:               "policy-checkpoint",
		CandidateStatus:              input.Confidence.CandidateStatus,
		CandidateDigest:              input.Confidence.CandidateDigest,
		CandidateSource:              input.Confidence.CandidateSource,
		RevisionSource:               input.Confidence.RevisionSource,
		ReviewDisposition:            input.Confidence.ReviewDisposition,
		ConfidenceBand:               input.Confidence.ConfidenceBand,
		AutomationMode:               input.Confidence.AutomationMode,
		RequestedCapability:          input.RequestedCapability,
		GuardianDecision:             decision,
		PermissionState:              "not-authorized",
		GuardianPolicyEvidenceDigest: input.GuardianPolicyEvidenceDigest,
		FallbackStage:                input.Confidence.FallbackStage,
		EvidencePrefixDigest:         input.Confidence.EvidencePrefixDigest,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-guardian-checkpoint-evidence")
	}
	return output
}

func guardianDecisionForGoooExtendedLineageCandidateRevision(reviewDisposition, automationMode string) string {
	switch {
	case reviewDisposition == "review-rejected":
		return "review-rejected"
	case reviewDisposition == "review-abstained":
		return "human-or-deterministic-policy"
	case automationMode == "shadow-only":
		return "shadow-observation-only"
	case automationMode == "reversible-review-only":
		return "deterministic-reversible-policy-review"
	case automationMode == "deterministic-policy-required":
		return "deterministic-policy-required"
	default:
		return ""
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpoint(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGuardianCheckpointBinding) string {
	digest, err := Digest(struct {
		Status                       string
		GuardianStatus               string
		CandidateStatus              string
		CandidateDigest              string
		CandidateSource              string
		RevisionSource               string
		ReviewDisposition            string
		ConfidenceBand               string
		AutomationMode               string
		RequestedCapability          string
		GuardianDecision             string
		PermissionState              string
		GuardianPolicyEvidenceDigest string
		FallbackStage                string
		EvidencePrefixDigest         string
		NonExecuting                 bool
		NonAuthorizing               bool
	}{
		Status:                       b.Status,
		GuardianStatus:               b.GuardianStatus,
		CandidateStatus:              b.CandidateStatus,
		CandidateDigest:              b.CandidateDigest,
		CandidateSource:              b.CandidateSource,
		RevisionSource:               b.RevisionSource,
		ReviewDisposition:            b.ReviewDisposition,
		ConfidenceBand:               b.ConfidenceBand,
		AutomationMode:               b.AutomationMode,
		RequestedCapability:          b.RequestedCapability,
		GuardianDecision:             b.GuardianDecision,
		PermissionState:              b.PermissionState,
		GuardianPolicyEvidenceDigest: b.GuardianPolicyEvidenceDigest,
		FallbackStage:                b.FallbackStage,
		EvidencePrefixDigest:         b.EvidencePrefixDigest,
		NonExecuting:                 b.NonExecuting,
		NonAuthorizing:               b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}