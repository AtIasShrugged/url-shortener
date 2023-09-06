package shortening

import (
	"context"
	"sync"
	"time"
	"url-shortener/internal/model"
)

type InMemory struct {
	m sync.Map
}

func NewInMemory() *InMemory {
	return &InMemory{}
}

func (s *InMemory) Put(_ context.Context, shortening model.Shortening) (*model.Shortening, error) {
	if _, exists := s.m.Load(shortening.Identifier); exists {
		return nil, model.ErrIdentifierExists
	}

	shortening.CreatedAt = time.Now().UTC()
	s.m.Store(shortening.Identifier, shortening)

	return &shortening, nil
}

func (s *InMemory) Get(_ context.Context, identifier string) (*model.Shortening, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, model.ErrNotFound
	}

	shortening := v.(model.Shortening)

	return &shortening, nil
}

func (s *InMemory) IncrementVisits(_ context.Context, identifier string) error {
	v, ok := s.m.Load(identifier)
	if !ok {
		return model.ErrNotFound
	}

	shortening := v.(model.Shortening)
	shortening.Visits++
	s.m.Store(identifier, shortening)

	return nil
}
