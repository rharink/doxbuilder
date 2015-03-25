package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

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
	file, _, err := req.FormFile("file")
	if err != nil {
		app.Render.JSON(res, 400, "No file found in request.")
		return
	}
	defer file.Close()
	tempName := fmt.Sprintf("%d-document", time.Now().Unix())

	// Save contents to tempfile
	path := fmt.Sprintf("%s/%s", app.Config.OutputDir, tempName)
	tempOut, err := os.Create(path)
	defer tempOut.Close()
	if err != nil {
		app.Render.JSON(res, 500, "Can't create temporary file path")
		return
	}
	_, err = io.Copy(tempOut, file)
	if err != nil {
		app.Render.JSON(res, 500, "Can't copy contents to temporary file path")
		return
	}
	defer os.Remove(path)

	format, err := app.getFormatFromRequest(req)
	if err != nil {
		app.Render.JSON(res, 400, "Unsupported format")
		return
	}
	cmd := exec.Command("soffice", "--headless", "--convert-to", format, "--outdir", app.Config.OutputDir, path)
	_, err = cmd.Output()
	if err != nil {
		app.Render.JSON(res, 500, fmt.Sprintf("Error during conversion: %s", err))
		return
	}

	// File has been converted now serve it back
	newPath := fmt.Sprintf("%s/%s.%s", app.Config.OutputDir, tempName, format)
	//delete file from server once it has been served
	defer os.Remove(newPath)

	res.Header().Set("Content-type", fmt.Sprintf("application/%s", format))
	http.ServeFile(res, req, newPath)
}

//Check the format from query string to match the ones in the config
func (app *App) getFormatFromRequest(req *http.Request) (string, error) {
	format := req.URL.Query().Get("format")
	format = strings.ToLower(format)
	if a := stringInSlice(format, app.Config.AllowedFormats); a == true {
		return format, nil
	} else {
		err := errors.New("Unsupported format")
		return "", err
	}
}

// Check if a string exists in array
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
