package handler

import (
	"fmt"
	"net/http"
)

// http.Handlerインターフェースを満たす構造体の定義
// exportするために頭文字を大文字にしなければならない
type DoPanicHandler struct{}

// ファイル内で定義しておくとコードジャンプ出来る
func NewDoPanicHandler() *DoPanicHandler {
	return &DoPanicHandler{}
}

func (d *DoPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("次の行でpanic発火します。")
	panic("panic発火")
	fmt.Println("panic発火したためこの行は実行されません。")
}
