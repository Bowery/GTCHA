//
// tests for gitcha.go functions
//

package gtcha

import (
	"net/http"
	"testing"
	"time"

	"code.google.com/p/appengine-go/appengine/aetest"
)

var g = &gtcha{
	Tag: "cute puppy",
	In: []gimg{
		gimg{"in_1", "DgzJFvt6StyFi", "http://media0.giphy.com/media/DgzJFvt6StyFi/100w.gif"},
		gimg{"in_2", "13OWYvLosSoK2I", "http://media0.giphy.com/media/13OWYvLosSoK2I/100w.gif"},
	},
	Out: []gimg{
		gimg{"out_1", "14f3BPP6SCc0Mw", "http://media1.giphy.com/media/14f3BPP6SCc0Mw/100w.gif"},
		gimg{"out_2", "OrkjamOz6caA0", "http://media0.giphy.com/media/OrkjamOz6caA0/100w.gif"},
	},
	Maybe: []gimg{
		gimg{"maybe_1", "K2BQcrA30rUMU", "http://media0.giphy.com/media/K2BQcrA30rUMU/100w.gif"},
		gimg{"maybe_2", "f4i4IpVQVhtu0", "http://media4.giphy.com/media/f4i4IpVQVhtu0/100w.gif"},
	},
}

var img = gimg{"in_1", "DgzJFvt6StyFi", "http://media0.giphy.com/media/DgzJFvt6StyFi/100w.gif"}

func TestNewGtcha(t *testing.T)    { t.Log("TODO") }
func TestVerifyGtcha(t *testing.T) { t.Log("TODO") }

func TestToGImg(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	_, err = img.toGImg(c, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}

	// TODO(r-medina): actually write tests
	t.Log("TODO(r-medina): actually write tests")
}

func TestToCaptcha(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	start := time.Now()
	captcha1, err := g.toCaptcha(c, http.DefaultClient)
	d1 := time.Since(start)
	t.Logf("without caching %+v\n", d1) // output for debug
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(500 * time.Millisecond) // make sure images are cached

	start = time.Now()
	captcha2, err := g.toCaptcha(c, http.DefaultClient)
	d2 := time.Since(start)
	t.Logf("with caching %+v\n", d2) // output for debug
	if err != nil {
		t.Fatal(err)
	}

	if d2 > d1 {
		t.Logf("with caching it took longer %v vs %v", d2, d1)
	}

	isIn := false
	for _, img := range captcha1.Images {
		if captcha2.Images[0].ID == img.ID {
			isIn = true
			break
		}
	}
	if !isIn {
		t.Fatalf("image %s not in first captcha", captcha2.Images[0].ID)
	}

	n := 16
	start = time.Now()
	for range make([]struct{}, n) {
		g.toCaptcha(c, http.DefaultClient)
	}
	d3 := time.Since(start)
	t.Logf("%d times %+v\n", n, d3) // output for debug
}

func TestIsImageIn(t *testing.T) { t.Log("TODO") }
