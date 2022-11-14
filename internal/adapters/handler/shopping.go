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
		ctx.SetCookie(userCartCookie, "", -1, "/", domain, false, true)
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
	productId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to parse product id to int: %v", err).Error())
		return
	}

	err = h.services.Shopping.DeleteFromCart(productId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      productId,
		Message: "product deleted",
	})
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
// @Router       /api/client/favourites [post]
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
// @Success      200  {array}  	product.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/favourites [get]
func (h *Handler) getFavourites(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	products, err := h.services.Product.GetFavourites(userId)
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
// @Param 		 id query string true "Product id"
// @Success      200  {array}  	handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/favourites [delete]
func (h *Handler) deleteFromFavourites(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "user id was not found "+err.Error())
		return
	}

	err = h.services.Shopping.DeleteFromFavourites(userId, id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "product deleted from favourites",
	})
}
