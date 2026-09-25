package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func main() {
    log.Println("Job collector running...")
    for {
        evt := map[string]any{
            "id":      "job-" + time.Now().Format("20060102150405"),
            "kind":    "job",
            "payload": map[string]any{"job_id": "job-1", "status": "running"},
            "ts":      time.Now(),
            "source":  "collector-jobs",
        }
        b, _ := json.Marshal(evt)
        _, err := http.Post("http://localhost:9000/events", "application/json", bytes.NewBuffer(b))
        if err != nil {
            log.Println("collector-jobs error:", err)
        }
        time.Sleep(25 * time.Second)
    }
}
