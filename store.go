package main

import (
	"errors"
	"sync"
	"time"
)

type Store struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewStore() *Store {
	return &Store{orders: make(map[string]Order)}
}

func (s *Store) Save(o Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[o.ID] = o
}

func (s *Store) Get(id string) (Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[id]
	if !ok {
		return Order{}, errors.New("not found")
	}
	return o, nil
}

func NewOrderFromRequest(req OrderRequest, id string) Order {
	return Order{
		ID:         id,
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Status:     "created",
		CreatedAt:  time.Now().UTC(),
	}
}
