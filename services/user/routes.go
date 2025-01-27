package user

import (
	"fmt"
	"net/http"

	"github.com/FarheenParvez/GOECOMAPI/services/auth"
	"github.com/FarheenParvez/GOECOMAPI/types"
	"github.com/FarheenParvez/GOECOMAPI/utils"

	//"github.com/goccy/go-json"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore


}

func NewHandler(store types.UserStore) *Handler {

	return &Handler{
		store: store,
	}

}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleLogin).Methods("POST")

}

func (h* Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h* Handler) handleRegistery(w http.ResponseWriter, r *http.Request) {
	// 1. get jason payload from request
	if r.Body == nil {

	}
    var payload types.RegisterUserPayload
	if err:= utils.ParseJSON(r, payload); err!= nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	//err:= json.NewDecoder(r.Body).Decode(payload)

	// 2. validate the payload i.e cheecck if user already exists

	 _, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user  with email %s already exists", payload.Email))
		return 
	}
	
	hashedpassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// 3. if user does not exist, create a new user
	err = h.store.CreatedUser(types.User{
		FirstName: payload.FirstName,
		LastName: payload.Lastname,
	    Email: payload.Email,
		Password: hashedpassword,	
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
	 

}