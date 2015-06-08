package giphy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unrolled/render"
)

var (
	renderer = render.New(render.Options{
		IndentJSON:    true,
		IsDevelopment: true,
	})
)

func TestGetTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(testGetTagHandler))
	defer server.Close()
	giphyAPI = server.URL

	_, err := GetTag(http.DefaultClient)
	if err != nil {
		t.Error(err)
	}
}

func testGetTagHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func TestGetOtherTag(t *testing.T)                                     {}
func testGetOtherTagHandler(rw http.ResponseWriter, req *http.Request) {}

func TestGetImagesTagged(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(testGetImagesTaggedHandler))
	defer server.Close()
	giphyAPI = server.URL

	images, err := GetImagesTagged(http.DefaultClient, "test-tag", 0)
	if err != nil {
		t.Error(err)
	}

	if len(images) != 2 {
		t.Error("unexpected result")
	}
}

func testGetImagesTaggedHandler(rw http.ResponseWriter, req *http.Request) {
	images := []*Image{
		&Image{ID: "test-id-1"},
		&Image{ID: "test-id-2"},
	}

	renderer.JSON(rw, http.StatusOK, images)
}

func TestGetImagesNotTagged(t *testing.T)                                     {}
func testGetImagesNotTaggedHandler(rw http.ResponseWriter, req *http.Request) {}

func TestGetImagesMaybeTagged(t *testing.T)                                     {}
func testGetImagesMaybeTaggedHandler(rw http.ResponseWriter, req *http.Request) {}
