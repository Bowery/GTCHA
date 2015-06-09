//
// utility functions for parsing domains from user input
//

package gtcha

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var errEmptyDomain = errors.New("domain empty")

// parseDomains takes the raw user input string for their app origins
// and makes it a slice of strings that are just the host.
// Properly formatted domains should include scheme and be separated by newlines
// eg:
//     http://bowery.io
//
//     http://localhost:8080
//
// Empty lines will be removed.
func parseDomains(rawDomains []string) ([]string, error) {
	var domains []string
	for _, domain := range rawDomains {
		url, err := parseDomain(domain)
		if err == errEmptyDomain {
			continue
		}
		if err != nil {
			return nil, err
		}

		domains = append(domains, url)
	}

	return domains, nil
}

// parseDomain parses and individual line of user input. See documentation for `parseDomains`.
func parseDomain(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errEmptyDomain
	}
	origin, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	if domain := origin.Host; domain != "" { // handles cases like http://bowery.io
		return domain, nil
	} else if domain = origin.Path; domain != "" { // handles cases like google.com
		return domain, nil
	}

	return "", fmt.Errorf("bad origin '%s'", raw)
}
