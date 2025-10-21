package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	errorspkg "reservations/errors"
	"reservations/models"
	"reservations/services"
	"reservations/utils"
)

// ReservationHandler holds service reference
type ReservationHandler struct {
	svc services.ReservationService
}

func NewReservationHandler(svc services.ReservationService) *ReservationHandler {
	return &ReservationHandler{svc: svc}
}

func (h *ReservationHandler) Register(r *mux.Router) {
	r.HandleFunc("/reservations", h.create).Methods("POST")
	r.HandleFunc("/reservations/{id}", h.get).Methods("GET")
	r.HandleFunc("/reservations/{id}", h.update).Methods("PUT")
	r.HandleFunc("/reservations/{id}", h.delete).Methods("DELETE")
}

func (h *ReservationHandler) create(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var r models.Reservation
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if r.OwnerID == "" {
		utils.WriteError(w, http.StatusBadRequest, "ownerId required")
		return
	}
	id, err := h.svc.Create(ctx, &r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *ReservationHandler) get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := mux.Vars(req)["id"]
	r, err := h.svc.Get(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == errorspkg.ErrNotFound {
			utils.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, r)
}

func (h *ReservationHandler) update(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := mux.Vars(req)["id"]
	var payload models.Reservation
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	if err := h.svc.Update(ctx, id, &payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (h *ReservationHandler) delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := mux.Vars(req)["id"]
	if err := h.svc.Delete(ctx, id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"id": id})
}
