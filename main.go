package main

import (
	"encode/json"
	"math/rand"
	"net/http"
)

type Status struct {
	Windows     string `json:"windows"`
	Wet         string `json:"wet"`
	Temperature int    `json:"temperature"`
}

// type Temperature struct{
// 	Temperature int `json:"temperature"`

// }

func generateTemperature() (temp int) {
	temp = rand.Intn(35-17+1) + 17
	return temp
}

func generateWet() (wet int) {
	wet = rand.Intn(100-20+1) + 20
	return wet
}

func getTemperature(w http.ResponseWriter, r *http.Request) {
	// temp := rand.Intn(35);
	temperature := rand.Intn(29)
	windows := "close"
	wet := rand.Intn(67)
	if temp > 28 || wet > 65 {
		windows = "open"
	}
	temp := Status{
		Windows:     windows,
		Wet:         wet,
		Temperature: temperature,
	}
	json.NewEncoder(w).Encode(temp)
}

func main() {
	http.HandleFunc("/temperature", getTemperature)
	http.ListenAndServe(":8080", nil)
}
