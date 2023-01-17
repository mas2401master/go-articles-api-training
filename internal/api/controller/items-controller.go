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

// CreateItem godoc
// @Summary Crea el registro de un articulo
// @Description El articulo se crea en status activo; Valida autorizacion de usuario
// @Tags Items
// @Accept	json
// @Produce json
// @Param ItemDTO body dto.ItemDTOCreate true "create item"
// @Success 200 {object} dto.MessageInfo
// @Router /api/v1/items/{id} [post]
// @Security Authorization Token
func CreateItem(c *gin.Context) {
	var ItemDTO dto.ItemDTOCreate
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	err := c.ShouldBindJSON(&ItemDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		log.Println(err)
		return
	}
	//price, _ := strconv.ParseFloat(ItemDTO.Price, 64)
	ItemNew := models.Item{
		NameItem:    ItemDTO.NameItem,
		Description: ItemDTO.Description,
		Price:       ItemDTO.Price,
	}

	err = repository.AddItem(ItemNew)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusCreated, dto.MessageInfo{Status: http.StatusCreated, Message: "item created!!"})
	}
}

// GetItemById godoc
// @Summary Recupera un articuulos segun su ID
// @Description Obtiene informacion del articulo segun su ID Valida autorizacion de usuario; existencia del articulo a mostrar; usuario cliente sin acceso al articulo con status falso
// @Tags Items
// @Produce json
// @Param id path integer true "Item ID"
// @Success 200 {object} models.Item
// @Router /api/v1/items/{id} [get]
// @Security Authorization Token
func GetItemById(c *gin.Context) {
	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	Items, err := repository.FindByItemId(ID)
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))

	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
	}
	if Items == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item not found"))
		return
	}
	if rolename == "CLIENTE" && !Items.Status {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	c.JSON(http.StatusOK, Items)

}

// GetItems godoc
// @Summary Recupera un listado de Articulos registrados
// @Description Lista articulos registrados con status ACTIVO valida autorizacion de usuario
// @Tags Items
// @Produce json
// @Param name query string true "name"
// @Param available query string true "available"
// @Success 200 {object} []models.Item
// @Router /api/v1/items [get]
// @Security Authorization Token
func GetItems(c *gin.Context) {
	filter := dto.ItemFilter{
		NameItem:  c.Query("name"),
		Available: c.Query("available"),
	}
	Items, err := repository.FindAllItem(filter)
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("items not found"))
		return
	}
	if Items == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("items not found"))
		return
	}
	c.JSON(http.StatusOK, Items)
}

// UpdateItem godoc
// @Summary Modifica un articulo segun su ID
// @Description modifica el articulo segun su ID valida autorizacion de usuario; existencia del articulo a editar;
// @Tags Items
// @Produce json
// @Param id path integer true "Item ID"
// @Param itemdto body dto.ItemDTOUpdate true "create item"
// @Success 200 {object} models.Item
// @Router /api/v1/items/{id} [put]
// @Security Authorization Token
func UpdateItem(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	id := c.Params.ByName("id")
	var ItemDTO dto.ItemDTOUpdate

	err := c.ShouldBindJSON(&ItemDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}

	ID, _ := strconv.ParseUint(id, 10, 64)
	ItemValid, err := repository.FindByItemId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if ItemValid == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item not found"))
		return
	}
	ItemUpdate, err := repository.UpdateItem(ID, ItemDTO)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, ItemUpdate)
}

// DeleteItem godoc
// @Summary Elimina un articulo segun su ID
// @Description Elimina un arituculo  segun su ID valida autorizacion de usuario; existencia del articulo a eliminar; Articulo no asociado a una orden
// @Tags Items
// @Produce json
// @Param id path integer true "Item ID"
// @Success 204 {object} dto.MessageInfo
// @Router /api/v1/items/{id} [delete]
// @Security Authorization Token
func DeleteItem(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	item, err := repository.FindByItemId(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if item == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("item not found"))
		return
	}

	if !item.Available {
		httperror.NewError(c, http.StatusBadRequest, errors.New("item is associated with an order, it cannot be deleted"))
		return
	}

	err = repository.DeleteItem(ID)
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusNoContent, dto.MessageInfo{Status: http.StatusNoContent, Message: "items deleted"})
}
