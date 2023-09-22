package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/writeFile", writeFile)
	http.HandleFunc("/readFile", readFile)

	http.ListenAndServe(":8080", nil)
}

func writeFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.FormValue("filePath")
	fileContent := r.FormValue("fileContent")

	err := ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to write to file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File written successfully"))
}

func readFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.FormValue("filePath")
	fmt.Printf("path: %s", filePath)
	data, err := ioutil.ReadFile(filePath)
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
