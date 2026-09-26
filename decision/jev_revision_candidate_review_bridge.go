package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const jevRevisionCandidateReviewBridgeReady = "revision-candidate-selection-ready"
const jevRevisionCandidateReviewBridgeUnknown = "UNKNOWN"

type JEVRevisionCandidateReviewBridgeInput struct {
    Candidate      JEVImprovementRevisionCandidate
    NonAuthorizing bool
}

type JEVRevisionCandidateReviewBridge struct {
    Status                 string
    MissingStage           string
    Selection              JEVImprovementCandidateSelection
    ParentCandidateDigest  string
    RevisionChangeDigest   string
    CandidateEvidenceDigest string
    LineageEvidenceDigest  string
    EvidenceDigest         string
    NonExecuting           bool
    NonAuthorizing         bool
}

func (b JEVRevisionCandidateReviewBridge) Validate() error {
    if b.Status == "" || b.ParentCandidateDigest == "" || b.RevisionChangeDigest == "" || b.CandidateEvidenceDigest == "" || b.LineageEvidenceDigest == "" || b.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV revision candidate review bridge")
    }
    if b.Status != jevRevisionCandidateReviewBridgeReady {
        return fmt.Errorf("invalid JEV revision candidate review bridge status")
    }
    if err := b.Selection.Validate(); err != nil {
        return fmt.Errorf("invalid review selection: %w", err)
    }
    if !b.NonExecuting {
        return fmt.Errorf("JEV revision candidate review bridge must be non-executing")
    }
    if !b.NonAuthorizing {
        return fmt.Errorf("JEV revision candidate review bridge must be non-authorizing")
    }
    expectedLineage := digestJEVRevisionCandidateLineage(b.ParentCandidateDigest, b.RevisionChangeDigest, b.Selection.CandidateDigest, b.CandidateEvidenceDigest)
    if b.LineageEvidenceDigest != expectedLineage {
        return fmt.Errorf("JEV revision candidate lineage evidence mismatch")
    }
    expected := digestJEVRevisionCandidateReviewBridge(b.Status, b.Selection.EvidenceDigest, b.LineageEvidenceDigest)
    if b.EvidenceDigest != expected {
        return fmt.Errorf("JEV revision candidate review bridge evidence mismatch")
    }
    return nil
}

func BridgeJEVImprovementRevisionCandidateForReview(input JEVRevisionCandidateReviewBridgeInput) JEVRevisionCandidateReviewBridge {
    output := JEVRevisionCandidateReviewBridge{
        Status: jevRevisionCandidateReviewBridgeUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Candidate.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Candidate.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Candidate.Validate(); err != nil {
        output.MissingStage = "revision-candidate"
        return output
    }
    output.Status = jevRevisionCandidateReviewBridgeReady
    output.ParentCandidateDigest = input.Candidate.ParentCandidateDigest
    output.RevisionChangeDigest = input.Candidate.RevisionChangeDigest
    output.CandidateEvidenceDigest = input.Candidate.EvidenceDigest
    output.Selection = JEVImprovementCandidateSelection{
        Status: jevImprovementCandidateReady,
        CandidateDigest: input.Candidate.CandidateDigest,
        CandidateSource: input.Candidate.RevisionSource,
        ObservationEvidenceDigest: input.Candidate.EvidenceDigest,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    output.Selection.EvidenceDigest = digestJEVImprovementCandidateSelection(output.Selection.Status, output.Selection.CandidateDigest, output.Selection.CandidateSource, output.Selection.ObservationEvidenceDigest)
    output.LineageEvidenceDigest = digestJEVRevisionCandidateLineage(output.ParentCandidateDigest, output.RevisionChangeDigest, output.Selection.CandidateDigest, output.CandidateEvidenceDigest)
    output.EvidenceDigest = digestJEVRevisionCandidateReviewBridge(output.Status, output.Selection.EvidenceDigest, output.LineageEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevRevisionCandidateReviewBridgeUnknown
        output.MissingStage = "revision-review-bridge-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVRevisionCandidateLineage(parentCandidateDigest, revisionChangeDigest, candidateDigest, candidateEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s", parentCandidateDigest, revisionChangeDigest, candidateDigest, candidateEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}

func digestJEVRevisionCandidateReviewBridge(status, selectionEvidenceDigest, lineageEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", status, selectionEvidenceDigest, lineageEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
