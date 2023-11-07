package metadata

import (
	"github.com/Saffica/image-storage/pkg/models"
)

type metaDataRepository struct {
}

func New() *metaDataRepository {
	return &metaDataRepository{}
}

func (r *metaDataRepository) Get(url string) (*models.MetaData, error) {
	return nil, nil
}

func (r *metaDataRepository) Insert(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, nil
}

func (r *metaDataRepository) Update(metaData *models.MetaData) (*models.MetaData, error) {
	return nil, nil
}
