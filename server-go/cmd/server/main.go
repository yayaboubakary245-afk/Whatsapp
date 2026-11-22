package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/groupe12-supptic/whatsapp-campus/internal/call"
	"github.com/groupe12-supptic/whatsapp-campus/internal/config"
	"github.com/groupe12-supptic/whatsapp-campus/internal/database"
	httpapi "github.com/groupe12-supptic/whatsapp-campus/internal/http"
	"github.com/groupe12-supptic/whatsapp-campus/internal/http/handlers"
	"github.com/groupe12-supptic/whatsapp-campus/internal/network/tcp"
	"github.com/groupe12-supptic/whatsapp-campus/internal/network/udp"
	"github.com/groupe12-supptic/whatsapp-campus/internal/websocket"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[DB] Connexion PostgreSQL échouée: %v", err)
	}
	defer pool.Close()

	if err := database.RunMigrations(ctx, pool, "../database/migrations"); err != nil {
		log.Printf("[DB] Avertissement migrations: %v", err)
	}

	// 1. Sockets TCP (port 5000) & UDP (port 5001) 100 % Golang
	tcpRegistry := tcp.NewClientRegistry()
	tcpServer := tcp.NewServer(cfg.TCPPort, tcpRegistry)
	if err := tcpServer.Start(); err != nil {
		log.Fatalf("[TCP] Impossible de démarrer sur :%s: %v", cfg.TCPPort, err)
	}
	defer tcpServer.Stop()

	udpServer := udp.NewServer(cfg.UDPPort)
	if err := udpServer.Start(); err != nil {
		log.Fatalf("[UDP] Impossible de démarrer sur :%s: %v", cfg.UDPPort, err)
	}
	defer udpServer.Stop()

	// 2. Hub WebSocket & Signalisation WebRTC (gorilla/websocket)
	hub := websocket.NewHub()
	signaling := call.NewSignalingService(pool, hub)

	// 3. API REST HTTP (gorilla/mux)
	h := &handlers.Handler{
		Cfg:       cfg,
		DB:        pool,
		Hub:       hub,
		TCPServer: tcpServer,
		UDPServer: udpServer,
		Signaling: signaling,
	}
	router := httpapi.NewRouter(h)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("[HTTP] API REST & WebSocket à l'écoute sur :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HTTP] Erreur serveur: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("[SERVER] Arrêt propre effectué.")
}
