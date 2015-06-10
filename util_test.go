package gtcha

import (
	"os"
	"testing"
)

func TestParseDomains(t *testing.T) {
	domains, err := parseDomains(
		[]string{"http://bowery.io", "\ngoogle.com    ", "   https://bing.no"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if n, a := len(domains), 3; n != a {
		t.Fatalf("expected %d domains, got %d", n, a)
	}

	if n, a := domains[0], "bowery.io"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if n, a := domains[1], "google.com"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if n, a := domains[2], "bing.no"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	domains, err = parseDomains([]string{"a", "\nb", "\nhttp://abc.com:8080"})
	if err != nil {
		t.Fatal(err)
	}

	if n, a := len(domains), 3; n != a {
		t.Fatalf("expected %d domains, got %d", n, a)
	}

	if n, a := domains[0], "a"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if n, a := domains[1], "b"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if n, a := domains[2], "abc.com:8080"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	domains, err = parseDomains([]string{"abc.com:8080"})
	if err == nil {
		t.Fatal("expected error")
	}

	domains, err = parseDomains([]string{
		"a", "\n\n", "\n\n", " \n\n", "b\n", "http://abc.com:8080",
	})
	if err != nil {
		t.Fatal(err)
	}

	if n, a := len(domains), 3; n != a {
		t.Fatalf("expected %d domains, got %d", n, a)
	}
}

func TestParseDomain(t *testing.T) {
	if domain, err := parseDomain("     http://bowery.io      "); err != nil {
		t.Fatal(err)
	} else if n, a := domain, "bowery.io"; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}

	if _, err := parseDomain("     spark.io:90      "); err == nil {
		t.Fatal("expected error")
	}

	if domain, err := parseDomain("     http://spark.io:90      "); err != nil {
		t.Fatal(err)
	} else if n, a := "spark.io:90", domain; n != a {
		t.Fatalf("expected '%s', got '%s'", n, a)
	}
}

func TestDataURI(t *testing.T) {
	exp := "data:image/gif;base64,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"

	f, err := os.Open("giphy_logo_laser.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, finfo.Size())

	_, err = f.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	uri := dataURI(buf, gifType)
	if uri != exp {
		t.Fatalf("data uri different than expected")
	}
}

func TestGet(t *testing.T) { t.Log("TODO") }
