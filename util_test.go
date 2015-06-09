package gtcha

import "testing"

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

	domains, err = parseDomains([]string{"a", "\n\n", "\n\n", " \n\n", "b\n", "http://abc.com:8080"})
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
