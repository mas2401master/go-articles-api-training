package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/repository"
	"github.com/mas2401master/go-articles-api-training/pkg/encryption"
	"github.com/mas2401master/go-articles-api-training/pkg/httperror"
)

// CreateUser godoc
// @Summary Crea token de autorizacion de usuario
// @Description Crea genera token jwt
// @Tags Login
// @Accept	json
// @Produce json
// @Param userdto body dto.LoginDTO true "Login"
// @Success 200 {object} dto.Logedin
// @Router /api/v1/login [post]
func Login(c *gin.Context) {
	var logindto dto.LoginDTO
	_ = c.BindJSON(&logindto)

	user, err := repository.FindByUsername(logindto.Username)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	if user == nil {
		httperror.NewError(c, http.StatusNotFound, errors.New("user not found"))
		return
	}

	bb, _ := encryption.FromBase64(user.Password)
	decryptedPassword, _ := encryption.Decrypt(bb)
	if string(decryptedPassword) != logindto.Password {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user and password not match"))
		return
	}

	if !user.Status {
		httperror.NewError(c, http.StatusUnauthorized, errors.New("user inactive"))
		return
	}
	xuserid := strconv.FormatUint(user.ID, 10)
	xrole := strconv.FormatUint(user.RoleID, 10)
	token, err := encryption.SignedLoginToken(xuserid, xrole)
	if err != nil {
		httperror.NewError(c, http.StatusInternalServerError, errors.New("unternal server error"))
		return
	}

	encryptedUserid, _ := encryption.Encrypt(encryption.Uint64toByte(user.ID))
	Idbb := encryption.ToBase64(encryptedUserid)
	c.Writer.Header().Set("X-USERID", Idbb)

	c.JSON(http.StatusOK, dto.Logedin{
		UserID:    Idbb,
		FirstName: user.Firstname,
		Lastname:  user.Lastname,
		Token:     token,
	})
}
