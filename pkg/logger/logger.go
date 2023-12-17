package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		Log.Fatal("Failed to get home directory:", err)
	}

	logFilePath := filepath.Join(homeDir, "logs", "riddles.log")

	if _, err := os.Stat(filepath.Dir(logFilePath)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm)
		if err != nil {
			Log.Fatal("Failed to create log directory")
		}
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.Fatal("Failed to open log file:", err)
	}

	Log.SetOutput(file)
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}
