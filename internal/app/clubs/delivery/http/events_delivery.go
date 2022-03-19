package delivery

import (
	"encoding/json"
	clubs "github.com/dantedoyl/car-life-api/internal/app/clubs"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ClubsHandler struct {
	clubsUcase clubs.IClubsUsecase
}

func NewClubsHandler(clubsUcase clubs.IClubsUsecase) *ClubsHandler {
	return &ClubsHandler{
		clubsUcase: clubsUcase,
	}
}

func (ch *ClubsHandler) Configure(r *mux.Router) {
	r.HandleFunc("/club/create", ch.CreateClub).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}", ch.GetClubByID).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs", ch.GetClubs).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/upload}", ch.UploadAvatarHandler).Methods(http.MethodPost, http.MethodOptions)
}

// CreateEvent godoc
// @Summary      create a club
// @Description  Handler for creating a club
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        body body models.Club true "Club"
// @Success      200  {object}  models.Club
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /club/create [post]
func (ch *ClubsHandler) CreateClub(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	clubsData := &models.Club{}
	err := json.NewDecoder(r.Body).Decode(&clubsData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	err = ch.clubsUcase.CreateClub(clubsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(clubsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetClubs godoc
// @Summary      get clubs list
// @Description  Handler for getting clubs list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Club
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs [get]
func (ch *ClubsHandler) GetClubs(w http.ResponseWriter, r *http.Request) {
	event, err := ch.clubsUcase.GetClubs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(event) == 0 {
		event = []*models.Club{}
	}

	body, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetClubByID godoc
// @Summary      get club by id
// @Description  Handler for getting a club by id
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Success      200  {object}  models.Club
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id} [get]
func (ch *ClubsHandler) GetClubByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	club, err := ch.clubsUcase.GetClubByID(clubID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(club)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// UploadAvatarHandler godoc
// @Summary      upload avatar for club
// @Description  Handler for uploading a club's avatar
// @Tags         Clubs
// @Accept       mpfd
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Success      200  {object}  models.Club
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/upload [post]
func (ch *ClubsHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	clubID, _ := strconv.ParseInt(vars["id"], 10, 64)

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
	club, errE := ch.clubsUcase.UpdateAvatar(clubID, file)
	if errE != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(club)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}