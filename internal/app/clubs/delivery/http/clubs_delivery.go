package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/dantedoyl/car-life-api/internal/app/clients/vk"
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
	vk         *vk.VKClient
}

func NewClubsHandler(clubsUcase clubs.IClubsUsecase, vkCl *vk.VKClient) *ClubsHandler {
	return &ClubsHandler{
		clubsUcase: clubsUcase,
		vk:         vkCl,
	}
}

func (ch *ClubsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/club/create", mw.CheckAuthMiddleware(ch.CreateClub)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}", mw.CheckAuthMiddleware(ch.GetClubByID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs", mw.CheckAuthMiddleware(ch.GetClubs)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/tags", mw.CheckAuthMiddleware(ch.GetTags)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/upload", mw.CheckAuthMiddleware(ch.UploadAvatarHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/{type:participate|subscribe}", mw.CheckAuthMiddleware(ch.SetUserStatusByClubID)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{cid:[0-9]+}/participate/{uid:[0-9]+}/{type:approve|reject}", mw.CheckAuthMiddleware(ch.ApproveRejectUserParticipateInClub)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/{type:participant|participant_request|subscriber}", mw.CheckAuthMiddleware(ch.GetClubsUsersByType)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/leave", mw.CheckAuthMiddleware(ch.LeaveClub)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/cars", mw.CheckAuthMiddleware(ch.GetClubsCars)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/events", mw.CheckAuthMiddleware(ch.GetClubsEvents)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/chat_link", mw.CheckAuthMiddleware(ch.GetClubChatLink)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/delete", mw.CheckAuthMiddleware(ch.DeleteClub)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/clubs/{id:[0-9]+}/complain", mw.CheckAuthMiddleware(ch.ComplainClub)).Methods(http.MethodPost, http.MethodOptions)
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
		Owner:       models.UserCard{
			VKID:      userID,
		},
	}

	err = ch.clubsUcase.CreateClub(clubsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	id, err := ch.vk.CreatChat(clubsData.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	err = ch.clubsUcase.SetClubChatID(int64(clubsData.ID), int64(id))
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
			SubscribersCount: club.SubscribersCount,
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

	club := &models.Club{}
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		club.UserStatus = "unknown"
	}

	club, err := ch.clubsUcase.GetClubByID(clubID, userID)
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
// @Param 		 file-upload formData file true "Image to upload"
// @Success      200  {object}  models.Club
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/upload [post]
func (ch *ClubsHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseInt(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

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

	chatID, err := ch.clubsUcase.GetClubChatID(clubID, int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	err = ch.vk.UploadChatPhoto(int(chatID), file)
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

// GetClubsUsersByType godoc
// @Summary      get clubs users list
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Param        type path string true "Type" Enums(participant, participant_request, subscriber)
// @Success      200  {object}  []models.UserCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/{type} [get]
func (ch *ClubsHandler) GetClubsUsersByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)
	role := vars["type"]

	if role == "participant_request" {
		userID, ok := r.Context().Value("userID").(uint64)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
			return
		}

		userClubSatus, err := ch.clubsUcase.GetUserStatusInClub(int64(clubID), int64(userID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
			return
		}

		if userClubSatus == nil || userClubSatus.Status != "admin" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
			return
		}
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

	users, err := ch.clubsUcase.GetClubsUserByStatus(int64(clubID), role, query.IdGt, query.IdLte, query.Limit)
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
// @Router       /clubs/{id}/events [get]
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

// SetUserStatusByClubID godoc
// @Summary      set user role in club
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        type path string true "Type" Enums(participate, subscribe)
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/{type} [post]
func (ch *ClubsHandler) SetUserStatusByClubID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)
	decision := vars["type"]

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	userClubSatus, err := ch.clubsUcase.GetUserStatusInClub(int64(clubID), int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if decision == "participate" && userClubSatus != nil && (userClubSatus.Status == "participant" || userClubSatus.Status == "admin" || userClubSatus.Status == "moderator") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	status := "participant_request"
	if decision == "subscribe" {
		status = "subscriber"
	}

	err = ch.clubsUcase.SetUserStatusByClubID(int64(clubID), int64(userID), status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if status == "participant_request" {
		club, err := ch.clubsUcase.GetClubByID(clubID, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
			return
		}

		clubUrl := "https://vk.com/app8099557"
		err = ch.vk.CreatMessage(int(club.Owner.VKID),
			fmt.Sprintf("Привет! Новый пользователь хочет поучаствовать в %s: %s\n", club.Name, clubUrl),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// LeaveClub godoc
// @Summary      leave club
// @Description  Handler for leaving club
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/leave [post]
func (ch *ClubsHandler) LeaveClub(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	userClubSatus, err := ch.clubsUcase.GetUserStatusInClub(int64(clubID), int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if userClubSatus == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	err = ch.clubsUcase.DeleteUserFromClub(int64(clubID), int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ApproveRejectUserParticipateInClub godoc
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
func (ch *ClubsHandler) ApproveRejectUserParticipateInClub(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["cid"], 10, 64)
	userID, _ := strconv.ParseUint(vars["uid"], 10, 64)
	decision := vars["type"]

	ownerID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	userClubSatus, err := ch.clubsUcase.GetUserStatusInClub(int64(clubID), int64(ownerID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if userClubSatus == nil || userClubSatus.Status != "admin" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	err = ch.clubsUcase.ApproveRejectUserParticipateInClub(int64(clubID), int64(userID), decision)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	club, err := ch.clubsUcase.GetClubByID(clubID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	clubUrl := "https://vk.com/app8099557"

	msg := fmt.Sprintf("Привет! Администратор принял вас в %s: %s\n", club.Name, clubUrl)
	if decision == "reject" {
		msg = fmt.Sprintf("Привет! К сожалению, администратор отклонил ваш запрос на участие в %s: %s\n", club.Name, clubUrl)
	}

	err = ch.vk.CreatMessage(int(userID), msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetClubChatLink godoc
// @Summary      get club chat link
// @Description  Handler for getting tags list
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Success      200  {object}  models.ChatLink
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/chat_link [get]
func (ch *ClubsHandler) GetClubChatLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	userClubSatus, err := ch.clubsUcase.GetUserStatusInClub(int64(clubID), int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if userClubSatus == nil || (userClubSatus.Status != "admin" && userClubSatus.Status != "participant") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	chatID, err := ch.clubsUcase.GetClubChatID(int64(clubID), int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if chatID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "no chat for this club"}))
		return
	}

	chatLink, err := ch.vk.GetChatLink(int(chatID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(models.ChatLink{ChatLink: chatLink})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// DeleteClub godoc
// @Summary      delete club
// @Description  Handler for deleting club
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/delete [post]
func (ch *ClubsHandler) DeleteClub(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	club, err := ch.clubsUcase.GetClubByID(clubID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if club.Owner.VKID != userID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	err = ch.clubsUcase.DeleteClubByID(int64(clubID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}


// ComplainClub godoc
// @Summary      complain club
// @Description  Handler for complaining club
// @Tags         Clubs
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Club ID"
// @Param        body body models.ComplaintReq true "Club"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /clubs/{id}/complain [post]
func (ch *ClubsHandler) ComplainClub(w http.ResponseWriter, r *http.Request) {
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

	err = ch.clubsUcase.ComplainByID(models.Complaint{
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