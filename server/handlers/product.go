package handlers

import (
	"context"
	productdto "dumbmerch/dto/product"
	dto "dumbmerch/dto/result"
	"dumbmerch/models"
	"dumbmerch/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var ctx = context.Background()
var CLOUD_NAME = os.Getenv("CLOUD_NAME")
var API_KEY = os.Getenv("API_KEY")
var API_SECRET = os.Getenv("API_SECRET")

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

func (h *handlerProduct) FindProducts(c echo.Context) error {
	products, err := h.ProductRepository.FindProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// delete this

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: products})
}

func (h *handlerProduct) GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var product models.Product
	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// delete this

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)})
}

func (h *handlerProduct) CreateProduct(c echo.Context) error {
	var err error
	filepath := c.Get("dataFile").(string)

	price, _ := strconv.Atoi(c.FormValue("price"))
	qty, _ := strconv.Atoi(c.FormValue("qty"))

	categoryIdString := c.FormValue("category_id")
	if categoryIdString == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Error: category_id form value is missing."})
	}

	var categoriesId []int
	err = json.Unmarshal([]byte(categoryIdString), &categoriesId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if len(categoriesId) == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Error: category_id form value is missing."})
	}

	request := productdto.CreateProductRequest{
		Name:       c.FormValue("name"),
		Desc:       c.FormValue("desc"),
		Price:      price,
		Image:      filepath,
		Qty:        qty,
		CategoryID: categoriesId,
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	// code here
	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbmerchb44"})

	if err != nil {
		fmt.Println(err.Error())
	}

	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	categories, _ := h.ProductRepository.FindCategoriesById(request.CategoryID)

	product := models.Product{
		Name:     request.Name,
		Desc:     request.Desc,
		Price:    request.Price,
		Image:    resp.SecureURL,
		Qty:      request.Qty,
		Category: categories,
		UserID:   int(userId),
	}

	product, err = h.ProductRepository.CreateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	product, _ = h.ProductRepository.GetProduct(product.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: convertResponseProduct(product)})
}

func (h *handlerProduct) UpdateProduct(c echo.Context) error {
	var err error
	filepath := c.Get("dataFile").(string)

	price, _ := strconv.Atoi(c.FormValue("price"))
	qty, _ := strconv.Atoi(c.FormValue("qty"))

	var categoriesId []int
	categoryIdString := c.FormValue("category_id")
	err = json.Unmarshal([]byte(categoryIdString), &categoriesId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	request := productdto.UpdateProductRequest{
		Name:       c.FormValue("name"),
		Desc:       c.FormValue("desc"),
		Price:      price,
		Image:      filepath,
		Qty:        qty,
		CategoryID: categoriesId,
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	// code here
	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbmerchb44"})

	if err != nil {
		fmt.Println(err.Error())
	}

	id, _ := strconv.Atoi(c.Param("id"))

	product, err := h.ProductRepository.GetProduct(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Name != "" {
		product.Name = request.Name
	}

	if request.Desc != "" {
		product.Desc = request.Desc
	}

	if request.Price != 0 {
		product.Price = request.Price
	}

	if request.Image != "" {
		product.Image = resp.SecureURL // change this
	}

	if request.Qty != 0 {
		product.Qty = request.Qty
	}

	if len(request.CategoryID) == 0 {
		data, err := h.ProductRepository.DeleteProductCategoryByProductId(product)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
	}

	categories, _ := h.ProductRepository.FindCategoriesById(request.CategoryID)
	product.Category = categories

	data, err := h.ProductRepository.UpdateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}

func (h *handlerProduct) DeleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.ProductRepository.DeleteProduct(product, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: http.StatusOK, Data: data})
}

func convertResponseProduct(u models.Product) models.ProductResponse {
	return models.ProductResponse{
		ID:       u.ID,
		Name:     u.Name,
		Desc:     u.Desc,
		Price:    u.Price,
		Image:    u.Image,
		Qty:      u.Qty,
		User:     u.User,
		Category: u.Category,
	}
}
