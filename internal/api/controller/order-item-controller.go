package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/models"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/repository"
	"github.com/mas2401master/go-articles-api-training/pkg/encryption"
	"github.com/mas2401master/go-articles-api-training/pkg/httperror"
)

// GetOrderItem godoc
// @Summary Recupera listado de articulos de detalle de la orden
// @Description Lista detalle de la orden valida  usuario autorizado; existencia de orden; detalle de orden a mostrar asociado a usuario de session para el rol cliente
// @Tags Order Detail
// @Produce json
// @Param id path integer true "Order ID"
// @Success 200 {object} []models.OrderItems
// @Router /api/v1/orders/details/{id} [get]
// @Security Authorization Token
func GetOrderItem(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	id := c.Param("id")
	idorder, _ := strconv.ParseUint(id, 10, 64)
	orderHead, err := repository.FindByOrderId(idorder)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	if orderHead == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("order not found"))
		return
	}

	if rolename == "CLIENTE" && orderHead.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	orders, err := repository.FindAllOrderItem(idorder)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetOrderItemById godoc
// @Summary Recupera una Orden segun su ID
// @Description Obtiene el detalle de la order segun su ID Valida usuario autorizado; existencia del detalla; detalle a mostrar asociada a usuario de session para el rol cliente
// @Tags Order Detail
// @Produce json
// @Param id path integer true "Detail ID"
// @Success 200 {object} models.OrderItems
// @Router /api/v1/orders/detail/{id} [get]
// @Security Authorization Token
func GetOrderItemById(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	orderItem, err := repository.FindByOrderItemId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	if orderItem == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("order not found"))
		return
	}

	orderHead, _ := repository.FindByOrderId(orderItem.OrderID)
	if rolename == "CLIENTE" && orderHead.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	c.JSON(http.StatusOK, orderItem)
}

// CreateOrderItem godoc
// @Summary Crea registro de las ordenes de los
// @Description Crea el detalle de la orden valida  usuario autorizado; existencia de la orden a agregar el detalle; orden asociada a usuario de session para el rol cliente; existencia y disponibilidad del artiulo a ser agregado al detalle de la orden
// @Tags Order Detail
// @Accept	json
// @Produce json
// @Param id path integer true "Order ID"
// @Param orderitemdto body dto.OrderItemDTO true "Create detail order"
// @Success 200 {object} dto.MessageInfo
// @Router /api/v1/orders/detail/{id} [post]
// @Security Authorization Token
func CreateOrderItem(c *gin.Context) {
	var orderitemDTO dto.OrderItemDTO
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "ADMIN" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	id := c.Params.ByName("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	err := c.ShouldBindJSON(&orderitemDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}

	orderHead, _ := repository.FindByOrderId(ID)
	if orderHead == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("order not found"))
		return
	}

	if rolename == "CLIENTE" && orderHead.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	if orderHead.Status != "R" {
		httperror.NewError(c, http.StatusNotFound, errors.New("order with status completed, cannot add items"))
		return
	}

	item, _ := repository.FindByItemId(orderitemDTO.ItemID)
	if item == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item not found"))
		return
	}

	if !item.Available {
		httperror.NewError(c, http.StatusBadRequest, errors.New("item not available"))
		return
	}

	newOrderItem := models.OrderItems{
		OrderID:  orderHead.ID,
		ItemID:   item.ID,
		Price:    item.Price,
		Quantity: orderitemDTO.Quantity,
		Total:    item.Price * float64(orderitemDTO.Quantity),
		Status:   orderHead.Status,
	}
	err = repository.AddOrderItem(newOrderItem)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, dto.MessageInfo{Status: http.StatusCreated, Message: "item added to your order!!"})
}

// UpdateOrder godoc
// @Summary Modifica una Orden segun su ID
// @Description modifica el detalle de la orden segun su ID valida  usuario autorizado; existencia del detalle; orden asociada a usuario de session para el rol cliente; existencia  del artiulo a ser editado al detalle de la orden
// @Tags Order Detail
// @Accept	json
// @Produce json
// @Param id path integer true "Detail ID"
// @Param orderitemdto body dto.OrderItemDTO true "Create detail order"
// @Success 200 {object} models.OrderItems
// @Router /api/v1/orders/detail/{id} [put]
// @Security Authorization Token
func UpdateOrderItem(c *gin.Context) {
	var orderitemDTO dto.OrderItemDTO
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	id := c.Params.ByName("id")
	err := c.ShouldBindJSON(&orderitemDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}
	ID, _ := strconv.ParseUint(id, 10, 64)
	orderItem, err := repository.FindByOrderItemId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if orderItem == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item does not exist in your order"))
		return
	}
	orderHead, _ := repository.FindByOrderId(orderItem.OrderID)
	if rolename == "CLIENTE" && orderHead.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	item, _ := repository.FindByItemId(orderitemDTO.ItemID)
	if item == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item not found"))
		return
	}
	changeOrderItem := models.OrderItems{
		ID:       orderItem.ID,
		OrderID:  orderItem.OrderID,
		Price:    item.Price,
		Quantity: orderitemDTO.Quantity,
		Total:    item.Price * float64(orderitemDTO.Quantity),
	}
	ItemUpdate, err := repository.UpdateOrderItem(changeOrderItem)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ItemUpdate)

}

// DeleteOrderItem godoc
// @Summary Elimina un detalle de la Orden segun su ID
// @Description Elimina una un detalle de la orden, libera disponibilidad de articulo asociado. valida  usuario autorizado; existencia del detalle; orden asociada a usuario de session para el rol cliente; existencia  del artiulo a ser eliminado al detalle de la orden.
// @Tags Order Detail
// @Produce json
// @Param id path integer true "Detail ID"
// @Success 204 {object} dto.MessageInfo
// @Router /api/v1/orders/detail/{id} [delete]
// @Security Authorization Token
func DeleteOrderItem(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	orderItem, err := repository.FindByOrderItemId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if orderItem == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item does not exist in your order"))
		return
	}
	orderHead, _ := repository.FindByOrderId(orderItem.OrderID)
	if rolename == "CLIENTE" && orderHead.UserID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	err = repository.DeleteOrderItem(orderItem.ID, orderItem.OrderID)
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, err)
		log.Println(err)
		return
	}
	c.JSON(http.StatusNoContent, dto.MessageInfo{Status: http.StatusNoContent, Message: "order detail deleted"})
}
