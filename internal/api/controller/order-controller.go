package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/models"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/repository"
	"github.com/mas2401master/go-articles-api-training/pkg/encryption"
	"github.com/mas2401master/go-articles-api-training/pkg/httperror"
)

// GetOrder godoc
// @Summary Recupera un listado Ordenes
// @Description Lista Ordenes registradas por los usuarios clientes
// @Tags Order
// @Produce json
// @Param status query string true "status"
// @Param code query string true "code"
// @Param username query string true "username"
// @Success 200 {object} []models.Order
// @Router /api/v1/orders [get]
// @Security Authorization Token
func GetOrder(c *gin.Context) {
	/*xuserid, err := encryption.ValidHeaderXUserId(c)
	if err != nil || xuserid == 0 {
		httperror.NewError(c, http.StatusUnauthorized, err)
		return
	}*/
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	var filter dto.OrderFilter
	if rolename == "CLIENTE" {
		filter = dto.OrderFilter{
			UserId: strconv.FormatUint(userid, 10),
			Code:   strings.ToUpper(c.Query("code")),
			Status: strings.ToUpper(c.Query("status")),
		}
	} else {
		filter = dto.OrderFilter{
			Code:     strings.ToUpper(c.Query("code")),
			Username: c.Query("username"),
			Status:   strings.ToUpper(c.Query("status")),
		}
	}
	Orders, err := repository.FindAllOrder(filter)
	if err != nil {
		fmt.Println("error:", err)
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if Orders == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("orders not found"))
		return
	}
	c.JSON(http.StatusOK, Orders)
}

// GetOrderById godoc
// @Summary Recupera una Orden segun su ID
// @Description Obtiene informacion de la orden segun su ID Valida orden a mostrar asociada a usuario de session para el rol cliente
// @Tags Order
// @Produce json
// @Param id path integer true "Order ID"
// @Success 200 {object} models.Order
// @Router /api/v1/orders/{id} [get]
// @Security Authorization Token
func GetOrderById(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	id := c.Param("id")

	ID, _ := strconv.ParseUint(id, 10, 64)
	order, err := repository.FindByOrderId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	if order == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("order not found"))
		return
	}
	if rolename == "CLIENTE" && order.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	c.JSON(http.StatusOK, order)
}

// CreateOrder godoc
// @Summary Crea registro de las ordenes de los
// @Description Se crean ordenes con status en Revision, se crea el detalle valida Autorizacion del usuario;Que no exista orden ya creada (status R); Que exista el articulo y que este disponible para crear el detalle; Que exista el codigo de promocion
// @Tags Order
// @Tags Promotion
// @Accept	json
// @Produce json
// @Param orderdto body dto.OrderDTOCreate true "Create order"
// @Success 200 {object} dto.MessageInfo
// @Router /api/v1/orders [post]
// @Security Authorization Token
func CreateOrder(c *gin.Context) {
	var orderDTO dto.OrderDTOCreate
	var promotion_id uint64
	var discount float64
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	user, _ := repository.FindById(userid)
	if rolename == "ADMIN" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized to perform this action"))
		return
	}

	err := c.ShouldBindJSON(&orderDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}

	orderUser, err := repository.FindByOrderUserStatus(user.ID, "R")
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if orderUser != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("there is an order created, try to add detail or edit the order"))
		return
	}
	if orderUser == nil {
		item, _ := repository.FindByItemId(orderDTO.ItemID)
		if item == nil {
			httperror.NewError(c, http.StatusNotFound, errors.New("item not found"))
			return
		}

		if !item.Available {
			httperror.NewError(c, http.StatusBadRequest, errors.New("item not available"))
			return
		}
		orderDTO.Price = item.Price
		if orderDTO.Code != "" {
			promo, _ := repository.FindByCodePromotion(strings.ToUpper(orderDTO.Code))
			if promo == nil {
				httperror.NewError(c, http.StatusNotFound, errors.New("promotion code not found"))
				return
			}
			if promo.Used {
				httperror.NewError(c, http.StatusBadRequest, errors.New("promotion code used"))
				return
			}
			promotion_id = promo.ID
			discount = item.Price * float64(orderDTO.Quantity) * promo.Discount / 100
		}

		order := models.Order{
			UserID:        user.ID,
			PromotionID:   promotion_id,
			Subtotal:      item.Price * float64(orderDTO.Quantity),
			TotalDiscount: discount,
			Total:         item.Price*float64(orderDTO.Quantity) - discount,
			Quantity:      orderDTO.Quantity,
			Status:        "R",
		}
		err := repository.AddOrder(order, orderDTO)
		if err != nil {
			httperror.NewError(c, http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, dto.MessageInfo{Status: http.StatusCreated, Message: "order created!!"})
	}
}

