package decision

import (
    "errors"
    "strings"
)

type ImprovementCandidateDecision string

const (
    CandidateProposed ImprovementCandidateDecision = "proposed"
    CandidateReview   ImprovementCandidateDecision = "review"
    CandidateRejected ImprovementCandidateDecision = "rejected"
)

type ImprovementCandidate struct {
    CycleDigest       string
    ReplayDigest      string
    CandidateReference string
    SourceReference   string
    Decision          ImprovementCandidateDecision
    NonExecuting      bool
    CandidateDigest   string
}

func ProposeImprovementCandidate(
    cycle ImprovementCycle,
    replay ImprovementCycleReplay,
    candidateReference,
    sourceReference string,
    decision ImprovementCandidateDecision,
) (ImprovementCandidate, error) {
    if err := cycle.Validate(); err != nil {
        return ImprovementCandidate{}, err
    }
    if err := replay.Validate(); err != nil {
        return ImprovementCandidate{}, err
    }
    if replay.CycleDigest != cycle.CycleDigest {
        return ImprovementCandidate{}, errors.New("improvement candidate cycle and replay digests do not match")
    }
    candidate := ImprovementCandidate{
        CycleDigest:        cycle.CycleDigest,
        ReplayDigest:       replay.ReplayDigest,
        CandidateReference: candidateReference,
        SourceReference:    sourceReference,
        Decision:           decision,
        NonExecuting:       true,
    }
    if err := candidate.validateShape(); err != nil {
        return ImprovementCandidate{}, err
    }
    digest, err := Digest(candidate)
    if err != nil {
        return ImprovementCandidate{}, err
    }
    candidate.CandidateDigest = digest
    return candidate, nil
}

func (candidate ImprovementCandidate) Validate() error {
    if err := candidate.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(candidate.CandidateDigest) == "" {
        return errors.New("improvement candidate digest is missing")
    }
    copy := candidate
    copy.CandidateDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != candidate.CandidateDigest {
        return errors.New("improvement candidate digest does not match its evidence")
    }
    return nil
}

func (candidate ImprovementCandidate) validateShape() error {
    if strings.TrimSpace(candidate.CycleDigest) == "" ||
        strings.TrimSpace(candidate.ReplayDigest) == "" ||
        strings.TrimSpace(candidate.CandidateReference) == "" ||
        strings.TrimSpace(candidate.SourceReference) == "" {
        return errors.New("improvement candidate is incomplete")
    }
    if !candidate.NonExecuting {
        return errors.New("improvement candidate must remain non-executing")
    }
    switch candidate.Decision {
    case CandidateProposed, CandidateReview, CandidateRejected:
        return nil
    default:
        return errors.New("unsupported improvement candidate decision")
    }
}
