package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func main() {
    log.Println("Device collector running...")
    for {
        evt := map[string]any{
            "id":      "device-" + time.Now().Format("20060102150405"),
            "kind":    "device",
            "payload": map[string]any{"device_id": "dev-1", "status": "active"},
            "ts":      time.Now(),
            "source":  "collector-devices",
        }
        b, _ := json.Marshal(evt)
        _, err := http.Post("http://localhost:9000/events", "application/json", bytes.NewBuffer(b))
        if err != nil {
            log.Println("collector-devices error:", err)
        }
        time.Sleep(20 * time.Second)
    }
}
