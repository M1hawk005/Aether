package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"aether-daemon"
)

var (
	PeerId   string
	HttpPort int
	WsPort   int
)

func init() {
	flag.StringVar(&PeerId, "peer", "peer-a", "Peer ID for the daemon")
	flag.IntVar(&HttpPort, "port", 5000, "HTTP REST API port")
}

func main() {
	flag.Parse()
	WsPort = HttpPort + 1000

	err := daemon.LoadNetworkConfig()
	if err != nil {
		log.Fatalf("Failed to load network config: %v", err)
	}
	defer daemon.CleanupNetworkConfig()

	// Setup Libp2p
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if err := daemon.InitLibp2p(ctx, WsPort); err != nil {
		log.Fatalf("Failed to initialize libp2p: %v", err)
	}

	// Setup HTTP API
	mux := http.NewServeMux()
	mux.HandleFunc("/api/publish", daemon.HandlePublish)
	mux.HandleFunc("/api/resolve/", daemon.HandleResolve)
	mux.HandleFunc("/api/fetch/", daemon.HandleFetch)
	mux.HandleFunc("/api/availability/", daemon.HandleAvailability)
	mux.HandleFunc("/api/delete/", daemon.HandleDeleteLocal)
	mux.HandleFunc("/api/revoke/", daemon.HandleRevokeGlobal)
	mux.HandleFunc("/api/claim-username", daemon.HandleRegistryClaim)
	mux.HandleFunc("/api/release-username", daemon.HandleRegistryRelease)
	mux.HandleFunc("/api/status", daemon.HandleStatus)
	mux.HandleFunc("/api/reputation", daemon.HandleReputation)
	mux.HandleFunc("/api/feed", daemon.HandleFeed)
	// V1 SDK backward compatibility
	mux.HandleFunc("/api/v1/status", daemon.HandleStatus)

	// BitTorrent Tracker Bridge
	mux.HandleFunc("/announce", daemon.HandleAnnounce)

	// Web UI (Delphi Tracker App)
	uiDir := filepath.Join("..", "..", "apps", "delphi")
	fs := http.FileServer(http.Dir(uiDir))
	mux.Handle("/ui/", http.StripPrefix("/ui/", fs))

	// Start HTTP Server
	serverAddr := fmt.Sprintf("127.0.0.1:%d", HttpPort)
	log.Printf("[Aether Daemon Go] Node '%s' listening on HTTP port %d and WS port %d", PeerId, HttpPort, WsPort)
	
	go func() {
		if err := http.ListenAndServe(serverAddr, mux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down daemon...")
}
