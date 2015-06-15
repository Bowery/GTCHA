package GTCHA

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	giphyAPI = "https://api.giphy.com/v1"
)

const (
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
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New("giphy API error")
		}
		return errors.New("giphy API error " + string(buf))
	}

	if out == nil {
		return nil
	}

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(out); err != nil {
		return err
	}

	return nil
}

// GetTag gets a random tag from the giphy API.
func GetTag(c *http.Client) (string, error) {
	url := genGiphyURL("tag/random", "")
	res := new(tagResult)
	if err := makeGiphyCall(c, url, "GET", res); err != nil {
		return "", err
	}

	return res.Data, nil
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
	// TODO(r-medina): add pagination

	url := genGiphyURL("gifs/search", "q="+tag)
	res := new(searchResult)
	if err := makeGiphyCall(c, url, "GET", res); err != nil {
		return nil, err
	}

	return res.Data, nil
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
	url := genGiphyURL("gifs/search", "maybe="+tag)
	res := new(searchResult)
	if err := makeGiphyCall(c, url, "GET", res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

// ConfirmTag tells the giphy API that an image was tagged by a
// user we know to be human.
func ConfirmTag(c *http.Client, tag, img string) error {
	url := genGiphyURL("confirm", "i="+img+"&q="+tag)

	if err := makeGiphyCall(c, url, "PUT", nil); err != nil {
		return err
	}

	return nil
}
