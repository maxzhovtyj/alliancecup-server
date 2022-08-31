package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server/pkg/models"
	"net/http"
	"strconv"
)

type ProductIdInput struct {
	Id int `json:"id"`
}

// getProducts godoc
// @Summary      GetProducts
// @Tags         api
// @Description  get products from certain category with params
// @ID 			 gets products
// @Produce      json
// @Param 		 category query string true "Category"
// @Param 	   	 size query int false "Size"
// @Param 		 type query string false "Type"
// @Param 		 search query string false "Search"
// @Param 		 price query string false "Price"
// @Param 		 characteristic query string false "characteristic"
// @Param		 created_at query string false "Created At"
// @Success      200  {array}   server.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/get-products [post]
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

	products, err := h.services.Products.GetWithParams(server.SearchParams{
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
// @Description  Adds a new product
// @ID 			 adds product
// @Accept 	     json
// @Produce      json
// @Param        input body server.ProductInfoDescription true "product info"
// @Success      201  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/add-product [post]
func (h *Handler) addProduct(ctx *gin.Context) {
	var input server.ProductInfoDescription

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Products.AddProduct(input.Info, input.Description)
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
// @Description  get product full info by its id
// @ID 			 gets full product info
// @Produce      json
// @Param 		 id query int true "Product id"
// @Success      200  {object}  server.ProductInfoDescription
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

	product, err := h.services.Products.GetProductById(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// updateProduct godoc
// @Summary      UpdateProduct
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Updates product
// @ID 			 updates product
// @Accept 	     json
// @Produce      json
// @Param        input body server.ProductInfoDescription true "product info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/update-product [put]
func (h *Handler) updateProduct(ctx *gin.Context) {
	var input server.ProductInfoDescription

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Products.Update(input)
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
// @Description  Deletes product
// @ID 			 delete product
// @Accept 	     json
// @Produce      json
// @Param        id query int true "product id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/delete-product [delete]
func (h *Handler) deleteProduct(ctx *gin.Context) {
	// TODO "pq: update or delete on table \"products\" violates foreign key constraint \"orders_products_product_id_fkey\" on table \"orders_products\""

	productId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Products.Delete(productId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      productId,
		Message: "product deleted",
	})
}
