package gooo

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// WorkloadIdentity is a parsed SPIFFE-shaped identity value, not an authorization.
type WorkloadIdentity struct {
	Status         string
	MissingStage   string
	Scheme         string
	TrustDomain    string
	Path           string
	Raw            string
	SourceDigest   string
	EvidenceDigest string
	IdentityDigest string
	NonExecuting   bool
	NonAuthorizing bool
}

var trustDomainPattern = regexp.MustCompile("^[a-z0-9](?:[a-z0-9.-]*[a-z0-9])?$")

// ParseWorkloadIdentity validates a SPIFFE-shaped identity and its external evidence.
func ParseWorkloadIdentity(raw, evidenceDigest string) (WorkloadIdentity, error) {
	identity := WorkloadIdentity{
		Status:         "UNKNOWN",
		MissingStage:   "identity-parse",
		Raw:            raw,
		SourceDigest:   digestString(raw),
		EvidenceDigest: evidenceDigest,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !validDigest(evidenceDigest) {
		identity.MissingStage = "identity-evidence"
		return identity, fmt.Errorf("gooo identity: evidence digest is required")
	}
	if !strings.HasPrefix(raw, "spiffe://") {
		identity.MissingStage = "identity-scheme"
		return identity, fmt.Errorf("gooo identity: scheme must be spiffe")
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "spiffe" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		identity.MissingStage = "identity-syntax"
		return identity, fmt.Errorf("gooo identity: URI syntax is invalid")
	}
	if !trustDomainPattern.MatchString(parsed.Host) {
		identity.MissingStage = "identity-trust-domain"
		return identity, fmt.Errorf("gooo identity: trust domain is invalid")
	}
	if !validIdentityPath(parsed.Path) {
		identity.MissingStage = "identity-path"
		return identity, fmt.Errorf("gooo identity: workload path is invalid")
	}

	identity.Status = "BOUND"
	identity.MissingStage = ""
	identity.Scheme = parsed.Scheme
	identity.TrustDomain = parsed.Host
	identity.Path = parsed.Path
	identity.IdentityDigest = digestWorkloadIdentity(identity)
	return identity, nil
}

// Validate ensures the identity remains source and evidence bound without authorizing it.
func (i WorkloadIdentity) Validate() error {
	if i.Status != "BOUND" {
		return fmt.Errorf("identity status must be BOUND")
	}
	if i.MissingStage != "" {
		return fmt.Errorf("identity missing stage must be empty")
	}
	if i.Scheme != "spiffe" || !trustDomainPattern.MatchString(i.TrustDomain) || !validIdentityPath(i.Path) {
		return fmt.Errorf("identity fields are invalid")
	}
	if i.Raw == "" || i.SourceDigest != digestString(i.Raw) {
		return fmt.Errorf("identity source digest does not match")
	}
	if !validDigest(i.EvidenceDigest) || !validDigest(i.IdentityDigest) {
		return fmt.Errorf("identity evidence digest is invalid")
	}
	if !i.NonExecuting || !i.NonAuthorizing {
		return fmt.Errorf("identity must remain non-executing and non-authorizing")
	}
	if expected := digestWorkloadIdentity(i); expected != i.IdentityDigest {
		return fmt.Errorf("identity digest does not match its fields")
	}
	return nil
}

func validIdentityPath(path string) bool {
	if path == "" || !strings.HasPrefix(path, "/") {
		return false
	}
	segments := strings.Split(path, "/")[1:]
	if len(segments) == 0 {
		return false
	}
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

func digestWorkloadIdentity(identity WorkloadIdentity) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s",
		identity.Scheme,
		identity.TrustDomain,
		identity.Path,
		identity.SourceDigest,
		identity.EvidenceDigest,
	))
}
