package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"tokyn/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db      *gorm.DB
	rclient *redis.Client
}

func NewAPIKeyRepository(db *gorm.DB, rclient *redis.Client) *APIKeyRepository {
	return &APIKeyRepository{
		db:      db,
		rclient: rclient,
	}
}

func (repo *APIKeyRepository) Create(ak *models.APIKey) error {
	if err := repo.db.Create(ak).Error; err != nil {
		return err
	}

	return nil
}

func (repo *APIKeyRepository) Update(ak *models.APIKey) error {
	if err := repo.db.Save(ak).Error; err != nil {
		return err
	}

	return nil
}

func (repo *APIKeyRepository) FindByHash(hash string) (*models.APIKey, error) {
	var ak models.APIKey
	if err := repo.db.Where("key_hash = ? AND revoked = false", hash).First(&ak).Error; err != nil {
		return nil, err
	}

	return &ak, nil
}

func (repo *APIKeyRepository) FindByID(id string) (*models.APIKey, error) {
	var ak models.APIKey
	if err := repo.db.Where("id = ? AND revoked = false", id).First(&ak).Error; err != nil {
		return nil, err
	}

	return &ak, nil
}

func (repo *APIKeyRepository) Revoke(id string) error {
	var ak models.APIKey
	if err := repo.db.First(&ak, "id = ?", id).Error; err != nil {
		return err
	}

	now := time.Now()
	ak.Revoked = true
	ak.RevokedAt = &now

	if err := repo.db.Save(&ak).Error; err != nil {
		return err
	}

	return nil
}

func (repo *APIKeyRepository) InsertRedis(ctx context.Context, keyData models.APIKeyRedis) error {
	data, err := json.Marshal(keyData)
	if err != nil {
		return err
	}

	err = repo.rclient.Set(ctx, keyData.KeyHash, data, time.Hour*24).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *APIKeyRepository) FindByHashRedis(ctx context.Context, hash string) (*models.APIKeyRedis, error) {
	val, err := repo.rclient.Get(ctx, hash).Result()
	if err != nil {
		return nil, err
	}

	var apiKey models.APIKeyRedis
	if err := json.Unmarshal([]byte(val), &apiKey); err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (repo *APIKeyRepository) AllowRequest(ctx context.Context, hash string) (bool, error) {
	counterKey := fmt.Sprintf("rate:%s", hash)

	count, err := repo.rclient.Incr(ctx, counterKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = repo.rclient.Expire(ctx, counterKey, time.Hour*24*7).Err()
		if err != nil {
			return false, err
		}
	}

	if count > 100 {
		return false, nil
	}

	return true, nil
}
