package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var globalCounter int64 // compteur de requêtes

func main() {
	targetURL := getEnv("TARGET_URL", "")
	if targetURL == "" {
		log.Fatal("❌ TARGET_URL must be set")
	}

	thresholdMS := getEnvInt("THRESHOLD_MS", 100)
	workerCount := getEnvInt("WORKER_COUNT", 1000)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	log.Printf("🚀 Starting attack on %s with %d workers (threshold: %dms)", targetURL, workerCount, thresholdMS)

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for {
				start := time.Now()
				resp, err := client.Get(targetURL)
				duration := time.Since(start)

				reqNum := atomic.AddInt64(&globalCounter, 1) // incrément sécurisé

				if err != nil {
					log.Printf("❌ Worker %d | #%d: Request failed: %v", id, reqNum, err)
					continue
				}
				resp.Body.Close()

				if duration.Milliseconds() > int64(thresholdMS) {
					log.Printf("⚠️ Worker %d | #%d: Slow response: %dms (status: %d)", id, reqNum, duration.Milliseconds(), resp.StatusCode)
				}
			}
		}(i)
	}

	select {} // Block forever
}

// Helpers
func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if valStr, ok := os.LookupEnv(key); ok {
		val, err := strconv.Atoi(valStr)
		if err == nil {
			return val
		}
	}
	return fallback
}
