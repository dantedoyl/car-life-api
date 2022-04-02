package delivery

import (
	"encoding/json"
	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/users"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

type UsersHandler struct {
	usersUcase users.IUsersUsecase
}

func NewUserssHandler(usersUcase users.IUsersUsecase) *UsersHandler {
	return &UsersHandler{
		usersUcase: usersUcase,
	}
}

func (uh *UsersHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/signup", uh.SignUp).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/me", mw.CheckAuthMiddleware(uh.MyProfile)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}", uh.UserProfile).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/garage", uh.UserGarage).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/events/{type:admin|participant|spectator}", uh.UserEvents).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/clubs/{type:admin|participant|subscriber}", uh.UserClubs).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/own_clubs", mw.CheckAuthMiddleware(uh.UserOwnClubs)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/login", uh.Login).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/garage/{id:[0-9]+}/upload", uh.UploadAvatarHandler).Methods(http.MethodPost, http.MethodOptions)
}

// SignUp godoc
// @Summary      sign uo new user
// @Description  Handler for signing up new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body models.SignUpRequest true "User"
// @Success      200 {object} models.CarIDResponse
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /signup [post]
func (uh *UsersHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	signUp := models.SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&signUp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "unable to decode data"}))
		return
	}

	user := &models.User{
		VKID:        signUp.VKID,
		Tags:        signUp.Tags,
		Name:        signUp.Name,
		Surname:     signUp.Surname,
		AvatarUrl:   signUp.AvatarUrl,
		Description: signUp.Description,
	}

	var car *models.CarCard
	if len(signUp.Garage) != 0 {
		car = &models.CarCard{
				Brand:       signUp.Garage[0].Brand,
				Model:       signUp.Garage[0].Model,
				Date:        signUp.Garage[0].Date,
				Description: signUp.Garage[0].Description,
				Body:        signUp.Garage[0].Body,
				Engine:      signUp.Garage[0].Engine,
				HorsePower:  signUp.Garage[0].HorsePower,
				Name:        signUp.Garage[0].Name,
		}
	}

	user, err = uh.usersUcase.Create(user, car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	session := models.CreateSession(user.VKID)
	err = uh.usersUcase.CreateSession(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	body, err := json.Marshal(models.CarIDResponse{
		CarID: user.CarID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Login godoc
// @Summary      login user
// @Description  Handler for signing up new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body models.LoginRequest true "User"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /login [post]
func (uh *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	login := &models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "unable to decode data"}))
		return
	}

	user, err := uh.usersUcase.GetByID(login.VKID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session := models.CreateSession(login.VKID)
	err = uh.usersUcase.CreateSession(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Value,
		Expires:  session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

// UploadAvatarHandler godoc
// @Summary      upload avatar for car
// @Description  Handler for creating an event
// @Tags         Users
// @Accept       mpfd
// @Produce      json
// @Param        id path int64 true "Car ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /garage/{id}/upload [post]
func (uh *UsersHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	carID, _ := strconv.ParseInt(vars["id"], 10, 64)

	r.Body = http.MaxBytesReader(w, r.Body, 3*1024*1024)
	err := r.ParseMultipartForm(3 * 1024 * 1024)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't parse data"}))
		return
	}

	if len(r.MultipartForm.File["file-upload"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "no photo"}))
		return
	}

	file := r.MultipartForm.File["file-upload"][0]
	user, err := uh.usersUcase.UpdateAvatar(uint64(carID), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UserProfile godoc
// @Summary      get user by id
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /user/{id} [get]
func (uh *UsersHandler) UserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.ParseUint(vars["id"], 10, 64)

	user, err := uh.usersUcase.GetByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// MyProfile godoc
// @Summary      get user by id
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /me [get]
func (uh *UsersHandler) MyProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	user, err := uh.usersUcase.GetByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UserOwnClubs godoc
// @Summary      get clubs where user is owner
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.ClubCard
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /user/own_clubs [get]
func (uh *UsersHandler) UserOwnClubs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	clubs, err := uh.usersUcase.GetClubsByUserStatus(int64(userID), "admin", query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(clubs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UserGarage godoc
// @Summary      get user garage
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "User ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.CarCard
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /user/{id}/garage [get]
func (uh *UsersHandler) UserGarage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	cars, err := uh.usersUcase.SelectCarByUserID(int64(userID), query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(cars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UserClubs godoc
// @Summary      get clubs where user is in status
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "User ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Param        type path string true "Type" Enums(admin, participant, subscriber)
// @Success      200  {object}  []models.ClubCard
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /user/{id}/clubs/{type} [get]
func (uh *UsersHandler) UserClubs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.ParseUint(vars["id"], 10, 64)
	role := vars["type"]

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	clubs, err := uh.usersUcase.GetClubsByUserStatus(int64(userID), role, query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(clubs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UserEvents godoc
// @Summary      get events where user is in status
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "User ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Param        type path string true "Type" Enums(admin, participant, spectator)
// @Success      200  {object}  []models.ClubCard
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /user/{id}/events/{type} [get]
func (uh *UsersHandler) UserEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.ParseUint(vars["id"], 10, 64)
	role := vars["type"]

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	clubs, err := uh.usersUcase.GetEventsByUserStatus(int64(userID), role, query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(clubs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
