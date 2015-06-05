//
// HTTP and appengine code
//

package gtcha

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
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

	// clean up origin domains
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
	domains, err := ParseDomains(r.PostForm.Get("domains"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// GetApp returns a GtchaApp entity from the appengine datastore.
func GetApp(c context.Context, id, origin string) (*GtchaApp, error) {
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

// ParseDomains takes the raw user input string for their app origins
// and makes it a slice of strings that are just the host.
func ParseDomains(raw string) ([]string, error) {
	domains := strings.Split(raw, "\n")
	for i, domain := range domains {
		origin, err := url.Parse(domain)
		if err != nil {
			return nil, err
		}
		domains[i] = origin.Host
	}

	return domains, nil
}

// GetGtcha looks for the captcha with the given ID in memcache and datastore.
func GetGtcha(c context.Context, id string) (*gtcha, error) {
	return nil, nil
}

// SaveGtcha saves a gtcha in the database and memcache.
func SaveGtcha(c context.Context, g *gtcha, id string) error {
	// put g in memcache
	go func() {
		buf, err := json.Marshal(g)
		if err != nil {
			return
		}

		memcache.Add(c, &memcache.Item{
			Key:   id,
			Value: buf,
		})
	}()

	key := datastore.NewKey(c, "Captcha", id, 0, nil)
	if _, err := datastore.Put(c, key, g); err != nil {
		return err
	}

	return nil
}
