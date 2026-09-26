package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevImprovementRevisionCandidateReady = "revision-candidate-ready-for-review"
const jevImprovementRevisionCandidateUnknown = "UNKNOWN"

type JEVImprovementRevisionCandidateInput struct {
    Directive           JEVImprovementDirectionDirective
    RevisionSource     string
    RevisionChangeDigest string
    NonAuthorizing     bool
}

type JEVImprovementRevisionCandidate struct {
    Status                 string
    MissingStage           string
    ParentCandidateDigest  string
    RevisionSource         string
    RevisionChangeDigest   string
    DirectiveEvidenceDigest string
    CandidateDigest        string
    EvidenceDigest         string
    NonExecuting           bool
    NonAuthorizing         bool
}

func (c JEVImprovementRevisionCandidate) Validate() error {
    if c.Status == "" || c.ParentCandidateDigest == "" || c.RevisionSource == "" || c.RevisionChangeDigest == "" || c.DirectiveEvidenceDigest == "" || c.CandidateDigest == "" || c.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement revision candidate")
    }
    if c.Status != jevImprovementRevisionCandidateReady {
        return fmt.Errorf("invalid JEV improvement revision candidate status")
    }
    if !c.NonExecuting {
        return fmt.Errorf("JEV improvement revision candidate must be non-executing")
    }
    if !c.NonAuthorizing {
        return fmt.Errorf("JEV improvement revision candidate must be non-authorizing")
    }
    expectedCandidate := digestJEVImprovementRevisionCandidate(c.ParentCandidateDigest, c.RevisionSource, c.RevisionChangeDigest, c.DirectiveEvidenceDigest)
    if c.CandidateDigest != expectedCandidate {
        return fmt.Errorf("JEV improvement revision candidate digest mismatch")
    }
    expected := digestJEVImprovementRevisionCandidateEvidence(c.Status, c.CandidateDigest, c.DirectiveEvidenceDigest, c.RevisionChangeDigest)
    if c.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement revision candidate evidence digest mismatch")
    }
    return nil
}

func GenerateJEVImprovementRevisionCandidate(input JEVImprovementRevisionCandidateInput) JEVImprovementRevisionCandidate {
    output := JEVImprovementRevisionCandidate{
        Status: jevImprovementRevisionCandidateUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Directive.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Directive.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Directive.Validate(); err != nil {
        output.MissingStage = "improvement-direction-directive"
        return output
    }
    if input.Directive.Directive != jevImprovementDirectiveRevision {
        output.MissingStage = "revision-directive"
        return output
    }
    if input.RevisionSource == "" {
        output.MissingStage = "revision-source"
        return output
    }
    if input.RevisionChangeDigest == "" {
        output.MissingStage = "revision-change"
        return output
    }
    output.Status = jevImprovementRevisionCandidateReady
    output.ParentCandidateDigest = input.Directive.CandidateDigest
    output.RevisionSource = input.RevisionSource
    output.RevisionChangeDigest = input.RevisionChangeDigest
    output.DirectiveEvidenceDigest = input.Directive.EvidenceDigest
    output.CandidateDigest = digestJEVImprovementRevisionCandidate(output.ParentCandidateDigest, output.RevisionSource, output.RevisionChangeDigest, output.DirectiveEvidenceDigest)
    output.EvidenceDigest = digestJEVImprovementRevisionCandidateEvidence(output.Status, output.CandidateDigest, output.DirectiveEvidenceDigest, output.RevisionChangeDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementRevisionCandidateUnknown
        output.MissingStage = "revision-candidate-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementRevisionCandidate(parentCandidateDigest, revisionSource, revisionChangeDigest, directiveEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", parentCandidateDigest, revisionSource, revisionChangeDigest, directiveEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}

func digestJEVImprovementRevisionCandidateEvidence(status, candidateDigest, directiveEvidenceDigest, revisionChangeDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", status, candidateDigest, directiveEvidenceDigest, revisionChangeDigest)))
    return hex.EncodeToString(sum[:])
}
