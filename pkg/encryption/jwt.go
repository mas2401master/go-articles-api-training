package encryption

import (
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/config"
)

func SignedLoginToken(xuserid, xrole string) (string, error) {
	config := config.GetConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iduser": xuserid,
		"role":   xrole,
	})

	return token.SignedString([]byte(config.Server.Secret))
}

func ParseLoginJWT(tokenString string) (jwt.MapClaims, error) {
	config := config.GetConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func ValidateToken(tokenString string) bool {
	config := config.GetConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Server.Secret), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

func ClaimsFromToken(tokenString string) (uint64, uint64, string) {
	claims, _ := ParseLoginJWT(tokenString)
	xuserid := claims["iduser"].(string)
	role := claims["role"].(string)
	userid, _ := strconv.ParseUint(xuserid, 10, 64)
	roleid, _ := strconv.ParseUint(role, 10, 64)
	rolename := "ADMIN"
	if roleid == 2 {
		rolename = "CLIENTE"
	}
	return userid, roleid, rolename
}
