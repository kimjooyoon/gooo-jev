package gooo

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
)

// CapabilityDiscoveryFeedback records user feedback about a suggested discovery question.
// It is observational only: it never executes or authorizes the suggested operation.
type CapabilityDiscoveryFeedback struct {
    Version           string `json:"version"`
    OriginalQuery     string `json:"original_query"`
    SuggestedQuery    string `json:"suggested_query"`
    Outcome           string `json:"outcome"`
    CapabilityStatus  string `json:"capability_status"`
    FirstMissingStage string `json:"first_missing_stage,omitempty"`
    EvidenceDigest    string `json:"evidence_digest"`
    NonExecuting      bool   `json:"non_executing"`
    NonAuthorizing    bool   `json:"non_authorizing"`
}

// Validate enforces the evidence and safety boundary for discovery feedback.
func (f CapabilityDiscoveryFeedback) Validate() error {
    if f.Version == "" {
        return fmt.Errorf("feedback version is required")
    }
    if f.OriginalQuery == "" || f.SuggestedQuery == "" {
        return fmt.Errorf("original and suggested queries are required")
    }
    switch f.Outcome {
    case "useful", "not_useful", "unresolved":
    default:
        return fmt.Errorf("unsupported feedback outcome %q", f.Outcome)
    }
    if f.CapabilityStatus == "" {
        return fmt.Errorf("capability status is required")
    }
    if f.EvidenceDigest == "" {
        return fmt.Errorf("evidence digest is required")
    }
    if !f.NonExecuting || !f.NonAuthorizing {
        return fmt.Errorf("discovery feedback must remain non-executing and non-authorizing")
    }
    return nil
}

// CanonicalDigest returns the digest of validated feedback for later replay.
func (f CapabilityDiscoveryFeedback) CanonicalDigest() (string, error) {
    if err := f.Validate(); err != nil {
        return "", err
    }
    encoded, err := json.Marshal(f)
    if err != nil {
        return "", fmt.Errorf("marshal feedback: %w", err)
    }
    digest := sha256.Sum256(encoded)
    return hex.EncodeToString(digest[:]), nil
}