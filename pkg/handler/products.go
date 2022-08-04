package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server"
	"net/http"
	"strconv"
	"strings"
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
// @Router       /api/get-products [get]
func (h *Handler) getProducts(ctx *gin.Context) {
	// https://localhost:8080/products?category=Одноразові-Стакани&size=110&type=Гофрований-А&price=5.44:15.2

	searchBar := ctx.Query("search")

	category := ctx.Query("category")
	if category != "" {
		if strings.Index(category, "-") != -1 {
			category = strings.Replace(category, "-", " ", -1)
		}
	}

	size := ctx.Query("size")

	price := ctx.Query("price")
	if price != "" {
		if strings.Index(price, ":") != -1 {
			split := strings.Split(price, ":")
			gt := split[0]
			lt := split[1]
			price = gt + " " + lt
		}
	} else {
		price = "0.00 1000.00"
	}

	createdAt := ctx.Query("created_at")

	characteristic := ctx.Query("characteristic")

	params := &server.SearchParams{
		CategoryTitle: category,
		Size:          size, Price: price,
		Characteristic: characteristic,
	}

	products, err := h.services.Products.GetWithParams(*params, createdAt, searchBar)
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
// @Success      200  {object}  handler.ItemProcessedResponse
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
// @Param        input body handler.ProductIdInput true "product id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/delete-product [delete]
func (h *Handler) deleteProduct(ctx *gin.Context) {
	var input ProductIdInput

	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Products.Delete(input.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      input.Id,
		Message: "product deleted",
	})
}
