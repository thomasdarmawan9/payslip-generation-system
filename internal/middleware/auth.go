package middleware

import (
	"strings"

	"payslip-generation-system/config"
	errorUc "payslip-generation-system/internal/error"
	"payslip-generation-system/pkg/log"
	"payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthHandler struct {
	cfg *config.Config
	log *log.LogCustom
}

func New(r *gin.RouterGroup, c *config.Config, l *log.LogCustom) {
	handler := AuthHandler{
		cfg: c,
		log: l,
	}
	r.Use(handler.AuthJwt)
}

func (au *AuthHandler) AuthJwt(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		au.log.Error(log.LogData{Err: utils.MakeError(errorUc.ErrUnauthorized)})
		utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrUnauthorized))))
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrUnauthorized))))
		c.Abort()
		return
	}
	tokenString := parts[1]

	// Ambil secret dari config (sesuaikan path-nya dengan config kamu)
	secret := "" // contoh: isi di config kamu
	if secret == "" {
		// fallback kalau belum diisi; tapi sebaiknya DIHAPUS untuk production
		secret = "A7M+TXRMxdz0N3nFLjGaxVKgkELowtbxWipS+IFZkVE="
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// pastikan HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			au.log.Error(log.LogData{Err: utils.MakeError(errorUc.ErrUnauthorized)})
			return nil, utils.MakeError(errorUc.ErrUnauthorized)
		}
		return []byte(secret), nil
	})
	if err != nil || token == nil || !token.Valid {
		au.log.Error(log.LogData{Err: utils.MakeError(errorUc.ErrUnauthorized)})
		utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrUnauthorized))))
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrUnauthorized))))
		c.Abort()
		return
	}

	// Ambil claims dengan aman
	var userID uint
	if v, ok := claims["user_id"]; ok {
		switch t := v.(type) {
		case float64:
			userID = uint(t)
		case int:
			userID = uint(t)
		}
	}
	name, _ := claims["name"].(string)
	role, _ := claims["role"].(string)

	if userID == 0 || name == "" || role == "" {
		utils.Failed(c, utils.CustomError(errorUc.ErrorCustom(utils.MakeError(errorUc.ErrUnauthorized))))
		c.Abort()
		return
	}

	c.Set("user_id", userID)
	c.Set("name", name)
	c.Set("role", role)

	c.Next()
}
