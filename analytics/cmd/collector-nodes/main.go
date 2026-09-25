package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

func main() {
    log.Println("Node collector running...")
    for {
        evt := map[string]any{
            "id":      "node-" + time.Now().Format("20060102150405"),
            "kind":    "node",
            "payload": map[string]any{"node_id": "node-1", "status": "online"},
            "ts":      time.Now(),
            "source":  "collector-nodes",
        }
        b, _ := json.Marshal(evt)
        _, err := http.Post("http://localhost:9000/events", "application/json", bytes.NewBuffer(b))
        if err != nil {
            log.Println("collector-nodes error:", err)
        }
        time.Sleep(15 * time.Second)
    }
}
