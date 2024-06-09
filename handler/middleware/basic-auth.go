package middleware

import (
	"fmt"
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// basic認証で入力されたidとpassと入力があるかないかを判定するflagokが戻り値
		id, pass, ok := r.BasicAuth()

		correctId := os.Getenv("BASIC_AUTH_USER_ID")
		correctPass := os.Getenv("BASIC_AUTH_PASSWORD")

		fmt.Println("id:", id, "pass:", pass, "ok:", ok)

		if !ok || !(correctId == id) || !(correctPass == pass) {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			h.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}
