package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound = "revision-candidate-application-candidate-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady = "application-candidate-ready"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold = "application-candidate-hold"
const jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected = "application-candidate-rejected"

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeInput struct {
    CandidateGate             JEVExternalApplyCapabilityReviewRevisionCandidateGate
    ApplicationTarget         string
    ApplicationSource         string
    ApplicationEvidenceDigest string
    NonAuthorizing            bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge struct {
    Status                    string
    MissingStage              string
    CandidateGateStatus       string
    CandidateDecision         string
    CandidateDigest           string
    RevisionSource            string
    CandidateGateDigest       string
    ApplicationCandidateStatus string
    ApplicationTarget         string
    ApplicationSource         string
    ApplicationEvidenceDigest string
    BridgeDigest              string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge) Validate() error {
    if b.Status == "" || b.CandidateGateStatus == "" || b.CandidateDecision == "" ||
        b.CandidateGateDigest == "" || b.ApplicationCandidateStatus == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV candidate application candidate bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound {
        return fmt.Errorf("invalid JEV candidate application candidate bridge status")
    }
    if b.CandidateGateStatus != jevExternalApplyCapabilityReviewRevisionCandidateGated {
        return fmt.Errorf("invalid JEV candidate gate status")
    }
    switch b.CandidateDecision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if b.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady ||
            b.CandidateDigest == "" || b.RevisionSource == "" || b.ApplicationTarget == "" ||
            b.ApplicationSource == "" || b.ApplicationEvidenceDigest == "" {
            return fmt.Errorf("ready JEV candidate application candidate bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        if b.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold ||
            b.ApplicationTarget != "" || b.ApplicationSource != "" || b.ApplicationEvidenceDigest != "" {
            return fmt.Errorf("held JEV candidate application candidate bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        if b.ApplicationCandidateStatus != jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected ||
            b.ApplicationTarget != "" || b.ApplicationSource != "" || b.ApplicationEvidenceDigest != "" {
            return fmt.Errorf("rejected JEV candidate application candidate bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid JEV candidate gate decision")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV candidate application candidate bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate application candidate bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(
        b.Status,
        b.CandidateGateStatus,
        b.CandidateDecision,
        b.CandidateDigest,
        b.RevisionSource,
        b.CandidateGateDigest,
        b.ApplicationCandidateStatus,
        b.ApplicationTarget,
        b.ApplicationSource,
        b.ApplicationEvidenceDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate application candidate bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.CandidateGate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.CandidateGate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.CandidateGate.Validate(); err != nil {
        output.MissingStage = "revision-candidate-gate"
        return output
    }

    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeBound
    output.CandidateGateStatus = input.CandidateGate.Status
    output.CandidateDecision = input.CandidateGate.Decision
    output.CandidateDigest = input.CandidateGate.CandidateDigest
    output.RevisionSource = input.CandidateGate.CandidateSource
    output.CandidateGateDigest = input.CandidateGate.GateDigest
    switch input.CandidateGate.Decision {
    case jevExternalApplyCapabilityReviewRevisionCandidateReady:
        if input.ApplicationTarget == "" {
            output.MissingStage = "application-target"
            return output
        }
        if input.ApplicationSource == "" {
            output.MissingStage = "application-source"
            return output
        }
        if input.ApplicationEvidenceDigest == "" {
            output.MissingStage = "application-evidence"
            return output
        }
        output.ApplicationCandidateStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateReady
        output.ApplicationTarget = input.ApplicationTarget
        output.ApplicationSource = input.ApplicationSource
        output.ApplicationEvidenceDigest = input.ApplicationEvidenceDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateHold:
        output.ApplicationCandidateStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateHold
    case jevExternalApplyCapabilityReviewRevisionCandidateRejected:
        output.ApplicationCandidateStatus = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateRejected
    default:
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown
        output.MissingStage = "candidate-gate-decision"
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(
        output.Status,
        output.CandidateGateStatus,
        output.CandidateDecision,
        output.CandidateDigest,
        output.RevisionSource,
        output.CandidateGateDigest,
        output.ApplicationCandidateStatus,
        output.ApplicationTarget,
        output.ApplicationSource,
        output.ApplicationEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridgeUnknown
        output.MissingStage = "candidate-application-candidate-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateApplicationCandidateBridge(status, candidateGateStatus, candidateDecision, candidateDigest, revisionSource, candidateGateDigest, applicationCandidateStatus, applicationTarget, applicationSource, applicationEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", status, candidateGateStatus, candidateDecision, candidateDigest, revisionSource, candidateGateDigest, applicationCandidateStatus, applicationTarget, applicationSource, applicationEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
