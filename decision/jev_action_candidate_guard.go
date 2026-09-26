package decision

import (
	"path"
	"strings"
	"time"
)

// JEVActionCandidateGuardInput describes a typed candidate and the observation
// boundary from which it may be considered, without executing it.
type JEVActionCandidateGuardInput struct {
	CandidateID         string
	ObservationDigest   string
	ObservationAt       time.Time
	Now                 time.Time
	MaxAge              time.Duration
	RequestedPath       string
	AllowedPathPrefix   string
	WriteRequested      bool
	WriteAllowed        bool
	ConfirmationRequired bool
	ConfirmationPresent bool
	NonAuthorizing      bool
}

// JEVActionCandidateGuard records a non-executing disposition for a candidate.
type JEVActionCandidateGuard struct {
	Status            string
	CandidateID       string
	ObservationDigest string
	EvidenceDigest    string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

// GuardJEVActionCandidate validates freshness and capability boundaries before
// a separate executor is allowed to consider a candidate.
func GuardJEVActionCandidate(input JEVActionCandidateGuardInput) JEVActionCandidateGuard {
	output := JEVActionCandidateGuard{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if strings.TrimSpace(input.CandidateID) == "" || strings.TrimSpace(input.ObservationDigest) == "" {
		output.MissingStage = "candidate-observation"
		return output
	}
	if input.ObservationAt.IsZero() || input.Now.IsZero() || input.MaxAge <= 0 {
		output.MissingStage = "freshness"
		return output
	}
	requestedPath, allowedPathPrefix, ok := normalizeJEVPathScope(input.RequestedPath, input.AllowedPathPrefix)
	if !ok {
		output.MissingStage = "path-scope"
		return output
	}
	evidenceDigest, err := Digest(struct {
		CandidateID          string
		ObservationDigest    string
		ObservationAt        time.Time
		Now                  time.Time
		MaxAgeNanoseconds    int64
		RequestedPath        string
		AllowedPathPrefix    string
		WriteRequested       bool
		WriteAllowed         bool
		ConfirmationRequired bool
		ConfirmationPresent  bool
	}{
		CandidateID:          input.CandidateID,
		ObservationDigest:    input.ObservationDigest,
		ObservationAt:        input.ObservationAt.UTC(),
		Now:                  input.Now.UTC(),
		MaxAgeNanoseconds:    int64(input.MaxAge),
		RequestedPath:        requestedPath,
		AllowedPathPrefix:    allowedPathPrefix,
		WriteRequested:       input.WriteRequested,
		WriteAllowed:         input.WriteAllowed,
		ConfirmationRequired: input.ConfirmationRequired,
		ConfirmationPresent:  input.ConfirmationPresent,
	})
	if err != nil {
		output.MissingStage = "candidate-evidence"
		return output
	}
	output.CandidateID = input.CandidateID
	output.ObservationDigest = input.ObservationDigest
	output.EvidenceDigest = evidenceDigest
	age := input.Now.Sub(input.ObservationAt)
	switch {
	case age < 0:
		output.MissingStage = "observation-future"
		return output
	case age > input.MaxAge:
		output.Status = "review"
		output.MissingStage = "observation-stale"
		return output
	case input.WriteRequested && !input.WriteAllowed:
		output.Status = "rejected"
		output.MissingStage = "write-permission"
		return output
	case input.ConfirmationRequired && !input.ConfirmationPresent:
		output.Status = "review"
		output.MissingStage = "confirmation"
		return output
	default:
		output.Status = "admitted"
		return output
	}
}

func normalizeJEVPathScope(requestedPath, allowedPathPrefix string) (string, string, bool) {
	requested := path.Clean(strings.TrimSpace(requestedPath))
	allowed := path.Clean(strings.TrimSpace(allowedPathPrefix))
	if requested == "." || allowed == "." || requested == ".." || allowed == ".." ||
		strings.HasPrefix(requested, "../") || strings.HasPrefix(allowed, "../") {
		return "", "", false
	}
	if allowed == "/" {
		return requested, allowed, strings.HasPrefix(requested, "/")
	}
	if requested != allowed && !strings.HasPrefix(requested, allowed+"/") {
		return "", "", false
	}
	return requested, allowed, true
}
