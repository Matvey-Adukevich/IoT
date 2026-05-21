package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Status struct {
	Windows     string         `json:"windows"`
	Wet         int            `json:"wet"`
	Temperature int            `json:"temperature"`
	PCs         map[int]string `json:"pcs"`
}

type Config struct {
	Addr        string
	Password    string
	User        string
	DB          int
	MaxRetries  int
	DialTimeout time.Duration
	Timeout     time.Duration
}

var db *redis.Client

func generateTemperature() int { return rand.Intn(35-17+1) + 17 }
func generateWet() int         { return rand.Intn(100-20+1) + 20 }

func generateStatusPC() map[int]string {
	m := make(map[int]string)
	for i := 1; i < 10; i++ {
		if rand.Intn(2) == 0 {
			m[i] = "off"
		} else {
			m[i] = "on"
		}
	}
	return m
}

func getRoomStatus(ctx context.Context) (*Status, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	val, err := db.Get(ctx, "room_status").Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("no data found in redis")
	} else if err != nil {
		return nil, err
	}

	var status Status
	if err := json.Unmarshal([]byte(val), &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	temperature := generateTemperature()
	wet := generateWet()
	pcs := generateStatusPC()
	windows := "close"

	if temperature > 28 || wet > 65 {
		windows = "open"
	}

	temp := Status{
		Windows:     windows,
		Wet:         wet,
		Temperature: temperature,
		PCs:         pcs,
	}

	jsonData, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("server error")
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if db != nil {
		ctx := r.Context()
		if err := db.Set(ctx, "room_status", jsonData, 0).Err(); err != nil {
			fmt.Printf("failed to save status to redis: %v\n", err)
		}

		logEntry := fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), string(jsonData))
		if err := db.LPush(ctx, "room_logs", logEntry).Err(); err != nil {
			fmt.Printf("failed to push log to redis: %v\n", err)
		}

		db.LTrim(ctx, "room_logs", 0, 49)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if db == nil {
		http.Error(w, "database not initialized", http.StatusInternalServerError)
		return
	}

	logs, err := db.LRange(r.Context(), "room_logs", 0, 19).Result()
	if err != nil {
		http.Error(w, "failed to fetch logs from redis", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(logs)
}

func getTemperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getRoomStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Temperature int `json:"temperature"`
	}{
		Temperature: status.Temperature,
	})
}

func getWet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getRoomStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Wet int `json:"wet"`
	}{
		Wet: status.Wet,
	})
}

func getWindows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getRoomStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Windows string `json:"windows"`
	}{
		Windows: status.Windows,
	})
}

func getPCStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getRoomStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(status.PCs)
}

func main() {
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.DefaultServeMux.ServeHTTP(w, r)
	})
	db = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "1234",
		DB:           0,
		MaxRetries:   5,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return
	}
	defer db.Close()

	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/logs", getLogs)
	http.HandleFunc("/temperature", getTemperature)
	http.HandleFunc("/wet", getWet)
	http.HandleFunc("/windows", getWindows)
	http.HandleFunc("/pcstatus", getPCStatus)

	fmt.Println("HTTP-server starts on :8080...")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		panic(err)
	}
}
