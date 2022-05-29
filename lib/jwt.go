package lib

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type MyClaims struct {
	UserTokenList
	jwt.StandardClaims
}

//用户token参数
type UserTokenList struct {
	UserToken    int    `json:"user_token"`
	ShopToken    int    `json:"shop_token"`
	UserPhone    string `json:"user_phone"`
	CompanyToken int    `json:"company_token"`
	Device       string `json:"device"`
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(viper.GetString("jwt.screat_key")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
