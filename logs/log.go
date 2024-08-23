package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

type Log struct {
	mu      sync.RWMutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) AddRegister(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var record Record
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Printf("Hay un error")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Hay un error")
		return
	}

	l.records = append(l.records, record)
	fmt.Printf("correcto")

	w.WriteHeader(http.StatusCreated)

}

func (l *Log) ShowRegister(w http.ResponseWriter, r *http.Request) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Printf("Hay un error")
		return
	}

	var offsetOnly struct {
		Offset uint64 `json:"offset"`
	}

	err := json.NewDecoder(r.Body).Decode(&offsetOnly)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Hay un error")
		return
	}

	if offsetOnly.Offset < uint64(len(l.records)) {
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

	file, ErrFile := os.Open(filename)
	if ErrFile != nil {
		return ErrFile
	}
	defer file.Close()

	FileJson, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err3 := json.Unmarshal(FileJson, &records)
	if err3 != nil {
		return err3
	}

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
