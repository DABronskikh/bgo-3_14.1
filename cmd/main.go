package main

import (
	"context"
	app "github.com/DABronskikh/bgo-3_14.1/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "8080"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	log.Println(host)
	log.Println(port)

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

func execute(addr string) (err error) {
	dsn := "postgres://postgres:example@localhost:5432/db"
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Release()

	mux := http.NewServeMux()
	application := app.NewServer(mux, ctx, conn)
	application.Init()

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
