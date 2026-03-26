package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Appouse/HostPeek/internal/collector"
	"github.com/Appouse/HostPeek/internal/config"
)

// Server wraps the HTTP server and configuration.
type Server struct {
	cfg *config.Config
	mux *http.ServeMux
}

// New creates a new Server with all routes registered.
func New(cfg *config.Config) *Server {
	s := &Server{cfg: cfg, mux: http.NewServeMux()}
	s.routes()
	return s
}

// Handler returns the HTTP handler (with middleware).
func (s *Server) Handler() http.Handler {
	var h http.Handler = s.mux

	// API key authentication middleware
	if s.cfg.Auth.Enabled && s.cfg.Auth.APIKey != "" {
		h = s.authMiddleware(h)
	}

	return h
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler().ServeHTTP(w, r)
}

// routes registers all endpoint handlers.
func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /metrics", s.handleMetricsAll)
	s.mux.HandleFunc("GET /metrics/os", s.handleMetricsOS)
	s.mux.HandleFunc("GET /metrics/cpu", s.handleMetricsCPU)
	s.mux.HandleFunc("GET /metrics/memory", s.handleMetricsMemory)
	s.mux.HandleFunc("GET /metrics/disk", s.handleMetricsDisk)
	s.mux.HandleFunc("GET /metrics/network", s.handleMetricsNetwork)
	s.mux.HandleFunc("GET /metrics/uptime", s.handleMetricsUptime)
}

// authMiddleware checks the X-API-Key header.
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow health endpoint without auth
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("X-API-Key")
		if key != s.cfg.Auth.APIKey {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized: invalid or missing API key",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": collector.Version,
	})
}

func (s *Server) handleMetricsAll(w http.ResponseWriter, _ *http.Request) {
	metrics, err := collector.CollectAll(s.cfg)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, metrics)
}

func (s *Server) handleMetricsOS(w http.ResponseWriter, _ *http.Request) {
	data, err := collector.CollectOS()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleMetricsCPU(w http.ResponseWriter, _ *http.Request) {
	data, err := collector.CollectCPU()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleMetricsMemory(w http.ResponseWriter, _ *http.Request) {
	data, err := collector.CollectMemory()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleMetricsDisk(w http.ResponseWriter, _ *http.Request) {
	data, err := collector.CollectDisk()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleMetricsNetwork(w http.ResponseWriter, _ *http.Request) {
	data, err := collector.CollectNetwork()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (s *Server) handleMetricsUptime(w http.ResponseWriter, _ *http.Request) {
	data, err := collector.CollectUptime()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Printf("error encoding JSON: %v", err)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": fmt.Sprintf("failed to collect metrics: %v", err),
	})
}
