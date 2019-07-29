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

// Init Web Application
func NewApplication(config Config, app gocli.Application, connState func(net.Conn, http.ConnState)) *Application {
	return &Application{
		config: config,
		app:    app,
		server: &http.Server{
			Addr:         config.Host + ":" + strconv.Itoa(config.Port),
			ReadTimeout:  time.Second * time.Duration(config.Timeout.Read),
			WriteTimeout: time.Second * time.Duration(config.Timeout.Write),
			IdleTimeout:  time.Second * time.Duration(config.Timeout.Idle),
			ConnState:    connState,
		},
	}
}

// Make server and listen
func (a *Application) Listen(routes http.Handler) {
	// Set routes
	a.server.Handler = routes
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			a.app.GetLogger(gocli.LogLevelDebug).Error("Can't listen: ", err.Error())
		}
	}()

	a.app.GetLogger(gocli.LogLevelDebug).Infof("Web server started at %s", a.server.Addr)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(a.config.Timeout.Read))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := a.server.Shutdown(ctx)
	if err != nil {
		a.app.FatalError(err)
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your Application should wait for other services
	// to finalize based on context cancellation.
	a.app.GetLogger(gocli.LogLevelDebug).Warn("Server shutting down")
	os.Exit(0)

	return
}
