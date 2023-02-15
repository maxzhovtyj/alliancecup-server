package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/product"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/shopping"
	"net/http"
	"strconv"
	"strings"
)

type ProductIdInput struct {
	Id int `json:"id"`
}

type ProductVisibility struct {
	Id       int  `json:"id"`
	IsActive bool `json:"isActive"`
}

// getProducts godoc
// @Summary      GetProducts
// @Tags         api
// @Description  	 get products from certain category with params
// @ID 			 gets products
// @Produce      json
// @Param 		 category query string false "Category"
// @Param 	   	 size query int false "Size"
// @Param 		 search query string false "Search"
// @Param 		 price query string false "Price"
// @Param 		 characteristic query string false "Characteristic"
// @Param		 created_at query string false "Created At"
// @Success      200  {array}  product.Product
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/products [post]
func (h *Handler) getProducts(ctx *gin.Context) {
	// Example http://localhost:8000/api/products?characteristic=Колір:Білий|Розмір:110мл
	var params shopping.SearchParams

	category := ctx.Query("category")
	var err error
	if category != "" {
		params.CategoryId, err = strconv.Atoi(category)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	params.PriceRange = ctx.Query("priceRange")
	params.CreatedAt = ctx.Query("createdAt")
	params.Search = ctx.Query("search")

	isActive := ctx.Query("isActive")
	if isActive != "" {
		parseBool, err := strconv.ParseBool(isActive)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to parse isActive field %v", err).Error())
			return
		}
		params.IsActive = &parseBool
	}

	characteristic := ctx.Query("characteristic")
	if characteristic != "" {
		arr := strings.Split(characteristic, "|")
		for _, e := range arr {
			var paramChar shopping.CharacteristicParam

			eArr := strings.Split(e, ":")
			if len(eArr) == 2 {
				paramChar.Name = eArr[0]
				paramChar.Value = eArr[1]
			}

			params.Characteristic = append(params.Characteristic, paramChar)
		}
	}
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.services.Product.GetWithParams(params)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// addProduct godoc
// @Summary      Create
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Adds a new product
// @ID 			 adds product
// @Accept 	     json
// @Produce      json
// @Param        input body product.Product true "product info" // TODO
// @Success      201  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product [post]
func (h *Handler) addProduct(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "form/json")
	err := ctx.Request.ParseMultipartForm(fileMaxSize)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var dto product.CreateDTO

	dto.Article = ctx.Request.Form.Get("article")
	if dto.Article == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "article is empty")
		return
	}

	dto.CategoryTitle = ctx.Request.Form.Get("categoryTitle")
	if dto.CategoryTitle == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "category title is empty")
		return
	}

	dto.ProductTitle = ctx.Request.Form.Get("productTitle")
	if dto.ProductTitle == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "product title is empty")
		return
	}

	amountInStock := ctx.Request.Form.Get("amountInStock")
	if amountInStock != "" {
		dto.AmountInStock, err = strconv.ParseFloat(amountInStock, 64)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	dto.Price, err = strconv.ParseFloat(ctx.Request.Form.Get("price"), 64)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid price: %f", dto.Price))
		return
	}

	characteristics := ctx.Request.Form.Get("characteristic")
	if characteristics != "" {
		char := types.JSONText(characteristics)
		dto.Characteristics = &char
	}

	packaging := ctx.Request.Form.Get("packaging")
	if packaging != "" {
		p := types.JSONText(packaging)
		dto.Packaging = &p
	}

	description := ctx.Request.Form.Get("description")
	if description != "" {
		dto.Description = &description
	}

	dto.Img, err = parseFile(ctx.Request.MultipartForm.File)
	if err != nil {
		if !errors.Is(err, ErrEmptyFile) {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	id, err := h.services.Product.Add(dto)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, ItemProcessedResponse{
		Id:      id,
		Message: "product successfully created",
	})
}

// getProductById godoc
// @Summary      GetProductById
// @Tags         api
// @Description  get product full info by its id
// @ID 			 gets full product info
// @Produce      json
// @Param 		 id query int true "Product id"
// @Success      200  {object}  product.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/product [get]
func (h *Handler) getProductById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "id was not found in request path")
		return
	}

	selectedProduct, err := h.services.Product.GetProductById(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, selectedProduct)
}

// updateProduct godoc
// @Summary      UpdateProduct
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Updates product
// @ID 			 updates product
// @Accept 	     json
// @Produce      json
// @Param        input body product.Product true "product info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product-image [put]
func (h *Handler) updateProduct(ctx *gin.Context) {
	var input product.Product
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Product.Update(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "product updated",
	})
}

// updateProductImage godoc
// @Summary      Update product
// @Tags         api
// @Description  Update product image
// @ID 			 updates product image
// @Produce      json
// @Param 		 id query int true "Product id"
// @Success      200  {object}
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product [put]
func (h *Handler) updateProductImage(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "form/json")
	err := ctx.Request.ParseMultipartForm(fileMaxSize)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var dto product.UpdateImageDTO

	dto.Id, err = strconv.Atoi(ctx.Request.Form.Get("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid id: %d", dto.Id))
		return
	}

	dto.Img, err = parseFile(ctx.Request.MultipartForm.File)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Product.UpdateImage(dto)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "product image updated",
	})
}

// deleteProductImage godoc
// @Summary      Delete product image (Minio)
// @Tags         api
// @Description  Delete product image (Minio)
// @ID 			 delete product image (minio)
// @Produce      json
// @Param 		 id query int true "Product id"
// @Success      200  {object}
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product-image [delete]
func (h *Handler) deleteProductImage(ctx *gin.Context) {
	productId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid id, %v", err).Error())
		return
	}

	err = h.services.Product.DeleteImage(productId)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"id":      productId,
		"message": "product image (Minio) deleted",
	})
}

// updateProductVisibility godoc
// @Summary      Update product visibility
// @Tags         api
// @Description  Update product image
// @ID 			 updates product image
// @Produce      json
// @Param 		 id body bool true "Product id"
// @Success      200  {object}
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/product-visibility [put]
func (h *Handler) updateProductVisibility(ctx *gin.Context) {
	var input ProductVisibility
	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Product.UpdateVisibility(input.Id, input.IsActive)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      input.Id,
		Message: "product visibility changed",
	})
}

// deleteProduct godoc
// @Summary      DeleteProduct
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  	 Deletes product
// @ID 			 delete product
// @Accept 	     json
// @Produce      json
// @Param        id query int true "product id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product [delete]
func (h *Handler) deleteProduct(ctx *gin.Context) {
	productId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Product.Delete(productId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      productId,
		Message: "product deleted",
	})
}
