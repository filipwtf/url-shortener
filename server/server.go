package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/filipwtf/url-longer/postgres"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Server struct
type Server struct {
	*mux.Router
	*postgres.Queries
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	Dev          bool
}

//NewServer creates a new server
func NewServer(db *postgres.Queries, flag bool) *Server {
	s := &Server{
		Router:       mux.NewRouter(),
		Queries:      db,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Dev:          flag,
	}
	s.Routes()
	return s
}

// Routes Creating routes
func (s Server) Routes() {
	s.HandleFunc("/url/make", s.CreateLongerURL()).Methods("POST")
	if s.Dev {
		s.HandleFunc("/url/all", s.getAllUrls()).Methods("GET")
	}
	s.HandleFunc("/{id}", s.getURL()).Methods("GET")
}

// CreateLongerURL adds new url to database
func (s *Server) CreateLongerURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go LogRequest(w, r)
		var u postgres.URL
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}

		u.ID = uuid.New().String()

		_, err := s.Queries.CreateLonger(context.Background(), postgres.CreateLongerParams{
			ID:       u.ID,
			Original: u.Original,
		})
		if err != nil {
			HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		full := fmt.Sprintf("https://000000000000000000000000000000000000000000000000000000000000000.000000000000000000000000000000000000000000000000000000000000000.xyz/%s", u.ID)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(full); err != nil {
			HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go LogRequest(w, r)
		idStr, _ := mux.Vars(r)["id"]
		fmt.Println(idStr)
		url, err := s.Queries.GetOriginal(context.Background(), idStr)
		if err != nil {
			HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, url.Original, http.StatusPermanentRedirect)
	}
}

// This would be locked during production
func (s Server) getAllUrls() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go LogRequest(w, r)
		w.Header().Set("content-type", "application/json")

		rows, err := s.Queries.GetUrls(context.Background())
		if err != nil {
			HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(rows); err != nil {
			HandleError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// HandleError handles any http errors
func HandleError(w http.ResponseWriter, err string, status int) {
	http.Error(w, err, status)
}

// LogRequest logs http request
func LogRequest(w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()
	uri := r.RequestURI
	method := r.Method
	ip := r.RemoteAddr
	cLog := color.New(color.BgBlack, color.FgHiBlue)
	cLog.Printf("🍔 %s %s %s %s\n", uri, agent, method, ip)
}
