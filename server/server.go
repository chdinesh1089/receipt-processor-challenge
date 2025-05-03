package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router   *mux.Router
	receipts map[string]Receipt
	log      *logrus.Entry
}

func NewServer() *Server {
	s := &Server{
		router:   mux.NewRouter(),
		receipts: map[string]Receipt{},
		log:      getLogger(context.Background()),
	}

	s.router.Use(s.LoggingMiddleware)
	s.router.HandleFunc("/receipts/process", s.processReceipts).Methods("POST")
	s.router.HandleFunc("/receipts/{id}/points", s.getPoints).Methods("GET")
	return s
}

func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		request_id := uuid.New().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDKey, request_id)
		r = r.WithContext(ctx)

		s.log.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"request_id": request_id,
		}).Info("request received")

		next.ServeHTTP(w, r)

		s.log.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"duration":   time.Since(start).String(),
			"request_id": request_id,
		}).Info("request processed")
	})
}

func (s *Server) processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	json.NewDecoder(r.Body).Decode(&receipt)
	// TODO: Validate receipt
	id := uuid.New()
	receipt_id := id.String()
	s.receipts[receipt_id] = receipt
	json.NewEncoder(w).Encode(map[string]string{"id": receipt_id})
}

func (s *Server) getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receipt_id := vars["id"]
	receipt, ok := s.receipts[receipt_id]
	if !ok {
		s.log.WithField("request_id", r.Context().Value(RequestIDKey)).Warn("Receipt not found, receipt_id:", receipt_id, "receipts:", s.receipts)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Perhaps Cache this so we don't compute it every time
	json.NewEncoder(w).Encode(map[string]int64{"points": receipt.Points(r.Context())})
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(":80", s.router))
}
