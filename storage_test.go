package gtcha

import (
	"encoding/json"
	"testing"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func TestGetApp(t *testing.T) {
	domains := []string{"bowery.io"}

	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	key := datastore.NewKey(c, "GtchaApp", uuid.New(), 0, nil)
	app := &GtchaApp{
		Name:    "testApp",
		Secret:  key.StringID(),
		APIKey:  uuid.New(),
		Domains: domains,
	}
	if _, err := datastore.Put(c, key, app); err != nil {
		t.Fatal(err)
	}

	// increases likelihood that app goes in the datastore before call to Get
	<-time.After(500 * time.Millisecond)

	retApp, err := GetApp(c, app.APIKey, "bowery.io")
	if err != nil {
		t.Fatal(err)
	}

	if retApp.Secret != app.Secret {
		t.Fatalf(
			"returned app differs from expected app:\nreturned %#v\nexpected %#v", retApp, app,
		)
	}
}

func TestGetGtcha(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	g := &gtcha{
		Tag:   "dog",
		In:    []string{"a", "b", "c"},
		Out:   []string{"d", "e", "f"},
		Maybe: []string{"g", "h", "i"},
	}

	id := "testID"
	if err = SaveGtcha(c, g, id); err != nil {
		t.Fatal(err)
	}

	retG, err := GetGtcha(c, id)
	if err != nil {
		t.Fatal(err)
	}

	for i := range g.In {
		if g.In[i] != retG.In[i] {
			t.Fatalf("expected same string %s, got %s", g.In[i], retG.In[i])
		}
	}
}

func TestSaveGtcha(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	g := &gtcha{
		Tag:   "dog",
		In:    []string{"a", "b", "c"},
		Out:   []string{"d", "e", "f"},
		Maybe: []string{"g", "h", "i"},
	}

	id := "testID"
	if err = SaveGtcha(c, g, id); err != nil {
		t.Fatal(err)
	}

	item, err := memcache.Get(c, id)
	if err != nil {
		t.Fatal(err)
	}

	retG := new(gtcha)
	if err = json.Unmarshal(item.Value, retG); err != nil {
		t.Fatal(err)
	}

	for i := range g.In {
		if g.In[i] != retG.In[i] {
			t.Fatalf("expected same string %s, got %s", g.In[i], retG.In[i])
		}
	}

	key := datastore.NewKey(c, "Gtcha", id, 0, nil)
	if err := datastore.Get(c, key, retG); err != nil {
		t.Fatal(err)
	}

	for i := range g.In {
		if g.In[i] != retG.In[i] {
			t.Fatalf("expected same string %s, got %s", g.In[i], retG.In[i])
		}
	}
}
