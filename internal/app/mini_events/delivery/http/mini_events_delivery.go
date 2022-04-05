package delivery

import (
	"encoding/json"
	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	"github.com/dantedoyl/car-life-api/internal/app/mini_events"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
	"time"
)

type MiniEventsHandler struct {
	miniEventsUcase mini_events.IMiniEventsUsecase
}

func NewMiniEventsHandler(eventsUcase mini_events.IMiniEventsUsecase) *MiniEventsHandler {
	return &MiniEventsHandler{
		miniEventsUcase: eventsUcase,
	}
}

func (mh *MiniEventsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/mini_event/create", mw.CheckAuthMiddleware(mh.CreateMiniEvent)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/mini_events/{id:[0-9]+}", mw.CheckAuthMiddleware(mh.GetMiniEventByID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/mini_events", mw.CheckAuthMiddleware(mh.GetMiniEvents)).Methods(http.MethodGet, http.MethodOptions)
}

// CreateMiniEvent godoc
// @Summary      create a mini event
// @Description  Handler for creating an event
// @Tags         MiniEvents
// @Accept       json
// @Produce      json
// @Param        body body models.CreateMiniEventRequest true "Event"
// @Success      200  {object}  models.MiniEvent
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /mini_event/create [post]
func (mh *MiniEventsHandler) CreateMiniEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	miniEvent := &models.CreateMiniEventRequest{}
	err := json.NewDecoder(r.Body).Decode(&miniEvent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	miniEventsData := &models.MiniEvent{
		Type:        models.MiniEventType{
			ID:               uint64(miniEvent.TypeID),
		},
		User:        models.UserCard{
			VKID:      userID,
		},
		Description: miniEvent.Description,
		CreatedAt:   time.Now().Local(),
		EndedAt:     miniEvent.EndedAt,
		Latitude:    miniEvent.Latitude,
		Longitude:   miniEvent.Longitude,
	}

	if miniEventsData.CreatedAt.After(miniEvent.EndedAt) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "ended time can't be before started time"}))
		return
	}

	err = mh.miniEventsUcase.CreateMiniEvent(miniEventsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	body, err := json.Marshal(miniEventsData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: "can't marshal data"}))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// GetMiniEvents godoc
// @Summary      get mini events list
// @Description  Handler for getting events list
// @Tags         MiniEvents
// @Accept       json
// @Produce      json
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.MiniEvent
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /mini_events [get]
func (mh *MiniEventsHandler) GetMiniEvents(w http.ResponseWriter, r *http.Request) {
	query := &models.EventQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	events, err := mh.miniEventsUcase.GetMiniEvents(query.IdGt, query.IdLte, query.Limit, query.Query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(events) == 0 {
		events = []*models.MiniEvent{}
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

// GetMiniEventByID godoc
// @Summary      get mini event by id
// @Description  Handler for creating an event
// @Tags         MiniEvents
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Mini Event ID"
// @Success      200  {object}  models.MiniEvent
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /mini_events/{id} [get]
func (mh *MiniEventsHandler) GetMiniEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	miniEventID, _ := strconv.ParseUint(vars["id"], 10, 64)

	event, err := mh.miniEventsUcase.GetMiniEventByID(miniEventID)
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