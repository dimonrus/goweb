package goweb

import (
	"bytes"
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/porterr"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

// Middleware Collection
type middlewareCollection struct {
	config         Config
	app            gocli.Application
	maxLogBodySize int64
}

// NewMiddlewareCollection New Middleware Collection Init Method
func NewMiddlewareCollection(config Config, app gocli.Application, maxLogBodySize int64) *middlewareCollection {
	return &middlewareCollection{
		config: config,
		app:    app,
		// if maxLogBodySize == -1 then do not need log body
		// if maxLogBodySize == 0 then log all body
		// if maxLogBodySize > 0 then read maxBodySize bytes from body for logging
		maxLogBodySize: maxLogBodySize,
	}
}

// Get requested url
func (m *middlewareCollection) getRequestedUrl(r *http.Request) string {
	return m.config.Url + ":" + strconv.Itoa(m.config.Port) + r.URL.Path + "?" + r.URL.RawQuery
}

// Logging request
func (m *middlewareCollection) loggingRequest(r *http.Request) porterr.IError {
	var logHeaders string
	for k, v := range r.Header {
		logHeaders += "-H '" + k + ": " + strings.Join(v, ",") + "' "
	}
	var body []byte
	var message string
	// log body with limits or without
	if r.ContentLength > 0 && m.maxLogBodySize > -1 {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return porterr.New(porterr.PortErrorBody, "IO error: "+err.Error()).HTTP(http.StatusBadRequest)
		}
		buf := bytes.NewReader(data)
		if m.maxLogBodySize > 0 && m.maxLogBodySize <= r.ContentLength {
			body = data[:m.maxLogBodySize-1]
		} else {
			body = data
		}
		r.Body = ioutil.NopCloser(buf)
	}
	message = gohelp.AnsiYellow + "REQUESTED: " + gohelp.AnsiBlue + "curl -X " + r.Method + " '" + m.getRequestedUrl(r) + "' " + logHeaders + " -d" + "'" + strings.Join(strings.Fields(string(body)), " ") + "'" + gohelp.AnsiReset
	m.app.GetLogger().Info(message)
	return nil
}

// LoggingMiddleware Logging middleware
func (m *middlewareCollection) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.Debug {
			e := m.loggingRequest(r)
			if e != nil {
				gorest.Send(w, gorest.NewErrorJsonResponse(e))
				return
			}
		}
		defer func() {
			if r := recover(); r != nil {
				var e porterr.IError
				switch r.(type) {
				case porterr.IError:
					e = r.(porterr.IError)
				case error:
					e = porterr.New(porterr.PortErrorSystem, "Critical issue: "+r.(error).Error())
				case string:
					e = porterr.New(porterr.PortErrorSystem, "Critical issue: "+r.(string))
				default:
					e = porterr.New(porterr.PortErrorSystem, "Critical issue: "+fmt.Sprintf("unsupported error: %T", r))
				}
				e = e.PushDetail("stack", "callback", string(debug.Stack()))
				m.app.GetLogger().Error(e.Error())
				gorest.Send(w, gorest.NewErrorJsonResponse(e))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// NotFoundHandler Not found handler
func (m *middlewareCollection) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := m.getRequestedUrl(r) + " not found on server"
	m.app.GetLogger().Error(gohelp.AnsiYellow + "REQUESTED: " + gohelp.AnsiRed + message + gohelp.AnsiReset)
	e := porterr.New(porterr.PortErrorHandler, message).HTTP(http.StatusNotFound)
	gorest.Send(w, gorest.NewErrorJsonResponse(e))
	return
}
