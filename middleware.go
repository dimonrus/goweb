package goweb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/rest"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
)

// Middleware Collection
type middlewareCollection struct {
	config         Config
	app            gocli.Application
	maxLogBodySize int64
}

// New Middleware Collection Init Method
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
	return fmt.Sprintf("%s:%v%s", m.config.Url, m.config.Port, r.URL.Path+"?"+r.URL.RawQuery)
}

// Logging request
func (m *middlewareCollection) loggingRequest(r *http.Request) rest.IError {
	var logHeaders string
	for k, v := range r.Header {
		logHeaders += fmt.Sprintf("-H '%s: %s' ", k, strings.Join(v, ","))
	}
	var body []byte

	if r.ContentLength > 0 && m.maxLogBodySize > -1 {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return rest.NewRestError("IO error: "+err.Error(), http.StatusBadRequest)
		}
		buf := bytes.NewReader(data)
		if m.maxLogBodySize > 0 {
			body = data[:m.maxLogBodySize]
		} else {
			body = data
		}
		r.Body = ioutil.NopCloser(buf)
	}

	format := fmt.Sprintf("\x1b[33;1mREQUESTED: \x1b[34;1mcurl -X %s '%s' %s -d '%s'\x1b[0m", r.Method, m.getRequestedUrl(r), logHeaders, strings.Join(strings.Fields(string(body)), " "))
	m.app.GetLogger(gocli.LogLevelDebug).Info(format)

	return nil
}

// Logging middleware
func (m *middlewareCollection) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := m.loggingRequest(r)
		if e != nil {
			rest.ErrorResponse(w, e)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r.(type) {
				case error:
					err = r.(error)
				case string:
					err = errors.New(r.(string))
				case gocli.IError:
					err = errors.New(r.(gocli.IError).Error())
				default:
					err = errors.New("some unsupported error")
				}
				e := rest.NewRestError("​​​​Critical issue. Please send it to technical support: "+err.Error(), http.StatusInternalServerError)
				key := "stack"
				message := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
				e = e.AppendDetail(message, &key, nil)
				m.app.GetLogger(gocli.LogLevelDebug).Error(message)
				rest.ErrorResponse(w, e)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Not found handler
func (m *middlewareCollection) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s not found on server", m.getRequestedUrl(r))
	m.app.GetLogger(gocli.LogLevelDebug).Error("\x1b[33;1mREQUESTED: \x1b[31;1m", message, "\x1b[0m")
	rest.NotFoundResponse(w, message)
	return
}
