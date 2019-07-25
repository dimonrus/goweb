package goweb

import (
	"context"
	"github.com/dimonrus/gocli"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// Make server and listen
func PrepareServerAndListen(app gocli.Application, config Config, routes http.Handler, connState func(net.Conn, http.ConnState)) {
	server := &http.Server{
		Addr:    config.Host + ":" + strconv.Itoa(config.Port),
		Handler: routes,
	}

	if connState != nil {
		server.ConnState = connState
	} else {
		server.ReadTimeout = time.Second * time.Duration(config.Timeout.Read)
		server.WriteTimeout = time.Second * time.Duration(config.Timeout.Write)
		server.IdleTimeout = time.Second * time.Duration(config.Timeout.Idle)
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			app.GetLogger(gocli.LogLevelDebug).Error()
		}
	}()

	app.GetLogger(gocli.LogLevelDebug).Infof("Web server started at %s", server.Addr)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(config.Timeout.Read))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := server.Shutdown(ctx)
	if err != nil {
		app.FatalError(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	app.GetLogger(gocli.LogLevelDebug).Warn("Server shutting down")
	os.Exit(0)

	return
}
