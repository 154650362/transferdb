package jwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/wentaojin/transferdb/conf"

	"time"
)

// todo 这里这样写会有个循环引入的问题。。import cycle not allowed， 需要解决
// 好像通过分层解决了
var jwtSecret = []byte(conf.Gcfg.AppConfig.JwtSecret)

//var jwtSecret = []byte("123456")

//这里用个函数来获取

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(90 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "transdb",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
