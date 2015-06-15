package GTCHA

import (
	"net/http"
	"testing"
)

func TestGenGiphyURL(t *testing.T) {
	expected := "https://api.giphy.com/v1/gifs/search?q=puppy&api_key=dc6zaTOxFJmzC"
	if url := genGiphyURL("gifs/search", "q=puppy"); url != expected {
		t.Fatalf("expected %s, got %s", expected, url)
	}
}

func TestMakeGiphyCall(t *testing.T) {
	url := "https://api.giphy.com/v1/gifs/search?q=puppy&api_key=dc6zaTOxFJmzC"
	res := new(searchResult)
	if err := makeGiphyCall(http.DefaultClient, url, "GET", res); err != nil {
		t.Fatal(err)
	}
}

func TestGetTag(t *testing.T)      { t.Log("TODO") }
func TestGetOtherTag(t *testing.T) { t.Log("TODO") }

func TestGetImagesTagged(t *testing.T) {
	_, err := GetImagesTagged(http.DefaultClient, "puppy", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetImagesNotTagged(t *testing.T)   { t.Log("TODO") }
func TestGetImagesMaybeTagged(t *testing.T) { t.Log("TODO") }
func TestConfirmTagged(t *testing.T)        { t.Log("TODO") }
