package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

func GetOS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// コンテキスト -> 関数をまたいで値の転送が出来る箱みたいなイメージ
		ua := useragent.Parse(r.UserAgent())

		ctx := r.Context()

		// key重複の可能性があるのでkeyをstringで定義するのは良くないらしい
		ctx = context.WithValue(ctx, "osName", ua.OS)

		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