// UpdateOrder godoc
// @Summary Modifica una Orden segun su ID
// @Description modifica la Orden (promocion y status) segun su ID valida Autorizacion del usuario; Que el status de entrada sea R/C; Que exista la orden a editar; Que la orden pertenezca al usuario (cliente); Que exista el codigo de promocion
// @Tags Order
// @Accept	json
// @Produce json
// @Param id path integer true "Order ID"
// @Param orderdto body dto.OrderDTOUpdate true "Edit order"
// @Success 200 {object} models.Order
// @Router /api/v1/orders/{id} [put]
// @Security Authorization Token
func UpdateOrder(c *gin.Context) {
	var orderDTO dto.OrderDTOUpdate
	var promotion_id uint64
	var discount float64
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	id := c.Params.ByName("id")
	err := c.ShouldBindJSON(&orderDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}
	if strings.ToUpper(orderDTO.Status) != "R" && strings.ToUpper(orderDTO.Status) != "C" {
		httperror.NewError(c, http.StatusNotFound, errors.New("status order invalid"))
		return
	}
	ID, _ := strconv.ParseUint(id, 10, 64)
	order, err := repository.FindByOrderId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if order == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("order not found"))
		return
	}
	if rolename == "CLIENTE" && order.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	if orderDTO.Code != "" {
		promo, _ := repository.FindByCodePromotion(strings.ToUpper(orderDTO.Code))
		if promo == nil {
			httperror.NewError(c, http.StatusBadRequest, errors.New("promotion code not found"))
			return
		}

		orderpromoused, _ := repository.OrderPromotionUsed(promo.ID)
		if orderpromoused != 0 && orderpromoused != ID {
			httperror.NewError(c, http.StatusBadRequest, errors.New("promotion code used"))
			return
		}
		promotion_id = promo.ID
		discount = promo.Discount
	}

	ItemUpdate, err := repository.UpdateOrder(ID, promotion_id, discount, strings.ToUpper(orderDTO.Status))
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ItemUpdate)
}

// DeleteOrder godoc
// @Summary Elimina una Orden segun su ID
// @Description Elimina una orden y su detalla, libera disponibilidad de articulos y promociones asociadas. Valida: 1.- que el usuario este autorizado;2.-Que exista la orden registrada. 3.-Que no se elimine la orden si no pertenece al usuario (cliente).
// @Tags Order
// @Produce json
// @Param id path integer true "Order ID"
// @Success 204 {object} dto.MessageInfo
// @Router /api/v1/orders/{id} [delete]
// @Security Authorization Token
func DeleteOrder(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	order, err := repository.FindByOrderId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
	}
	if order == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("order not found"))
		return
	}

	if rolename == "CLIENTE" && order.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	repository.ReleaseItemOrder(ID)
	err = repository.DeleteOrder(ID)
	if order.PromotionID != 0 {
		repository.UpdateUsedPromotion(order.PromotionID, false)
	}
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, err)
		log.Println(err)
		return
	}
	c.JSON(http.StatusNoContent, dto.MessageInfo{Status: http.StatusNoContent, Message: "order deleted"})
}
