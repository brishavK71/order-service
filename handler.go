package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type App struct {
	Store    *Store
	Producer *KafkaProducer
	Validate *validator.Validate
}

func NewApp(store *Store, producer *KafkaProducer) *App {
	return &App{
		Store:    store,
		Producer: producer,
		Validate: validator.New(),
	}
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middlewareLogger)
	r.Post("/orders", a.CreateOrder)
	r.Get("/orders/{id}", a.GetOrder)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	return r
}

func (a *App) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := a.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	order := NewOrderFromRequest(req, id)

	// Save locally
	a.Store.Save(order)

	// Publish event
	evt := OrderEvent{
		EventType: "order.created",
		OrderID:   id,
		Payload:   order,
		Time:      order.CreatedAt,
	}

	if err := a.Producer.PublishOrderCreated(r.Context(), evt); err != nil {
		// log error but still let client know. You may choose to roll back storage in real app.
		http.Error(w, "failed to publish event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (a *App) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	o, err := a.Store.Get(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}
