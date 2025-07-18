package repository

import (
	"tokyn/internal/models"

	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func (repo *APIKeyRepository) Create(ak *models.APIKey) error {
	if err := repo.db.Create(ak).Error; err != nil {
		return err
	}

	return nil
}

func (repo *APIKeyRepository) Update(ak *models.APIKey) error {

}

func (repo *APIKeyRepository) FindByHash(hash string) (*models.APIKey, error) {

}
