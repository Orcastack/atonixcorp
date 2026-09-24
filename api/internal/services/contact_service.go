package services

import (
	"atonixcorp/api/internal/models"
	"atonixcorp/api/internal/repositories"
)

type ContactService interface {
	Submit(msg *models.ContactMessage) error
	GetAll() ([]models.ContactMessage, error)
}

type contactService struct {
	repo repositories.ContactRepository
}

func NewContactService(repo repositories.ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) Submit(msg *models.ContactMessage) error {
	return s.repo.Save(msg)
}

func (s *contactService) GetAll() ([]models.ContactMessage, error) {
	return s.repo.GetAll()
}
