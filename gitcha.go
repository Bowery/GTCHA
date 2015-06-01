package gitcha

import "path"

type captcha struct {
	Certain  image   `json:"certain"`   // image for which description is certain
	TestImgs []image `json:"test_imgs"` // array of images thought to be related to `Certain`
}

type image struct {
	URI  string `json:"uri"`            // path to the image
	Desc string `json:"desc,omitempty"` // description of the image
}

func (img *image) ID() string {
	return path.Base(img.URI)
}

// NewCaptcha generates a new captcha struct to send back to the user
func NewCaptcha() (*captcha, error) {
	return nil, nil
}

// AssociateImages attempts to associate `test` with `cert`.
func AssociateImages(cert, test string) (bool, error) {
	return true, nil
}
