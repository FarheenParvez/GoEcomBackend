package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/FarheenParvez/GOECOMAPI/services/user"
	"github.com/gorilla/mux"
)


type APIServer struct {
	addr string 
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer{
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	//1. Initialise a router 
	//2. Register the routes
	//3. Start the server
	//4. Return any error

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter() // api/v1 is the base path for all the routes. It is important because if in future we change the version of the API, we can easily change it.
   	
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	
   	userHandler.RegisterRoutes(subrouter)


	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router);
}