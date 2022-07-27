package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) getProducts(ctx *gin.Context) {
	// https://localhost:8080/products?category=1&size=110&type=Гофрований-А&price=5.44:15.2

	searchBar := ctx.Query("search")

	category := ctx.Query("category")
	if category != "" {
		if strings.Index(category, "-") != -1 {
			category = strings.Replace(category, "-", " ", -1)
		}
	}

	size, _ := strconv.Atoi(ctx.Query("size"))
	if size != 0 {
		//todo some validation
	}

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

	params := &server.SearchParams{CategoryTitle: category, Size: size, Price: price}

	products, err := h.services.Products.GetWithParams(*params, createdAt, searchBar)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": products,
	})
}

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

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "product added",
	})
}

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

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"product_id": id,
		"message":    "product updated",
	})
}

func (h *Handler) deleteProduct(ctx *gin.Context) {
	type ProductIdInput struct {
		Id int `json:"id"`
	}

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

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      input.Id,
		"message": "deleted",
	})
}
