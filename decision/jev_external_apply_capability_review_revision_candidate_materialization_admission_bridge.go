package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeBound = "candidate-materialization-admission-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown = "UNKNOWN"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBound = "candidate-admission-bound"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionHeld = "candidate-admission-held"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionRejected = "candidate-admission-rejected"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit = "admit"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionHold = "hold"
const jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionReject = "reject"

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeInput struct {
    Materialization             JEVExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterializationBridge
    AdmissionDecision           string
    AdmissionDigest             string
    AdmissionSource             string
    AdmissionEvidenceDigest     string
    QualityMetricDigest         string
    QualityMetricSource         string
    QualityMetricEvidenceDigest string
    ReviewDigest                string
    ReviewSource                string
    ReviewEvidenceDigest        string
    NonAuthorizing              bool
}

type JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge struct {
    Status                      string
    MissingStage                string
    CandidateDigest             string
    CandidateSource             string
    CandidateGateDigest         string
    GeneratedCandidateDigest    string
    GeneratedCandidateSource    string
    GenerationInputDigest       string
    GenerationEvidenceDigest   string
    GeneratorIdentity           string
    MaterializationStatus       string
    MaterializationBridgeDigest string
    AdmissionDecision           string
    AdmissionStatus             string
    AdmissionDigest             string
    AdmissionSource             string
    AdmissionEvidenceDigest    string
    QualityMetricDigest         string
    QualityMetricSource         string
    QualityMetricEvidenceDigest string
    ReviewDigest                string
    ReviewSource                string
    ReviewEvidenceDigest        string
    BridgeDigest                string
    NonExecuting                bool
    NonAuthorizing              bool
}

