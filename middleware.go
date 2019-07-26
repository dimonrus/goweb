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
	config Config
	app    gocli.Application
}

// New Middleware Collection Init Method
func NewMiddlewareCollection(config Config, app gocli.Application) *middlewareCollection {
	return &middlewareCollection{
		config: config,
		app:    app,
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
	switch r.Method {
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		fallthrough
	case http.MethodPost:
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return rest.NewRestError("IO error: "+err.Error(), http.StatusBadRequest)
		}
		err = r.Body.Close()
		if err != nil {
			return rest.NewRestError("IO error: "+err.Error(), http.StatusBadRequest)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		body = b
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
				message := fmt.Sprintf("%s", debug.Stack())
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
