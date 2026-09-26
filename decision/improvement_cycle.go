package decision

import (
    "errors"
    "strings"
)

type ImprovementCycleOutcome string

const (
    ImprovementObserved  ImprovementCycleOutcome = "observed"
    ImprovementRegressed ImprovementCycleOutcome = "regressed"
    ImprovementUnknown   ImprovementCycleOutcome = "unknown"
    ImprovementReview     ImprovementCycleOutcome = "review"
)

type ImprovementCycle struct {
    BaselineLedgerDigest string
    ResultLedgerDigest   string
    CandidateReference   string
    SourceReference      string
    Outcome              ImprovementCycleOutcome
    NonAuthorizing       bool
    CycleDigest          string
}

func BuildImprovementCycle(
    baselineLedgerDigest,
    resultLedgerDigest,
    candidateReference,
    sourceReference string,
    outcome ImprovementCycleOutcome,
) (ImprovementCycle, error) {
    cycle := ImprovementCycle{
        BaselineLedgerDigest: baselineLedgerDigest,
        ResultLedgerDigest:   resultLedgerDigest,
        CandidateReference:   candidateReference,
        SourceReference:      sourceReference,
        Outcome:              outcome,
        NonAuthorizing:       true,
    }
    if err := cycle.validateShape(); err != nil {
        return ImprovementCycle{}, err
    }
    digest, err := Digest(cycle)
    if err != nil {
        return ImprovementCycle{}, err
    }
    cycle.CycleDigest = digest
    return cycle, nil
}

func (cycle ImprovementCycle) Validate() error {
    if err := cycle.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(cycle.CycleDigest) == "" {
        return errors.New("improvement cycle digest is missing")
    }
    copy := cycle
    copy.CycleDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != cycle.CycleDigest {
        return errors.New("improvement cycle digest does not match its evidence")
    }
    return nil
}

func (cycle ImprovementCycle) validateShape() error {
    if strings.TrimSpace(cycle.BaselineLedgerDigest) == "" ||
        strings.TrimSpace(cycle.ResultLedgerDigest) == "" ||
        strings.TrimSpace(cycle.CandidateReference) == "" ||
        strings.TrimSpace(cycle.SourceReference) == "" {
        return errors.New("improvement cycle is incomplete")
    }
    if !cycle.NonAuthorizing {
        return errors.New("improvement cycle must remain non-authorizing")
    }
    switch cycle.Outcome {
    case ImprovementObserved, ImprovementRegressed, ImprovementUnknown, ImprovementReview:
        return nil
    default:
        return errors.New("unsupported improvement cycle outcome")
    }
}
