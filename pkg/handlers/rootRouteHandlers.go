package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ionutinit/riddles-api/models"
	"github.com/ionutinit/riddles-api/pkg/db"
	"github.com/ionutinit/riddles-api/pkg/logger"
)

func GetAllRiddlesHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Executing GetAllRiddlesHandler")
	database := db.GetDB()

	rows, err := database.Query("SELECT id, riddle, solution, synonyms FROM riddles WHERE published = TRUE")
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err,
			"handler": "GetAllRiddlesHandler",
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var riddles []models.RiddleBase
	for rows.Next() {
		var rdl models.RiddleBase
		if err := rows.Scan(&rdl.ID, &rdl.Riddle, &rdl.Solution, &rdl.Synonyms); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"error":   err,
				"handler": "GetAllRiddlesHandler",
			}).Error("Error scanning riddles table rows")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		riddles = append(riddles, rdl)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var riddlesResponse []models.RiddleResponse
	for _, rdlBase := range riddles {
		rdlResponse := models.RiddleResponse{
			RiddleBase: rdlBase,
			Links: []models.Link{
				{Rel: "self", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", rdlBase.ID))},
			},
		}
		riddlesResponse = append(riddlesResponse, rdlResponse)
	}

	logger.Log.WithFields(logrus.Fields{
		"count": len(riddles),
	}).Info("Successful query for GetAllRiddlesHandler")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(riddlesResponse)
}

func PostRiddleHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Executing PostRiddleHandler")

	var riddle models.Riddle
	err := json.NewDecoder(r.Body).Decode(&riddle)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err,
			"handler": "PostRiddleHandler",
		}).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if riddle.Riddle == "" || riddle.Solution == "" {
		logger.Log.WithFields(logrus.Fields{
			"handler": "PostRiddleHandler",
		}).Warn("Missing required fields in request")
		http.Error(w, "Missing required fields: riddle or solution", http.StatusBadRequest)
		return
	}

	id, err := db.InsertNewRiddle(riddle)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err,
			"handler": "PostRiddleHandler",
		}).Error("Error inserting new riddle in the database")
		http.Error(w, "Error inserting new riddle", http.StatusInternalServerError)
		return
	}

	riddleResponse := models.RiddleResponse{
		RiddleBase: models.RiddleBase{
			ID:       id,
			Riddle:   riddle.Riddle,
			Solution: riddle.Solution,
			Synonyms: riddle.Synonyms,
		},
		Links: []models.Link{
			{Rel: "view", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", id))},
			{Rel: "patch", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", id))},
			{Rel: "delete", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", id))},
		},
	}

	logger.Log.WithFields(logrus.Fields{
		"id":      id,
		"handler": "PostRiddleHandler",
	}).Info("Succesfully executed PostRiddleHandler")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(riddleResponse)
}
