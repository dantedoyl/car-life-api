package delivery

import (
	"encoding/json"
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
}

func NewEventsHandler(eventsUcase events.IEventsUsecase) *EventsHandler {
	return &EventsHandler{
		eventsUcase: eventsUcase,
	}
}

func (eh *EventsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/event/create", mw.CheckAuthMiddleware(eh.CreateEvent)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}", mw.CheckAuthMiddleware(eh.GetEventByID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events", mw.CheckAuthMiddleware(eh.GetEvents)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/upload", mw.CheckAuthMiddleware(eh.UploadAvatarHandler)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/participants", mw.CheckAuthMiddleware(eh.GetEventsParticipants)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/participants/requests", mw.CheckAuthMiddleware(eh.GetEventsParticipantsRequests)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events/{id:[0-9]+}/participate", mw.CheckAuthMiddleware(eh.ParticipateByEventID)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/events/{eid:[0-9]+}/participate/{uid:[0-9]+}/{type:approve|reject}", mw.CheckAuthMiddleware(eh.ApproveRejectUserParticipateInEvent)).Methods(http.MethodPost, http.MethodOptions)
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

	event := &models.CreateEventRequest{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
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
	}

	err = eh.eventsUcase.CreateEvent(eventsData)
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

	events, err := eh.eventsUcase.GetEvents(query.IdGt, query.IdLte, query.Limit, query.Query)
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
			ID:        event.ID,
			Name:      event.Name,
			EventDate: event.EventDate,
			AvatarUrl: event.AvatarUrl,
			Latitude:  event.Latitude,
			Longitude: event.Longitude,
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

	event, err := eh.eventsUcase.GetEventByID(eventID)
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
// @Success      200  {object}  models.Event
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/upload [post]
func (eh *EventsHandler) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	eventID, _ := strconv.ParseInt(vars["id"], 10, 64)

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

// GetEventsParticipants godoc
// @Summary      get events participants list
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.UserCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/participants [get]
func (eh *EventsHandler) GetEventsParticipants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	users, err := eh.eventsUcase.GetEventsUserByStatus(int64(eventID), "participant", query.IdGt, query.IdLte, query.Limit)
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

// GetEventsParticipantsRequests godoc
// @Summary      get events participants requestlist
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.UserCard
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/participants/requests [get]
func (eh *EventsHandler) GetEventsParticipantsRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)

	query := &models.ClubQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	users, err := eh.eventsUcase.GetEventsUserByStatus(int64(eventID), "participant_request", query.IdGt, query.IdLte, query.Limit)
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

// ParticipateByEventID godoc
// @Summary      request participate in event
// @Description  Handler for getting tags list
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Event ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events/{id}/participate [post]
func (eh *EventsHandler) ParticipateByEventID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	err := eh.eventsUcase.SetUserStatusByEventID(int64(eventID), int64(userID), "participant_request")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
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

	w.WriteHeader(http.StatusOK)
}
