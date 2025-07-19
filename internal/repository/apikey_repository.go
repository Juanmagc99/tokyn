package repository

import (
	"time"
	"tokyn/internal/models"

	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{
		db: db,
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
