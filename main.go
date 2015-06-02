//
// HTTP and appengine logic
//

package gitcha

import (
	"fmt"
	"io"
	"net/http"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/", captchaHndlr)
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
	name := r.PostForm.Get("name")
	switch err := CheckName(c, name); err {
	case nil:
		break
	case ErrNameExists:
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case ErrNameTooLong:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app, err := NewApp(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = SaveApp(c, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func captchaHndlr(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != "GET" && m != "POST" {
		http.Error(w, fmt.Sprintf("Method %s not allowed", m), http.StatusMethodNotAllowed)
		return
	}

	c := appengine.NewContext(r)
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	app, err := GetApp(c, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var res *http.Response
	switch m {
	case "GET":
		res, err = http.Get(fmt.Sprintf("%s/%s", giphyAPI, captchaEndpoint))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default: // "POST"
		if err = r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if sec := r.Form.Get("secret"); sec != app.Secret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req *http.Request
		req, err = http.NewRequest(
			"POST", fmt.Sprintf("%s/%s", giphyAPI, captchaEndpoint), r.Body,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, res.Body)
}

// GetApp returns an app entity from the appengine datastore.
func GetApp(c appengine.Context, id string) (*gitchaApp, error) {
	q := datastore.NewQuery("App").Filter("id =", id).Limit(1)
	app := new(gitchaApp)
	_, err := q.Run(c).Next(app)

	return app, err
}

// SaveApp commits an app to a new entity in the appengine datastore.
func SaveApp(c appengine.Context, app *gitchaApp) error {
	key := datastore.NewKey(c, "gitchaApp", app.Secret, 0, nil)
	_, err := datastore.Put(c, key, app)

	return err
}

// CheckName checks that the name is unique and of the right length.
func CheckName(c appengine.Context, name string) error {
	if n := len(name); n > keyLen {
		return ErrNameTooLong
	}

	q := datastore.NewQuery("App").Filter("name =", name).Limit(1)
	app := new(gitchaApp)
	if _, err := q.Run(c).Next(app); err != nil {
		return err
	}
	if app != nil {
		return ErrNameExists
	}

	return nil
}
