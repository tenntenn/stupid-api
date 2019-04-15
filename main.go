package main

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay())
		switch random([]int{90, 5, 5}) {
		case 0: // success
			writeJSON(w, map[string]string{
				"degree": strconv.Itoa(rand.Intn(40) - 10),
			})
		case 1: // error1
			writeJSON(w, map[string]string{
				"error": "internal server error",
			})
		case 2: // error2
			const status = http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort("", port)
	http.ListenAndServe(addr, nil)
}

func random(dist []int) int {
	if len(dist) == 0 {
		return -1
	}

	var sum int
	for _, d := range dist {
		sum += d
	}

	r := rand.Intn(sum)
	for i, d := range dist {
		r -= d
		if r <= 0 {
			return i
		}
	}

	return dist[len(dist)-1]
}

func delay() time.Duration {
	return []time.Duration{
		0,
		100 * time.Millisecond,
		1 * time.Second,
		10 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}[random([]int{80, 3, 3, 3, 3})]
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		const status = http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
	}
}
