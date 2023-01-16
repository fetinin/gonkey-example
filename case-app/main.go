package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	app "gonkey-example/case-app/internal"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	db, err := app.NewDB(ctx, "postgres://service:service@localhost:6543/service?sslmode=disable")
	if err != nil {
		panic(err)
	}

	const addr = ":7700"
	const externalService = "https://names.drycodes.com"
	fmt.Printf("Starting server listening on %s", addr)
	srv := http.Server{Addr: addr, Handler: app.NewAPI(db, externalService)}
	defer srv.Close()

	go func() {
		srv.ListenAndServe()
	}()
	<-ctx.Done()
}
