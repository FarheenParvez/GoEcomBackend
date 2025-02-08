package user

import (
	"fmt"
	"net/http"

	"github.com/FarheenParvez/GOECOMAPI/config"
	"github.com/FarheenParvez/GOECOMAPI/services/auth"
	"github.com/FarheenParvez/GOECOMAPI/types"
	"github.com/FarheenParvez/GOECOMAPI/utils"
	"github.com/go-playground/validator/v10"

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
	router.HandleFunc("/register", h.handleRegistery).Methods("POST")

}

func (h* Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//err:= json.NewDecoder(r.Body).Decode(payload)

	// 2. validate the payload i.e cheecck if user already exists

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err:= h.store.GetUserByEmail(user.Email) 
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)

	token, err := auth.CreateJWT(secret, u.ID)

	if err!= nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token}) 




}

func (h* Handler) handleRegistery(w http.ResponseWriter, r *http.Request) {
	// 1. get jason payload from request
    var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//err:= json.NewDecoder(r.Body).Decode(payload)

	// 2. validate the payload i.e cheecck if user already exists

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	 _, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user  with email %s already exists", user.Email))
		return 
	}
	
	hashedpassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// 3. if user does not exist, create a new user
	err = h.store.CreatedUser(types.User{
		FirstName: user.FirstName,
		LastName: user.Lastname,
	    Email: user.Email,
		Password: hashedpassword,	
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
	 

}