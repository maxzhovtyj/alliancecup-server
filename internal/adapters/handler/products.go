package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"net/http"
	"strconv"
)

type ProductIdInput struct {
	Id int `json:"id"`
}

// getProducts godoc
// @Summary      GetProducts
// @Tags         api
// @Product  get products from certain category with params
// @ID 			 gets products
// @Produce      json
// @Param 		 category query string true "Category"
// @Param 	   	 size query int false "Size"
// @Param 		 type query string false "Type"
// @Param 		 search query string false "Search"
// @Param 		 price query string false "Price"
// @Param 		 characteristic query string false "characteristic"
// @Param		 created_at query string false "Created At"
// @Success      200  {array}   product.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/products [post]
func (h *Handler) getProducts(ctx *gin.Context) {
	category := ctx.Query("category")
	var categoryId int
	var err error
	if category != "" {
		categoryId, err = strconv.Atoi(category)
	}

	price := ctx.Query("priceRange") // TODO price validation
	createdAt := ctx.Query("createdAt")
	characteristic := ctx.Query("characteristic")
	search := ctx.Query("search")

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.services.Product.GetWithParams(server.SearchParams{
		CategoryId:     categoryId,
		PriceRange:     price,
		CreatedAt:      createdAt,
		Characteristic: characteristic,
		Search:         search,
	})

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": products,
	})
}

// addProduct godoc
// @Summary      AddProduct
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Product  Adds a new product
// @ID 			 adds product
// @Accept 	     json
// @Produce      json
// @Param        input body product.Info true "product info"
// @Success      201  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product [post]
func (h *Handler) addProduct(ctx *gin.Context) {
	var input product.Info

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Product.AddProduct(input.Product, input.Description)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, ItemProcessedResponse{
		Id:      id,
		Message: "product added",
	})
}

// getProductById godoc
// @Summary      GetProductById
// @Tags         api
// @Product  get product full info by its id
// @ID 			 gets full product info
// @Produce      json
// @Param 		 id query int true "Product id"
// @Success      200  {object}  product.Info
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
// @Product  Updates product
// @ID 			 updates product
// @Accept 	     json
// @Produce      json
// @Param        input body product.Info true "product info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/product [put]
func (h *Handler) updateProduct(ctx *gin.Context) {
	var input product.Info

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

// deleteProduct godoc
// @Summary      DeleteProduct
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Product  Deletes product
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
	// TODO "pq: update or delete on table \"products\" violates foreign key constraint \"orders_products_product_id_fkey\" on table \"orders_products\""

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
