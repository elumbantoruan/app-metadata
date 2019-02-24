package repository

import (
	"errors"

	"github.com/elumbantoruan/app-metadata/metadata"
)

// MetadataRepository defines an interface to store Metadata
type MetadataRepository interface {
	Create(appID string, data *metadata.ApplicationMetadata) error
	Update(appID string, data *metadata.ApplicationMetadata) error
	Get(appID string) (*metadata.ApplicationMetadata, error)
	GetAll() ([]metadata.ApplicationMetadata, error)
	Delete(appID string) error
}

var (
	ErrIDNotFound = errors.New("id not found")
)
