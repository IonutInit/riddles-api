package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ionutinit/riddles-api/pkg/config"
	"github.com/ionutinit/riddles-api/pkg/db"
	"github.com/ionutinit/riddles-api/pkg/handlers"
	"github.com/ionutinit/riddles-api/pkg/logger"
	"github.com/ionutinit/riddles-api/pkg/middleware"
)

func rootRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetAllPublishedRiddles(w, r)
	case "POST":
		handlers.PostRiddleHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func singleRiddleHandler(allowedIPs []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetRiddleByIdHandler(w, r)
		case "DELETE":
			middleware.IPWhitelistMiddleware(http.HandlerFunc(handlers.DeleteRiddleHandler), allowedIPs).ServeHTTP(w, r)
		case "PATCH":
			middleware.IPWhitelistMiddleware(http.HandlerFunc(handlers.PatchRiddleHandler), allowedIPs).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func main() {

	config.LoadConfig("config.json")

	db.InitDB()
	defer db.GetDB().Close()

	allowedIPs := config.AppConfig.AllowedIPs

	// GET all and POST
	http.HandleFunc("/api/riddles", rootRouteHandler)

	// GET random riddle
	http.HandleFunc("/api/riddles/random", handlers.RandomRiddleHandler)

	// GET, DELETE, PATCH single riddle by id
	// DELETE and PATCH methods are IP-protected
	http.HandleFunc("/api/riddles/", singleRiddleHandler(allowedIPs))

	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/api", handlers.ApiPageHandler)

	server := &http.Server{
		Addr:    ":" + config.AppConfig.ServerPort,
		Handler: nil, //a nil handler defaults to http.DefaultServeMux
	}

	// starting the server in a go routine
	go startServer(server)

	// channel listening for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// blocks until a signal is received
	<-quit
	logger.Log.Info("Shutting down server...")

	// creates a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// doesn't block if no connection, but will otherwise wait until deadline
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Server forced to shutdown")
	}

	logger.Log.Info("Server exiting")
}

func startServer(server *http.Server) {
	logger.Log.WithFields(logrus.Fields{
		"port": config.AppConfig.ServerPort,
	}).Info("Starting server")
	log.Printf("Starting server on :%s\n", config.AppConfig.ServerPort)

	if err := server.ListenAndServe(); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Server start failed")
	}
}
