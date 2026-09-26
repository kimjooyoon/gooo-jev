package decision

import (
    "errors"
    "math"
    "strings"
)

type MetricDeltaStatus string

const (
    MetricObserved MetricDeltaStatus = "observed"
    MetricUnknown  MetricDeltaStatus = "unknown"
)

type MetricDeltaDirection string

const (
    MetricImproved         MetricDeltaDirection = "improved"
    MetricRegressed        MetricDeltaDirection = "regressed"
    MetricUnchanged        MetricDeltaDirection = "unchanged"
    MetricDirectionUnknown MetricDeltaDirection = "unknown"
)

type MetricDelta struct {
    MetricName     string
    Before         float64
    After          float64
    EvidenceDigest string
    Status         MetricDeltaStatus
    Direction      MetricDeltaDirection
    DeltaDigest    string
}

func NewMetricDelta(metricName string, before, after float64, evidenceDigest string, status MetricDeltaStatus) (MetricDelta, error) {
    delta := MetricDelta{
        MetricName:     metricName,
        Before:         before,
        After:          after,
        EvidenceDigest: evidenceDigest,
        Status:         status,
    }
    switch status {
    case MetricObserved:
        switch {
        case after > before:
            delta.Direction = MetricImproved
        case after < before:
            delta.Direction = MetricRegressed
        default:
            delta.Direction = MetricUnchanged
        }
    case MetricUnknown:
        delta.Direction = MetricDirectionUnknown
    }
    if err := delta.validateShape(); err != nil {
        return MetricDelta{}, err
    }
    digest, err := Digest(delta)
    if err != nil {
        return MetricDelta{}, err
    }
    delta.DeltaDigest = digest
    return delta, nil
}

func (delta MetricDelta) Validate() error {
    if err := delta.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(delta.DeltaDigest) == "" {
        return errors.New("metric delta digest is missing")
    }
    copy := delta
    copy.DeltaDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != delta.DeltaDigest {
        return errors.New("metric delta digest does not match its evidence")
    }
    return nil
}

func (delta MetricDelta) validateShape() error {
    if strings.TrimSpace(delta.MetricName) == "" ||
        strings.TrimSpace(delta.EvidenceDigest) == "" {
        return errors.New("metric delta is incomplete")
    }
    switch delta.Status {
    case MetricObserved:
        if math.IsNaN(delta.Before) || math.IsNaN(delta.After) ||
            math.IsInf(delta.Before, 0) || math.IsInf(delta.After, 0) {
            return errors.New("observed metric delta must be finite")
        }
        switch delta.Direction {
        case MetricImproved, MetricRegressed, MetricUnchanged:
        default:
            return errors.New("observed metric delta has unsupported direction")
        }
    case MetricUnknown:
        if delta.Direction != MetricDirectionUnknown {
            return errors.New("unknown metric delta must retain unknown direction")
        }
    default:
        return errors.New("unsupported metric delta status")
    }
    return nil
}
