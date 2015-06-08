// https://cloud.google.com/appengine/docs/go/tools/localunittesting/#Go_Introducing_the_aetest_package

package gtcha

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"code.google.com/p/appengine-go/appengine/aetest"

	"appengine"
)

func TestRegisterApp(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)

	}
	defer inst.Close()

	ts := httptest.NewServer(http.HandlerFunc(registerApp))
	defer ts.Close()

	u := url.Values{}
	u.Set("name", "bizzle")
	u.Set("domains", "http://bowery.io/wut/up")
	body := strings.NewReader(u.Encode())
	req, err := inst.NewRequest("POST", ts.URL, body)
	if err != nil {
		t.Fatal(err)
	}
	c := appengine.NewContext(req)

	rec := httptest.NewRecorder()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", res) // output for debug
}

func TestGetCaptcha(t *testing.T)    {}
func TestVerifySession(t *testing.T) {}