func (b JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge) Validate() error {
    if b.Status == "" || b.MaterializationBridgeDigest == "" ||
        b.AdmissionStatus == "" || b.BridgeDigest == "" {
        return fmt.Errorf("incomplete JEV candidate materialization admission bridge")
    }
    if b.Status != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeBound {
        return fmt.Errorf("invalid JEV candidate materialization admission bridge status")
    }
    switch b.MaterializationStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized:
        if b.CandidateDigest == "" || b.CandidateSource == "" || b.CandidateGateDigest == "" ||
            b.GeneratedCandidateDigest == "" || b.GeneratedCandidateSource == "" ||
            b.GenerationInputDigest == "" || b.GenerationEvidenceDigest == "" ||
            b.GeneratorIdentity == "" ||
            b.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit ||
            b.AdmissionStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBound ||
            b.AdmissionDigest == "" || b.AdmissionSource == "" || b.AdmissionEvidenceDigest == "" ||
            b.QualityMetricDigest == "" || b.QualityMetricSource == "" ||
            b.QualityMetricEvidenceDigest == "" ||
            b.ReviewDigest == "" || b.ReviewSource == "" || b.ReviewEvidenceDigest == "" {
            return fmt.Errorf("materialized candidate admission bridge is incomplete")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateHeld:
        if b.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionHold ||
            b.AdmissionStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionHeld ||
            b.AdmissionDigest != "" || b.AdmissionSource != "" || b.AdmissionEvidenceDigest != "" ||
            b.QualityMetricDigest != "" || b.QualityMetricSource != "" ||
            b.QualityMetricEvidenceDigest != "" ||
            b.ReviewDigest != "" || b.ReviewSource != "" || b.ReviewEvidenceDigest != "" {
            return fmt.Errorf("held candidate admission bridge has inconsistent evidence")
        }
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateRollback:
        if b.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionReject ||
            b.AdmissionStatus != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionRejected ||
            b.AdmissionDigest != "" || b.AdmissionSource != "" || b.AdmissionEvidenceDigest != "" ||
            b.QualityMetricDigest != "" || b.QualityMetricSource != "" ||
            b.QualityMetricEvidenceDigest != "" ||
            b.ReviewDigest != "" || b.ReviewSource != "" || b.ReviewEvidenceDigest != "" {
            return fmt.Errorf("rollback candidate admission bridge has inconsistent evidence")
        }
    default:
        return fmt.Errorf("invalid candidate materialization status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV candidate materialization admission bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV candidate materialization admission bridge must be non-authorizing")
    }
    expected := digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(
        b.Status,
        b.CandidateDigest,
        b.CandidateSource,
        b.CandidateGateDigest,
        b.GeneratedCandidateDigest,
        b.GeneratedCandidateSource,
        b.GenerationInputDigest,
        b.GenerationEvidenceDigest,
        b.GeneratorIdentity,
        b.MaterializationStatus,
        b.MaterializationBridgeDigest,
        b.AdmissionDecision,
        b.AdmissionStatus,
        b.AdmissionDigest,
        b.AdmissionSource,
        b.AdmissionEvidenceDigest,
        b.QualityMetricDigest,
        b.QualityMetricSource,
        b.QualityMetricEvidenceDigest,
        b.ReviewDigest,
        b.ReviewSource,
        b.ReviewEvidenceDigest,
    )
    if b.BridgeDigest != expected {
        return fmt.Errorf("JEV candidate materialization admission bridge digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(input JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeInput) JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge {
    output := JEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge{
        Status:         jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown,
        NonExecuting:   true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Materialization.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Materialization.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Materialization.Validate(); err != nil {
        output.MissingStage = "candidate-materialization"
        return output
    }
    output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeBound
    output.CandidateDigest = input.Materialization.CandidateDigest
    output.CandidateSource = input.Materialization.CandidateSource
    output.CandidateGateDigest = input.Materialization.CandidateGateDigest
    output.GeneratedCandidateDigest = input.Materialization.GeneratedCandidateDigest
    output.GeneratedCandidateSource = input.Materialization.GeneratedCandidateSource
    output.GenerationInputDigest = input.Materialization.GenerationInputDigest
    output.GenerationEvidenceDigest = input.Materialization.GenerationEvidenceDigest
    output.GeneratorIdentity = input.Materialization.GeneratorIdentity
    output.MaterializationStatus = input.Materialization.MaterializationStatus
    output.MaterializationBridgeDigest = input.Materialization.BridgeDigest

    switch input.Materialization.MaterializationStatus {
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateMaterialized:
        if input.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionAdmit {
            output.MissingStage = "admission-decision-consistency"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.AdmissionDigest == "" {
            output.MissingStage = "admission-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.AdmissionSource == "" {
            output.MissingStage = "admission-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.AdmissionEvidenceDigest == "" {
            output.MissingStage = "admission-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.QualityMetricDigest == "" {
            output.MissingStage = "quality-metric-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.QualityMetricSource == "" {
            output.MissingStage = "quality-metric-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.QualityMetricEvidenceDigest == "" {
            output.MissingStage = "quality-metric-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.ReviewDigest == "" {
            output.MissingStage = "review-digest"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.ReviewSource == "" {
            output.MissingStage = "review-source"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        if input.ReviewEvidenceDigest == "" {
            output.MissingStage = "review-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        output.AdmissionDecision = input.AdmissionDecision
        output.AdmissionStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBound
        output.AdmissionDigest = input.AdmissionDigest
        output.AdmissionSource = input.AdmissionSource
        output.AdmissionEvidenceDigest = input.AdmissionEvidenceDigest
        output.QualityMetricDigest = input.QualityMetricDigest
        output.QualityMetricSource = input.QualityMetricSource
        output.QualityMetricEvidenceDigest = input.QualityMetricEvidenceDigest
        output.ReviewDigest = input.ReviewDigest
        output.ReviewSource = input.ReviewSource
        output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateHeld:
        if input.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionHold ||
            input.AdmissionDigest != "" || input.AdmissionSource != "" || input.AdmissionEvidenceDigest != "" ||
            input.QualityMetricDigest != "" || input.QualityMetricSource != "" ||
            input.QualityMetricEvidenceDigest != "" ||
            input.ReviewDigest != "" || input.ReviewSource != "" || input.ReviewEvidenceDigest != "" {
            output.MissingStage = "held-admission-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        output.AdmissionDecision = input.AdmissionDecision
        output.AdmissionStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionHeld
    case jevExternalApplyCapabilityReviewRevisionCandidateCapabilityScopeFeedbackRevisionProposalCandidateRollback:
        if input.AdmissionDecision != jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionReject ||
            input.AdmissionDigest != "" || input.AdmissionSource != "" || input.AdmissionEvidenceDigest != "" ||
            input.QualityMetricDigest != "" || input.QualityMetricSource != "" ||
            input.QualityMetricEvidenceDigest != "" ||
            input.ReviewDigest != "" || input.ReviewSource != "" || input.ReviewEvidenceDigest != "" {
            output.MissingStage = "rollback-admission-evidence"
            output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
            return output
        }
        output.AdmissionDecision = input.AdmissionDecision
        output.AdmissionStatus = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionRejected
    default:
        output.MissingStage = "materialization-status"
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
        return output
    }
    output.BridgeDigest = digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(
        output.Status,
        output.CandidateDigest,
        output.CandidateSource,
        output.CandidateGateDigest,
        output.GeneratedCandidateDigest,
        output.GeneratedCandidateSource,
        output.GenerationInputDigest,
        output.GenerationEvidenceDigest,
        output.GeneratorIdentity,
        output.MaterializationStatus,
        output.MaterializationBridgeDigest,
        output.AdmissionDecision,
        output.AdmissionStatus,
        output.AdmissionDigest,
        output.AdmissionSource,
        output.AdmissionEvidenceDigest,
        output.QualityMetricDigest,
        output.QualityMetricSource,
        output.QualityMetricEvidenceDigest,
        output.ReviewDigest,
        output.ReviewSource,
        output.ReviewEvidenceDigest,
    )
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridgeUnknown
        output.MissingStage = "candidate-materialization-admission-evidence"
        output.BridgeDigest = ""
    }
    return output
}

func digestJEVExternalApplyCapabilityReviewRevisionCandidateMaterializationAdmissionBridge(status, candidateDigest, candidateSource, candidateGateDigest, generatedCandidateDigest, generatedCandidateSource, generationInputDigest, generationEvidenceDigest, generatorIdentity, materializationStatus, materializationBridgeDigest, admissionDecision, admissionStatus, admissionDigest, admissionSource, admissionEvidenceDigest, qualityMetricDigest, qualityMetricSource, qualityMetricEvidenceDigest, reviewDigest, reviewSource, reviewEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(strings.Join([]string{
        status, candidateDigest, candidateSource, candidateGateDigest,
        generatedCandidateDigest, generatedCandidateSource, generationInputDigest,
        generationEvidenceDigest, generatorIdentity, materializationStatus,
        materializationBridgeDigest, admissionDecision, admissionStatus, admissionDigest,
        admissionSource, admissionEvidenceDigest, qualityMetricDigest, qualityMetricSource,
        qualityMetricEvidenceDigest, reviewDigest, reviewSource, reviewEvidenceDigest,
    }, "|")))
    return hex.EncodeToString(sum[:])
}
