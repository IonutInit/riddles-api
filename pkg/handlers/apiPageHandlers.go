package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/ionutinit/riddles-api/pkg/logger"
	"github.com/sirupsen/logrus"
)

func ApiPageHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"handler": "ApiPageHandler",
		}).Error("Error parsing HTML template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"handler": "ApiPageHandler",
		}).Error("Error executing HTML template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}