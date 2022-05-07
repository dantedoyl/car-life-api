package delivery

import (
	"encoding/json"
	"github.com/dantedoyl/car-life-api/internal/app/events_posts"
	"github.com/dantedoyl/car-life-api/internal/app/middleware"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"github.com/dantedoyl/car-life-api/internal/app/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

type EventsPostsHandler struct {
	eventsUcase events_posts.IEventsPostsUsecase
}

func NewEventsPostsHandler(eventsUcase events_posts.IEventsPostsUsecase) *EventsPostsHandler {
	return &EventsPostsHandler{
		eventsUcase: eventsUcase,
	}
}

func (eph *EventsPostsHandler) Configure(r *mux.Router, mw *middleware.Middleware) {
	r.HandleFunc("/event_posts/{event_id:[0-9]+}/create", mw.CheckAuthMiddleware(eph.CreateEventPost)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/event_posts/{event_id:[0-9]+}", mw.CheckAuthMiddleware(eph.GetEventsPostsByEventID)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/events_posts/{post_id:[0-9]+}/upload", mw.CheckAuthMiddleware(eph.UploadAttachments)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/event_posts/{post_id:[0-9]+}/delete", mw.CheckAuthMiddleware(eph.DeletePost)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/event_posts/{post_id:[0-9]+}/complain", mw.CheckAuthMiddleware(eph.ComplainPost)).Methods(http.MethodPost, http.MethodOptions)

}

// CreateEventPost godoc
// @Summary      create an event post
// @Description  Handler for creating an event post
// @Tags         EventsPosts
// @Accept       json
// @Produce      json
// @Param        event_id path int64 true "Event ID"
// @Param        body body models.CreatePostRequest true "EventPost"
// @Success      200  {object}  models.EventPost
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /event_posts/{event_id}/create [post]
func (eph *EventsPostsHandler) CreateEventPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["event_id"], 10, 64)
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	event := &models.CreatePostRequest{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "can't unmarshal data"}))
		return
	}

	eventsData := &models.EventPost{
		Text: event.Text,
		User: models.UserCard{
			VKID: userID,
		},
		EventID: eventID,
	}

	err = eph.eventsUcase.CreateEventPost(eventsData)
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

// GetEventsPostsByEventID godoc
// @Summary      get events posts list
// @Description  Handler for getting events posts list
// @Tags         EventsPosts
// @Accept       json
// @Produce      json
// @Param        event_id path int64 true "Event ID"
// @Param        IdGt query integer false "IdGt"
// @Param        IdLte query integer false "IdLte"
// @Param        Limit query integer false "Limit"
// @Success      200  {object}  []models.EventPost
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /event_posts/{event_id} [get]
func (eph *EventsPostsHandler) GetEventsPostsByEventID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, _ := strconv.ParseUint(vars["event_id"], 10, 64)

	query := &models.EventQuery{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(query, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	events, err := eph.eventsUcase.GetEventsPostsByEventID(eventID, query.IdGt, query.IdLte, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}
	if len(events) == 0 {
		events = []*models.EventPost{}
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

// UploadAttachments godoc
// @Summary      upload attachments for event post
// @Description  Handler for uploading an event post attachment
// @Tags         EventsPosts
// @Accept       mpfd
// @Produce      json
// @Param        post_id path int64 true "Post ID"
// @Success      200  {object}  models.EventPost
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /events_posts/{post_id}/upload [post]
func (eph *EventsPostsHandler) UploadAttachments(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	postID, _ := strconv.ParseUint(vars["post_id"], 10, 64)

	_, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	err := r.ParseMultipartForm(10 * 1024 * 1024)
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

	file := r.MultipartForm.File["file-upload"]
	event, err := eph.eventsUcase.UploadAttachments(postID, file)
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

// DeletePost godoc
// @Summary      delete post
// @Description  Handler for deleting post
// @Tags         EventsPosts
// @Accept       json
// @Produce      json
// @Param        post_id path int64 true "Post ID"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /event_posts/{post_id}/delete [post]
func (eph *EventsPostsHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.ParseUint(vars["post_id"], 10, 64)

	userID, ok := r.Context().Value("userID").(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(utils.JSONError(&utils.Error{Message: "you're unauthorized"}))
		return
	}

	post, err := eph.eventsUcase.GetEventPostByPostID(postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	if post.User.VKID != userID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.JSONError(&utils.Error{Message: "user has inappropriate status"}))
		return
	}

	err = eph.eventsUcase.DeletePostByID(int64(postID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.JSONError(&utils.Error{Message: err.Error()}))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ComplainPost godoc
// @Summary      complain post
// @Description  Handler for complaining post
// @Tags         EventsPosts
// @Accept       json
// @Produce      json
// @Param        post_id path int64 true "Post ID"
// @Param        body body models.ComplaintReq true "Event"
// @Success      200
// @Failure      400  {object}  utils.Error
// @Failure      401  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Failure      500  {object}  utils.Error
// @Router       /event_posts/{post_id}/complain [post]
func (eph *EventsPostsHandler) ComplainPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubID, _ := strconv.ParseUint(vars["post_id"], 10, 64)

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

	err = eph.eventsUcase.ComplainByID(models.Complaint{
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
