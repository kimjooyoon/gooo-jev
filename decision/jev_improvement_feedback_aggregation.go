package decision

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

const (
    jevImprovementFeedbackStableForReview = "stable-for-review"
    jevImprovementFeedbackNeedsRevision = "needs-revision"
    jevImprovementFeedbackHold = "hold"
    jevImprovementFeedbackAggregateUnknown = "UNKNOWN"
)

type JEVImprovementFeedbackAggregationInput struct {
    Feedback       []JEVImprovementReplayFeedback
    NonAuthorizing bool
}

type JEVImprovementFeedbackAggregation struct {
    Status             string
    MissingStage       string
    Total              int
    Confirmed          int
    Refuted            int
    Unknown            int
    InputEvidenceDigest string
    EvidenceDigest     string
    NonExecuting       bool
    NonAuthorizing     bool
}

func (a JEVImprovementFeedbackAggregation) Validate() error {
    if a.Status == "" || a.Total <= 0 || a.InputEvidenceDigest == "" || a.EvidenceDigest == "" {
        return fmt.Errorf("incomplete JEV improvement feedback aggregation")
    }
    if a.Confirmed < 0 || a.Refuted < 0 || a.Unknown < 0 || a.Confirmed+a.Refuted+a.Unknown != a.Total {
        return fmt.Errorf("invalid JEV improvement feedback aggregation counts")
    }
    switch a.Status {
    case jevImprovementFeedbackStableForReview, jevImprovementFeedbackNeedsRevision, jevImprovementFeedbackHold:
    default:
        return fmt.Errorf("invalid JEV improvement feedback aggregation status")
    }
    if !a.NonExecuting {
        return fmt.Errorf("JEV improvement feedback aggregation must be non-executing")
    }
    if !a.NonAuthorizing {
        return fmt.Errorf("JEV improvement feedback aggregation must be non-authorizing")
    }
    expected := digestJEVImprovementFeedbackAggregation(a.Status, a.Total, a.Confirmed, a.Refuted, a.Unknown, a.InputEvidenceDigest)
    if a.EvidenceDigest != expected {
        return fmt.Errorf("JEV improvement feedback aggregation evidence digest mismatch")
    }
    return nil
}

func AggregateJEVImprovementFeedback(input JEVImprovementFeedbackAggregationInput) JEVImprovementFeedbackAggregation {
    output := JEVImprovementFeedbackAggregation{
        Status: jevImprovementFeedbackAggregateUnknown,
        NonExecuting: true,
        NonAuthorizing: true,
    }
    if !input.NonAuthorizing {
        output.MissingStage = "authorization-boundary"
        return output
    }
    if len(input.Feedback) == 0 {
        output.MissingStage = "feedback-history"
        return output
    }
    evidenceParts := make([]string, 0, len(input.Feedback))
    for index, feedback := range input.Feedback {
        if !feedback.NonAuthorizing {
            output.MissingStage = fmt.Sprintf("feedback[%d]-authorization-boundary", index)
            return output
        }
        if !feedback.NonExecuting {
            output.MissingStage = fmt.Sprintf("feedback[%d]-execution-boundary", index)
            return output
        }
        if err := feedback.Validate(); err != nil {
            output.MissingStage = fmt.Sprintf("feedback[%d]", index)
            return output
        }
        evidenceParts = append(evidenceParts, feedback.EvidenceDigest)
        output.Total++
        switch feedback.FeedbackKind {
        case jevReplayFeedbackConfirmed:
            output.Confirmed++
        case jevReplayFeedbackRefuted:
            output.Refuted++
        case jevReplayFeedbackUnknown:
            output.Unknown++
        }
    }
    output.InputEvidenceDigest = digestJEVImprovementFeedbackInputs(evidenceParts)
    switch {
    case output.Refuted > 0:
        output.Status = jevImprovementFeedbackNeedsRevision
    case output.Unknown > 0:
        output.Status = jevImprovementFeedbackHold
    default:
        output.Status = jevImprovementFeedbackStableForReview
    }
    output.EvidenceDigest = digestJEVImprovementFeedbackAggregation(output.Status, output.Total, output.Confirmed, output.Refuted, output.Unknown, output.InputEvidenceDigest)
    if err := output.Validate(); err != nil {
        output.Status = jevImprovementFeedbackAggregateUnknown
        output.MissingStage = "feedback-aggregation-evidence"
        output.EvidenceDigest = ""
    }
    return output
}

func digestJEVImprovementFeedbackInputs(evidenceDigests []string) string {
    sum := sha256.Sum256([]byte(strings.Join(evidenceDigests, "|")))
    return hex.EncodeToString(sum[:])
}

func digestJEVImprovementFeedbackAggregation(status string, total, confirmed, refuted, unknown int, inputEvidenceDigest string) string {
    sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d|%d|%d|%s", status, total, confirmed, refuted, unknown, inputEvidenceDigest)))
    return hex.EncodeToString(sum[:])
}
