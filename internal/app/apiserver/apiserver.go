package apiserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GritselMaks/postgreSQL-api-server/internal/app/model"
	"github.com/GritselMaks/postgreSQL-api-server/internal/app/store"
	"github.com/gorilla/mux"
)

// Create new API server
type APIServer struct {
	router *mux.Router
	store  store.Store
}

// Start API server
func NewServer(store store.Store) *APIServer {
	s := APIServer{
		router: mux.NewRouter(),
		store:  store,
	}
	s.configRouter()
	return &s
}

// Configurate router
func Start(config *Config) error {
	db, err := store.NewDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := store.New(db)
	srv := NewServer(*store)

	return http.ListenAndServe(config.BinAddr, srv.router)
}

func (s *APIServer) configRouter() {

	s.router.HandleFunc("/articles", s.HandleShowArticles()).Methods("GET")
	s.router.HandleFunc("/article/{id}", s.HandleShowArticle()).Methods("GET")
	s.router.HandleFunc("/articles/new", s.HandleCreate()).Methods("POST")

}

// Handles the request to receive data and returns datas in Json format
func (s *APIServer) HandleShowArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParam := strings.Split(r.URL.Query().Get("sort"), ",")
		for i, s := range sortParam {
			switch s {
			case "price":
				sortParam[i] = "price DESC"
			case "-price":
				sortParam[i] = "price ASC"
			case "-date":
				sortParam[i] = "date_at ASC"
			case "date":
				sortParam[i] = "date_at DESC"
			default:
				sortParam[i] = ""
			}
		}
		articles, err := s.store.Articles().ShowArticles(sortParam)
		if err != nil {
			s.respond(w, r, http.StatusNoContent, err)
		} else {
			s.respond(w, r, http.StatusOK, articles)
		}
	}
}

// Handles the request to receive data and returns datas in Json format
func (s *APIServer) HandleShowArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)["id"]
		fields := r.FormValue("fields")
		article, err := s.store.Articles().ShowArticle(vars, fields)
		if err != nil {
			s.respond(w, r, http.StatusNoContent, article)
		} else {
			s.respond(w, r, http.StatusOK, article)
		}
	}
}

// Creare article
func (s *APIServer) HandleCreate() http.HandlerFunc {

	type request struct {
		Title    string   `json:"title"`
		FullText string   `json:"fulltext"`
		Price    int      `json:"price"`
		URLFoto  []string `json:"urlfoto"`
		Date     string   `json:"date"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.respond(w, r, http.StatusBadRequest, "not valid request")
		}
		art := &model.Articles{
			Title:    req.Title,
			FullText: req.FullText,
			Price:    req.Price,
			Date:     req.Date,
			URLFoto:  req.URLFoto,
		}
		id, err := s.store.Articles().Save(art)
		if err != nil {
			s.respond(w, r, http.StatusExpectationFailed, err)
		}
		s.respond(w, r, http.StatusCreated, id)
	}

}

// Forms a response with http Status and Json data
func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if data != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(data)
	}
}
