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

var invalidReceiptErrStr = "The receipt is invalid."
var receiptNotFoundErrStr = "No receipt found for that ID."

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

func (s *Server) writeJsonResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		s.log.WithField("request_id", w.Header().Get("request_id")).
			Error("writeJsonResponse> Failed to encode response to json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
}

func (s *Server) processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)

	if err != nil {
		s.log.WithField("request_id", r.Context().Value(RequestIDKey)).Error(err)
		s.writeJsonResponse(w, http.StatusBadRequest, map[string]string{"error": invalidReceiptErrStr, "details": err.Error()})
		return
	}

	if err := receipt.Validate(r.Context()); err != nil {
		s.log.WithField("request_id", r.Context().Value(RequestIDKey)).Error(err)
		s.writeJsonResponse(w, http.StatusBadRequest, map[string]string{"error": invalidReceiptErrStr, "details": err.Error()})
		return
	}

	id := uuid.New()
	receipt_id := id.String()
	s.receipts[receipt_id] = receipt
	s.writeJsonResponse(w, http.StatusOK, map[string]string{"id": receipt_id})
}

func (s *Server) getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receipt_id := vars["id"]
	receipt, ok := s.receipts[receipt_id]
	if !ok {
		s.log.WithField("request_id", r.Context().Value(RequestIDKey)).
			Error("getPoints> Receipt not found, receipt_id:", receipt_id, "receipts:", s.receipts)
		s.writeJsonResponse(w, http.StatusNotFound, map[string]string{"error": receiptNotFoundErrStr})
		return
	}

	// Perhaps Cache this so we don't compute it every time
	s.writeJsonResponse(w, http.StatusOK, map[string]int64{"points": receipt.Points(r.Context())})
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(":80", s.router))
}
