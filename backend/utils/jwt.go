package utils

import (
	"errors"
	"time"

	"MyBlog/global"

	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	AdminID  uint32 `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成管理员token
func GenerateAdminToken(adminID uint32, username string) (string, error) {
	jwtKey := []byte(global.Config.Jwt.Secret)

	expirationTime := time.Now().Add(time.Duration(global.Config.Jwt.ExpireHours) * time.Hour)

	claims := &AdminClaims{
		AdminID:  adminID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "admin_auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

// 解析管理员token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	jwtKey := []byte(global.Config.Jwt.Secret)

	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		//防止别人改签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
