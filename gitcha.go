//
//
//

package gitcha

import (
	"errors"
	"fmt"
)

const (
	giphyAPI        = "https://api.giphy.com"
	captchaEndpoint = "captcha"
	keyLen          = 128
	secLen          = keyLen >> 1
)

var (
	ErrNameTooLong = fmt.Errorf("name too long; expected length %d", keyLen)
	ErrNameExists  = errors.New("name exists")
)

type gitchaApp struct {
	Name string `json:"name,omitempty"`

	// ID is the public key of an app for use in its front-end.
	ID string `json:"id"`

	// Secret is an app's secret key.
	Secret string `json:"secret,omitempty"`
}

// NewApp allocates a new app and sets its ID and secret based on the name.
// The logic for checking that the name is new and the right length is outside of this
// function, as it simplifies testing.
func NewApp(name string) (*gitchaApp, error) {
	app := new(gitchaApp)
	app.Name = name

	if err := app.genSecret(); err != nil {
		return nil, err
	}

	if err := app.genID(); err != nil {
		return nil, err
	}

	return app, nil
}

// genSecret resets an app's secret based on its name.
func (app *gitchaApp) genSecret() error {
	if app.Name == "" {
		return errors.New("app has no name")
	}

	sec, err := genSecret(app.Name)
	if err != nil {
		return err
	}

	app.Secret = sec

	return nil
}

// genID resets an app's ID based on it's secret key.
func (app *gitchaApp) genID() error {
	if app.Secret == "" {
		return errors.New("app has no secret")
	}

	id, err := genSecret(app.Secret)
	if err != nil {
		return err
	}

	app.ID = id

	return nil
}
