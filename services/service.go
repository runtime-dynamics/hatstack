package services

import (
	"context"
)

// Service is the base interface that all services should implement
// for common functionality like context handling
type Service interface {
	// Context returns the service's context
	Context() context.Context
}

// BaseService provides common functionality for all services
type BaseService struct {
	ctx context.Context
}

// NewBaseService creates a new base service with the given context
func NewBaseService(ctx context.Context) *BaseService {
	return &BaseService{ctx: ctx}
}

// Context returns the service's context
func (s *BaseService) Context() context.Context {
	return s.ctx
}
