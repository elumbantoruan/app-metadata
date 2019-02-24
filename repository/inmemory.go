package repository

import (
	"github.com/elumbantoruan/app-metadata/metadata"
)

// InMemoryMetadataRepository is a concrete implementation of MetadataRepository interface in memory
type InMemoryMetadataRepository struct {
	Storage map[string]*metadata.ApplicationMetadata
}

// NewInMemoryMetadataRepository creates a new instance of InMemoryMetadataRepository
func NewInMemoryMetadataRepository() MetadataRepository {
	data := make(map[string]*metadata.ApplicationMetadata)
	return &InMemoryMetadataRepository{
		Storage: data,
	}
}

// Create adds an application metadata into a repository
func (im *InMemoryMetadataRepository) Create(appID string, data *metadata.ApplicationMetadata) error {
	im.Storage[appID] = data
	return nil
}

// Update updates the application metadata for a given appID
func (im *InMemoryMetadataRepository) Update(appID string, data *metadata.ApplicationMetadata) error {
	if _, ok := im.Storage[appID]; !ok {
		return ErrIDNotFound
	}
	im.Storage[appID] = data
	return nil
}

// Get returns application metadata for a given appID
func (im *InMemoryMetadataRepository) Get(appID string) (*metadata.ApplicationMetadata, error) {
	ret, ok := im.Storage[appID]
	if ok {
		return ret, nil
	}
	return nil, nil
}

// GetAll returns all application metadata
func (im *InMemoryMetadataRepository) GetAll() ([]metadata.ApplicationMetadata, error) {
	var results []metadata.ApplicationMetadata
	for k, v := range im.Storage {
		v.ApplicationID = k
		results = append(results, *v)
	}

	return results, nil
}

// Delete removes the application metadata for a given an appID
func (im *InMemoryMetadataRepository) Delete(appID string) error {
	if _, ok := im.Storage[appID]; !ok {
		return ErrIDNotFound
	}
	delete(im.Storage, appID)
	return nil
}
