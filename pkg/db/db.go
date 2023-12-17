package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/ionutinit/riddles-api/models"
	"github.com/ionutinit/riddles-api/pkg/config"
	"github.com/ionutinit/riddles-api/pkg/logger"
)

var db *sql.DB

func InitDB() {
	cfg := config.AppConfig

	configFile, err := os.Open("config.json")
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"file":  "config.json",
		}).Fatal("Failed to open config file")
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
			"file":  "config.json",
		}).Fatal("Failed to decode config file")
	}
	defer configFile.Close()

	psqlCreds := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Dbname,
		cfg.Database.Sslmode,
	)

	db, err = sql.Open("postgres", psqlCreds)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error connecting to the database with current credentials")
	}

	err = db.Ping()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error pinging the database")
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxConnLifetime))

	logger.Log.Info("Successfully connected to the database")
}

func GetDB() *sql.DB {
	return db
}

func InsertNewRiddle(riddle models.Riddle) (int, error) {
	query := "INSERT INTO riddles (riddle, solution, synonyms, username, user_email) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id int
	err := db.QueryRow(query, riddle.Riddle, riddle.Solution, riddle.Synonyms, riddle.Username, riddle.UserEmail).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteRiddle(id int) (int64, error) {
	query := "DELETE FROM riddles WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func UpdateRiddle(id int, riddle models.Riddle) error {
	query := "UPDATE riddles SET "
	args := []interface{}{}
	argID := 1

	if riddle.Riddle != "" {
		query += fmt.Sprintf("riddle = $%d, ", argID)
		args = append(args, riddle.Riddle)
		argID++
	}

	if riddle.Solution != "" {
		query += fmt.Sprintf("solution = $%d, ", argID)
		args = append(args, riddle.Solution)
		argID++
	}

	if riddle.Synonyms != nil {
		query += fmt.Sprintf("synonyms = $%d, ", argID)
		args = append(args, riddle.Synonyms)
		argID++
	}

	if riddle.Username.Valid {
		query += fmt.Sprintf("username = $%d, ", argID)
		args = append(args, riddle.Username)
		argID++
	}

	if riddle.UserEmail.Valid {
		query += fmt.Sprintf("user_email = $%d, ", argID)
		args = append(args, riddle.UserEmail)
		argID++
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)

	// log.Printf("Executing query: %s with args: %v\n", query, args)

	_, err := db.Exec(query, args...)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"id":    id,
			"error": err,
			"query": query,
			"args":  args,
		}).Error("Error executing patch query")
		return err
	}
	return nil
}
