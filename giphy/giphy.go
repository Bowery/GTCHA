//
// giphy API interactions
//

package giphy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	giphyAPI        = "https://api.giphy.com/v1"
	giphyKey        = "dc6zaTOxFJmzC"
	captchaEndpoint = "captcha"
)

func genGiphyURL(endpoint, query string) string {
	return fmt.Sprintf("%s/%s?%s&api_key=%s", giphyAPI, endpoint, query, giphyKey)
}

func makeGiphyCall(c *http.Client, url, method string, out interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(out); err != nil {
		return err
	}

	return nil
}

// GetTag gets a random tag from the giphy API.
func GetTag(c *http.Client) (string, error) {
	return "", nil
}

// GetOtherTag gets a tag that isn't `tag` from the giphyAPI.
func GetOtherTag(c *http.Client, tag string) (string, error) {
	i := 0
	otag := tag
	for otag == tag {
		var err error
		otag, err = GetTag(c)
		if err != nil {
			return "", err
		}
		i++
		if i == 5 {
			return "", errors.New("could not get new tag")
		}
	}

	return otag, nil
}

// GetImagesTagged returns a slice of images that definitely match `tag`.
func GetImagesTagged(c *http.Client, tag string, page int) ([]*Image, error) {
	url := genGiphyURL("gifs/search", "q="+tag)
	var imgs []*Image
	if err := makeGiphyCall(c, url, "GET", imgs); err != nil {
		return nil, err
	}

	return nil, nil
}

// GetImagesNotTagged gets images that do not match a specified tag.
func GetImagesNotTagged(c *http.Client, tag string, page int) ([]*Image, error) {
	otag, err := GetOtherTag(c, tag)
	if err != nil {
		return nil, err
	}

	return GetImagesTagged(c, otag, page)
}

// GetImagesMaybeTagged returns images that *might* match `tag`.
func GetImagesMaybeTagged(c *http.Client, tag string, page int) ([]*Image, error) {
	return nil, nil
}
