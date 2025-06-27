package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/responses"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

func GetSectionHandler() SectionHandlerI {
	return &SectionHandler{
		service: services.GetSectionService(),
	}
}

type SectionHandlerI interface {
	GetAll() http.HandlerFunc
}

type SectionHandler struct {
	service services.SectionServiceI
}

func (h *SectionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			response         *responses.DataResponse = &responses.DataResponse{}
			sectionsResponse []*responses.SectionResponse
			sections         []*models.Section
		)

		w.Header().Set("Content-Type", "application/json")

		sections = h.service.GetAll()
		sectionsResponse = mappers.GetListSectionResponseFromListModel(sections)
		response.Data = sectionsResponse

		w.WriteHeader(http.StatusOK)
		encodeErr := json.NewEncoder(w).Encode(response)
		if encodeErr != nil {
			response.SetError(encodeErr.Error())
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
