package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevExternalApplyOriginBindingReady = "origin-bound-external-apply"
const jevExternalApplyOriginBindingUnknown = "UNKNOWN"

type JEVExternalApplyOriginBindingInput struct {
    Origin         JEVImprovementOriginResolution
    Request        JEVExternalApplyRequest
    NonAuthorizing bool
}

type JEVExternalApplyOriginBinding struct {
    Status                  string
    MissingStage            string
    OriginDigest            string
    CandidateDigest         string
    RequestEvidenceDigest   string
    BindingEvidenceDigest   string
    EvidenceDigest          string
    NonExecuting            bool
    NonAuthorizing          bool
}

func (b JEVExternalApplyOriginBinding) Validate() error {
    if b.Status == "" || b.OriginDigest == "" || b.CandidateDigest == "" || b.RequestEvidenceDigest == "" || b.BindingEvidenceDigest == "" || b.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV external apply origin binding")
    }
    if b.Status != jevExternalApplyOriginBindingReady {
        return fmt.Errorf("invalid JEV external apply origin binding status")
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV external apply origin binding must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV external apply origin binding must be non-authorizing")
    }
    expectedBinding := digestJEVExternalApplyOriginBinding(b.OriginDigest, b.CandidateDigest, b.RequestEvidenceDigest)
    if b.BindingEvidenceDigest != expectedBinding {
        return fmt.Errorf("JEV external apply origin binding evidence mismatch")
    }
    expected := digestJEVExternalApplyOriginBindingEvidence(b.Status, b.OriginDigest, b.BindingEvidenceDigest)
    if b.EvidenceDigest != expected {
        return fmt.Errorf("JEV external apply origin binding digest mismatch")
    }
    return nil
}

func BindJEVExternalApplyToOrigin(input JEVExternalApplyOriginBindingInput) JEVExternalApplyOriginBinding {
    output := JEVExternalApplyOriginBinding{
        Status: jevExternalApplyOriginBindingUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Origin.NonAuthorizing || !input.Request.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Origin.NonExecuting || !input.Request.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Origin.Validate(); err != nil {
        output.MissingStage = "origin-resolution"
        return output
    }
    if err := input.Request.Validate(); err != nil {
        output.MissingStage = "external-apply-request"
        return output
    }
    if input.Origin.Status != jevImprovementOriginResolved {
        output.MissingStage = "origin-status"
        return output
    }
    if input.Request.Status != jevExternalApplyRequestReady {
        output.MissingStage = "external-apply-status"
        return output
    }
    if input.Origin.CandidateDigest != input.Request.CandidateDigest {
        output.MissingStage = "origin-candidate-binding"
        return output
    }
    output.Status = jevExternalApplyOriginBindingReady
    output.OriginDigest = input.Origin.OriginDigest
    output.CandidateDigest = input.Request.CandidateDigest
    output.RequestEvidenceDigest = input.Request.EvidenceDigest
    output.BindingEvidenceDigest = digestJEVExternalApplyOriginBinding(output.OriginDigest, output.CandidateDigest, output.RequestEvidenceDigest)
    output.EvidenceDigest = digestJEVExternalApplyOriginBindingEvidence(output.Status, output.OriginDigest, output.BindingEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevExternalApplyOriginBindingUnknown
        output.MissingStage = "origin-apply-binding-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVExternalApplyOriginBinding(originDigest, candidateDigest, requestEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", originDigest, candidateDigest, requestEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}

func digestJEVExternalApplyOriginBindingEvidence(status, originDigest, bindingEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", status, originDigest, bindingEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
