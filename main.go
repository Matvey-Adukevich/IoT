package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type Status struct {
	Windows        string         `json:"windows"`
	Wet            int            `json:"wet"`
	Temperature    int            `json:"temperature"`
	NumberStudents int            `json:"numberStudents"`
	PCs            map[int]string `json:"pcs"`
}

func generateTemperature() (temp int) {
	temp = rand.Intn(35-17+1) + 17
	return temp
}

func generateWet() (wet int) {
	wet = rand.Intn(100-20+1) + 20
	return wet
}

// func getNumberStudents() (number int) {
// 	number = rand.Intn(11)
// 	return number
// }

func generateStatusPC() map[int]string {
	m := make(map[int]string)
	// num := getNumberStudents()
	for i := 1; i < 10; i++ {
		flag := rand.Intn(2)
		if flag == 0 {
			m[i] = "off"
		} else {
			m[i] = "on"
		}
	}
	return m
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	temperature := generateTemperature()
	windows := "close"
	wet := generateWet()
	pcs := generateStatusPC()
	if temperature > 28 || wet > 65 {
		windows = "open"
	}
	temp := Status{
		Windows:     windows,
		Wet:         wet,
		Temperature: temperature,
		PCs:         pcs,
	}
	json.NewEncoder(w).Encode(temp)
}

// func getWet(w http.ResponseWriter, r *http.Request){

// }

func main() {
	http.HandleFunc("/status", getStatus)
	// http.HandleFunc("/temperature", getTemperature)
	// http.HandleFunc("/wet", getWet)
	// http.HandleFunc("/windows", getWindows)
	// http.HandleFunc("/pcstatus", getPCStatus)
	// http.HandleFunc("/logs", getLogs)
	http.ListenAndServe(":8080", nil)
}
