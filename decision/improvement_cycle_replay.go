package decision

import (
    "errors"
    "strings"
)

type ImprovementReplayStatus string

const (
    ImprovementReplayed  ImprovementReplayStatus = "replayed"
    ImprovementDiverged  ImprovementReplayStatus = "diverged"
    ImprovementReplayUnknown ImprovementReplayStatus = "unknown"
)

type ImprovementCycleReplay struct {
    CycleDigest                    string
    RecordedResultLedgerDigest     string
    RecomputedResultLedgerDigest   string
    EvidenceDigest                 string
    Status                         ImprovementReplayStatus
    NonAuthorizing                 bool
    ReplayDigest                   string
}

func ReplayImprovementCycle(
    cycle ImprovementCycle,
    recomputedResultLedgerDigest,
    evidenceDigest string,
) (ImprovementCycleReplay, error) {
    if err := cycle.Validate(); err != nil {
        return ImprovementCycleReplay{}, err
    }
    replay := ImprovementCycleReplay{
        CycleDigest:                  cycle.CycleDigest,
        RecordedResultLedgerDigest:   cycle.ResultLedgerDigest,
        RecomputedResultLedgerDigest: recomputedResultLedgerDigest,
        EvidenceDigest:               evidenceDigest,
        NonAuthorizing:               true,
    }
    switch {
    case strings.TrimSpace(recomputedResultLedgerDigest) == "":
        replay.Status = ImprovementReplayUnknown
    case recomputedResultLedgerDigest == cycle.ResultLedgerDigest:
        replay.Status = ImprovementReplayed
    default:
        replay.Status = ImprovementDiverged
    }
    if err := replay.validateShape(); err != nil {
        return ImprovementCycleReplay{}, err
    }
    digest, err := Digest(replay)
    if err != nil {
        return ImprovementCycleReplay{}, err
    }
    replay.ReplayDigest = digest
    return replay, nil
}

func (replay ImprovementCycleReplay) Validate() error {
    if err := replay.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(replay.ReplayDigest) == "" {
        return errors.New("improvement cycle replay digest is missing")
    }
    copy := replay
    copy.ReplayDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != replay.ReplayDigest {
        return errors.New("improvement cycle replay digest does not match its evidence")
    }
    return nil
}

func (replay ImprovementCycleReplay) validateShape() error {
    if strings.TrimSpace(replay.CycleDigest) == "" ||
        strings.TrimSpace(replay.RecordedResultLedgerDigest) == "" ||
        strings.TrimSpace(replay.EvidenceDigest) == "" {
        return errors.New("improvement cycle replay is incomplete")
    }
    if !replay.NonAuthorizing {
        return errors.New("improvement cycle replay must remain non-authorizing")
    }
    switch replay.Status {
    case ImprovementReplayUnknown:
        if strings.TrimSpace(replay.RecomputedResultLedgerDigest) != "" {
            return errors.New("unknown replay must not claim a recomputed digest")
        }
    case ImprovementReplayed:
        if replay.RecomputedResultLedgerDigest != replay.RecordedResultLedgerDigest {
            return errors.New("replayed cycle digest does not match the recorded result")
        }
    case ImprovementDiverged:
        if strings.TrimSpace(replay.RecomputedResultLedgerDigest) == "" ||
            replay.RecomputedResultLedgerDigest == replay.RecordedResultLedgerDigest {
            return errors.New("diverged replay must retain a different recomputed digest")
        }
    default:
        return errors.New("unsupported improvement replay status")
    }
    return nil
}
