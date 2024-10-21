package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/scorcism/go-auth/service"
	"github.com/scorcism/go-auth/service/secrets"
	"github.com/scorcism/go-auth/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

// sort of constructor
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := store.NewStore(s.db)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	secretsStore := secrets.NewStore(s.db)
	secretsHandler := secrets.NewHandler(secretsStore, userStore)
	secretsHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
