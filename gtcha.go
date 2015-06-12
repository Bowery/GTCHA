//
// gitcha application logic
//

package gtcha

import (
	"bytes"
	"errors"
	"net/http"
	"sync"

	"code.google.com/p/appengine-go/appengine"
	"code.google.com/p/go-uuid/uuid"

	"github.com/Bowery/gtcha/giphy"
)

const gifType = "image/gif"

// GtchaApp represents an app that is using our GIF captcha.
type GtchaApp struct {
	Name string `json:"name,omitempty"`

	// Secret is an app's secret key.
	Secret string `json:"secret,omitempty"`

	// public key of an app for use in its front-end.
	APIKey string `json:"api_key,omitempty"`

	Domains []string `json:"domains,omitempty"`
}

// Captcha represents our user-facing GIF captcha.
type Captcha struct {
	ID     string `json:"id,omitempty"`
	Tag    string `json:"tag"`
	Images []GImg `json:"images"`
}

type GImg struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

type gtcha struct {
	Tag     string
	In      []gimg
	Out     []gimg
	Maybe   []gimg
	IsHuman bool
}

type gimg struct {
	ID   string `json:"id"`
	GID  string `json:"gid"` // the ID on the giphy backend
	GURI string `json:"guri"`
}

// newGtcha returns the internal representation of the GIF captcha.
// This function works by making a lot of API calls in parrallel.
func newGtcha(c *http.Client) (*gtcha, error) {
	tag, err := giphy.GetTag(c)
	if err != nil {
		return nil, err
	}

	var (
		in      []gimg
		out     []gimg
		maybe   []gimg
		wg      sync.WaitGroup
		errOnce sync.Once
		errCh   = make(chan error)
	)

	// closes over some variables, so that we can get all the imges in parrallel
	processImages := func(fn func(*http.Client, string, int) ([]*giphy.Image, error)) []gimg {
		apiImgs, err := fn(c, tag, 0)
		if err != nil {
			errOnce.Do(func() { errCh <- err })
		}
		imgs := make([]gimg, len(apiImgs))
		for i, img := range apiImgs {
			imgs[i] = gimg{
				ID:   uuid.New(),
				GID:  img.ID,
				GURI: img.Images.FixedWidthSmall.URL,
			}
		}

		return imgs
	}

	wg.Add(3)

	go func() {
		in = processImages(giphy.GetImagesTagged)
		wg.Done()
	}()

	go func() {
		out = processImages(giphy.GetImagesNotTagged)
		wg.Done()
	}()

	go func() {
		maybe = processImages(giphy.GetImagesMaybeTagged)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		errOnce.Do(func() { errCh <- nil })
	}()

	err = <-errCh
	if err != nil {
		return nil, err
	}

	g := &gtcha{
		Tag:   tag,
		In:    in,
		Out:   out,
		Maybe: maybe,
	}

	return g, nil
}

func verifyGtcha(c *http.Client, g *gtcha, in []string) bool {
	isHuman := false
	for _, img := range in {
		if isImageIn(g.Out, img) {
			return false
		}

		if isImageIn(g.In, img) {
			isHuman = true
		}
	}

	if !isHuman {
		return false
	}

	// check the images that might be tagged g's tag against the user's submitted images
	// to let the giphy API know that a human verif
	for _, img := range in {
		go func(img string) {
			if isImageIn(g.Maybe, img) {
				giphy.ConfirmTag(c, g.Tag, img)
			}
		}(img)
	}

	return true
}

func (img *gimg) toGImg(c appengine.Context, httpC *http.Client) (*GImg, error) {
	gotten := false
	var mtx sync.Mutex

	f1 := func() interface{} {
		mtx.Lock()
		defer mtx.Unlock()
		uri, err := GetImageURI(c, img.GID)
		if err != nil {
			return nil
		}

		gotten = true
		return &GImg{img.ID, uri}
	}

	f2 := func() (interface{}, error) {
		mtx.Lock()
		if gotten {
			mtx.Unlock()
			return nil, errors.New("done")
		}
		mtx.Unlock()

		req, err := http.NewRequest("GET", img.GURI, nil)
		if err != nil {
			return nil, err
		}
		res, err := httpC.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(res.Body); err != nil {
			return nil, err
		}

		uri := dataURI(buf.Bytes(), gifType)
		go CacheImageURI(c, img.GID, uri)
		return &GImg{img.ID, uri}, nil
	}

	i, err := Get(f1, f2)
	if err != nil {
		return nil, err
	}

	return i.(*GImg), nil
}

func (g *gtcha) toCaptcha(c appengine.Context, httpC *http.Client) (*Captcha, error) {
	var (
		o     sync.Once
		wg    sync.WaitGroup
		errCh = make(chan error)
		imgCh = make(chan *GImg) // saves overhead to use pointer
		n     = len(g.In) + len(g.Out) + len(g.Maybe)
	)

	wg.Add(n)

	for _, l := range [][]gimg{g.In, g.Out, g.Maybe} {
		for _, img := range l {
			go func(img gimg) {
				defer wg.Done()

				i, err := img.toGImg(c, httpC)
				if err != nil {
					o.Do(func() { errCh <- err })
					return
				}

				imgCh <- i
			}(img)
		}
	}

	go func() {
		wg.Wait()
		close(imgCh)
	}()

	imgs := make([]GImg, 0, n)
LOOP:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				return nil, err
			}
		case img, ok := <-imgCh:
			if !ok {
				break LOOP
			}
			imgs = append(imgs, *img)
		}
	}

	captcha := &Captcha{
		ID:     uuid.New(),
		Tag:    g.Tag,
		Images: imgs,
	}

	return captcha, nil
}

func isImageIn(imgs []gimg, img string) bool {
	for _, i := range imgs {
		if i.ID == img {
			return true
		}
	}

	return false
}
