//
// gitcha application logic
//

package gtcha

import (
	"net/http"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"github.com/Bowery/gtcha/giphy"
)

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
	ID     string   `json:"id,omitempty"`
	Tag    string   `json:"tag"`
	Images []string `json:"images"`
}

type gtcha struct {
	Tag     string
	In      []string
	Out     []string
	Maybe   []string
	IsHuman bool
}

// newGtcha returns the internal representation of the GIF captcha.
// This function works by making a lot of API calls in parrallel.
func newGtcha(c *http.Client) (*gtcha, error) {
	tag, err := giphy.GetTag(c)
	if err != nil {
		return nil, err
	}

	var (
		in      []string
		out     []string
		maybe   []string
		wg      sync.WaitGroup
		errOnce sync.Once
		errCh   = make(chan error)
	)

	// closes over some variables, so that we can get all the imges in parrallel
	processImages := func(fn func(*http.Client, string, int) ([]*giphy.Image, error)) []string {
		apiImgs, err := fn(c, tag, 0)
		if err != nil {
			errOnce.Do(func() { errCh <- err })
		}
		imgs := make([]string, len(apiImgs))
		for i, img := range apiImgs {
			imgs[i] = img.ID // TODO: maybe change this to one of the URLs?
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
		if checkImageIn(g.Out, img) {
			return false
		}

		if checkImageIn(g.In, img) {
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
			if checkImageIn(g.Maybe, img) {
				giphy.ConfirmTag(c, g.Tag, img)
			}
		}(img)
	}

	return true
}

func (g *gtcha) toCaptcha() *Captcha {
	imgs := make([]string, 0, len(g.In)+len(g.Out)+len(g.Maybe))
	for _, l := range [][]string{g.In, g.Out, g.Maybe} {
		for _, img := range l {
			imgs = append(imgs, img)
		}
	}

	return &Captcha{
		ID:     uuid.New(),
		Tag:    g.Tag,
		Images: imgs,
	}
}

func checkImageIn(imgs []string, img string) bool {
	for _, i := range imgs {
		if i == img {
			return true
		}
	}

	return false
}
