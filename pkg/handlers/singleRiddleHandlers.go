package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ionutinit/riddles-api/models"
	"github.com/ionutinit/riddles-api/pkg/db"
	"github.com/ionutinit/riddles-api/pkg/logger"
)

func GetRiddleByIdHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Executing GetRiddleBtIdHandler")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		logger.Log.WithFields(logrus.Fields{
			"path":    path,
			"handler": "GetRiddleByIdHandler",
		}).Error("Invalid request path")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"path":    path,
			"error":   err,
			"handler": "GetRiddleByIdHandler",
		}).Error("Invalid riddle ID in path")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	database := db.GetDB()

	row := database.QueryRow("SELECT id, riddle, solution, synonyms FROM riddles WHERE id = $1", id)

	var rdlBase models.RiddleBase
	if err := row.Scan(&rdlBase.ID, &rdlBase.Riddle, &rdlBase.Solution, &rdlBase.Synonyms); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"id":      id,
			"error":   err,
			"handler": "GetRiddleBtIdHandler",
		}).Error("Error scanning riddles table rows")
		http.Error(w, "Riddle not found", http.StatusNotFound)
		return
	}

	rdlResponse := models.RiddleResponse{
		RiddleBase: rdlBase,
		Links: []models.Link{
			{Rel: "update", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", rdlBase.ID))},
			{Rel: "delete", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", rdlBase.ID))},
		},
	}

	logger.Log.WithFields(logrus.Fields{
		"id":      id,
		"handler": "GetRiddleByIdHandler",
	}).Info("Successfully executed GetRiddleByIdHandler")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rdlResponse)
}

func RandomRiddleHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Executing RandomRiddleHandler")
	database := db.GetDB()

	row := database.QueryRow("SELECT id, riddle, solution, synonyms FROM riddles WHERE published = TRUE ORDER BY RANDOM() LIMIT 1")

	var rdlBase models.RiddleBase
	if err := row.Scan(&rdlBase.ID, &rdlBase.Riddle, &rdlBase.Solution, &rdlBase.Synonyms); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err,
			"handler": "RandomRiddleHandler",
		}).Error("Error scanning riddles table rows")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rdlResponse := models.RiddleResponse{
		RiddleBase: rdlBase,
		Links: []models.Link{
			{Rel: "update", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", rdlBase.ID))},
			{Rel: "delete", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", rdlBase.ID))},
		},
	}

	logger.Log.WithFields(logrus.Fields{
		"riddleID": rdlBase.ID,
	}).Info("Successful query for RandomRiddleHandler")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rdlResponse)
}

func DeleteRiddleHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Executing DeleteRiddleHandler")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		logger.Log.WithFields(logrus.Fields{
			"path":    path,
			"handler": "DeleteRiddleHandler",
		}).Error("Invalid request path")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"path":    path,
			"error":   err,
			"handler": "DeleteRiddleHandler",
		}).Error("Invalid riddle ID")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.DeleteRiddle(id)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"id":      id,
			"error":   err,
			"handler": "DeleteRiddleHandler",
		}).Error("Error deleting riddle from database")
		http.Error(w, "Error deleting riddle", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		logger.Log.WithFields(logrus.Fields{
			"id":      id,
			"handler": "DeleteRiddleHandler",
		}).Warn("ID not matching any riddle for deletion")
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "Riddle deleted successfully",
		"links": []models.Link{
			{Rel: "all-riddles", Href: constructURL(r, "/api/riddles")},
		},
	}

	logger.Log.WithFields(logrus.Fields{
		"id":      id,
		"handler": "DeleteRiddleHandler",
	}).Info("Successfully executed DeleteRiddleHandler")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func PatchRiddleHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Executing PatchRiddleHandler")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		logger.Log.WithFields(logrus.Fields{
			"path":    path,
			"handler": "PatchRiddleHandler",
		}).Error("Invalid request path")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"path":    path,
			"error":   err,
			"handler": "PatchRiddleHandler",
		}).Error("Invalid riddle ID")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// creating a map in order to check that only allowed fields exist in the request body
	// more efficient than creating a json directly and then searching through it, is not efficient and it involves further parsing
	var requestBodyMap map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBodyMap); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err,
			"handler": "PatchRiddleHandler",
		}).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allowedFields := map[string]bool{
		"riddle": true, "solution": true, "synonyms": true, "username": true, "user_email": true,
	}

	for field := range requestBodyMap {
		if !allowedFields[field] {
			logger.Log.WithFields(logrus.Fields{
				"invalidField": field,
				"handler":      "PatchRiddleHandler",
			}).Warn("Invalid field in request body")
			http.Error(w, fmt.Sprintf("Invalid field: %s", field), http.StatusBadRequest)
			return
		}
	}

	// marshalling the body into json, as it a the type of the map cannot be directly converted into byte
	// the decoder cannot be referred to again, as it is a read-once stream, which has been exhausted in the map
	requestBodyJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updatedRiddle models.Riddle
	if err := json.Unmarshal(requestBodyJSON, &updatedRiddle); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateRiddle(id, updatedRiddle); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"id":      id,
			"error":   err,
			"handler": "PatchRiddleHandler",
		}).Error("Error updating riddle")
		http.Error(w, "Error updating riddle", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Riddle updated successfully",
		"links": []models.Link{
			{Rel: "view", Href: constructURL(r, fmt.Sprintf("/api/riddles/%d", id))},
			{Rel: "all-riddles", Href: constructURL(r, "/api/riddles")},
		},
	}

	logger.Log.WithFields(logrus.Fields{
		"id":      id,
		"handler": "PatchRiddleHandler",
	}).Info("Successfully executed PatchRiddleHandler")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
