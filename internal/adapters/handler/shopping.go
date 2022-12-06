package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"net/http"
	"strconv"
)

type CartProductsResponse struct {
	Products []shopping.CartProduct `json:"products"`
	Sum      float64                `json:"sum"`
}

type AddToFavouritesInput struct {
	ProductId int `json:"product_id"`
}

// addToCart godoc
// @Summary      AddToCart
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  adds a product to a cart
// @ID adds a product to a cart
// @Accept       json
// @Produce      json
// @Param        input body shopping.CartProduct true "product info"
// @Success      200  {string}  string
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/cart [post]
func (h *Handler) addToCart(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	userCartId, err := getCartId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var input shopping.CartProduct
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Shopping.AddToCart(input, userCartId, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, "product added")
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
// @Router       /api/client/cart [get]
func (h *Handler) getFromCartById(ctx *gin.Context) {
	cartId, err := getCartId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	products, sum, err := h.services.Shopping.GetCart(cartId)
	if err != nil {
		ctx.SetCookie(userCartCookie, "", -1, "/", h.cfg.Domain, false, true)
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
// @Param 		 id query string true "Product id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/cart [delete]
func (h *Handler) deleteFromCart(ctx *gin.Context) {
	cartId, err := getCartId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to find cart id, %v", err).Error())
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		userId = 0
	}

	productId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to parse product id to int: %v", err).Error())
		return
	}

	err = h.services.Shopping.DeleteFromCart(productId, userId, cartId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      productId,
		Message: "product deleted",
	})
}
