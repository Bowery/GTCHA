package gtcha

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"google.golang.org/appengine/aetest"
)

func TestRegisterApp(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)

	}
	defer inst.Close()

	req, err := inst.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = url.Values{}
	req.PostForm.Set("name", "bizzle")
	req.PostForm.Add("domain", "http://bowery.io/wut/up")

	rec := httptest.NewRecorder()

	registerApp(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatal("bad status code")
	}

	decoder := json.NewDecoder(rec.Body)
	app := new(GtchaApp)
	if err = decoder.Decode(app); err != nil {
		t.Fatal(err)
	}

	if app.Name != "bizzle" {
		t.Fatalf("expected app name %s, got %s", "bizzle", app.Name)
	}

	if n, a := len(app.Domains), 1; n != a {
		t.Fatalf("expected %d domains, got %d", n, a)
	}

	if app.Secret == "" {
		t.Fatal("empty app secret")
	}

	if app.APIKey == "" {
		t.Fatal("empty api key")
	}

	req, err = inst.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = url.Values{}
	req.PostForm.Add("domain", "http://bowery.io/")

	rec = httptest.NewRecorder()

	registerApp(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatal("should have returned error")
	}

	req, err = inst.NewRequest("POST", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.PostForm = url.Values{}
	req.PostForm.Set("name", "dat app")

	rec = httptest.NewRecorder()

	registerApp(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatal("should have returned error")
	}
}

// TODO: these two can't be written until we have the endpoints from giphy

func TestGetCaptcha(t *testing.T)    {}
func TestVerifySession(t *testing.T) {}
