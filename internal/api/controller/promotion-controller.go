package controller

import (
	"errors"
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

// CreatePromotion godoc
// @Summary Crea registro de promociones de los
// @Description Las promociones se crean con Usado en falso Valida autorizacion de usuario; que el codigo de promocion no este registrado
// @Tags Promotion
// @Accept	json
// @Produce json
// @Param PromoDTO body dto.PromotionDTOCreate true "create promotion"
// @Success 200 {object} dto.MessageInfo
// @Router /api/v1/promotion [post]
// @Security Authorization Token
func CreatePromotion(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	var PromoDTO dto.PromotionDTOCreate
	err := c.ShouldBindJSON(&PromoDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	PromoValid, _ := repository.FindByCodePromotion(strings.ToUpper(PromoDTO.Code))
	if PromoValid != nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("promotion code not available"))
		return
	}

	PromoNew := models.Promotion{
		Name:     PromoDTO.Name,
		Code:     strings.ToUpper(PromoDTO.Code),
		Discount: PromoDTO.Discount,
	}

	err = repository.AddPromotion(PromoNew)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, dto.MessageInfo{Status: http.StatusCreated, Message: "promotion created!!"})
}

// GetPromotionById godoc
// @Summary Recupera una Promocion segun su ID
// @Description Obtiene informacion de la Promocion segun su ID Valida autorizacion de usuario; existencia de promocion a mostrar
// @Tags Promotion
// @Produce json
// @Param id path integer true "Promotion ID"
// @Success 200 {object} models.Promotion
// @Router /api/v1/promotion/{id} [get]
// @Security Authorization Token
func GetPromotionById(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	Promo, err := repository.FindByIdPromotion(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
	}
	if Promo == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("promotion not found"))
		return
	}
	c.JSON(http.StatusOK, Promo)
}

// GetPromotions godoc
// @Summary Recupera un listado de las Promociones registrados
// @Description Lista Promociones registradas. Valida autorizacion de usuario
// @Tags Promotion
// @Produce json
// @Param name query string true "name"
// @Param code query string true "code"
// @Param used query string true "code"
// @Success 200 {object} []models.Promotion
// @Router /api/v1/promotion [get]
// @Security Authorization Token
func GetPromotions(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	filter := dto.PromotionFilter{
		Name: c.Query("name"),
		Code: strings.ToUpper(c.Query("code")),
		Used: c.Query("used"),
	}
	Promo, err := repository.FindAllIPromotion(filter)
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("promotions not found"))
		return
	}
	if Promo == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("promotions not found"))
		return
	}
	c.JSON(http.StatusOK, Promo)
}

// UpdatePromotion godoc
// @Summary Modifica una Promocion segun su ID
// @Description modifica la promocion segun su ID Valida autorizacion de usuario; que exista la promocion a editar; que el codigo de promocion no este registrado
// @Tags Promotion
// @Accept	json
// @Produce json
// @Param id path integer true "Promotion ID"
// @Param PromoDTO body dto.PromotionDTOUpdate true "update promotion"
// @Success 200 {object} models.Promotion
// @Router /api/v1/promotion/{id} [put]
// @Security Authorization Token
func UpdatePromotion(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	id := c.Params.ByName("id")
	var PromoDTO dto.PromotionDTOUpdate

	err := c.ShouldBindJSON(&PromoDTO)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}
	ID, _ := strconv.ParseUint(id, 10, 64)
	PromoValid, err := repository.FindByIdPromotion(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if PromoValid == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("promotion not found"))
		return
	}

	if PromoDTO.Code != "" {
		PromoCodeValid, _ := repository.FindByCodePromotion(strings.ToUpper(PromoDTO.Code))
		if PromoCodeValid != nil && PromoCodeValid.ID != PromoValid.ID {
			httperror.NewError(c, http.StatusBadRequest, errors.New("promotion code not available"))
			return
		}
	}

	PromoUpdate, err := repository.UpdatePromotion(ID, PromoDTO)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, PromoUpdate)
}

// DeletePromotion godoc
// @Summary Elimina una promocion segun su ID
// @Description Elimina una promocion segun su ID Valida Autorizacion de usuario; promocion a eliminar exista; promocion asociada a una orden
// @Tags Promotion
// @Produce json
// @Param id path integer true "Promotion ID"
// @Success 204 {object} dto.MessageInfo
// @Router /api/v1/promotion/{id} [delete]
// @Security Authorization Token
func DeletePromotion(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	Promo, err := repository.FindByIdPromotion(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
	}
	if Promo == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("promotion not found"))
		return
	}

	if Promo.Used {
		httperror.NewError(c, http.StatusBadRequest, errors.New("promotion is associated with an order, it cannot be deleted"))
		return
	}
	err = repository.DeletePromotion(ID)
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, err)
		log.Println(err)
		return
	}
	c.JSON(http.StatusNoContent, dto.MessageInfo{Status: http.StatusNoContent, Message: "promotion deleted"})
}
