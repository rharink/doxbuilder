package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

//App represents current context
type App struct {
	Config  Configuration
	Version string
	Render  *render.Render
}

//VersionAction Prints the current version
func (app *App) VersionAction(res http.ResponseWriter, req *http.Request) {
	m := make(map[string]string)
	m["version"] = app.Version

	app.Render.JSON(res, 200, m)
}

//ConvertAction Converts any docx to given format
func (app *App) ConvertAction(res http.ResponseWriter, req *http.Request) {
	// Handle incoming file (upload)
	file, header, err := req.FormFile("file")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	tempName := fmt.Sprintf("%d-%s", time.Now().Nanosecond(), header.Filename)

	// Save contents to tempfile
	path := fmt.Sprintf("%s/%s", app.Config.OutputDir, tempName)
	tempOut, err := os.Create(path)
	defer tempOut.Close()
	if err != nil {
		panic(err)
	}

	// Write contents to file
	_, err = io.Copy(tempOut, file)
	if err != nil {
		panic(err)
	}

	vars := mux.Vars(req)
	format := vars["format"]
	out, err := exec.Command("soffice", "--headless", "--convert-to", format, "--outdir", app.Config.OutputDir, tempName).Output()
	if err != nil {
		panic(err)
	}
	app.Render.JSON(res, 200, out)
}

// Run a Command
func (app *App) runCommand(cmd string) {
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	log.Printf("%s", out)
}
