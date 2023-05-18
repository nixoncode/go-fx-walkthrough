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
	"net"
	"net/http"
)

func main() {
	fx.New(
		fx.Provide(NewHttpServer),
		fx.Invoke(func(*http.Server) {}), // need to find out why we need this and if it's possible to avoid
	).Run()
}

func NewHttpServer(lc fx.Lifecycle) *http.Server {
	srv := &http.Server{Addr: ":8080"}

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
