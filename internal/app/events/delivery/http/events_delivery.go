package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/dantedoyl/car-life-api/internal/app/clients/vk"
	clubs "github.com/dantedoyl/car-life-api/internal/app/clubs"
	"github.com/dantedoyl/car-life-api/internal/app/events"
	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

type EventsHandler struct {
	eventsUcase events.IEventsUsecase
	clubUcase   clubs.IClubsUsecase
	vk          *vk.VKClient
}

func NewEventsHandler(eventsUcase events.IEventsUsecase, clubUcase clubs.IClubsUsecase, vk *vk.VKClient) *EventsHandler {
	return &EventsHandler{
		eventsUcase: eventsUcase,
		clubUcase:   clubUcase,
		vk:          vk,
	}
}

func (eh *EventsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/event/create", mw.CheckAuthMiddleware(eh.CreateEvent)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}", mw.CheckAuthMiddleware(eh.GetEventByID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events", mw.CheckAuthMiddleware(eh.GetEvents)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/upload", mw.CheckAuthMiddleware(eh.UploadAvatarHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/{type:participant|participant_request|spectator}", mw.CheckAuthMiddleware(eh.GetEventsUsersByType)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/{type:participate|spectate}", mw.CheckAuthMiddleware(eh.SetUserStatusByEventID)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{eid:[0-9]+}/participate/{uid:[0-9]+}/{type:approve|reject}", mw.CheckAuthMiddleware(eh.ApproveRejectUserParticipateInEvent)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/chat_link", mw.CheckAuthMiddleware(eh.GetEventChatLink)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/leave", mw.CheckAuthMiddleware(eh.LeaveEvent)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/delete", mw.CheckAuthMiddleware(eh.DeleteEvent)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/complain", mw.CheckAuthMiddleware(eh.ComplainEvent)).Methods(http.MethodPost, http.MethodOptions)
}

// CreateEvent godoc
// @Summary      create an event
// @Description  Handler for creating an event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        body body models.CreateEventRequest true "Event"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /event/create [post]
func (eh *EventsHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	event := &models.CreateEventRequest{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	userClubSatus, err := eh.clubUcase.GetUserStatusInClub(int64(event.ClubID), int64(userID))
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

	eventsData := &models.Event{
		Name:        event.Name,
		Club:        models.Club{ID: event.ClubID},
		Description: event.Description,
		EventDate:   event.EventDate,
		Latitude:    event.Latitude,
		Longitude:   event.Longitude,
		AvatarUrl:   event.AvatarUrl,
		Creator:   models.UserCard{
			VKID:      userID,
		},
	}

	err = eh.eventsUcase.CreateEvent(eventsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	id, err := eh.vk.CreatChat(eventsData.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	err = eh.eventsUcase.SetEventChatID(int64(eventsData.ID), int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(eventsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetEvents godoc
// @Summary      get events list
// @Description  Handler for getting events list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Param        Query query string false "Query"
// @Param        UpperRightLatitude query number false "UpperRightLatitude"
// @Param        UpperRightLongitude query number false "UpperRightLongitude"
// @Param        DownLeftLatitude query number false "DownLeftLatitude"
// @Param        DownLeftLongitude query number false "DownLeftLongitude"
// @Success      200  {object}  []models.EventCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events [get]
func (eh *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	query := &models.EventQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	events, err := eh.eventsUcase.GetEvents(query.IdGt, query.IdLte, query.Limit, query.Query, query.DownLeftLongitude, query.DownLeftLatitude, query.UpperRightLongitude, query.UpperRightLatitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(events) == 0 {
		events = []*models.Event{}
	}

	eventCards := make([]models.EventCard, 0, len(events))
	for _, event := range events {
		eventCards = append(eventCards, models.EventCard{
			ID:                event.ID,
			Name:              event.Name,
			EventDate:         event.EventDate,
			AvatarUrl:         event.AvatarUrl,
			Latitude:          event.Latitude,
			Longitude:         event.Longitude,
			ParticipantsCount: event.ParticipantsCount,
			SpectatorsCount:   event.SpectatorsCount,
		})
	}

	body, err := json.Marshal(eventCards)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetEventByID godoc
// @Summary      get event by id
// @Description  Handler for creating an event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id} [get]
func (eh *EventsHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)

	event := &models.Event{}
	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		event.UserStatus = "unknown"
	}

	event, err := eh.eventsUcase.GetEventByID(eventID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
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

// UploadAvatarHandler godoc
// @Summary      upload avatar for event
// @Description  Handler for creating an event
// @Tags         Events
// @Accept       mpfd
// @Produce      json
// @Param        id path int64 true "Account ID"
// @Param 		 file-upload formData file true "Image to upload"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/upload [post]
func (eh *EventsHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	eventID, _ := strconv.ParseInt(vars["id"], 10, 64)

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
	event, err := eh.eventsUcase.UpdateAvatar(eventID, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	chatID, err := eh.eventsUcase.GetEventChatID(eventID, int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	err = eh.vk.UploadChatPhoto(int(chatID), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
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

// GetEventsUsersByType godoc
// @Summary      get events users list
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Param        type path string true "Type" Enums(participant, participant_request, spectator)
// @Success      200  {object}  []models.UserCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/{type} [get]
func (eh *EventsHandler) GetEventsUsersByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)
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

	users, err := eh.eventsUcase.GetEventsUserByStatus(int64(eventID), role, query.IdGt, query.IdLte, query.Limit)
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

// SetUserStatusByEventID godoc
// @Summary      set user role in event
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Param        type path string true "Type" Enums(participate, spectate)
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/{type} [post]
func (eh *EventsHandler) SetUserStatusByEventID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)
	decision := vars["type"]

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	status := "participant_request"
	if decision == "spectate" {
		status = "spectator"
	}

	err := eh.eventsUcase.SetUserStatusByEventID(int64(eventID), int64(userID), status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if status == "participant_request" {
		event, err := eh.eventsUcase.GetEventByID(eventID, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
			return
		}

		eventUrl := "https://vk.com/app8099557"
		err = eh.vk.CreatMessage(int(event.Creator.VKID),
			fmt.Sprintf("????????????! ?????????? ???????????????? ?????????? ?????????????????????????? ?? %s: %s", event.Name, eventUrl),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// ApproveRejectUserParticipateInEvent godoc
// @Summary      approve/reject participate in event
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eid path int64 true "Event ID"
// @Param        uid path int64 true "User ID"
// @Param        type path string true "Type" Enums(approve, reject)
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{eid}/participate/{uid}/{type} [post]
func (eh *EventsHandler) ApproveRejectUserParticipateInEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["eid"], 10, 64)
	userID, _ := strconv.ParseUint(vars["uid"], 10, 64)
	decision := vars["type"]

	_, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	err := eh.eventsUcase.ApproveRejectUserParticipateInEvent(int64(eventID), int64(userID), decision)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	event, err := eh.eventsUcase.GetEventByID(eventID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	eventUrl := "https://vk.com/app8099557"

	msg := fmt.Sprintf("????????????! ?????????????????????????? ???????????? ?????? ?? %s: %s\n", event.Name, eventUrl)
	if decision == "reject" {
		msg = fmt.Sprintf("????????????! ?? ??????????????????, ?????????????????????????? ???????????????? ?????? ???????????? ???? ?????????????? ?? %s: %s\n", event.Name, eventUrl)
	}

	err = eh.vk.CreatMessage(int(userID), msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetEventChatLink godoc
// @Summary      get event chat link
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Success      200  {object}  models.ChatLink
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/chat_link [get]
func (eh *EventsHandler) GetEventChatLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	chatID, err := eh.eventsUcase.GetEventChatID(int64(eventID), int64(userID))
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

	chatLink, err := eh.vk.GetChatLink(int(chatID))
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

// LeaveEvent godoc
// @Summary      leave event
// @Description  Handler for leaving event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/leave [post]
func (eh *EventsHandler) LeaveEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	err := eh.eventsUcase.DeleteUserFromEvent(int64(clubID), int64(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteEvent godoc
// @Summary      delete event
// @Description  Handler for deleting event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/delete [post]
func (eh *EventsHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	event, err := eh.eventsUcase.GetEventByID(clubID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if event.Creator.VKID != userID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	err = eh.eventsUcase.DeleteEventByID(int64(clubID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ComplainEvent godoc
// @Summary      complain event
// @Description  Handler for complaining event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Param        body body models.ComplaintReq true "Event"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/complain [post]
func (ch *EventsHandler) ComplainEvent(w http.ResponseWriter, r *http.Request) {
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

	err = ch.eventsUcase.ComplainByID(models.Complaint{
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
