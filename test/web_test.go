package test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gorest"
	"github.com/dimonrus/goweb"
	"github.com/dimonrus/porterr"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type WebConfig struct {
	Web       goweb.Config
	Arguments gocli.ArgumentMap
}

var config WebConfig

func TestDecomposeCommand(t *testing.T) {
	command := gocli.ParseCommand([]byte("web_unknown"))
	_, _, e := goweb.DecomposeCommand(command)
	if e == nil {
		t.Fatal("must be an error")
	}
	command = gocli.ParseCommand([]byte("web"))
	_, _, e = goweb.DecomposeCommand(command)
	if e == nil {
		t.Fatal("must be an error")
	}
	command = gocli.ParseCommand([]byte("web sfsdf"))
	_, _, e = goweb.DecomposeCommand(command)
	if e == nil {
		t.Fatal("must be an error")
	}
	command = gocli.ParseCommand(nil)
	_, _, e = goweb.DecomposeCommand(command)
	if e == nil {
		t.Fatal("must be an error")
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	gorest.Send(w, gorest.NewOkJsonResponse("Alive", nil, nil))
}

func PostBodyHandler(w http.ResponseWriter, r *http.Request) {
	message := struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}{
		Key:   "information",
		Value: "notification",
	}
	gorest.Send(w, gorest.NewOkJsonResponse("Alive", message, nil))
}

func GetErrHandler(w http.ResponseWriter, r *http.Request) {
	errType := r.URL.Query().Get("type")
	t, err := strconv.ParseInt(errType, 10, 64)
	if err != nil {
		panic(err)
	}
	switch t {
	case 500:
		panic("500 error")
	case 501:
		panic(errors.New("501 error"))
	case 502:
		e := porterr.New(porterr.PortErrorSystem, "502 porterr")
		panic(e)
	case 503:
		newErr := struct {
			Hello string
		}{
			Hello: "hello",
		}
		panic(newErr)
	}
	gorest.Send(w, gorest.NewOkJsonResponse("Error", nil, nil))
}

func getWebApplication() *goweb.Application {
	rootPath, err := filepath.Abs("")
	if err != nil {
		panic(err)
	}
	app := gocli.NewApplication("global", rootPath+"/config", &config)
	//app.ParseFlags(&config.Arguments)
	config.Web.Security.HTTP2Enable = false
	config.Web.Security.ServerCert = "cert.crt"
	config.Web.Security.ServerKey = "key.key"

	return goweb.NewApplication(config.Web, app, nil)
}

func TestWebServer(t *testing.T) {
	wa := getWebApplication()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthCheckHandler)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				wa.FailMessage(r.(error).Error())
			}
		}()
		func() {
			time.Sleep(time.Second)
			command := gocli.ParseCommand([]byte("web stop"))
			wa.WebCommander(command)

			// check error
			command = gocli.ParseCommand([]byte("web csdvsdv"))
			wa.WebCommander(command)
		}()
	}()
	wa.Listen(mux)
	fmt.Println("Graceful shutdown")
}

func TestLoggingMiddleware(t *testing.T) {
	wa := getWebApplication()
	mux := http.NewServeMux()
	middleware := goweb.NewMiddlewareCollection(config.Web, wa.Application, 1024)
	mux.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(HealthCheckHandler)))
	mux.Handle("/information", middleware.LoggingMiddleware(http.HandlerFunc(PostBodyHandler)))
	mux.Handle("/error", middleware.LoggingMiddleware(http.HandlerFunc(GetErrHandler)))
	mux.Handle("/", http.HandlerFunc(middleware.NotFoundHandler))

	go func() {
		defer func() {
			if r := recover(); r != nil {
				wa.FailMessage(r.(error).Error())
			}
		}()
		func() {
			resp1, err := http.Get("http://0.0.0.0:8080/health")
			if err != nil {
				wa.FatalError(err)
			}
			defer resp1.Body.Close()
			data, err := ioutil.ReadAll(resp1.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(string(data))

			body := []byte("{\"information\": true}")
			b := bytes.NewBuffer(body)
			resp2, err := http.Post("http://0.0.0.0:8080/information", "application/json", b)
			if err != nil {
				wa.FatalError(err)
			}
			data, err = ioutil.ReadAll(resp2.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer resp2.Body.Close()
			fmt.Println(string(data))

			for i := 0; i < 4; i++ {
				errType := 500 + i
				resp3, err := http.Post("http://0.0.0.0:8080/error?type="+strconv.Itoa(errType), "application/json", nil)
				if err != nil {
					wa.FatalError(err)
				}
				data, err = ioutil.ReadAll(resp3.Body)
				if err != nil {
					t.Fatal(err)
				}
				resp3.Body.Close()
				fmt.Println(string(data))
			}

			resp4, err := http.Get("http://0.0.0.0:8080/notknown")
			if err != nil {
				wa.FatalError(err)
			}
			defer resp4.Body.Close()
			data, err = ioutil.ReadAll(resp1.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(string(data))

			time.Sleep(time.Second)
			command := gocli.ParseCommand([]byte("web stop"))
			wa.WebCommander(command)
		}()
	}()
	wa.Listen(mux)
	fmt.Println("Graceful shutdown")
}

func BenchmarkDecomposeCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		command := gocli.ParseCommand([]byte("web stop"))
		_, _, e := goweb.DecomposeCommand(command)
		if e != nil {
			b.Fatal(e)
		}
	}
	// TODO reduce allocation
	b.ReportAllocs()
}
