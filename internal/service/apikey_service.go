package service

import (
	"context"
	"log"
	"time"
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

func (s *APIKeyService) Create(akc dto.APIKeyCreate) (*dto.APIKeyCreateResponse, error) {
	var apikey models.APIKey

	tok, hash, err := generator.GenerateToken()
	if err != nil {
		return nil, err
	}

	apikey.ID = uuid.NewString()
	apikey.Name = akc.Name
	apikey.KeyHash = hash

	if akc.Hours > 0 {
		exp := time.Now().Add(time.Hour * time.Duration(akc.Hours))
		apikey.ExpiresAt = &exp
	}

	if err := s.akrepo.Create(&apikey); err != nil {
		return nil, err
	}

	s.cacheAPIKey(context.Background(), &apikey)

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

func (s *APIKeyService) Delete(id string) error {
	apikey, err := s.akrepo.Delete(id)
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

	if !apikey.IsValid() {
		_, err = s.akrepo.Revoke(apikey.ID)
		if err != nil {
			return nil, err
		}

		err = s.akrepo.DeleteRedis(ctx, apikey.KeyHash)
		if err != nil {
			return nil, err
		}

		return nil, models.ErrAPIKeyInvalid
	}

	s.cacheAPIKey(ctx, apikey)

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

func (s *APIKeyService) Details(id string) (*dto.APIKeyDetailsResponse, error) {
	apikey, err := s.akrepo.Details(id)
	if err != nil {
		return nil, err
	}

	apikeyDetails := &dto.APIKeyDetailsResponse{
		ID:        apikey.ID,
		Name:      apikey.Name,
		CreatedAt: apikey.CreatedAt,
		Revoked:   apikey.Revoked,
		RevokedAt: apikey.RevokedAt,
		ExpiresAt: apikey.ExpiresAt,
	}

	return apikeyDetails, nil
}

func (s *APIKeyService) cacheAPIKey(ctx context.Context, ak *models.APIKey) {
	go func() {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		data := models.APIKeyRedis{
			ID:      ak.ID,
			KeyHash: ak.KeyHash,
			Name:    ak.Name,
		}

		if err := s.akrepo.InsertRedis(ctx, data); err != nil {
			log.Printf("warning: failed to cache API key in Redis: %v", err)
		}
	}()
}
