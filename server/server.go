package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	router   *mux.Router
	receipts map[string]Receipt
}

func NewServer() *Server {
	s := &Server{
		router:   mux.NewRouter(),
		receipts: map[string]Receipt{},
	}
	s.router.HandleFunc("/receipts/process", s.processReceipts).Methods("POST")
	s.router.HandleFunc("/receipts/{id}/points", s.getPoints).Methods("GET")
	return s
}

func (s *Server) processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	json.NewDecoder(r.Body).Decode(&receipt)
	// TODO: Validate receipt
	id := uuid.New()
	receipt_id := id.String()
	s.receipts[receipt_id] = receipt
	json.NewEncoder(w).Encode(receipt_id)
}

func (s *Server) getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receipt_id := vars["id"]
	receipt, ok := s.receipts[receipt_id]
	if !ok {
		fmt.Println("Receipt not found, receipt_id:", receipt_id, "receipts:", s.receipts)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Perhaps Cache this so we don't compute it every time
	json.NewEncoder(w).Encode(receipt.Points())
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(":80", s.router))
}
