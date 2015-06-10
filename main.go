//
// HTTP endpoints
//

package gtcha

import (
	"encoding/json"
	"errors"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"

	"github.com/Bowery/gopackages/web"
)

type verificationResponse struct {
	IsHuman bool `json:"is_human"`
}

var routes = []web.Route{
	{"POST", "/register", registerApp, false},
	{"GET", "/captcha", getCaptcha, false},
	{"PUT", "/verify", verifySession, false},
	{"GET", "/verify", isVerified, false},
	{"GET", "/dummy", dummyHandler, false},
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(domains) == 0 {
		http.Error(w, "origins cannot be empty", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// TODO(r-medina): cache captchas and get from a cache instead of generating each time

	g, err := newGtcha(urlfetch.Client(c))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	captcha := g.toCaptcha()
	if err = SaveGtcha(c, g, captcha.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(captcha)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func verifySession(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := verifyRequest(c, r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.PostForm.Get("id")
	tag := r.PostForm.Get("tag")
	in := r.PostForm["in"]
	if len(in) == 0 {
		http.Error(w, "no images selected", http.StatusBadRequest)
		return
	}

	g, err := GetGtcha(c, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tag != g.Tag {
		http.Error(w, "tag incorrect for given id", http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO(r-medina): abstract out everything above this

	g, err := GetGtcha(c, r.PostForm.Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(verificationResponse{g.IsHuman})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	c := &Captcha{
		ID:  uuid.New(),
		Tag: "#beyonce",
		Images: []string{
			"https://media1.giphy.com/media/yFNA1ndGA5ZuM/200.gif",
			"https://media4.giphy.com/media/10H8p7oa4LUSB2/200.gif",
			"https://media2.giphy.com/media/skYmSo5tpORr2/200.gif",
			"https://media1.giphy.com/media/HfGqchLEK2WFq/200.gif",
		},
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
