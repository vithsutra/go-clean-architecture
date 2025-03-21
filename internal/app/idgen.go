package app

import "id-generator/internal/domain"

type IDService struct {
	generator domain.IDGenerator
}

func NewIDService(gen domain.IDGenerator) *IDService {
	return &IDService{generator: gen}
}

func (s *IDService) GetID() (string, error) {
	return s.generator.GenerateID()
}
