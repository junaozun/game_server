package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("873dhc7sfs09d")

type MyClaims struct {
	Uid uint64
	jwt.StandardClaims
}

// SetToken 生成Token
func SetToken(uid uint64) (string, error) {
	// 过期时间 默认7天
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &MyClaims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, error) {
	setToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || setToken == nil {
		return nil, err
	}
	if key, _ := setToken.Claims.(*MyClaims); setToken.Valid {
		return key, nil
	} else {
		return nil, fmt.Errorf("setToken.Valid")
	}
}
