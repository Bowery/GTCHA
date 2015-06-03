//
// giphy API logic
//

package gitcha

import (
	"errors"
	"fmt"
)

const (
	giphyAPI        = "https://api.giphy.com"
	giphyVer        = "v1"
	captchaEndpoint = "captcha"
	keyLen          = 128
	secLen          = keyLen >> 1
)

var (
	ErrNameTooLong = fmt.Errorf("name too long; expected length %d", keyLen)
	ErrNameExists  = errors.New("name exists")
)

type GtchaApp struct {
	Name string `json:"name,omitempty"`

	// Secret is an app's secret key.
	Secret string `json:"secret,omitempty"`

	// public key of an app for use in its front-end.
	APIKey string `json:"api_key,omitempty"`

	Domains []string `json:"domains,omitempty"`
}

type Captcha struct {
	ID     string   `json:"id,omitempty"`
	Tag    string   `json:"tag"`
	Images []string `json:"images"`
}

type gtcha struct {
	id    string
	tag   string
	in    []string
	out   []string
	maybe []string
}

func newGtcha() (*gtcha, error) {
	tag, err := genTag()
	if err != nil {
		return nil, err
	}

	in, err := getImagesTagged(tag)
	if err != nil {
		return nil, err
	}

	out, err := getImagesNotTagged(tag)
	if err != nil {
		return nil, err
	}

	maybe, err := getImagesMaybeTagged(tag)
	if err != nil {
		return nil, err
	}

	g := &gtcha{
		tag:   tag,
		in:    in,
		out:   out,
		maybe: maybe,
	}

	return g, nil
}

func (g *gtcha) toCaptcha() *Captcha {
	imgs := make([]string, 0, len(g.in)+len(g.out)+len(g.maybe))
	for _, l := range [][]string{g.in, g.out, g.maybe} {
		for _, img := range l {
			imgs = append(imgs, img)
		}
	}

	return &Gtcha{
		ID:     g.id,
		Tag:    g.tag,
		Images: imgs,
	}
}

func getImagesTagged(tag string) ([]string, error) {
	return nil, nil
}

func getImagesNotTagged(tag string) ([]string, error) {
	otag, err := getOtherTag(tag)
	if err != nil {
		return nil, err
	}

	return getImagesTagged(otag)
}

func getOtherTag(tag string) (string, error) {
	otag := tag
	for otag == tag {
		otag, err := genTag()
		if err != nil {
			return "", err
		}
	}

	return otag, nil
}
