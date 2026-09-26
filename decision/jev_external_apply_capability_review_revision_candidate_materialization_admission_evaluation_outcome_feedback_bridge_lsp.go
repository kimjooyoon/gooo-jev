package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPBoundCode = "jev.provenance.candidate-evaluation-outcome-feedback-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPUnknownCode = "jev.provenance.unknown"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic struct {
	JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic
	FeedbackDirection     string
	FeedbackTarget        string
	FeedbackDigest        string
	FeedbackSource        string
	FeedbackEvidenceDigest string
	FeedbackStatus        string
	FeedbackBridgeDigest  string
	ProjectionDigest      string
	Publishable           bool
	NonExecuting          bool
	NonAuthorizing        bool
}

func (d JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic) Validate() error {
	if d.Severity == "" || d.Code == "" || d.Message == "" || d.Status == "" || d.ProjectionDigest == "" {
		return fmt.Errorf("incomplete JEV candidate evaluation feedback LSP diagnostic")
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("JEV candidate evaluation feedback LSP diagnostic must be non-executing and non-authorizing")
	}
	switch d.Status {
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown:
		if d.Severity != "warning" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPUnknownCode ||
			d.MissingStage == "" || d.Publishable || d.FeedbackBridgeDigest != "" {
			return fmt.Errorf("unknown JEV candidate evaluation feedback LSP diagnostic is inconsistent")
		}
	case jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeBound:
		if d.Severity != "info" || d.Code != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPBoundCode ||
			!d.Publishable || d.FeedbackBridgeDigest == "" || d.FeedbackStatus == "" {
			return fmt.Errorf("bound JEV candidate evaluation feedback LSP diagnostic is incomplete")
		}
		if d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.ProjectionDigest == "" ||
			d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.OutcomeBridgeDigest == "" ||
			!d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.Publishable {
			return fmt.Errorf("incomplete prior candidate outcome LSP diagnostic")
		}
		switch d.AdmissionDecision {
		case "admit":
			if d.FeedbackStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeFeedbackBound ||
				d.FeedbackDirection == "" || d.FeedbackTarget == "" || d.FeedbackDigest == "" ||
				d.FeedbackSource == "" || d.FeedbackEvidenceDigest == "" {
				return fmt.Errorf("admitted candidate feedback LSP diagnostic lost evidence")
			}
			switch d.FeedbackDirection {
			case "improve", "retain", "rollback":
			default:
				return fmt.Errorf("invalid candidate feedback direction")
			}
		case "hold", "reject":
			if d.FeedbackDirection != "" || d.FeedbackTarget != "" || d.FeedbackDigest != "" ||
				d.FeedbackSource != "" || d.FeedbackEvidenceDigest != "" {
				return fmt.Errorf("held or rejected candidate feedback LSP diagnostic contains evidence")
			}
		default:
			return fmt.Errorf("invalid candidate feedback admission decision")
		}
	default:
		return fmt.Errorf("invalid candidate evaluation feedback LSP status")
	}
	expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic(
		d.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.ProjectionDigest,
		d.FeedbackDirection, d.FeedbackTarget, d.FeedbackDigest, d.FeedbackSource,
		d.FeedbackEvidenceDigest, d.FeedbackStatus, d.FeedbackBridgeDigest, d.Publishable,
	)
	if d.ProjectionDigest != expected {
		return fmt.Errorf("JEV candidate evaluation feedback LSP projection digest mismatch")
	}
	return nil
}

func ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSP(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridge) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic {
	prior := ProjectJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSP(
		input.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridge,
	)
	output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic{
		JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic: prior,
		FeedbackDirection:      input.FeedbackDirection,
		FeedbackTarget:         input.FeedbackTarget,
		FeedbackDigest:         input.FeedbackDigest,
		FeedbackSource:         input.FeedbackSource,
		FeedbackEvidenceDigest: input.FeedbackEvidenceDigest,
		FeedbackStatus:         input.FeedbackStatus,
		FeedbackBridgeDigest:   input.BridgeDigest,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	if input.Status == jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown ||
		input.MissingStage != "" || !prior.Publishable {
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPUnknownCode
		output.Message = "candidate evaluation feedback provenance is incomplete"
		output.Publishable = false
		output.FeedbackBridgeDigest = ""
		if output.MissingStage == "" {
			output.MissingStage = "evaluation-outcome-feedback-bridge"
		}
	} else {
		output.Severity = "info"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPBoundCode
		output.Message = "candidate evaluation feedback provenance is bound"
		output.Publishable = true
	}
	output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic(
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.ProjectionDigest,
		output.FeedbackDirection, output.FeedbackTarget, output.FeedbackDigest,
		output.FeedbackSource, output.FeedbackEvidenceDigest, output.FeedbackStatus,
		output.FeedbackBridgeDigest, output.Publishable,
	)
	if err := output.Validate(); err != nil {
		output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.Status =
			jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationBridgeUnknown
		output.Severity = "warning"
		output.Code = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPUnknownCode
		output.Message = "candidate evaluation feedback provenance is incomplete"
		output.MissingStage = "lsp-projection"
		output.Publishable = false
		output.FeedbackBridgeDigest = ""
		output.ProjectionDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic(
			output.JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeBridgeLSPDiagnostic.ProjectionDigest,
			output.FeedbackDirection, output.FeedbackTarget, output.FeedbackDigest,
			output.FeedbackSource, output.FeedbackEvidenceDigest, output.FeedbackStatus,
			output.FeedbackBridgeDigest, output.Publishable,
		)
	}
	return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionEvaluationOutcomeFeedbackBridgeLSPDiagnostic(
	priorProjectionDigest, feedbackDirection, feedbackTarget, feedbackDigest, feedbackSource,
	feedbackEvidenceDigest, feedbackStatus, feedbackBridgeDigest string, publishable bool,
) string {
	values := []string{
		priorProjectionDigest, feedbackDirection, feedbackTarget, feedbackDigest, feedbackSource,
		feedbackEvidenceDigest, feedbackStatus, feedbackBridgeDigest, strconv.FormatBool(publishable),
	}
	sum := sha256.Sum256([]byte(strings.Join(values, "|")))
	return hex.EncodeToString(sum[:])
}

