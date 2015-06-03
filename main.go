//
// HTTP and appengine code
//

package gitcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", getCaptcha)
	http.HandleFunc("/verify", verifySession)
	http.HandleFunc("/register", registerApp) // requires unique name. returns id and secret
}

func registerApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(
			w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed,
		)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
	fmt.Fprintf(w, "key %#v\n", key)

	app := &GtchaApp{
		Name:    r.PostForm.Get("name"),
		Secret:  key.StringID(),
		APIKey:  uuid.New(),
		Domains: strings.Split(r.PostForm.Get("domains"), "\n"),
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
	corsHeaders(w)
	m := r.Method
	if m == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if m != "GET" && m != "POST" {
		http.Error(w, fmt.Sprintf("Method %s not allowed", m), http.StatusMethodNotAllowed)
		return
	}

	c := appengine.NewContext(r)
	id := r.URL.Query().Get("api_key")

	if id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	app, err := GetApp(c, id, r.Header.Get("Origin"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	httpC := urlfetch.Client(c)

	// TODO(r-medina): actually make a captcha
	//
	// Requests:
	//     - GET /tag/random
	//     - GET /images/tag
	//     - GET /images/NOTtag
	//     - GET /images/tag/maybe
	//     - POST tag + picture id if human
	//
	// Process:
	//     -

	url := fmt.Sprintf(
		"%s/%s/%s", giphyAPI, giphyVer, "/gifs/search?q=funny+cat&api_key=dc6zaTOxFJmzC",
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := httpC.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, res.Body)
}

func verifySession(w http.ResponseWriter, r *http.Request) {}

func corsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

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
		return nil, errors.New("invalid origin")
	}

	return app, nil
}
