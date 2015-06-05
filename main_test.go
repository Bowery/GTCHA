package gtcha

import (
	"fmt"
	"testing"
)

func TestRegisterApp(t *testing.T)   {}
func TestGetCaptcha(t *testing.T)    {}
func TestVerifySession(t *testing.T) {}
func TestGetApp(t *testing.T)        {}
func TestParseDomains(t *testing.T) {
	//"a\nb\nhttp://abc.com"
	domains, err := ParseDomains("http://bowery.io\ngoogle.com    \n   https://bing.no")
	if err != nil {
		t.Fatal(err)
	}

	if n, a := len(domains), 3; n != a {
		t.Fatalf("expected %d domains, got %d", n, a)
	}

	fmt.Println(domains)

	if n, a := domains[0], "bowery.io"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if n, a := domains[1], "google.com"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if n, a := domains[1], "bing.no"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}
}
