package common

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"serve/model"
	"time"
)

var jwtKey = []byte("create_key")

type Claims struct {
	UserId uint
	jwt.RegisteredClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := Claims{
		user.ID,
		jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// 发放时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 发放人
			Issuer: "oceanlearn.tech",
			// 	主题
			Subject: "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("token生成失败")
	}
	return tokenString, err
}

// ParseToken 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Println("token解析失败", err)
	}
	return token, claims, err
}
