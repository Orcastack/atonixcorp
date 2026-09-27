package backend

import (
	"log"
	"net/http"
	"time"

	"atonixcorp/cloud/devtools/portal/backend/auth"
	"atonixcorp/cloud/devtools/portal/backend/handlers"
)

func Start(addr string) error {
	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/version", handlers.Version)

	// Authenticated endpoints
	mux.Handle("/projects", auth.WithAuth(http.HandlerFunc(handlers.ListProjects)))
	mux.Handle("/compute", auth.WithAuth(http.HandlerFunc(handlers.ListCompute)))
	mux.Handle("/network", auth.WithAuth(http.HandlerFunc(handlers.ListNetwork)))
	mux.Handle("/storage", auth.WithAuth(http.HandlerFunc(handlers.ListStorage)))
	mux.Handle("/telemetry", auth.WithAuth(http.HandlerFunc(handlers.Telemetry)))

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("DevTools Portal backend listening on %s", addr)
	return srv.ListenAndServe()
}
