package encryption

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ValidHeaderXUserId(c *gin.Context) (uint64, error) {
	if c.GetHeader("X-USERID") == "" {
		return 0, errors.New("user not authorized")
	}
	bb, err := FromBase64(c.GetHeader("X-USERID"))
	if err != nil {
		return 0, errors.New("user not authorized")
	}
	decrypted, err := Decrypt(bb)
	if err != nil {
		return 0, errors.New("user not authorized")
	}
	xuserid := BytetoUint64(decrypted)

	return xuserid, nil
}
