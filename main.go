package gitcha

import (
	"encoding/json"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHndlr)
	http.HandleFunc("/is_h", isHumanHndlr)
	http.HandleFunc("/isn_h", isNotHumanHndlr)
	http.HandleFunc("/c", captchaHndlr)
}

func rootHndlr(w http.ResponseWriter, r *http.Request) {
	captcha, err := NewCaptcha()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(captcha)

	if _, err := w.Write(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func isHumanHndlr(w http.ResponseWriter, r *http.Request)    {}
func isNotHumanHndlr(w http.ResponseWriter, r *http.Request) {}

type captchaResp struct {
	cert string
	test string
}

func captchaHndlr(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cr captchaResp
	err := decoder.Decode(&cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cr.cert == "" || cr.test == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ok, err := AssociateImages(cr.cert, cr.test)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Invalid input", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
