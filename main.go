package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// TODO: サーバーをlistenする
	fmt.Println("done前")
	// http.ListenAndServe(port, mux)
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	srv.ListenAndServe()
	// os.Interrupt か os.Kill を受け取ったらこれ以降のコードが実行される
	<-ctx.Done()
	fmt.Println("done後")
	// os.Interrupt か os.Kill を受け取ってから5秒待つ
	ctxTimeOut, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = srv.Shutdown(ctxTimeOut)

	if err != nil {
		return err
	}

	return nil
}
