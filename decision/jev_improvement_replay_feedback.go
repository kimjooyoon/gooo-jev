package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

const (
    jevReplayFeedbackConfirmed = "confirmed"
    jevReplayFeedbackRefuted = "refuted"
    jevReplayFeedbackUnknown = "unknown"
)

type JEVImprovementReplayFeedbackInput struct {
    ReplayObservation JEVImprovementReplayObservation
    MetricDigest      string
    NonAuthorizing    bool
}

type JEVImprovementReplayFeedback struct {
    Status                    string
    FeedbackKind              string
    MissingStage              string
    CandidateDigest           string
    ReplayObservationDigest   string
    MetricDigest              string
    FeedbackEvidenceDigest    string
    EvidenceDigest            string
    NonExecuting              bool
    NonAuthorizing            bool
}

func (f JEVImprovementReplayFeedback) Validate() error {
    if f.Status == "" || f.FeedbackKind == "" || f.CandidateDigest == "" || f.ReplayObservationDigest == "" || f.MetricDigest == "" || f.FeedbackEvidenceDigest == "" || f.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement replay feedback")
    }
    switch f.FeedbackKind {
    case jevReplayFeedbackConfirmed, jevReplayFeedbackRefuted, jevReplayFeedbackUnknown:
    default:
        return fmt.Errorf("invalid JEV improvement replay feedback kind")
    }
    if !f.NonExecuting {
        return fmt.Errorf("JEV improvement replay feedback must be non-executing")
    }
    if !f.NonAuthorizing {
        return fmt.Errorf("JEV improvement replay feedback must be non-authorizing")
    }
    expectedFeedbackEvidence := digestJEVImprovementReplayFeedbackEvidence(f.ReplayObservationDigest, f.MetricDigest, f.FeedbackKind)
    if f.FeedbackEvidenceDigest != expectedFeedbackEvidence {
        return fmt.Errorf("JEV improvement replay feedback evidence digest mismatch")
    }
    expected := digestJEVImprovementReplayFeedback(f.Status, f.FeedbackKind, f.CandidateDigest, f.ReplayObservationDigest, f.MetricDigest, f.FeedbackEvidenceDigest)
    if f.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement replay feedback digest mismatch")
    }
    return nil
}

func DeriveJEVImprovementReplayFeedback(input JEVImprovementReplayFeedbackInput) JEVImprovementReplayFeedback {
    output := JEVImprovementReplayFeedback{
        Status: jevReplayFeedbackUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing || !input.ReplayObservation.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if !input.ReplayObservation.NonExecuting {
        output.MissingStage = "execution-boundary"
        return output
    }
    if err := input.ReplayObservation.Validate(); err != nil {
        output.MissingStage = "reverse-observation"
        return output
    }
    if input.MetricDigest == "" {
        output.MissingStage = "metric"
        return output
    }
    output.CandidateDigest = input.ReplayObservation.CandidateDigest
    output.ReplayObservationDigest = input.ReplayObservation.EvidenceDigest
    output.MetricDigest = input.MetricDigest
    switch input.ReplayObservation.Status {
    case jevImprovementReplayReproduced:
        output.Status = jevReplayFeedbackConfirmed
        output.FeedbackKind = jevReplayFeedbackConfirmed
    case jevImprovementReplayCounterexample:
        output.Status = jevReplayFeedbackRefuted
        output.FeedbackKind = jevReplayFeedbackRefuted
    case jevImprovementReplayObservationUnknown:
        output.Status = jevReplayFeedbackUnknown
        output.FeedbackKind = jevReplayFeedbackUnknown
    default:
        output.MissingStage = "feedback-kind"
        return output
    }
    output.FeedbackEvidenceDigest = digestJEVImprovementReplayFeedbackEvidence(output.ReplayObservationDigest, output.MetricDigest, output.FeedbackKind)
    output.EvidenceDigest = digestJEVImprovementReplayFeedback(output.Status, output.FeedbackKind, output.CandidateDigest, output.ReplayObservationDigest, output.MetricDigest, output.FeedbackEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevReplayFeedbackUnknown
        output.MissingStage = "feedback-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementReplayFeedbackEvidence(replayObservationDigest, metricDigest, feedbackKind string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", replayObservationDigest, metricDigest, feedbackKind)))
    return hex.EncodeToString(sum[:])
}

func digestJEVImprovementReplayFeedback(status, feedbackKind, candidateDigest, replayObservationDigest, metricDigest, feedbackEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s", status, feedbackKind, candidateDigest, replayObservationDigest, metricDigest, feedbackEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
