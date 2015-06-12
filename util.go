//
// utility functions for parsing domains from user input
//

package gtcha

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
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

// http://blog.manki.in/2012/03/generating-data-uris-in-go-language.html
func dataURI(buf []byte, contentType string) string {
	return fmt.Sprintf(
		"data:%s;base64,%s",
		contentType, base64.StdEncoding.EncodeToString(buf),
	)
}

// Get has to solve the following problem:
//
// You have two routines trying to get data. You want to return the data returned on a
// successful call and you only want to send an error once both routines have errored.
//
// We mostly solve this by using a waitgroup, but this also requires some additional
// machinery. We send the completed data over a channel and then select on that
// channel and another channel that sends an error (if there is one) once the
// waitgroup is finished.
//
// f1 is a first routine for getting the data. Typically, you call memcache here.
// f2 is the second, more expensive routine. This can call stable storage or do a computation.
func Get(f1 func() interface{}, f2 func() (interface{}, error)) (interface{}, error) {
	var (
		o     sync.Once
		wg    sync.WaitGroup
		iCh   = make(chan interface{})
		errCh = make(chan error)
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		if i := f1(); i != nil {
			o.Do(func() { iCh <- i })
		}
	}()

	go func() {
		defer wg.Done()

		i, err := f2()
		if err != nil {
			errCh <- err
			return
		}

		o.Do(func() { iCh <- i })
		errCh <- nil
	}()

	eCh := make(chan error)
	go func() {
		err := <-errCh
		wg.Wait()
		o.Do(func() { eCh <- err })
	}()

	select {
	case i := <-iCh:
		return i, nil
	case err := <-eCh:
		return nil, err
	}
}
