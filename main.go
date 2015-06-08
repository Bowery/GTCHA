//
// HTTP endpoints
//

package gtcha

import (
	"encoding/json"
	"net/http"
	"net/url"

	"code.google.com/p/go-uuid/uuid"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

func init() {
<<<<<<< HEAD
	handle("/captcha", getCaptcha, "GET")
	handle("/verify", verifySession, "PUT")
	handle("/register", registerApp, "POST")
=======
	handle("/register", registerApp, "POST")
	handle("/captcha", getCaptcha, "GET")
	handle("/verify", verifySession, "PUT")
>>>>>>> ff2977960434562fbb161d039a0a9efd89d4e7d7
}

func registerApp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// clean up origin domains
	domains, err := parseDomains(r.PostForm.Get("domains"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
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

	c := appengine.NewContext(r)
	var app *GtchaApp
	if url := origin.Host; url != "" {
		app, err = GetApp(c, apiKey, origin.Host)
	} else {
		app, err = GetApp(c, apiKey, origin.Path)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	httpC := urlfetch.Client(c)

	// TODO(r-medina): cache captchas and get from a cache instead of generating each time

	g, err := newGtcha(httpC)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	captcha := g.toCaptcha()
	buf, err := json.Marshal(captcha)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = SaveGtcha(c, g, captcha.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}

// TODO(r-medina): write this function
func verifySession(w http.ResponseWriter, r *http.Request) {}
