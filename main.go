package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

const (
	saveDirectory = "./data"
)

func main() {
	err := godotenv.Load("app.env")
	if err != nil {
		fmt.Println("Error loading app.env")
	}

	// Save directory
	err = os.MkdirAll(saveDirectory, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Handler
	http.HandleFunc("/data", handleData)
	fmt.Println("Listening on port " + os.Getenv("API_PORT") + "...")
	http.ListenAndServe(":"+os.Getenv("API_PORT"), nil)
}

func handleData(w http.ResponseWriter, r *http.Request) {
	// Check Method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode body
	var sysInfo map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&sysInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get IP Client
	ip := r.RemoteAddr

	// filenames CSV and JSON
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s_%s", ip, date)
	// jsonFilename := fmt.Sprintf("%s_%s.json", ip, date)

	// Open CSV y JSON file
	csvFilePath := filepath.Join(saveDirectory, filename+".csv")
	file, err := os.Create(csvFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	jsonFilePath := filepath.Join(saveDirectory, filename+".json")
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	// Save info into JSON
	jsonEncoder := json.NewEncoder(jsonFile)
	err = jsonEncoder.Encode(sysInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Save info into CSV
	for key, value := range sysInfo {
		row := []string{key, fmt.Sprintf("%v", value)}
		err := writer.Write(row)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	fmt.Printf("Information saved into CSV & JSON files in: %s, %s\n", file.Name(), jsonFile.Name())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
