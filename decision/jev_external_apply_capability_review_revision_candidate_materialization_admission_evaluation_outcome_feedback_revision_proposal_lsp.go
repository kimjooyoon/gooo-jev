package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPBoundCode = "jev.provenance.candidate-evaluation-feedback-revision-proposal-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic
	RevisionProposalStatus string
	ProposalDecision       string
	ProposalTarget         string
	ProposalDigest         string
	ProposalSource         string
	ProposalEvidenceDigest string
	ProposalBridgeDigest   string
	ProjectionDigest       string
	Publishable            bool
	NonExecuting           bool
	NonAuthorizing         bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic) Validate() error {
	if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.ProjectionDigest == "" {
		return fmt.Errorf("incomplete JEV candidate revision proposal LSP diagnostic")
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("JEV candidate revision proposal LSP diagnostic must be non-executing and non-authorizing")
	}
	switch d.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown:
		if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPUnknownCode ||
			d.MissingStage == "" || d.Publishable || d.ProposalBridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate revision proposal LSP diagnostic is inconsistent")
		}
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeBound:
		if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPBoundCode ||
			!d.Publishable || d.ProposalBridgeDigest == "" || d.RevisionProposalStatus == "" {
			return fmt.Errorf("bound JEV candidate revision proposal LSP diagnostic is incomplete")
		}
		if d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.ProjectionDigest == "" ||
			d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.FeedbackBridgeDigest == "" ||
			!d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.Publishable {
			return fmt.Errorf("incomplete prior candidate feedback LSP diagnostic")
		}
		switch d.AdmissionDecision {
		case "admit":
			if d.RevisionProposalStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBound ||
				d.ProposalDecision == "" || d.ProposalTarget == "" || d.ProposalDigest == "" ||
				d.ProposalSource == "" || d.ProposalEvidenceDigest == "" ||
				!validJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalDecision(d.ProposalDecision, d.FeedbackDirection) {
				return fmt.Errorf("admitted candidate revision proposal LSP diagnostic lost evidence or consistency")
			}
		case "hold", "reject":
			if d.ProposalDecision != "" || d.ProposalTarget != "" || d.ProposalDigest != "" ||
				d.ProposalSource != "" || d.ProposalEvidenceDigest != "" {
				return fmt.Errorf("held or rejected candidate revision proposal LSP diagnostic contains evidence")
			}
		default:
			return fmt.Errorf("invalid candidate revision proposal admission decision")
		}
	default:
		return fmt.Errorf("invalid candidate revision proposal LSP status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic(
		d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.ProjectionDigest,
		d.RevisionProposalStatus, d.ProposalDecision, d.ProposalTarget, d.ProposalDigest,
		d.ProposalSource, d.ProposalEvidenceDigest, d.ProposalBridgeDigest, d.Publishable,
	)
	if d.ProjectionDigest != expected {
		return fmt.Errorf("JEV candidate revision proposal LSP projection digest mismatch")
	}
	return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic {
	prior := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSP(
		input.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge,
	)
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic: prior,
		RevisionProposalStatus: input.RevisionProposalStatus,
		ProposalDecision:       input.ProposalDecision,
		ProposalTarget:         input.ProposalTarget,
		ProposalDigest:         input.ProposalDigest,
		ProposalSource:         input.ProposalSource,
		ProposalEvidenceDigest: input.ProposalEvidenceDigest,
		ProposalBridgeDigest:   input.BridgeDigest,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		input.MissingStage != "" || !prior.Publishable {
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPUnknownCode
		output.Message = "candidate revision proposal provenance is incomplete"
		output.Publishable = false
		output.ProposalBridgeDigest = ""
		if output.MissingStage == "" {
			output.MissingStage = "evaluation-outcome-feedback-revision-proposal"
		}
	} else {
		output.Severity = "info"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPBoundCode
		output.Message = "candidate revision proposal provenance is bound"
		output.Publishable = true
	}
	output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.ProjectionDigest,
		output.RevisionProposalStatus, output.ProposalDecision, output.ProposalTarget,
		output.ProposalDigest, output.ProposalSource, output.ProposalEvidenceDigest,
		output.ProposalBridgeDigest, output.Publishable,
	)
	if err := output.Validate(); err != nil {
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.Status =
			jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPUnknownCode
		output.Message = "candidate revision proposal provenance is incomplete"
		output.MissingStage = "lsp-projection"
		output.Publishable = false
		output.ProposalBridgeDigest = ""
		output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic(
			output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic.ProjectionDigest,
			output.RevisionProposalStatus, output.ProposalDecision, output.ProposalTarget,
			output.ProposalDigest, output.ProposalSource, output.ProposalEvidenceDigest,
			output.ProposalBridgeDigest, output.Publishable,
		)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackRevisionProposalBridgeLSPDiagnostic(
	priorProjectionDigest, revisionProposalStatus, proposalDecision, proposalTarget, proposalDigest,
	proposalSource, proposalEvidenceDigest, proposalBridgeDigest string, publishable bool,
) string {
	values := []string{
		priorProjectionDigest, revisionProposalStatus, proposalDecision, proposalTarget,
		proposalDigest, proposalSource, proposalEvidenceDigest, proposalBridgeDigest,
		strconv.FormatBool(publishable),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

