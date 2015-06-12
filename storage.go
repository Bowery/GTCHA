//
// appengine datastore and memcache logic
//

package gtcha

import (
	"encoding/json"
	"errors"

	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

// GetApp returns a GtchaApp entity from the appengine datastore.
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

// GetGtcha looks for the captcha with the given ID in memcache and datastore.
func GetGtcha(c appengine.Context, id string) (*gtcha, error) {
	// memcache
	f1 := func() interface{} {
		item, err := memcache.Get(c, "Gtcha"+id)
		if err != nil {
			return nil
		}

		g := new(gtcha)
		if err = json.Unmarshal(item.Value, g); err != nil {
			return nil
		}

		return g
	}

	// datastore
	f2 := func() (interface{}, error) {
		key := datastore.NewKey(c, "Gtcha", id, 0, nil)
		g := new(gtcha)
		if err := datastore.Get(c, key, g); err != nil {
			return nil, err
		}

		return g, nil
	}

	g, err := Get(f1, f2)
	if err != nil {
		return nil, err
	}

	return g.(*gtcha), nil
}

// SaveGtcha saves a gtcha in the database and memcache.
func SaveGtcha(c appengine.Context, g *gtcha, id string) error {
	// put g in memcache
	go func() {
		buf, err := json.Marshal(g)
		if err != nil {
			return
		}

		memcache.Add(c, &memcache.Item{
			Key:   "Gtcha" + id,
			Value: buf,
		})
	}()

	key := datastore.NewKey(c, "Gtcha", id, 0, nil)
	if _, err := datastore.Put(c, key, g); err != nil {
		return err
	}

	return nil
}

// CacheImageURI caches the generated image uri.
func CacheImageURI(c appengine.Context, id, uri string) error {
	return memcache.Add(c, &memcache.Item{
		Key:   "uri" + id,
		Value: []byte(uri),
	})
}

// GetImageURI attempts to find the base64 encoded ata URI of an image in memcache.
func GetImageURI(c appengine.Context, id string) (string, error) {
	item, err := memcache.Get(c, "uri"+id)
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}
