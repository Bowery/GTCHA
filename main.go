//
// HTTP and appengine code
//

package gtcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

func init() {
	handle("/captcha", getCaptcha, "GET")
	handle("/verify", verifySession, "PUT")
	handle("/register", registerApp, "POST")
}

func registerApp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
	fmt.Fprintf(w, "key %#v\n", key)

	// clean up origin domains
	domains := strings.Split(r.PostForm.Get("domains"), "\n")
	for i, domain := range domains {
		origin, err := url.Parse(domain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		domains[i] = origin.Host
	}

	app := &GtchaApp{
		Name:    r.PostForm.Get("name"),
		Secret:  key.StringID(),
		APIKey:  uuid.New(),
		Domains: domains,
	}

	if _, err := datastore.Put(c, key, app); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func getCaptcha(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	apiKey := r.URL.Query().Get("api_key")

	if apiKey == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	origin, err := url.Parse(r.Header.Get("Origin"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app, err := GetApp(c, apiKey, origin.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	httpC := urlfetch.Client(c)
	g, err := newGtcha(httpC)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf, err := json.Marshal(g.toCaptcha())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

func verifySession(w http.ResponseWriter, r *http.Request) {}

// GetApp returns an app entity from the appengine datastore.
func GetApp(c appengine.Context, id, origin string) (*GtchaApp, error) {
	q := datastore.NewQuery("GtchaApp").Filter("APIKey =", id).Limit(1)
	app := new(GtchaApp)
	if _, err := q.Run(c).Next(app); err != nil {
		return nil, err
	}

	ok := false
	for _, domain := range app.Domains {
		if origin == domain {
			ok = true
		}
	}
	if !ok {
		return nil, errors.New("bad origin")
	}

	return app, nil
}
