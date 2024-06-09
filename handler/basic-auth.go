package handler

import (
	"fmt"
	"net/http"
	"os"
)

type BasicAuthHandler struct{}

func NewBasicAuthHandler() *BasicAuthHandler {
	return &BasicAuthHandler{}
}

func (b *BasicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("basicAuthハンドラーが発火した")

	// 起動時にコマンドライン引数としてenv渡すのはターミナルじゃできなかった.git bashならいけた
	// ex)BASIC_AUTH_USER_ID=testUser BASIC_AUTH_PASSWORD=testPass go run main.go
	fmt.Println("os.Getenv(BASIC_AUTH_USER_ID):", os.Getenv("BASIC_AUTH_USER_ID"))
	fmt.Println("os.Getenv(BASIC_AUTH_PASSWORD):", os.Getenv("BASIC_AUTH_PASSWORD"))

}
