//
// appengine datastore and memcache logic
//

package gtcha

import (
	"encoding/json"
	"errors"
	"sync"

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
	var (
		o     sync.Once
		wg    sync.WaitGroup
		gCh   = make(chan *gtcha)
		errCh = make(chan error)
	)

	//
	// This has to solve the following problem:
	//
	// You have two routines trying to get data. You want to return the data returned on a
	// successful call and you only want to send an error once both routines have errored.
	//
	// We mostly solve this by using a waitgroup, but this also requires some additional
	// machinery. We send the completed data over a channel and then select on that
	// channel and another channel that sends an error (if there is one) once the
	// waitgroup is finished.
	//

	wg.Add(2)

	// memcache
	go func() {
		defer wg.Done()
		item, err := memcache.Get(c, id)
		if err != nil {
			return
		}

		g := new(gtcha)
		if err = json.Unmarshal(item.Value, g); err != nil {
			return
		}

		o.Do(func() { gCh <- g })
	}()

	// datastore
	go func() {
		defer wg.Done()
		key := datastore.NewKey(c, "Gtcha", id, 0, nil)
		g := new(gtcha)
		if err := datastore.Get(c, key, g); err != nil {
			errCh <- err
			return
		}

		o.Do(func() { gCh <- g })
		errCh <- nil
	}()

	eCh := make(chan error)
	go func() {
		err := <-errCh
		wg.Wait()
		o.Do(func() { eCh <- err })
	}()

	select {
	case g := <-gCh:
		return g, nil
	case err := <-eCh:
		return nil, err
	}
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
			Key:   id,
			Value: buf,
		})
	}()

	key := datastore.NewKey(c, "Gtcha", id, 0, nil)
	if _, err := datastore.Put(c, key, g); err != nil {
		return err
	}

	return nil
}
