package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server"
	"net/http"
)

// addToCart godoc
// @Summary      AddToCart
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  adds a product to a cart
// @ID adds a product to a cart
// @Accept       json
// @Produce      json
// @Param        input body server.CartProduct true "product info"
// @Success      200  {object}  string 2
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/add-to-cart [post]
func (h *Handler) addToCart(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	fmt.Println(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	var input server.CartProduct

	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	price, err := h.services.Shopping.AddToCart(userId, input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"price_for_quantity": price,
		"message":            "product added",
	})
}

type CartProductsResponse struct {
	Products []server.CartProduct `json:"products"`
	Sum      float64              `json:"sum"`
}

// getFromCart godoc
// @Summary      GetProductsInCart
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  gets products from a cart
// @ID gets products from a cart
// @Produce      json
// @Success      200  {object}  handler.CartProductsResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/user-cart [get]
func (h *Handler) getFromCartById(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "no user's id")
		return
	}

	products, sum, err := h.services.Shopping.GetProductsInCart(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, CartProductsResponse{
		Products: products,
		Sum:      sum,
	})
}

// deleteFromCart godoc
// @Summary      DeleteFromCart
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  deletes a product from users cart
// @ID deletes from cart
// @Accept       json
// @Produce      json
// @Param        input body handler.ProductIdInput true "product id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/delete-from-cart [delete]
func (h *Handler) deleteFromCart(ctx *gin.Context) {
	var product ProductIdInput

	if err := ctx.BindJSON(&product); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "product_id to delete was not found: "+err.Error())
		return
	}

	err := h.services.Shopping.DeleteFromCart(product.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      product.Id,
		Message: "product deleted",
	})
}

type AddToFavouritesInput struct {
	ProductId int `json:"product_id"`
}

// addToFavourites godoc
// @Summary      AddToFavourites
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  adds a product to favourites
// @ID adds to favourites
// @Accept       json
// @Produce      json
// @Param        input body handler.ProductIdInput true "product id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/add-to-favourites [post]
func (h *Handler) addToFavourites(ctx *gin.Context) {
	var input ProductIdInput
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Shopping.AddToFavourites(userId, input.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      input.Id,
		Message: "item added to favourites",
	})
}

// getFavourites godoc
// @Summary      GetFavourites
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  gets user favourite products
// @ID get favourites
// @Produce      json
// @Success      200  {array}  	server.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/get-favourites [get]
func (h *Handler) getFavourites(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	products, err := h.services.Shopping.GetFavourites(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

// deleteFromFavourites godoc
// @Summary DeleteFromFavourites
// @Security ApiKeyAuth
// @Tags api/client
// @Description deletes product from favourites
// @ID deletes from favourites
// @Accepts json
// @Produce json
// @Param        input body handler.ProductIdInput true "product id"
// @Success      200  {array}  	handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/delete-from-favourites [delete]
func (h *Handler) deleteFromFavourites(ctx *gin.Context) {
	var input ProductIdInput

	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "user id was not found "+err.Error())
		return
	}

	err = h.services.Shopping.DeleteFromFavourites(userId, input.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      input.Id,
		Message: "product deleted from favourites",
	})
}
