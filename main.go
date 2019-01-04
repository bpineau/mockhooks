package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	version            = "0.3.0"
	mode               *string
	minDelay, maxDelay *int
	failurePct         *int
)

type ctxKey int

const (
	isFailing ctxKey = iota
	usernameID
)

type rObj struct {
	Kind     string   `json:"kind"`
	Metadata metadata `json:"metadata"`
}

type metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var name string
		var obj rObj

		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &obj); err != nil {
			name = "[buggy request object]"
		} else {
			name = fmt.Sprintf("[kind=%s name=%s/%s]",
				obj.Kind, obj.Metadata.Namespace, obj.Metadata.Name)
		}

		delay := randInRange(*minDelay, *maxDelay)
		fail := shouldFail(*failurePct)
		time.Sleep(time.Duration(delay) * time.Second)

		ctx := context.WithValue(r.Context(), isFailing, fail)
		next.ServeHTTP(w, r.WithContext(ctx))

		log.Printf("%s %s %-15s shouldFail=%-6t delay=%d %s",
			time.Now().Format(time.RFC3339), r.Method,
			r.URL.RequestURI(), fail, delay, name)
	})
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	failing := r.Context().Value(isFailing).(bool)

	if failing {
		http.Error(w, "500 - Synthetic error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func shouldFail(pct int) bool {
	return randInRange(0, 100) < pct
}

func randInRange(min, max int) int {
	if max <= 0 {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func main() {
	listen := flag.String("listen", ":8082", "listening address")
	minDelay = flag.Int("min-delay", 0, "minimum delay before responses (seconds)")
	maxDelay = flag.Int("max-delay", 0, "maximum delay before responses (seconds)")
	failurePct = flag.Int("failure-pct", 0, "percentage of failing responses")
	mode = flag.String("mode", "webhook", "mode: webhook, elasticsearch")

	flag.Parse()

	mux := http.NewServeMux()
	switch *mode {
	case "webhook":
		mux.HandleFunc("/", webhookHandler)
	case "elasticsearch":
		mux.HandleFunc("/_cat/health", esHealthHandler)
		mux.HandleFunc("/_cluster/settings", esSettingsHandler)
		mux.HandleFunc("/_flush/synced", esFlushHandler)
	default:
		log.Fatal("unknown mode", mode)
	}

	log.Println("starting mockhooks, version", version)

	contextedMux := AddContext(mux)
	log.Fatal(http.ListenAndServe(*listen, contextedMux))
}
