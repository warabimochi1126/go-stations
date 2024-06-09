package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func GetTime(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		beforeHandleTime := time.Now()

		h.ServeHTTP(w, r)

		afterHandleTime := time.Now()
		diffMill := afterHandleTime.Sub(beforeHandleTime).Milliseconds()

		// any型が返ってくるので型キャストしてエラー帰ってきたらpanic投げる
		osName, ok := r.Context().Value("osName").(string)
		if !ok {
			panic("Middleware:GetTimeの型キャストでエラーが発生した。")
		}

		var responseJson model.ResponseLatencyJson
		responseJson.Timestamp = beforeHandleTime
		responseJson.Latency = diffMill
		responseJson.Path = r.URL.Path
		responseJson.OS = osName

		json, err := json.Marshal(responseJson)
		if err != nil {
			fmt.Println("json変換でエラー発生した。")
		}
		fmt.Println("json変換後のresponseJson:", string(json))
	}

	return http.HandlerFunc(fn)
}
