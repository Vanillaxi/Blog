package utils

import (
	"errors"
	"log"
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

	tokenString, err := token.SignedString(jwtKey)
	log.Printf("[auth-debug] GenerateAdminToken admin_id=%d username=%s secret_len=%d expire_hours=%d expires_at=%s token_len=%d err=%v",
		adminID,
		username,
		len(jwtKey),
		global.Config.Jwt.ExpireHours,
		expirationTime.Format(time.RFC3339),
		len(tokenString),
		err,
	)

	return tokenString, err
}

// 解析管理员token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	jwtKey := []byte(global.Config.Jwt.Secret)
	log.Printf("[auth-debug] ParseAdminToken start token_len=%d secret_len=%d", len(tokenString), len(jwtKey))

	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		//防止别人改签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("[auth-debug] ParseAdminToken error=%v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		log.Printf("[auth-debug] ParseAdminToken invalid claims_ok=%v token_valid=%v", ok, token.Valid)
		return nil, errors.New("invalid token")
	}

	expiresAt := "<nil>"
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time.Format(time.RFC3339)
	}
	issuedAt := "<nil>"
	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Time.Format(time.RFC3339)
	}
	log.Printf("[auth-debug] ParseAdminToken success admin_id=%d username=%s subject=%s issued_at=%s expires_at=%s",
		claims.AdminID,
		claims.Username,
		claims.Subject,
		issuedAt,
		expiresAt,
	)

	return claims, nil
}
