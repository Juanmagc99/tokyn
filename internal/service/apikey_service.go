package service

import (
	"context"
	"tokyn/internal/api/dto"
	"tokyn/internal/models"
	"tokyn/internal/repository"
	"tokyn/pkg/generator"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

	if err := s.akrepo.Create(&apikey); err != nil {
		return nil, err
	}

	redisData := models.APIKeyRedis{
		ID:      apikey.ID,
		KeyHash: apikey.KeyHash,
		Name:    apikey.Name,
	}

	if err := s.akrepo.InsertRedis(context.Background(), redisData); err != nil {

	}

	akr := &dto.APIKeyCreateResponse{
		ID:    apikey.ID,
		Token: tok,
		Name:  apikey.Name,
	}

	return akr, nil
}

func (s *APIKeyService) Revoke(id string) error {
	apikey, err := s.akrepo.Revoke(id)
	if err != nil {
		return err
	}

	if err := s.akrepo.DeleteRedis(context.Background(), apikey.KeyHash); err != nil {
		return err
	}

	return nil
}

func (s *APIKeyService) CheckToken(ctx context.Context, token string) (*dto.APIKeyCreateResponse, error) {
	hashedToken := generator.GetHash(token)

	apikeyRedis, err := s.akrepo.FindByHashRedis(ctx, hashedToken)
	if err == nil {
		return &dto.APIKeyCreateResponse{
			ID:    apikeyRedis.ID,
			Token: token,
			Name:  apikeyRedis.Name,
		}, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	apikey, err := s.akrepo.FindByHash(hashedToken)
	if err != nil {
		return nil, err
	}

	err = s.akrepo.InsertRedis(ctx, models.APIKeyRedis{
		ID:      apikey.ID,
		KeyHash: apikey.KeyHash,
		Name:    apikey.Name,
	})
	if err != nil {
		//Silenced error ;P
	}

	return &dto.APIKeyCreateResponse{
		ID:    apikey.ID,
		Token: token,
		Name:  apikey.Name,
	}, nil
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
