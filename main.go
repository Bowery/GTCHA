//
// HTTP endpoints
//

package gtcha

import (
	"encoding/json"
	"net/http"
	"net/url"

	"code.google.com/p/go-uuid/uuid"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"

	"github.com/Bowery/gopackages/web"
)

func init() {
	handlers := []web.Handler{new(web.CorsHandler)}
	s := web.NewServer("", handlers, routes)
	s.Prestart()

	http.Handle("/", handlers[0])
}

var routes = []web.Route{
	{"POST", "/register", registerApp, false},
	{"GET", "/captcha", getCaptcha, false},
	{"PUT", "/verify", verifySession, false},
	// {"GET", "/verify", isVerified, false},
	{"GET", "/dummy", dummyHandler, false},
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
	domains, err := parseDomains(r.PostForm.Get("domains"))
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
	// w.WriteHeader(http.StatusOK)
}
