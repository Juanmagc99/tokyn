package service

import (
	"context"
	"tokyn/internal/api/dto"
	"tokyn/internal/models"
	"tokyn/internal/repository"
	"tokyn/pkg/generator"

	"github.com/google/uuid"
)

type APIKeyService struct {
	akrepo repository.APIKeyRepository
}

func NewAPIKeyService(akr repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{
		akrepo: akr,
	}
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

func (s *APIKeyService) Revoke(id string) error {
	err := s.akrepo.Revoke(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *APIKeyService) CheckToken(token string) (*dto.APIKeyCreateResponse, error) {

	hashedToken := generator.GetHash(token)
	apikey, err := s.akrepo.FindByHash(hashedToken)
	if err != nil {
		return nil, err
	}

	akr := &dto.APIKeyCreateResponse{
		ID:    apikey.ID,
		Token: token,
		Name:  apikey.Name,
	}

	return akr, nil
}

func (s *APIKeyService) CheckTokenRedis(ctx context.Context, token string) (*dto.APIKeyCreateResponse, error) {

	hashedToken := generator.GetHash(token)
	apikey, err := s.akrepo.FindByHashRedis(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	akr := &dto.APIKeyCreateResponse{
		ID:    apikey.ID,
		Token: token,
		Name:  apikey.Name,
	}

	return akr, nil
}

func (s *APIKeyService) CheckRateLimit(ctx context.Context, token string) (bool, error) {
	ok, err := s.akrepo.AllowRequest(ctx, token)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}
