package gtcha

import (
	"testing"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"appengine/aetest"
	"appengine/datastore"
)

func TestRegisterApp(t *testing.T)   {}
func TestGetCaptcha(t *testing.T)    {}
func TestVerifySession(t *testing.T) {}

func TestGetApp(t *testing.T) {
	domains, err := parseDomains("bowery.io")
	if err != nil {
		t.Fatal(err)
	}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
	app := &GtchaApp{
		Name:    "testApp",
		Secret:  key.StringID(),
		APIKey:  uuid.New(),
		Domains: domains,
	}
	if _, err := datastore.Put(c, key, app); err != nil {
		t.Fatal(err)
	}

	// increases likelihood that app goes in the datastore before call to Get
	<-time.After(500 * time.Millisecond)

	retApp, err := GetApp(c, app.APIKey, "bowery.io")
	if err != nil {
		t.Fatal(err)
	}

	if retApp.Secret != app.Secret {
		t.Fatalf(
			"returned app differs from expected app:\nreturned %#v\nexpected %#v", retApp, app,
		)
	}
}

func TestGetGtcha(t *testing.T) {

}

func TestSaveGtcha(t *testing.T) {

}

func TestParseDomains(t *testing.T) {
	domains, err := parseDomains("http://bowery.io\ngoogle.com    \n   https://bing.no")
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

	domains, err = parseDomains("a\nb\nhttp://abc.com:8080")
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

	domains, err = parseDomains("abc.com:8080")
	if err == nil {
		t.Fatal("expected error")
	}

	domains, err = parseDomains("a\n\n\n\n \n\nb\nhttp://abc.com:8080")
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
