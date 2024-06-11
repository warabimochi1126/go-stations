package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// Go基礎編Station6
	// os.Interrupt か os.Kill を受け取るまで待つ
	// 受け取ったらctx.Done()から先を発火させられる？
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// メインゴルーチンが終了したら強制終了されるからサーバ停止用のゴルーチンが終了するまで待つ
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// 平行で実行することでメインのゴルーチンの実行を止めない
	go func() {
		// サーバ停止してからWaitGroupのカウンタをデクリメントする
		defer wg.Done()

		// contextという平行実行される上でそもそも独立してる実行間で終了信号を受け渡している？
		<-ctx.Done()
		fmt.Println("goroutineの中で実行されている")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// サーバーを閉じる
		srv.Shutdown(ctx)
	}()

	// TODO: サーバーをlistenする
	fmt.Println("listen前")
	// http.ListenAndServe(port, mux)

	srv.ListenAndServe()

	fmt.Println("listen後")
	// WaitGroupのカウンタが0になるまでメインのゴルーチンを終了しない
	wg.Wait()

	return nil
}
