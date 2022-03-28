package delivery

import (
	"encoding/json"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/users"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
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

func (uh *UsersHandler) Configure(r *mux.Router) {
	r.HandleFunc("/signup", uh.SignUp).Methods(http.MethodPost, http.MethodOptions)
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
// @Success      200
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
		VKID:       signUp.VKID,
		Name:       signUp.Name,
		Surname:    signUp.Surname,
		AvatarUrl:  signUp.AvatarUrl,
	}

	err = uh.usersUcase.Create(user)
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
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
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

	user, errE := uh.usersUcase.GetByID(login.VKID)
	if errE != nil {
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
		SameSite: http.SameSiteLaxMode,
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


