package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server"
	"net/http"
)

// addToCart godoc
// @Summary      AddToCart
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

func (h *Handler) deleteFromCart(ctx *gin.Context) {
	type ProductInput struct {
		Id int `json:"product_id"`
	}

	var product ProductInput

	if err := ctx.BindJSON(&product); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "product_id to delete was not found: "+err.Error())
		return
	}

	err := h.services.Shopping.DeleteFromCart(product.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      product.Id,
		"message": "deleted from cart",
	})
}

type AddToFavouritesInput struct {
	ProductId int `json:"product_id"`
}

func (h *Handler) addToFavourites(ctx *gin.Context) {
	var input AddToFavouritesInput
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Shopping.AddToFavourites(userId, input.ProductId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "product added to favourites",
	})
}

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
