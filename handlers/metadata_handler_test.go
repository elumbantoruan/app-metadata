package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"gopkg.in/yaml.v2"

	"github.com/elumbantoruan/app-metadata/metadata"
	"github.com/elumbantoruan/app-metadata/repository"
	"github.com/stretchr/testify/assert"
)

func TestMetadataHandler_HandlePostMetadata_ResultedCreated(t *testing.T) {

	payload := createValidPayload()
	request, _ := http.NewRequest("POST", "app-metadata", strings.NewReader(payload))
	responseRecorder := httptest.NewRecorder()

	mr := repository.NewInMemoryMetadataRepository()
	mh := NewMetadataHandler(mr)
	mh.HandlePostMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
}

func TestMetadataHandler_HandlePostMetadataBadYamlFormat_ResultedBadRequest(t *testing.T) {

	payload := createInvalidPayloadBadYamlFormat()
	request, _ := http.NewRequest("POST", "app-metadata", strings.NewReader(payload))
	responseRecorder := httptest.NewRecorder()

	mr := repository.NewInMemoryMetadataRepository()
	mh := NewMetadataHandler(mr)
	mh.HandlePostMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestMetadataHandler_HandlePostMetadataRepositoryError_ResultedInternalServerError(t *testing.T) {

	payload := createValidPayload()
	request, _ := http.NewRequest("POST", "app-metadata", strings.NewReader(payload))
	responseRecorder := httptest.NewRecorder()

	mr := NewFakeMetadataRepository()
	mh := NewMetadataHandler(mr)
	mh.HandlePostMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestMetadataHandler_HandlePutMetadata_ResultedOK(t *testing.T) {

	im := repository.NewInMemoryMetadataRepository()
	payload := createValidPayload()
	appID := "appID1"
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)

	mtd.ApplicationID = appID
	im.Create(appID, &mtd)

	mtd.Company = "updated company"

	updatedPayload, _ := yaml.Marshal(mtd)

	request, _ := http.NewRequest("PUT", "app-metadata/appID1", strings.NewReader(string(updatedPayload)))
	responseRecorder := httptest.NewRecorder()
	mapper := map[string]string{
		"appID": "appID1",
	}
	request = mux.SetURLVars(request, mapper)

	mh := NewMetadataHandler(im)
	mh.HandlePutMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestMetadataHandler_HandlePutMetadata_ResultedInternalStatusError(t *testing.T) {

	mr := NewFakeMetadataRepository()
	payload := createValidPayload()
	appID := "appID1"
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)

	mtd.ApplicationID = appID
	mr.Create(appID, &mtd)

	mtd.Company = "updated company"

	updatedPayload, _ := yaml.Marshal(mtd)

	request, _ := http.NewRequest("PUT", "app-metadata/appID1", strings.NewReader(string(updatedPayload)))
	responseRecorder := httptest.NewRecorder()
	mapper := map[string]string{
		"appID": "appID1",
	}
	request = mux.SetURLVars(request, mapper)

	// mr is a FakeMetadataRepository that will return an error
	mh := NewMetadataHandler(mr)
	mh.HandlePutMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
}

func TestMetadataHandler_HandlePutMetadata_ResultedBadRequest(t *testing.T) {

	payload := createValidPayload()
	request, _ := http.NewRequest("PUT", "app-metadata", strings.NewReader(payload))
	responseRecorder := httptest.NewRecorder()

	im := repository.NewInMemoryMetadataRepository()
	mh := NewMetadataHandler(im)
	mh.HandlePutMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestMetadataHandler_HandleGetMetadata_ResultedOK(t *testing.T) {

	im := repository.NewInMemoryMetadataRepository()
	payload := createValidPayload()
	appID := "appID1"
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)

	mtd.ApplicationID = appID
	im.Create(appID, &mtd)
	request, _ := http.NewRequest("GET", "app-metadata/appID1", strings.NewReader(""))
	mapper := map[string]string{
		"appID": "appID1",
	}
	request = mux.SetURLVars(request, mapper)
	responseRecorder := httptest.NewRecorder()

	mh := NewMetadataHandler(im)
	mh.HandleGetMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	yaml.NewDecoder(responseRecorder.Body).Decode(&mtd)

	assert.Equal(t, appID, mtd.ApplicationID)
}

func TestMetadataHandler_HandleGetMetadata_ResultedNotFound(t *testing.T) {

	im := repository.NewInMemoryMetadataRepository()
	payload := createValidPayload()
	appID := "appID1"
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)

	mtd.ApplicationID = appID
	im.Create(appID, &mtd)
	request, _ := http.NewRequest("GET", "app-metadata/notfound", strings.NewReader(""))
	mapper := map[string]string{
		"appID": "notfound",
	}
	request = mux.SetURLVars(request, mapper)
	responseRecorder := httptest.NewRecorder()

	mh := NewMetadataHandler(im)
	mh.HandleGetMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}

func TestMetadataHandler_HandleGetMetadata_ResultedBadRequest(t *testing.T) {

	payload := createValidPayload()
	request, _ := http.NewRequest("GET", "app-metadata", strings.NewReader(payload))
	responseRecorder := httptest.NewRecorder()

	im := repository.NewInMemoryMetadataRepository()
	mh := NewMetadataHandler(im)
	mh.HandleGetMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestMetadataHandler_HandleGetAllMetadata_ResultedOK(t *testing.T) {

	im := repository.NewInMemoryMetadataRepository()
	appID1 := "appID1"
	appID2 := "appID2"

	// add 1st record
	payload := createValidPayload()
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)
	mtd.ApplicationID = appID1
	im.Create(appID1, &mtd)

	// add 2nd record
	payload = createValidPayload2()
	var mtd2 metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd2)
	mtd2.ApplicationID = appID2
	im.Create(appID2, &mtd2)

	request, _ := http.NewRequest("GET", "app-metadata", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()

	mh := NewMetadataHandler(im)
	mh.HandleGetAllMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	var mtds []metadata.ApplicationMetadata
	yaml.NewDecoder(responseRecorder.Body).Decode(&mtds)

	assert.True(t, len(mtds) == 2)

}

