package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/GritselMaks/postgreSQL-api-server/internal/app/model"
	"github.com/GritselMaks/postgreSQL-api-server/internal/app/store"
	"github.com/gorilla/mux"
)

type APIServer struct {
	router *mux.Router
	store  store.Store
}

// new API server
func NewServer(store store.Store) *APIServer {
	s := APIServer{
		router: mux.NewRouter(),
		store:  store,
	}
	s.configRouter()
	return &s
}

//start API server
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

// response all articles
func (s *APIServer) HandleShowArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := s.store.User().ShowArticles(r.URL.Query())
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, "Error in select From database")
		}
		s.respond(w, r, http.StatusOK, articles)
	}
}

//response one article
func (s *APIServer) HandleShowArticle() http.HandlerFunc {

	var article model.Articles
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if err := s.store.User().ShowArticle(article, vars); err != nil {

			s.respond(w, r, http.StatusInternalServerError, "Error in select From database")
		}
		s.respond(w, r, http.StatusOK, article)
	}
}

//creare article
func (s *APIServer) HandleCreate() http.HandlerFunc {

	type request struct {
		Title    string   `json:"title"`
		FullText string   `json:"fulltext"`
		Prise    float64  `json:"prise"`
		URLFoto  []string `json:"urlfoto"`
		Data     string   `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.respond(w, r, http.StatusBadRequest, err)
		}
		art := &model.Articles{
			Title:    req.Title,
			FullText: req.FullText,
			Price:    req.Prise,
			Data:     req.Data,
			URLFoto:  req.URLFoto,
		}
		id, err := s.store.User().Save(art)
		if err != nil {
			s.respond(w, r, http.StatusExpectationFailed, err)
		}
		s.respond(w, r, http.StatusCreated, id)
	}

}

//respond http status and Json string
func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

	// w.WriteHeader(code)
	// if result, err := json.MarshalIndent(data, "", " "); err == nil {
	// 	w.Write(result)
	// }
}
