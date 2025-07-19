package service

import (
	"tokyn/internal/api/dto"
	"tokyn/internal/models"
	"tokyn/internal/repository"
	"tokyn/pkg/generator"

	"github.com/google/uuid"
)

type APIKeyService struct {
	akrepo repository.APIKeyRepository
}

func (s *APIKeyService) Create(name string) (*dto.APIKeyCreateResponse, error) {
	var apikey models.APIKey

	tok, hash, err := generator.GenerateToken()
	if err != nil {
		return nil, err
	}

	apikey.ID = uuid.NewString()
	apikey.Name = name
	apikey.KeyHash = hash
	err = s.akrepo.Create(&apikey)
	if err != nil {
		return nil, err
	}

	akr := &dto.APIKeyCreateResponse{
		ID:    apikey.ID,
		Token: tok,
		Name:  apikey.Name,
	}

	return akr, nil
}
