package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/koind/action-log/internal/domain/repository"
	"github.com/koind/action-log/internal/domain/service"
	"net/http"
)

// HTTP сервер
type HTTPServer struct {
	http.Server
	router         http.Handler
	domain         string
	historyService *service.HistoryService
}

// Возвращает новый HTTP сервер
func NewHTTPServer(historyService *service.HistoryService, domain string) *HTTPServer {

	r := mux.NewRouter()
	hs := HTTPServer{router: r, domain: domain, historyService: historyService}

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
	history := repository.History{}

	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

		return
	}

	validate := validator.New()
	err = validate.StructCtx(r.Context(), history)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

		return
	}

	newHistory, err := s.historyService.Add(r.Context(), history)
	if err != nil {
		fmt.Fprint(w, err)

		return
	}

	err = json.NewEncoder(w).Encode(newHistory)
	if err != nil {
		fmt.Fprint(w, err)

		return
	}
}

// Возвращает список действий
func (s *HTTPServer) GetHistoriesHandle(w http.ResponseWriter, r *http.Request) {
	filter := repository.SearchFilter{}

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

		return
	}

	validate := validator.New()
	err = validate.StructCtx(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

		return
	}

	list, err := s.historyService.FindAllByFilter(r.Context(), filter)
	if err != nil {
		fmt.Fprint(w, err)

		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Fprint(w, err)

		return
	}
}
