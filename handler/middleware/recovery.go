package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		// defer -> deferが実行されている関数が終了した後に実行される
		// recover() -> 実行した1回だけpanicに渡した引数が戻り値として返される
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("panicの引数:", err)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
