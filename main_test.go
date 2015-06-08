package gtcha

import "testing"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"code.google.com/p/appengine-go/appengine/aetest"
)

func TestRegisterApp(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)

	}
	defer inst.Close()

	u := url.Values{}
	u.Set("name", "bizzle")
	u.Set("domains", "http://bowery.io/wut/up")
	body := strings.NewReader(u.Encode())
	req, err := inst.NewRequest("POST", "/register", body)
	if err != nil {
		t.Fatal(err)
	}
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

	fmt.Printf("%+v\n", app) // output for debug

}

func TestGetCaptcha(t *testing.T)    {}
func TestVerifySession(t *testing.T) {}
