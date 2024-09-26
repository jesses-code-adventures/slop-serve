package api

import (
	"encoding/json"
	"net/http"
)

func (a AppApi) preferCookieOverHeader(r *http.Request, key string) string {
	if cookie, err := r.Cookie(key); err == nil {
		return cookie.Value
	}
	return r.Header.Get(key)
}

func (a AppApi) setAuthorization(w http.ResponseWriter, r *http.Request, newToken string) {
	cookie := &http.Cookie{
		Name:  "Authorization",
		Value: newToken,
	}
	http.SetCookie(w, cookie)
}

func (a *AppApi) jsonResponse(w http.ResponseWriter, r *http.Request, key string, value interface{}) {
	w.Header().Add("Content-Type", "application/json")
	response := map[string]interface{}{key: value}
	marshalled, err := json.Marshal(response)
	if err != nil {
		http.Error(w, a.jsonErrorString("Json marshalling error"), http.StatusInternalServerError)
		return
	}
	w.Write(marshalled)
}

func (a *AppApi) jsonErrorString(errorMessage string) string {
	response := map[string]interface{}{"error": errorMessage}
	marshalled, err := json.Marshal(response)
	if err != nil {
		panic("couldn't marshal json message")
	}
	return string(marshalled)
}
