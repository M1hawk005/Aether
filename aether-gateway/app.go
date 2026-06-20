package main

import (
	"context"
	"log"
	"net/http"
	
	"aether-daemon"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Start the Aether Daemon in the background
	log.Println("Starting embedded Aether Daemon...")
	daemon.LoadNetworkConfig()
	
	err := daemon.InitLibp2p(ctx, 6000)
	if err != nil {
		log.Printf("Failed to init libp2p: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/publish", corsMiddleware(daemon.HandlePublish))
	mux.HandleFunc("/api/resolve/", corsMiddleware(daemon.HandleResolve))
	mux.HandleFunc("/api/fetch/", corsMiddleware(daemon.HandleFetch))
	mux.HandleFunc("/api/availability/", corsMiddleware(daemon.HandleAvailability))
	mux.HandleFunc("/api/delete/", corsMiddleware(daemon.HandleDeleteLocal))
	
	// Start REST API on 5000
	go func() {
		http.ListenAndServe("127.0.0.1:5000", mux)
	}()
}
