package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"

	"github.com/ionutinit/riddles-api/models"
	"github.com/ionutinit/riddles-api/pkg/db"
	"github.com/ionutinit/riddles-api/pkg/logger"
	"github.com/ionutinit/riddles-api/pkg/config"
)

type Response struct {
    Riddle models.RiddleBase `json:"riddle"`
    ImageURL string `json:"image.url"`
}

type ImageStyle struct {
    Style string `json:"style"`
}

func GenerateImageHandler(w http.ResponseWriter, r *http.Request) {
    logger.Log.Info("Executing GenerateImageHandler")

    path := r.URL.Path
    parts := strings.Split(path, "/")
    if len(parts) != 5 {
        logger.Log.WithFields(logrus.Fields{
            "path": path,
            "handler": "GenerateImageHandler",
        }).Error("Invalid request path")
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(parts[4])
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "extractedId": id,
            "error": err,
            "handler": "GenerateImageHandler",
        })
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    database := db.GetDB()
    row := database.QueryRow("SELECT id, riddle, solution, synonyms FROM riddles WHERE id = $1", id)

    var rdl models.RiddleBase
    if err := row.Scan(&rdl.ID, &rdl.Riddle, &rdl.Solution, &rdl.Synonyms); err != nil{
        logger.Log.WithFields(logrus.Fields{
            "id": id,
            "error": err,
            "handler": "GenerateImageHandler",
        }).Error("Error scanning riddles table rows")
        http.Error(w, "Riddle not found", http.StatusNotFound)
        return
    }

    var imageStyle ImageStyle

    if r.Body != nil {
        defer r.Body.Close()
        err := json.NewDecoder(r.Body).Decode(&imageStyle)
        if err != nil && err != io.EOF { //refers to no body
            logger.Log.WithFields(logrus.Fields{
                "error": err,
                "handler": "GenerateImageHandler",
            }).Error("Error parsing request body")
            http.Error(w, "Error parsing request body", http.StatusBadRequest)
            return
        }
    }

    client := openai.NewClient(config.AppConfig.OpenAiToken)
    ctx := context.Background()

    prompt := rdl.Riddle
    if imageStyle.Style != "" {
        prompt += imageStyle.Style + "style"
    }

    reqUrl := openai.ImageRequest{
        Prompt: prompt,
        Size: openai.CreateImageSize512x512,
        ResponseFormat: openai.CreateImageResponseFormatURL,
        N: 1,
    }

    respUrl, err := client.CreateImage(ctx, reqUrl)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "id": id,
            "error": err,
            "handler": "GenerateImageHandler",
        }).Error("Error generating image")
        http.Error(w, fmt.Sprintf("Error generating image: %v", err), http.StatusInternalServerError)
        return
    }
    imageUrl := respUrl.Data[0].URL


    response := Response {
        Riddle:  rdl,
        ImageURL: imageUrl,
    }

    logger.Log.WithFields(logrus.Fields{
        "riddleId": id,
        "handler": "GenerateImage Handler",
    }).Info("Succesfully sent riddle and image URL to client")

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    storeImage(imageUrl, rdl.ID)

    logger.Log.WithFields(logrus.Fields{
        "riddleId": id,
        "handler": "GenerateImageHandler",
    }).Info("Succesfully store image")
}


func storeImage(imageUrl string, riddleId int) {
    imgResp, err := http.Get(imageUrl)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error fetching image")
        return
    }
    defer imgResp.Body.Close()

    imgBytes, err := io.ReadAll(imgResp.Body)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error reading image bytes")
        return
    }

    base64Image := base64.StdEncoding.EncodeToString(imgBytes)

    insertImageQuery := "INSERT INTO images (riddleId, image) VALUES ($1, $2)"
    _, err = db.GetDB().Exec(insertImageQuery, riddleId, base64Image)
    if err != nil {
        logger.Log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error storing image")
    } 
}
