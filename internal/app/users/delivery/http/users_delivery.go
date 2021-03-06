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
	r.HandleFunc("/me/update", mw.CheckAuthMiddleware(uh.UpdateUserProfile)).Methods(http.MethodPut, http.MethodOptions)
	r.HandleFunc("/new_car", mw.CheckAuthMiddleware(uh.NewUserCar)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/garage", uh.UserGarage).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/complain", mw.CheckAuthMiddleware(uh.ComplainUser)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/events/{type:admin|participant|spectator}", uh.UserEvents).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{id:[0-9]+}/clubs/{type:admin|participant|subscriber}", uh.UserClubs).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/own_clubs", mw.CheckAuthMiddleware(uh.UserOwnClubs)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/login", uh.Login).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/garage/{id:[0-9]+}/upload", uh.UploadAvatarHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/garage/{id:[0-9]+}", uh.GetCarByID).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/garage/{id:[0-9]+}/delete", mw.CheckAuthMiddleware(uh.DeleteCar)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/garage/{id:[0-9]+}/complain", mw.CheckAuthMiddleware(uh.ComplainCar)).Methods(http.MethodPost, http.MethodOptions)
}

// SignUp godoc
// @Summary      sign uo new user
// @Description  Handler for signing up new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body models.SignUpRequest true "User"
// @Success      200 {object} models.SignUpResponse
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

	body, err := json.Marshal(models.SignUpResponse{
		CarID:   user.CarID,
		Session: session,
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
// @Success      200  {object}  models.Session
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

	body, err := json.Marshal(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UploadAvatarHandler godoc
// @Summary      upload avatar for car
// @Description  Handler for creating an event
// @Tags         Users
// @Accept       mpfd
// @Produce      json
// @Param        id path int64 true "Car ID"
// @Param 		 file-upload formData file true "Image to upload"
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

	if user.Tags == nil {
		user.Tags = []string{}
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

// UpdateUserProfile godoc
// @Summary      get user by id
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body models.UpdateRequest true "User"
// @Success      200  {object}  models.User
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /me/update [put]
func (uh *UsersHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	signUp := models.UpdateRequest{}
	err := json.NewDecoder(r.Body).Decode(&signUp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "unable to decode data"}))
		return
	}

	user := &models.User{
		VKID:        userID,
		Tags:        signUp.Tags,
		Description: signUp.Description,
	}

	user, err = uh.usersUcase.UpdateUserInfo(user)
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

	if user.Tags == nil {
		user.Tags = []string{}
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
	if len(clubs) == 0 {
		clubs = []*models.ClubCard{}
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
	if len(cars) == 0 {
		cars = []*models.CarCard{}
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
	if len(clubs) == 0 {
		clubs = []*models.ClubCard{}
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
// @Success      200  {object}  []models.EventCard
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

	events, err := uh.usersUcase.GetEventsByUserStatus(int64(userID), role, query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(events) == 0 {
		events = []*models.EventCard{}
	}

	body, err := json.Marshal(events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// NewUserCar godoc
// @Summary      add new car to user
// @Description  Handler for signing up new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body models.CarRequest true "Car"
// @Success      200 {object} models.CarCard
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /new_car [post]
func (uh *UsersHandler) NewUserCar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	car := models.CarRequest{}
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "unable to decode data"}))
		return
	}

	carData := &models.CarCard{
		Brand:       car.Brand,
		Model:       car.Model,
		Date:        car.Date,
		Description: car.Description,
		Body:        car.Body,
		Engine:      car.Engine,
		HorsePower:  car.HorsePower,
		Name:        car.Name,
		Owner:     	 models.UserCard{
			VKID:      userID,
		},
	}

	carData, err = uh.usersUcase.AddNewUserCar(carData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(carData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetCarByID godoc
// @Summary      get car
// @Description  Handler for getting a user by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Car ID"
// @Success      200  {object}  models.CarCard
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /garage/{id} [get]
func (uh *UsersHandler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carID, _ := strconv.ParseUint(vars["id"], 10, 64)

	cars, err := uh.usersUcase.SelectCarByID(int64(carID))
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

// DeleteCar godoc
// @Summary      delete car
// @Description  Handler for deleting car
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Car ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /garage/{id}/delete [post]
func (uh *UsersHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	car, err := uh.usersUcase.SelectCarByID(int64(carID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if car.Owner.VKID != userID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	err = uh.usersUcase.DeleteCarByID(int64(carID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}


// ComplainUser godoc
// @Summary      complain user
// @Description  Handler for complaining user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "User ID"
// @Param        body body models.ComplaintReq true "User"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /user/{id}/complain [post]
func (uh *UsersHandler) ComplainUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	req := &models.ComplaintReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	err = uh.usersUcase.ComplainByID("user", models.Complaint{
		UserID:   int64(userID),
		Text:     req.Text,
		TargetID: int64(clubID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ComplainCar godoc
// @Summary      complain car
// @Description  Handler for complaining car
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Car ID"
// @Param        body body models.ComplaintReq true "Car"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /garage/{id}/complain [post]
func (uh *UsersHandler) ComplainCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	req := &models.ComplaintReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	err = uh.usersUcase.ComplainByID("car", models.Complaint{
		UserID:   int64(userID),
		Text:     req.Text,
		TargetID: int64(clubID),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}
