package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeBound = "revision-proposal-generation-request-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBound = "generation-request-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestHeld = "generation-request-held"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestRejected = "generation-request-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeInput struct {
	RevisionProposal         JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge
	GenerationRequestDigest  string
	GenerationRequestSource  string
	GenerationRequestEvidenceDigest string
	GeneratorIdentity        string
	NonAuthorizing           bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge
	GenerationRequestStatus       string
	GenerationRequestDigest       string
	GenerationRequestSource       string
	GenerationRequestEvidenceDigest string
	GeneratorIdentity              string
	BridgeDigest                   string
	NonExecuting                   bool
	NonAuthorizing                 bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge) Validate() error {
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("JEV generation request bridge must be non-executing and non-authorizing")
	}
	if b.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		if b.MissingStage == "" || b.BridgeDigest != "" {
			return fmt.Errorf("unknown JEV generation request bridge is inconsistent")
		}
		return nil
	}
	if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeBound ||
		b.BridgeDigest == "" {
		return fmt.Errorf("incomplete JEV generation request bridge")
	}
	if b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge.BridgeDigest == "" {
		return fmt.Errorf("incomplete prior revision proposal bridge")
	}
	switch b.ProposalDecision {
	case "revise":
		if b.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBound ||
			b.GenerationRequestDigest == "" || b.GenerationRequestSource == "" ||
			b.GenerationRequestEvidenceDigest == "" || b.GeneratorIdentity == "" {
			return fmt.Errorf("revise generation request lost evidence")
		}
	case "retain":
		if b.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestHeld ||
			b.GenerationRequestDigest != "" || b.GenerationRequestSource != "" ||
			b.GenerationRequestEvidenceDigest != "" || b.GeneratorIdentity != "" {
			return fmt.Errorf("retain generation request contains evidence")
		}
	case "rollback":
		if b.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestRejected ||
			b.GenerationRequestDigest != "" || b.GenerationRequestSource != "" ||
			b.GenerationRequestEvidenceDigest != "" || b.GeneratorIdentity != "" {
			return fmt.Errorf("rollback generation request contains evidence")
		}
	default:
		return fmt.Errorf("invalid revision proposal decision")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge(
		b.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge.BridgeDigest,
		b.GenerationRequestStatus, b.GenerationRequestDigest, b.GenerationRequestSource,
		b.GenerationRequestEvidenceDigest, b.GeneratorIdentity,
	)
	if b.BridgeDigest != expected {
		return fmt.Errorf("JEV generation request bridge digest mismatch")
	}
	return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge {
	unknown := func(stage string, prior JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge {
		if stage == "" {
			stage = "revision-proposal-generation-request"
		}
		prior.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		prior.MissingStage = stage
		prior.BridgeDigest = ""
		prior.NonExecuting = true
		prior.NonAuthorizing = true
		return JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge{
			JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge: prior,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	prior := input.RevisionProposal
	if !input.NonAuthorizing || !prior.NonExecuting || !prior.NonAuthorizing {
		return unknown("capability-boundary", prior)
	}
	if err := prior.Validate(); err != nil {
		return unknown("revision-proposal-bridge", prior)
	}
	if prior.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown {
		return unknown(prior.MissingStage, prior)
	}
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge: prior,
		GenerationRequestDigest:         input.GenerationRequestDigest,
		GenerationRequestSource:         input.GenerationRequestSource,
		GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
		GeneratorIdentity:               input.GeneratorIdentity,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	switch prior.ProposalDecision {
	case "revise":
		output.GenerationRequestStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBound
		stages := []struct {
			value string
			stage string
		}{
			{output.GenerationRequestDigest, "generation-request-digest"},
			{output.GenerationRequestSource, "generation-request-source"},
			{output.GenerationRequestEvidenceDigest, "generation-request-evidence"},
			{output.GeneratorIdentity, "generator-identity"},
		}
		for _, item := range stages {
			if item.value == "" {
				return unknown(item.stage, prior)
			}
		}
	case "retain":
		output.GenerationRequestStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestHeld
		if input.GenerationRequestDigest != "" || input.GenerationRequestSource != "" ||
			input.GenerationRequestEvidenceDigest != "" || input.GeneratorIdentity != "" {
			return unknown("retain-generation-evidence", prior)
		}
	case "rollback":
		output.GenerationRequestStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestRejected
		if input.GenerationRequestDigest != "" || input.GenerationRequestSource != "" ||
			input.GenerationRequestEvidenceDigest != "" || input.GeneratorIdentity != "" {
			return unknown("rollback-generation-evidence", prior)
		}
	default:
		return unknown("revision-proposal-decision", prior)
	}
	output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeBound
	output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge.BridgeDigest,
		output.GenerationRequestStatus, output.GenerationRequestDigest,
		output.GenerationRequestSource, output.GenerationRequestEvidenceDigest,
		output.GeneratorIdentity,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-proposal-generation-request", prior)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge(
	revisionProposalBridgeDigest, generationRequestStatus, generationRequestDigest,
	generationRequestSource, generationRequestEvidenceDigest, generatorIdentity string,
) string {
	values := []string{
		revisionProposalBridgeDigest, generationRequestStatus, generationRequestDigest,
		generationRequestSource, generationRequestEvidenceDigest, generatorIdentity,
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

