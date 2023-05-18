// Package main
// Created by GoLand.
// User: nixon
// Date: 18/5/2023
// Time: 06:55
package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {
	fx.New(
		fx.Provide(NewHttpServer, NewEchoHandler, NewServerMux),
		fx.Invoke(func(*http.Server) {}), // need to find out why we need this and if it's possible to avoid
	).Run()
}

func NewHttpServer(lc fx.Lifecycle, mux *http.ServeMux) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Http Server starting at", srv.Addr)
			go func() {
				_ = srv.Serve(listener)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (*EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to handle request", err)
	}
}

func NewServerMux(echo *EchoHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/echo", echo)
	return mux
}
