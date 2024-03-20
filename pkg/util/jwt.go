package util

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// 该文件用于写jwt的基础操作
type Claims struct {
	User     string
	Password string
	jwt.StandardClaims
}

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(UserName, Password string) (string, error) {
	//生成token
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		User:     UserName,
		Password: Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "to-do-list",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtStr, err := tokenClaims.SignedString(JwtSecret)
	return jwtStr, err

}

func ParseToken(Token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
