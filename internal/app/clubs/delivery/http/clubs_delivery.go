package delivery

import (
	"encoding/json"
	clubs "github.com/dantedoyl/car-life-api/internal/app/clubs"
	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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

func (ch *ClubsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/club/create", mw.CheckAuthMiddleware(ch.CreateClub)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}", mw.CheckAuthMiddleware(ch.GetClubByID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs", mw.CheckAuthMiddleware(ch.GetClubs)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/tags", mw.CheckAuthMiddleware(ch.GetTags)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/upload", mw.CheckAuthMiddleware(ch.UploadAvatarHandler)).Methods(http.MethodPost, http.MethodOptions)
	//r.HandleFunc("/clubs/{id:[0-9]+}/subscribe", mw.CheckAuthMiddleware(ch.SubscribeByClubID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/participate", mw.CheckAuthMiddleware(ch.ParticipateByClubID)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{cid:[0-9]+}/participate/{uid:[0-9]+}/{type:approve|reject}", mw.CheckAuthMiddleware(ch.ApproveRejectUserParticipate)).Methods(http.MethodPost, http.MethodOptions)
	//r.HandleFunc("/clubs/{id:[0-9]+}/subscribers", mw.CheckAuthMiddleware(ch.GetClubsSubscribers)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/participants", mw.CheckAuthMiddleware(ch.GetClubsParticipants)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/cars", mw.CheckAuthMiddleware(ch.GetClubsCars)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/events", mw.CheckAuthMiddleware(ch.GetClubsEvents)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/participants/requests", mw.CheckAuthMiddleware(ch.GetClubsParticipantsRequests)).Methods(http.MethodGet, http.MethodOptions)

}

// CreateClub godoc
// @Summary      create a club
// @Description  Handler for creating a club
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        body body models.CreateClubRequest true "Club"
// @Success      200  {object}  models.Club
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /club/create [post]
func (ch *ClubsHandler) CreateClub(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// добавить проверку авторизации
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	club := &models.CreateClubRequest{}
	err := json.NewDecoder(r.Body).Decode(&club)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	clubsData := &models.Club{
		Name:        club.Name,
		Description: club.Description,
		AvatarUrl:   club.AvatarUrl,
		Tags:        club.Tags,
		OwnerID:     userID,
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
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Param        Query query string false "Query"
// @Success      200  {object}  []models.ClubCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs [get]
func (ch *ClubsHandler) GetClubs(w http.ResponseWriter, r *http.Request) {
	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	clubs, err := ch.clubsUcase.GetClubs(query.IdGt, query.IdLte, query.Limit, query.Query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(clubs) == 0 {
		clubs = []*models.Club{}
	}

	clubCards := make([]models.ClubCard, 0, len(clubs))
	for _, club := range clubs {
		clubCards = append(clubCards, models.ClubCard{
			ID:                club.ID,
			Name:              club.Name,
			AvatarUrl:         club.AvatarUrl,
			Tags:              club.Tags,
			ParticipantsCount: club.ParticipantsCount,
		})
	}

	body, err := json.Marshal(clubCards)
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
	club, err := ch.clubsUcase.UpdateAvatar(clubID, file)
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

// GetTags godoc
// @Summary      get tags list
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Success      200  {object}  []models.Tag
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/tags [get]
func (ch *ClubsHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := ch.clubsUcase.GetTags()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(tags) == 0 {
		tags = []models.Tag{}
	}

	body, err := json.Marshal(tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetClubsParticipants godoc
// @Summary      get clubs participants list
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.UserCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/participants [get]
func (ch *ClubsHandler) GetClubsParticipants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	users, err := ch.clubsUcase.GetClubsUserByStatus(int64(clubID), "participant", query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(users) == 0 {
		users = []*models.UserCard{}
	}

	body, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetClubsCars godoc
// @Summary      get clubs cars list
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.CarCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/cars [get]
func (ch *ClubsHandler) GetClubsCars(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	cars, err := ch.clubsUcase.GetClubsCars(int64(clubID), query.IdGt, query.IdLte, query.Limit)
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

// GetClubsParticipantsRequests godoc
// @Summary      get clubs participants request list
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.UserCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/participants/requests [get]
func (ch *ClubsHandler) GetClubsParticipantsRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	users, err := ch.clubsUcase.GetClubsUserByStatus(int64(clubID), "participant_request", query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(users) == 0 {
		users = []*models.UserCard{}
	}

	body, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}


// GetClubsEvents godoc
// @Summary      get clubs events list
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.EventCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/cars [get]
func (ch *ClubsHandler) GetClubsEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	events, err := ch.clubsUcase.GetClubsEvents(int64(clubID), query.IdGt, query.IdLte, query.Limit)
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

// ParticipateByClubID godoc
// @Summary      request participate
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/participate [post]
func (ch *ClubsHandler) ParticipateByClubID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	err := ch.clubsUcase.SetUserStatusByClubID(int64(clubID), int64(userID), "participant_request")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ApproveRejectUserParticipate godoc
// @Summary      request participate
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        cid path int64 true "Club ID"
// @Param        uid path int64 true "User ID"
// @Param        type path string true "Type" Enums(approve, reject)
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{cid}/participate/{uid}/{type} [post]
func (ch *ClubsHandler) ApproveRejectUserParticipate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["cid"], 10, 64)
	userID, _ := strconv.ParseUint(vars["uid"], 10, 64)
	decision := vars["type"]

	_, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	err := ch.clubsUcase.ApproveRejectUserParticipate(int64(clubID), int64(userID), decision)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}