package handler

import (
	"github.com/gorilla/mux"
	"github.com/koind/action-log/internal/domain/service"
	"net/http"
)

// HTTP сервер
type HTTPServer struct {
	http.Server
	router http.Handler
	domain string
	s      *service.HistoryService
}

// Возвращает новый HTTP сервер
func NewHTTPServer(historyService *service.HistoryService, domain string) *HTTPServer {

	r := mux.NewRouter()
	hs := HTTPServer{router: r, domain: domain, s: historyService}

	r.HandleFunc("/history", hs.AddHistoryHandle).Methods("POST")
	r.HandleFunc("/histories", hs.GetHistoriesHandle).Methods("GET")

	http.Handle("/", r)

	return &hs
}

// Запускает HTTP сервер
func (s *HTTPServer) Start() error {
	return http.ListenAndServe(s.domain, s.router)
}

// Добавляет новую историю действий
func (s *HTTPServer) AddHistoryHandle(w http.ResponseWriter, r *http.Request) {

}

// Возвращает список действий
func (s *HTTPServer) GetHistoriesHandle(w http.ResponseWriter, r *http.Request) {

}
