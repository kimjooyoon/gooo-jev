package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevImprovementReplayReady = "ready-for-replay-observation"
    jevImprovementReplayUnknown = "UNKNOWN"
)

type JEVImprovementReplayPreparationInput struct {
    Review              JEVImprovementCandidateReview
    ReplayScope         string
    ReplayEnvironment   string
    NonAuthorizing      bool
}

type JEVImprovementReplayPreparation struct {
    Status              string
    MissingStage        string
    CandidateDigest     string
    ReviewEvidenceDigest string
    ReplayScope         string
    ReplayEnvironment   string
    EvidenceDigest      string
    NonExecuting        bool
    NonAuthorizing      bool
}

func (p JEVImprovementReplayPreparation) Validate() error {
    if p.Status == "" || p.CandidateDigest == "" || p.ReviewEvidenceDigest == "" || p.ReplayScope == "" || p.ReplayEnvironment == "" || p.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement replay preparation")
    }
    if !p.NonExecuting {
        return fmt.Errorf("JEV improvement replay preparation must be non-executing")
    }
    if !p.NonAuthorizing {
        return fmt.Errorf("JEV improvement replay preparation must be non-authorizing")
    }
    expected := digestJEVImprovementReplayPreparation(p.Status, p.CandidateDigest, p.ReviewEvidenceDigest, p.ReplayScope, p.ReplayEnvironment)
    if p.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement replay preparation evidence digest mismatch")
    }
    return nil
}

func PrepareJEVImprovementReplay(input JEVImprovementReplayPreparationInput) JEVImprovementReplayPreparation {
    output := JEVImprovementReplayPreparation{
        Status: jevImprovementReplayUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.Review.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.Review.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.Review.Validate(); err != nil {
        output.MissingStage = "candidate-review"
        return output
    }
    if input.Review.Status != jevImprovementCandidateReviewConfirmed {
        output.MissingStage = "review-confirmation"
        return output
    }
    if input.ReplayScope == "" {
        output.MissingStage = "replay-scope"
        return output
    }
    if input.ReplayEnvironment == "" {
        output.MissingStage = "replay-environment"
        return output
    }
    output.Status = jevImprovementReplayReady
    output.CandidateDigest = input.Review.CandidateDigest
    output.ReviewEvidenceDigest = input.Review.EvidenceDigest
    output.ReplayScope = input.ReplayScope
    output.ReplayEnvironment = input.ReplayEnvironment
    output.EvidenceDigest = digestJEVImprovementReplayPreparation(output.Status, output.CandidateDigest, output.ReviewEvidenceDigest, output.ReplayScope, output.ReplayEnvironment)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementReplayUnknown
        output.MissingStage = "replay-preparation-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementReplayPreparation(status, candidateDigest, reviewEvidenceDigest, replayScope, replayEnvironment string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s", status, candidateDigest, reviewEvidenceDigest, replayScope, replayEnvironment)))
    return hex.EncodeToString(sum[:])
}
