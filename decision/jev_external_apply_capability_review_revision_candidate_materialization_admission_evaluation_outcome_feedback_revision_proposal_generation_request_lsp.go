package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPBoundCode = "jev.provenance.revision-proposal-generation-request-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic
	GenerationRequestStatus       string
	GenerationRequestDigest       string
	GenerationRequestSource       string
	GenerationRequestEvidenceDigest string
	GeneratorIdentity             string
	GenerationRequestBridgeDigest string
	ProjectionDigest              string
	Publishable                   bool
	NonExecuting                  bool
	NonAuthorizing                bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic) Validate() error {
	if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.ProjectionDigest == "" {
		return fmt.Errorf("incomplete JEV generation request LSP diagnostic")
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("JEV generation request LSP diagnostic must be non-executing and non-authorizing")
	}
	switch d.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown:
		if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPUnknownCode ||
			d.MissingStage == "" || d.Publishable || d.GenerationRequestBridgeDigest != "" {
			return fmt.Errorf("unknown JEV generation request LSP diagnostic is inconsistent")
		}
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeBound:
		if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPBoundCode ||
			!d.Publishable || d.GenerationRequestBridgeDigest == "" || d.GenerationRequestStatus == "" {
			return fmt.Errorf("bound JEV generation request LSP diagnostic is incomplete")
		}
		if d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.ProjectionDigest == "" ||
			d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.ProposalBridgeDigest == "" ||
			!d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.Publishable {
			return fmt.Errorf("incomplete prior revision proposal LSP diagnostic")
		}
		switch d.ProposalDecision {
		case "revise":
			if d.GenerationRequestStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBound ||
				d.GenerationRequestDigest == "" || d.GenerationRequestSource == "" ||
				d.GenerationRequestEvidenceDigest == "" || d.GeneratorIdentity == "" {
				return fmt.Errorf("revise generation request LSP diagnostic lost evidence")
			}
		case "retain", "rollback":
			if d.GenerationRequestDigest != "" || d.GenerationRequestSource != "" ||
				d.GenerationRequestEvidenceDigest != "" || d.GeneratorIdentity != "" {
				return fmt.Errorf("retain or rollback generation request LSP diagnostic contains evidence")
			}
		default:
			return fmt.Errorf("invalid revision proposal decision")
		}
	default:
		return fmt.Errorf("invalid generation request LSP status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic(
		d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.ProjectionDigest,
		d.GenerationRequestStatus, d.GenerationRequestDigest, d.GenerationRequestSource,
		d.GenerationRequestEvidenceDigest, d.GeneratorIdentity,
		d.GenerationRequestBridgeDigest, d.Publishable,
	)
	if d.ProjectionDigest != expected {
		return fmt.Errorf("JEV generation request LSP projection digest mismatch")
	}
	return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic {
	prior := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSP(
		input.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge,
	)
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic: prior,
		GenerationRequestStatus:        input.GenerationRequestStatus,
		GenerationRequestDigest:        input.GenerationRequestDigest,
		GenerationRequestSource:        input.GenerationRequestSource,
		GenerationRequestEvidenceDigest: input.GenerationRequestEvidenceDigest,
		GeneratorIdentity:              input.GeneratorIdentity,
		GenerationRequestBridgeDigest:  input.BridgeDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		input.MissingStage != "" || !prior.Publishable {
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPUnknownCode
		output.Message = "generation request provenance is incomplete"
		output.Publishable = false
		output.GenerationRequestBridgeDigest = ""
		if output.MissingStage == "" {
			output.MissingStage = "revision-proposal-generation-request"
		}
	} else {
		output.Severity = "info"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPBoundCode
		output.Message = "generation request provenance is bound"
		output.Publishable = true
	}
	output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.ProjectionDigest,
		output.GenerationRequestStatus, output.GenerationRequestDigest,
		output.GenerationRequestSource, output.GenerationRequestEvidenceDigest,
		output.GeneratorIdentity, output.GenerationRequestBridgeDigest, output.Publishable,
	)
	if err := output.Validate(); err != nil {
		output.Status =
			jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.Status =
			jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPUnknownCode
		output.Message = "generation request provenance is incomplete"
		output.MissingStage = "lsp-projection"
		output.Publishable = false
		output.GenerationRequestBridgeDigest = ""
		output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic(
			output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic.ProjectionDigest,
			output.GenerationRequestStatus, output.GenerationRequestDigest,
			output.GenerationRequestSource, output.GenerationRequestEvidenceDigest,
			output.GeneratorIdentity, output.GenerationRequestBridgeDigest, output.Publishable,
		)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalGenerationRequestBridgeLSPDiagnostic(
	priorProjectionDigest, generationRequestStatus, generationRequestDigest, generationRequestSource,
	generationRequestEvidenceDigest, generatorIdentity, generationRequestBridgeDigest string,
	publishable bool,
) string {
	values := []string{
		priorProjectionDigest, generationRequestStatus, generationRequestDigest,
		generationRequestSource, generationRequestEvidenceDigest, generatorIdentity,
		generationRequestBridgeDigest, strconv.FormatBool(publishable),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

