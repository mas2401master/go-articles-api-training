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

// CreateUser godoc
// @Summary Crea registro de los usuarios
// @Description Crea registro de usuarios tipo CLIENTE y ADMIN Valida que el nombre de usuario no exista
// @Tags User
// @Accept	json
// @Produce json
// @Param userdto body dto.UserDTO true "Create user"
// @Success 200 {object} dto.MessageInfo
// @Router /api/v1/users [post]
// @Security Authorization Token
func CreateUser(c *gin.Context) {
	//Authorization - token
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	var userdto dto.UserDTO
	var rol *models.Role
	var userValid *models.User
	err := c.ShouldBindJSON(&userdto)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		log.Println(err)
		return
	}
	user := models.User{
		Username:  userdto.Username,
		Firstname: userdto.Firstname,
		Lastname:  userdto.Lastname,
		Password:  userdto.Password,
		Status:    true,
		Email:     userdto.Email,
		RoleID:    userdto.RoleID,
	}
	userValid, _ = repository.FindByUsername(user.Username)
	if userValid != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("Username: "+user.Username+" already exists..."))
		return
	}
	rol, _ = repository.FindByIdRol(userdto.RoleID)
	if rol == nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("role id does not exist"))
		return
	}

	err = repository.AddUser(user)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, dto.MessageInfo{Status: http.StatusCreated, Message: "user created!!"})
}

// GetUserById godoc
// @Summary Recupera una Usuario segun su ID
// @Description Obtiene informacion del usuario segun su ID Valida autorizacion de usuario, usuario a mostrar asociado a usuario de session para el rol cliente
// @Tags User
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.User
// @Router /api/v1/users/{id} [get]
// @Security Authorization Token
func GetUserById(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	if rolename == "CLIENTE" && ID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	user, err := repository.FindById(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	if user == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("user not found"))
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary Recupera un listado de las Usuarios registrados
// @Description Lista Usuarios registradas Valida autorizacion de usuario
// @Tags User
// @Produce json
// @Param status query string true "status"
// @Param role query string true "role"
// @Success 200 {object} []models.User
// @Router /api/v1/users [get]
// @Security Authorization Token
func GetUsers(c *gin.Context) {
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}
	role, _ := strconv.ParseUint(c.Query("role"), 10, 64)
	filter := dto.UserFilter{
		Status: c.Query("status"),
		RoleID: role,
	}
	users, err := repository.FindAllUser(filter)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if users == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("users not found"))
		return
	}
	c.JSON(http.StatusOK, users)

}

// UpdateUser godoc
// @Summary Modifica un Usuario segun su ID
// @Description modifica usuario segun su ID Valida autorizacion de usuario; existencia de usuario a editar; existencia del username; usuario a editar asociado a usuario de session para el rol cliente
// @Tags User
// @Accept	json
// @Produce json
// @Param id path integer true "User ID"
// @Param userdto body dto.UserUpdateDTO true "Edit user"
// @Success 200 {object} models.User
// @Router /api/v1/users/{id} [put]
// @Security Authorization Token
func UpdateUser(c *gin.Context) {
	//Authorization - token
	userid, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	id := c.Params.ByName("id")
	ID, _ := strconv.ParseUint(id, 10, 64)

	var userNew dto.UserUpdateDTO
	var userValid *models.User
	err := c.ShouldBindJSON(&userNew)
	if err != nil {
		httperror.NewError(c, http.StatusBadRequest, err)
		return
	}

	userOld, err := repository.FindById(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if userOld == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("user not found"))
		return
	}

	if rolename == "CLIENTE" && userOld.ID != userid {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	if userOld.Username != userNew.Username {
		userValid, _ = repository.FindByUsername(userNew.Username)
		if userValid != nil && userOld.ID != userValid.ID {
			httperror.NewError(c, http.StatusBadRequest, errors.New("Username: "+userValid.Username+" already exists..."))
			return
		}
	}

	user, err := repository.UpdateUser(ID, userNew)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Elimina un usuario segun su ID
// @Description Elimina una usuario segun su ID Valida autorizacion de usuario; existencia de usuario a eliminar;usuario no asociado a una orden
// @Tags User
// @Produce json
// @Param id path integer true "User ID"
// @Success 204 {object} dto.MessageInfo
// @Router /api/v1/users/{id} [delete]
// @Security Authorization Token
func DeleteUser(c *gin.Context) {
	_, _, rolename := encryption.ClaimsFromToken(c.GetHeader("Authorization"))
	if rolename == "CLIENTE" {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user not authorized"))
		return
	}

	id := c.Param("id")
	ID, _ := strconv.ParseUint(id, 10, 64)
	user, err := repository.FindById(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("user not found"))
		return
	}

	order, _ := repository.FindByOrderUser(ID)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if order != nil {
		httperror.NewError(c, http.StatusBadRequest, errors.New("user associated with an order"))
		return
	}
	err = repository.DeleteUser(ID)
	if err != nil {
		httperror.NewError(c, http.StatusNotFound, err)
		log.Println(err)
		return
	}
	c.JSON(http.StatusNoContent, dto.MessageInfo{Status: http.StatusNoContent, Message: "user deleted"})
}
