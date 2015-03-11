package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/carbocation/handlers"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func main() {
	app := App{}
	app.Version = "0.1-alpha"
	app.Render = render.New()
	app.Config = LoadConfiguration()

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
	http.ListenAndServe(":3001", middle)
}

//GorillaLog Logs to handler
func GorillaLog(out io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.CombinedLoggingHandler(out, next).ServeHTTP(w, r)
		})
	}
}
