package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"flag"
	"fmt"

	"github.com/carbocation/handlers"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	configPath = flag.String("c", "configuration.yml", "Set the configuration file location")
	port       = flag.Int("p", 3000, "Set a port to run doxbuilder on, defaults to 3000")
)

func main() {
	flag.Parse()
	config, err := LoadConfiguration(*configPath)
	if err != nil {
		fmt.Println("Can't load configuration")
		os.Exit(2)
	}
	app := App{}
	app.Version = "0.1-alpha"
	app.Render = render.New()
	app.Config = config

	// Middleware
	middle := interpose.New()
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/version", app.VersionAction).Methods("GET")
	router.HandleFunc("/convert", app.ConvertAction).Methods("POST")

	// Logging
	multi := io.MultiWriter(os.Stdout)
	log.SetOutput(multi)
	middle.Use(GorillaLog(multi))

	// Add router
	middle.UseHandler(router)

	// Listen and Serve
	http.ListenAndServe(fmt.Sprintf(":%d", *port), middle)
}

//GorillaLog Logs to handler
func GorillaLog(out io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.CombinedLoggingHandler(out, next).ServeHTTP(w, r)
		})
	}
}
