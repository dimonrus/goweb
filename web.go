package goweb

import (
	"context"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/porterr"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	CommandStart   = "start"
	CommandStop    = "stop"
	CommandRestart = "restart"
	CommandStatus  = "status"
	CommandWeb     = "web"
)

var CommandActions = []string{
	CommandStart, CommandStop, CommandRestart, CommandStatus,
}

// Config Web configuration
type Config struct {
	// Web application port
	Port int
	// Web application server host
	Host string
	// Web application url. For logging
	Url string
	// Server timeouts
	Timeout struct {
		// Timeout read
		Read int
		// Timeout write
		Write int
		// Timeout idle
		Idle int
	}
}

// Application Web Application struct
type Application struct {
	// The console application base interface.
	// Required for start
	gocli.Application
	// Config of web application
	config Config
	// Http server type
	server *http.Server
	// Exit web server
	exit chan bool
}

// DecomposeCommand Parse gocli.Command
func DecomposeCommand(command *gocli.Command) (action string, arguments []gocli.Argument, e porterr.IError) {
	args := command.Arguments()
	if len(args) > 0 {
		if args[0].Name != CommandWeb {
			e = porterr.New(porterr.PortErrorArgument, "Web command must start with keyword 'web'")
			return
		}
		if len(args) < 2 {
			e = porterr.New(porterr.PortErrorArgument, "Web command must contain action")
			return
		}
		action = args[1].Name
		if !gohelp.ExistsInArrayString(action, CommandActions) {
			e = porterr.New(porterr.PortErrorArgument, "Web command action is unknown: "+action)
		}
	} else {
		e = porterr.New(porterr.PortErrorArgument, "Web command is empty")
	}
	return
}

// Graceful shutdown web application
func (a *Application) shutdown() <-chan bool {
	go func() {
		sig := make(chan os.Signal, 1)
		// Accept graceful shutdowns when quit (Ctrl+C)
		signal.Notify(sig, os.Interrupt)
		<-sig
		a.exit <- true
	}()
	return a.exit
}

// WebCommander Web command processor
func (a *Application) WebCommander(command *gocli.Command) {
	a.SuccessMessage("Receive command: "+command.String(), &gocli.Command{})
	action, _, e := DecomposeCommand(command)
	if e != nil {
		a.FatalError(e)
		return
	}
	switch action {
	case CommandStop:
		a.AttentionMessage("Stopping web server by command... " + command.String())
		go func() {
			a.exit <- true
		}()
	}
}

// Listen Make server and listen
func (a *Application) Listen(routes http.Handler) {
	// Set routes
	a.server.Handler = routes
	// Run server so that it doesn't block.
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			a.GetLogger().Error("Can't listen: ", err.Error())
		}
	}()
	// Log into console that server started
	a.GetLogger().Infof("Web server started at %s", a.server.Addr)
	// Block until program receive exit command or wait for OS interrupt
	<-a.shutdown()
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(a.config.Timeout.Read))
	// Call cancel immediately after func out of scope
	defer cancel()
	// Shut down the server
	err := a.server.Shutdown(ctx)
	if err != nil {
		a.FatalError(err)
	}
	// Log shutdown into console
	a.GetLogger().Warn("Server shutting down")
	// os.Exit(0)
	return
}

// NewApplication Init Web Application
func NewApplication(config Config, app gocli.Application, connState func(net.Conn, http.ConnState)) *Application {
	return &Application{
		config:      config,
		Application: app,
		exit:        make(chan bool),
		server: &http.Server{
			Addr:         config.Host + ":" + strconv.Itoa(config.Port),
			ReadTimeout:  time.Second * time.Duration(config.Timeout.Read),
			WriteTimeout: time.Second * time.Duration(config.Timeout.Write),
			IdleTimeout:  time.Second * time.Duration(config.Timeout.Idle),
			// We can handle our connections.
			// It is useful for web sockets or SSE or distributed transaction
			ConnState: connState,
		},
	}
}