func TestMetadataHandler_HandleDeleteMetadata_ResultedOK(t *testing.T) {

	im := repository.NewInMemoryMetadataRepository()
	payload := createValidPayload()
	appID := "appID1"
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)

	mtd.ApplicationID = appID
	im.Create(appID, &mtd)
	request, _ := http.NewRequest("GET", "app-metadata/appID1", strings.NewReader(""))
	mapper := map[string]string{
		"appID": "appID1",
	}
	request = mux.SetURLVars(request, mapper)
	responseRecorder := httptest.NewRecorder()

	mh := NewMetadataHandler(im)
	mh.HandleDeleteMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusNoContent, responseRecorder.Code) // 204

}

func TestMetadataHandler_HandleDeleteMetadata_ResultedConflict(t *testing.T) {

	im := repository.NewInMemoryMetadataRepository()
	payload := createValidPayload()
	appID := "appID1"
	var mtd metadata.ApplicationMetadata
	yaml.Unmarshal([]byte(payload), &mtd)

	mtd.ApplicationID = appID
	im.Create(appID, &mtd)
	request, _ := http.NewRequest("GET", "app-metadata/appID2", strings.NewReader(""))
	// aooID is not the same as what's inserted before, so delete will result as Conflict
	mapper := map[string]string{
		"appID": "appID2",
	}
	request = mux.SetURLVars(request, mapper)
	responseRecorder := httptest.NewRecorder()

	mh := NewMetadataHandler(im)
	mh.HandleDeleteMetadata(responseRecorder, request)

	assert.Equal(t, http.StatusConflict, responseRecorder.Code) // 409

}

func createValidPayload() string {
	return `
title: Valid App 1
version: 1.0.1
maintainers:
- name: First Maintainer App1
  email: firstmaintainer@hotmail.com
- name: Second Maintainer App1
  email: secondmaitainer@gmail.com
company: pellucid Computing
website: http://pellucidcomputing.com
source: https://github.com/elumbantoruan/app-metadata
license: Apache-2.0
description: |-
  ### Interesting title
  Some application content
`
}

func createValidPayload2() string {
	return `
title: Valid App 2
version: 1.0.1
maintainers:
- name: First Maintainer App1
  email: firstmaintainer@hotmail.com
- name: Second Maintainer App1
  email: secondmaitainer@gmail.com
company: pellucid Computing
website: http://pellucidcomputing.com
source: https://github.com/elumbantoruan/app-metadata
license: Apache-2.0
description: |-
  ### Interesting title
  Some application content
`
}

func createInvalidPayloadBadYamlFormat() string {
	return `
title: Valid App 1
			version: 1.0.1
maintainers:
- name: First Maintainer App1
  email: firstmaintainer@hotmail.com
- name: Second Maintainer App1
  email: secondmaitainer@gmail.com
company: pellucid Computing
website: http://pellucidcomputing.com
source: https://github.com/elumbantoruan/app-metadata
license: Apache-2.0
description: |-
  ### Interesting title
  Some application content
`
}

func createInvalidPayloadMissingVersion() string {
	return `
title: Valid App 1
maintainers:
- name: First Maintainer App1
  email: firstmaintainer@hotmail.com
- name: Second Maintainer App1
  email: secondmaitainer@gmail.com
company: pellucid Computing
website: http://pellucidcomputing.com
source: https://github.com/elumbantoruan/app-metadata
license: Apache-2.0
description: |-
  ### Interesting title
  Some application content
`
}

func createInvalidPayloadBadEmail() string {
	return `
title: Valid App 1
version: 1.0.1
maintainers:
- name: First Maintainer App1
  email: firstmaintainer@hotmail.com
- name: Second Maintainer App1
  email: secondmaitainergmail.com
company: pellucid Computing
website: http://pellucidcomputing.com
source: https://github.com/elumbantoruan/app-metadata
license: Apache-2.0
description: |-
  ### Interesting title
  Some application content
`
}

var (
	errInCreate = errors.New("error in create")
	errInUpdate = errors.New("error in update")
	errInGet    = errors.New("error in get")
	errInGetAll = errors.New("error in get all")
	errInDelete = errors.New("error in delete")
)

// FakeMetadataRepository is a concrete implementation of MetadataRepository interface in memory
type FakeMetadataRepository struct {
}

// NewFakeMetadataRepository creates a new instance of InMemoryMetadataRepository
func NewFakeMetadataRepository() repository.MetadataRepository {
	return &FakeMetadataRepository{}
}

// Create adds an application metadata into a repository
func (fm *FakeMetadataRepository) Create(appID string, data *metadata.ApplicationMetadata) error {
	return errInCreate
}

// Update updates the application metadata for a given appID
func (fm *FakeMetadataRepository) Update(appID string, data *metadata.ApplicationMetadata) error {
	return errInUpdate
}

// Get returns application metadata for a given appID
func (fm *FakeMetadataRepository) Get(appID string) (*metadata.ApplicationMetadata, error) {
	return nil, errInGet
}

// GetAll returns all application metadata
func (fm *FakeMetadataRepository) GetAll() ([]metadata.ApplicationMetadata, error) {
	return nil, errInGetAll
}

// Delete removes the application metadata for a given an appID
func (fm *FakeMetadataRepository) Delete(appID string) error {
	return errInDelete
}
