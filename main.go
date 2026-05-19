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
	Windows        string         `json:"windows"`
	Wet            int            `json:"wet"`
	Temperature    int            `json:"temperature"`
	NumberStudents int            `json:"numberStudents"`
	PCs            map[int]string `json:"pcs"`
}

type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

var db *redis.Client

func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return client, nil
}

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

func getLastStatus(ctx context.Context) (*Status, error) {
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
	status, err := getLastStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"temperature": status.Temperature})
}

func getWet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getLastStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"wet": status.Wet})
}

func getWindows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getLastStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"windows": status.Windows})
}

func getPCStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status, err := getLastStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(status.PCs)
}

func main() {
	var err error
	cfg := Config{
		Addr:        "localhost:6379",
		Password:    "1234",
		User:        "",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	db, err = NewClient(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/logs", getLogs)
	http.HandleFunc("/temperature", getTemperature)
	http.HandleFunc("/wet", getWet)
	http.HandleFunc("/windows", getWindows)
	http.HandleFunc("/pcstatus", getPCStatus)

	fmt.Println("HTTP-server starts on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
