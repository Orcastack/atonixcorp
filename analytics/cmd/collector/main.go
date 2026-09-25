package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func main() {
    log.Println("System collector running...")
    for {
        evt := map[string]any{
            "id":      time.Now().Format("20060102150405"),
            "kind":    "system",
            "payload": map[string]any{"heartbeat": "ok"},
            "ts":      time.Now(),
            "source":  "collector-system",
        }
        b, _ := json.Marshal(evt)
        _, err := http.Post("http://localhost:9000/events", "application/json", bytes.NewBuffer(b))
        if err != nil {
            log.Println("collector error:", err)
        }
        time.Sleep(10 * time.Second)
    }
}
