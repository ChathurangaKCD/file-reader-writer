package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var serverId string

func main() {
	serverId = generateUUID()
	http.HandleFunc("/writeFile", writeFile)
	http.HandleFunc("/readFile", readFile)
	http.HandleFunc("/listFiles", listFiles)
	http.HandleFunc("/deleteFile", deleteFile)

	http.ListenAndServe(":8080", nil)
}

func generateUUID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
func writeJSON(w http.ResponseWriter, msg string, requestId string, data interface{}) {
	res := map[string]interface{}{
		"message":   msg,
		"serverId":  serverId,
		"requestId": requestId,
		"data":      data,
	}
	responseData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create JSON response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(responseData))
}

func writeFile(w http.ResponseWriter, r *http.Request) {
	requestId := generateUUID()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.FormValue("filePath")
	fileContent := r.FormValue("fileContent")
	logrus.WithFields(logrus.Fields{
		"filePath":    filePath,
		"fileContent": fileContent,
		"requestId":   requestId,
		"serverId":    serverId,
	}).Info("Writing file")

	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to write to file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, "File written successfully", requestId, nil)
}

func readFile(w http.ResponseWriter, r *http.Request) {
	requestId := generateUUID()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.FormValue("filePath")
	logrus.WithFields(logrus.Fields{
		"filePath":  filePath,
		"requestId": requestId,
		"serverId":  serverId,
	}).Info("Reading file")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Unable to read file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	requestId := generateUUID()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dirPath := r.FormValue("dirPath")
	logrus.WithFields(logrus.Fields{
		"dirPath":   dirPath,
		"requestId": requestId,
		"serverId":  serverId,
	}).Info("Listing files")

	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read directory: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	// Return the file names as a JSON response
	responseData, err := json.Marshal(fileNames)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create JSON response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writeJSON(w, "Files listed successfully", requestId, responseData)
}

// New function to handle file deletion
func deleteFile(w http.ResponseWriter, r *http.Request) {
	requestId := generateUUID()
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.FormValue("filePath")
	logrus.WithFields(logrus.Fields{
		"filePath":  filePath,
		"requestId": requestId,
		"serverId":  serverId,
	}).Info("Deleting file")

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Unable to delete file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, "File deleted successfully", requestId, nil)
}
