package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

// Record struct of the logs
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// Log handler
type Log struct {
	mu      sync.RWMutex
	records []Record
}

// NewLog =Initialize the log
func NewLog() *Log {
	return &Log{}
}

func (l *Log) AddRegister(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var record Record
	// check if the HTTP request method is a POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Printf("Hay un error")
		return
	}

	// Decode the JSON request body into the record struct
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Hay un error")
		return
	}

	// Add the new register to the log struct
	l.records = append(l.records, record)
	fmt.Printf("correcto")

	w.WriteHeader(http.StatusCreated)

}

func (l *Log) ShowRegister(w http.ResponseWriter, r *http.Request) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	// check if the HTTP request method is a GET
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Printf("Hay un error")
		return
	}

	//We only request the offset when we want to display the user, so we create a
	//structure to handle that value specifically.
	var offsetOnly struct {
		Offset uint64 `json:"offset"`
	}

	// Decode the JSON request body into the offsetOnly struct
	err := json.NewDecoder(r.Body).Decode(&offsetOnly)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Hay un error")
		return
	}

	//check if the offset exists
	if offsetOnly.Offset < uint64(len(l.records)) {
		//if it exists,take the value
		value := l.records[offsetOnly.Offset]
		fmt.Printf("correcto")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(value)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Hay un error")
	}

}

func (l *Log) UploadJSON(filename string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var records []Record
	//open the json file
	file, ErrFile := os.Open(filename)
	if ErrFile != nil {
		return ErrFile
	}
	defer file.Close()

	//read the entire content of the file
	FileJson, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// converting json data into a structured data format
	err3 := json.Unmarshal(FileJson, &records)
	if err3 != nil {
		return err3
	}

	//save the data in the log struct
	l.records = records
	return nil
}

func main() {

	log := NewLog()
	err := log.UploadJSON("log.json")
	if err != nil {
		fmt.Print(err)
	}

	http.HandleFunc("/AddLog", log.AddRegister)
	http.HandleFunc("/ShowLog", log.ShowRegister)

	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		return
	}

}
