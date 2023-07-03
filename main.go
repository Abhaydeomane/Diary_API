package main

import (
	// "fmt"
	"encoding/json" //for json conversion
	"log"
	"math/rand" //for generating random value
	"net/http"  //for http server
	"strconv"   //for string - int coversion
	"time"
)

type Log struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"createAt"`
}

type DiaryEntry struct {
	ID   string `json:"id"`
	Date string `json:"date"`
	Logs []Log  `json:"log"`
}

type User struct {
	ID           string       `json:"id"`
	SecretCode   string       `json:"secretCode"`
	Name         string       `json:"name"`
	EmailAddress string       `json:"emailAddress"`
	DateOfBirth  string       `json:"dateOfBirth"`
	DiaryEntries []DiaryEntry `json:"diaryEntries"`
}

var users = make(map[string]*User) // Map to store User data

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct { //struct for storing request
		SecretCode string `json:"secretCode"`
	}
	// var req string

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[req.SecretCode] //searching in server memory
	if !ok {
		http.Error(w, " user not found", http.StatusNotFound)
		return
	}

	// here resp is response that we will send
	resp, err := json.Marshal(user) //Marshal returns the JSON encoding of user.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		DateOfBirth  string `json:"dateOfBirth"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate unique ID and secret code for the new user
	id := strconv.Itoa(rand.Intn(1000000)) // Example: generating ID as random integer between 0 and 999999
	secretCode := strconv.Itoa(rand.Intn(1000000))

	user := &User{
		ID:           id,
		SecretCode:   secretCode,
		Name:         req.Name,
		EmailAddress: req.EmailAddress,
		DateOfBirth:  req.DateOfBirth,
		DiaryEntries: []DiaryEntry{},
	}

	users[secretCode] = user //insert or adding user in server memory

	resp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func showDiaryOfMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct { //struct for storing request
		SecretCode string `json:"secretCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[req.SecretCode]
	if !ok {
		http.Error(w, "user not  found", http.StatusNotFound)
		return
	}

	// Get current month and year
	now := time.Now()
	month := now.Month()
	year := now.Year()

	// Filter diary entries for current month
	var filteredDiaryEntries []DiaryEntry
	for _, diaryEntry := range user.DiaryEntries {
		entryMonth, err := time.Parse("2006-01-02", diaryEntry.Date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if entryMonth.Month() == month && entryMonth.Year() == year {
			filteredDiaryEntries = append(filteredDiaryEntries, diaryEntry)
		}
	}

	resp, err := json.Marshal(filteredDiaryEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) 
	w.Write(resp)
}

// addEntry handles the "/addEntry" route
func addEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SecretCode string `json:"secretCode"`
		Log        Log    `json:"log"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[req.SecretCode]
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	now := time.Now().Format("2006-01-02")
	req.Log.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	req.Log.ID = strconv.FormatInt(time.Now().UnixNano(), 10)

	// Add log to current day's diary entry
	var found bool
	for i, diaryEntry := range user.DiaryEntries {
		if diaryEntry.Date == now {
			// Generate a unique ID based on current timestamp
			id := strconv.FormatInt(time.Now().UnixNano(), 10)
			req.Log.ID = id
			user.DiaryEntries[i].Logs = append(diaryEntry.Logs, req.Log)
			found = true
			break
		}
	}

	if !found {
		// If no diary entry exists for current day, create a new one
		// with a unique ID based on current timestamp
		id := strconv.FormatInt(time.Now().UnixNano(), 10)
		newDiaryEntry := DiaryEntry{
			ID:   id,
			Date: now,
			Logs: []Log{req.Log},
		}
		user.DiaryEntries = append(user.DiaryEntries, newDiaryEntry)
	}

	w.WriteHeader(http.StatusOK)
}

// }

func updateEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct { //for storing request
		SecretCode string `json:"secretCode"`
		Log        Log    `json:"log"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[req.SecretCode] //searching in server
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update log in current day's diary entry
	now := time.Now().Format("2006-01-02")
	var found bool
	for i, diaryEntry := range user.DiaryEntries {
		if diaryEntry.Date == now {
			for j, log := range diaryEntry.Logs {
				if log.ID == req.Log.ID {
					user.DiaryEntries[i].Logs[j] = req.Log
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		http.Error(w, "log not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct { // to store the request data
		SecretCode string `json:"secretCode"`
		ID         string `json:"id"` //here it is log id
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[req.SecretCode] //searching for the user
	if !ok {
		http.Error(w, "user not found in  the database", http.StatusNotFound)
		return
	}

	// Delete log from current day's diary entry
	now := time.Now().Format("2006-01-02")
	var found bool
	for i, diaryEntry := range user.DiaryEntries {
		if diaryEntry.Date == now {
			for j, log := range diaryEntry.Logs {
				if log.ID == req.ID {
					user.DiaryEntries[i].Logs = append(user.DiaryEntries[i].Logs[:j], user.DiaryEntries[i].Logs[j+1:]...)
					found = true
					break
				}
			}
			break
		}
	}

	if !found {
		http.Error(w, "log not found in server memoty.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func showEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SecretCode string `json:"secretCode"`
		Date       string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := users[req.SecretCode]
	if !ok {
		http.Error(w, "user not found ", http.StatusNotFound)
		return
	}

	var logs []Log
	for _, diaryEntry := range user.DiaryEntries {
		if diaryEntry.Date == req.Date {
			logs = diaryEntry.Logs
			break
		}
	}

	response := struct {
		Logs []Log `json:"logs"`
	}{
		Logs: logs,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {

	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/showDiaryOfMonth", showDiaryOfMonth)
	http.HandleFunc("/addEntry", addEntry)
	http.HandleFunc("/updateEntry", updateEntry)
	http.HandleFunc("/deleteEntry", deleteEntry)
	http.HandleFunc("/showEntry", showEntry)

	log.Fatal(http.ListenAndServe("localhost:3000", nil))

}
