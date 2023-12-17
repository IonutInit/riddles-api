package handlers

import (
	"net/http"

	"github.com/ionutinit/riddles-api/pkg/config"
)

func constructURL(req *http.Request, path string) string {
	baseURL := config.AppConfig.BaseURL
	if baseURL == "" {
		baseURL = "http://" + req.Host
	}
	return baseURL + path
}
