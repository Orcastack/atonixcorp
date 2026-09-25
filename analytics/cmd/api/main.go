package main

import (
    "encoding/json"
    "log"
    "net/http"

    "atonixcorp/analytics/internal/bus"
    "atonixcorp/analytics/internal/config"
    "atonixcorp/analytics/internal/events"
    "atonixcorp/analytics/internal/metrics"
    "atonixcorp/analytics/internal/storage"
)

func main() {
    cfg := config.Load()
    db := storage.Connect(cfg.DBURL)
    storage.InitSchema(db)

    eventBus := bus.New()

    eventBus.Subscribe(func(evt events.Event) {
        metrics.EventsIngested.Inc()
        _, err := db.NamedExec(
            `INSERT INTO events (id, kind, payload, ts, source)
             VALUES (:id, :kind, :payload, :ts, :source)
             ON CONFLICT (id) DO NOTHING`,
            evt,
        )
        if err != nil {
            log.Println("DB write error:", err)
        }
    })

    mux := http.NewServeMux()

    mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
        metrics.Requests.Inc()
        var evt events.Event
        json.NewDecoder(r.Body).Decode(&evt)
        eventBus.Publish(evt)
        w.Write([]byte(`{"status":"ok"}`))
    })

    mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        out := map[string]any{
            "requests_total":  metrics.Requests.Value(),
            "events_ingested": metrics.EventsIngested.Value(),
        }
        json.NewEncoder(w).Encode(out)
    })

    mux.HandleFunc("/analytics/usage", func(w http.ResponseWriter, r *http.Request) {
        metrics.Requests.Inc()
        rows := []map[string]any{}
        db.Select(&rows, `SELECT kind, COUNT(*) AS count FROM events GROUP BY kind`)
        json.NewEncoder(w).Encode(rows)
    })

    log.Println("AtonixCorp Analytics API running on :9000")
    http.ListenAndServe(":9000", mux)
}
