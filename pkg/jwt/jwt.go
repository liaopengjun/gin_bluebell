package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)
//过期时间
const TokenExpireDuration = time.Hour * 2

//加密盐
var MySecret = []byte("gogogo")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return MySecret, nil
}

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims jwt包自带的jwt.StandardClaims只包含了官方字段
type MyClaims struct {
	UserID int64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

//生成Token
func GenToken(UserID uint64,Username string) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		int64(UserID), // 自定义字段
		Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "bluebell",                                 // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(MySecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(), // 过期时间
		Issuer:    "bluebell",                              // 签发人
	}).SignedString(MySecret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}


//解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}
	// 从旧access token中解析出claims数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(uint64(claims.UserID),claims.Username)
	}
	return
}




