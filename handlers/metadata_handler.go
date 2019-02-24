package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/elumbantoruan/app-metadata/metadata"
	"github.com/elumbantoruan/app-metadata/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	yaml "gopkg.in/yaml.v2"
)

// MetadataHandler handles app-metadata resource
type MetadataHandler struct {
	Repository repository.MetadataRepository
}

// NewMetadataHandler returns an instance of MetadataHandler
func NewMetadataHandler(repo repository.MetadataRepository) *MetadataHandler {
	return &MetadataHandler{
		Repository: repo,
	}
}

// HandlePostMetadata handles POST operation
func (mh *MetadataHandler) HandlePostMetadata(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload metadata.ApplicationMetadata
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	err = yaml.Unmarshal(b, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		yaml.NewEncoder(w).Encode(err.Error())
		return
	}

	// validate the payload first
	valid, desc := payload.IsValid()
	if !valid {
		w.WriteHeader(http.StatusBadRequest) // 400
		yaml.NewEncoder(w).Encode(desc)
		return
	}

	id := uuid.New()
	payload.ApplicationID = id.String()

	err = mh.Repository.Create(id.String(), &payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		yaml.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	yaml.NewEncoder(w).Encode(payload)
}

// HandlePutMetadata handles PUT operation
func (mh *MetadataHandler) HandlePutMetadata(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload metadata.ApplicationMetadata
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	err = yaml.Unmarshal(b, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// validate the payload first
	valid, desc := payload.IsValid()
	if !valid {
		w.WriteHeader(http.StatusBadRequest) // 400
		yaml.NewEncoder(w).Encode(desc)
		return
	}

	vars := mux.Vars(r)
	var (
		appID string
		ok    bool
	)
	if appID, ok = vars["appID"]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	payload.ApplicationID = appID

	err = mh.Repository.Update(appID, &payload)
	if err != nil {
		if err == repository.ErrIDNotFound {
			w.WriteHeader(http.StatusConflict) // 409
		} else {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		yaml.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK) // 200
	yaml.NewEncoder(w).Encode(payload)
}

// HandleGetMetadata handles GET operation for specified applicationID
func (mh *MetadataHandler) HandleGetMetadata(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	var (
		appID string
		ok    bool
	)
	if appID, ok = vars["appID"]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	res, err := mh.Repository.Get(appID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	if res == nil {
		// no resource is found
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	w.WriteHeader(http.StatusOK) // 200
	yaml.NewEncoder(w).Encode(res)
}

// HandleGetAllMetadata handles all GET operation
func (mh *MetadataHandler) HandleGetAllMetadata(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	res, err := mh.Repository.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	if res == nil {
		// no resources found
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	w.WriteHeader(http.StatusOK) // 200
	yaml.NewEncoder(w).Encode(res)
}

// HandleDeleteMetadata handles DELETE operation
func (mh *MetadataHandler) HandleDeleteMetadata(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	var (
		appID string
		ok    bool
	)
	if appID, ok = vars["appID"]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	err := mh.Repository.Delete(appID)
	if err != nil {
		if err == repository.ErrIDNotFound {
			w.WriteHeader(http.StatusConflict) // 409
		} else {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		yaml.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}
