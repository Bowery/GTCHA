// Package gtcha uses Google Appengine and the Giphy API to make a front-end plugin to
// distinguish humans and robots.
//
// TODO
//
// 1. Proper error handling
// 2. Removing gtchas from datastore
// 3. Make a running task that makes captchas continuously
// 4. Move everything over so that all functions/methods use contexts instead of arguments like `*http.Client`
package gtcha
