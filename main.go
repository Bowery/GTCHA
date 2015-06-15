//
// HTTP endpoints
//

package GTCHA

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Bowery/gopackages/requests"
	"github.com/Bowery/gopackages/web"
	"github.com/pborman/uuid"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

type verificationResponse struct {
	IsHuman bool `json:"is_human"`
}

var routes = []web.Route{
	{"POST", "/register", registerApp, false},
	{"GET", "/captcha", getCaptcha, false},
	{"PUT", "/verify", verifySession, false},
	{"GET", "/verify", isVerified, false},
	{"GET", "/dummy_get", dummyGetHandler, false},
	{"PUT", "/dummy_put", dummyPutHandler, false},
}

func init() {
	handlers := []web.Handler{new(web.CorsHandler)}
	s := web.NewServer("", handlers, routes)
	s.Prestart()

	http.Handle("/", handlers[0])
}

func registerApp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.PostForm.Get("name")

	if name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	// clean up origin domains
	domains, err := parseDomains(r.PostForm["domain"])
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}
	if len(domains) == 0 {
		requests.ErrorJSON(
			w, http.StatusBadRequest, requests.StatusFailed, "origins cannot be empty",
		)
		return
	}

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
	app := &GtchaApp{
		Name:    name,
		Secret:  key.StringID(),
		APIKey:  uuid.New(),
		Domains: domains,
	}
	if _, err := datastore.Put(c, key, app); err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(app)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getCaptcha(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := verifyRequest(c, r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// TODO(r-medina): cache captchas and get from a cache
	// instead of generating each time

	httpC := urlfetch.Client(c)
	g, err := newGtcha(httpC)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}
	captcha, err := g.toCaptcha(c, httpC)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}
	if err = SaveGtcha(c, g, captcha.ID); err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(captcha)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func verifySession(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := verifyRequest(c, r); err != nil {
		requests.ErrorJSON(
			w, http.StatusUnauthorized, requests.StatusFailed, err.Error(),
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	id := r.PostForm.Get("id")
	tag := r.PostForm.Get("tag")
	in := r.PostForm["in"]
	if len(in) == 0 {
		requests.ErrorJSON(
			w, http.StatusBadRequest, requests.StatusFailed, "no images selected",
		)
		http.Error(w, "no images selected", http.StatusBadRequest)
		return
	}

	g, err := GetGtcha(c, id)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	if tag != g.Tag {
		requests.ErrorJSON(
			w, http.StatusInternalServerError,
			requests.StatusFailed, "tag incorrect for given id",
		)
		return
	}

	go func() {
		if verifyGtcha(urlfetch.Client(c), g, in) {
			SaveGtcha(c, g, id)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func isVerified(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := verifyRequest(c, r); err != nil {
		requests.ErrorJSON(
			w, http.StatusUnauthorized, requests.StatusFailed, err.Error(),
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	g, err := GetGtcha(c, r.PostForm.Get("id"))
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(verificationResponse{g.IsHuman})
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func dummyGetHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	httpC := urlfetch.Client(c)

	g, err := newGtcha(httpC)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	captcha, err := g.toCaptcha(c, httpC)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(captcha)
	if err != nil {
		requests.ErrorJSON(
			w, http.StatusInternalServerError, requests.StatusFailed, err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func dummyPutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func verifyRequest(c appengine.Context, r *http.Request) error {
	apiKey := r.URL.Query().Get("api_key")

	if apiKey == "" {
		return errors.New("unauthorized")
	}

	origin, err := parseDomain(r.Header.Get("Origin"))
	if err != nil {
		return err
	}

	app, err := GetApp(c, apiKey, origin)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("unauthorized")
	}

	return nil
}
