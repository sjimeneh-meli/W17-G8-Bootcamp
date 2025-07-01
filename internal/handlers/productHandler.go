package handlers

import (
	"errors"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/handlers/requests"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/mappers"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/validations"
	"net/http"
	"strconv"
	"strings"
)

type productHandler struct {
	service services.ProductServiceI
}

func NewProductHandler(service services.ProductServiceI) *productHandler {
	return &productHandler{
		service: service,
	}
}

func (ph *productHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := ph.service.GetAll()

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (ph *productHandler) Save(w http.ResponseWriter, r *http.Request) {
	var productRequest requests.ProductRequest

	err := request.JSON(r, &productRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	validation := validations.GetProductValidation()

	err = validation.ValidateProductRequestStruct(productRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	product := mappers.GetProductFromRequest(productRequest)

	newProduct, err := ph.service.Create(product)

	if err != nil {
		if errors.Is(err, error_message.ErrAlreadyExists) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	productResponse := mappers.GetProductResponseFromModel(&newProduct)

	response.JSON(w, http.StatusCreated, productResponse)

}

func (ph *productHandler) GetById(w http.ResponseWriter, r *http.Request) {

	idString := strings.TrimSpace(chi.URLParam(r, "id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "error: id is required")
		return
	}

	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "error: id is not a number")
		return
	}

	product, err := ph.service.GetByID(idInt)

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return

	}

	productResponse := mappers.GetProductResponseFromModel(product)

	response.JSON(w, http.StatusOK, productResponse)

}

func (ph *productHandler) Update(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimSpace(chi.URLParam(r, "id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "id is not a number")
		return
	}

	var productRequest requests.ProductRequest

	err = request.JSON(r, &productRequest)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	productToUpdate := mappers.GetProductFromRequest(productRequest)

	productUpdated, err := ph.service.UpdateById(idInt, productToUpdate)

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, productUpdated)

}

func (ph *productHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimSpace(chi.URLParam(r, "id"))

	if idString == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	idInt, err := strconv.Atoi(idString)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "id is not a number")
		return
	}

	err = ph.service.DeleteById(idInt)

	if err != nil {
		if errors.Is(err, error_message.ErrNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
