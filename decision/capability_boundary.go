package decision

import (
    "errors"
    "strings"
    "time"
)

type CapabilityBoundaryStatus string

const (
    CapabilityBoundaryActive   CapabilityBoundaryStatus = "active"
    CapabilityBoundaryExpired   CapabilityBoundaryStatus = "expired"
    CapabilityBoundaryNotYet    CapabilityBoundaryStatus = "not-yet-active"
    CapabilityBoundaryUnknown   CapabilityBoundaryStatus = "unknown"
)

type CapabilityBoundary struct {
    Subject       string
    Audience      string
    Capability    string
    EvidenceDigest string
    IssuedAt      time.Time
    ExpiresAt     time.Time
    BoundaryDigest string
}

func NewCapabilityBoundary(subject, audience, capability, evidenceDigest string, issuedAt, expiresAt time.Time) (CapabilityBoundary, error) {
    boundary := CapabilityBoundary{
        Subject:        subject,
        Audience:       audience,
        Capability:     capability,
        EvidenceDigest: evidenceDigest,
        IssuedAt:       issuedAt,
        ExpiresAt:      expiresAt,
    }
    if err := boundary.validateShape(); err != nil {
        return CapabilityBoundary{}, err
    }
    digest, err := Digest(boundary)
    if err != nil {
        return CapabilityBoundary{}, err
    }
    boundary.BoundaryDigest = digest
\n    return boundary, nil
}

func (boundary CapabilityBoundary) StatusAt(now time.Time) CapabilityBoundaryStatus {
    if boundary.validateShape() != nil || strings.TrimSpace(boundary.BoundaryDigest) == "" {
        return CapabilityBoundaryUnknown
    }
    if now.Before(boundary.IssuedAt) {
        return CapabilityBoundaryNotYet
    }
    if !now.Before(boundary.ExpiresAt) {
        return CapabilityBoundaryExpired
    }
\n    return CapabilityBoundaryActive
}

func (boundary CapabilityBoundary) Validate() error {
    if err := boundary.validateShape(); err != nil {
        return err
    }
    if strings.TrimSpace(boundary.BoundaryDigest) == "" {
        return errors.New("capability boundary digest is missing")
    }
    copy := boundary
    copy.BoundaryDigest = ""
    digest, err := Digest(copy)
    if err != nil {
        return err
    }
    if digest != boundary.BoundaryDigest {
        return errors.New("capability boundary digest does not match its evidence")
    }
    return nil
}

func (boundary CapabilityBoundary) validateShape() error {
    if strings.TrimSpace(boundary.Subject) == "" ||
        strings.TrimSpace(boundary.Audience) == "" ||
        strings.TrimSpace(boundary.Capability) == "" ||
\n        strings.TrimSpace(boundary.EvidenceDigest) == "" {
        return errors.New("capability boundary is incomplete")
    }
    if boundary.IssuedAt.IsZero() || boundary.ExpiresAt.IsZero() {
        return errors.New("capability boundary time window is incomplete")
    }
    if !boundary.ExpiresAt.After(boundary.IssuedAt) {
        return errors.New("capability boundary expiry must follow issuance")
    }
    return nil
}