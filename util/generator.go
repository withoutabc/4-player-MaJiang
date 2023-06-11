package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"majiang/model"
	"time"
)

var Secret = []byte("YJX")

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err

}

func HashWithSalt(password string, salt []byte) []byte {
	salted := append(salt, []byte(password)...)
	hashed := sha256.Sum256(salted)
	return hashed[:]
}

// GenerateRoomId 生成房间ID
func GenerateRoomId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GenToken GenToken生成aToken和rToken
func GenToken(userId int) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := model.MyClaims{
		ID: userId, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(), // 过期时间
			Issuer:    "YJX",                            // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(Secret)
	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 过期时间
		Issuer:    "YJX",                                 // 签发人
	}).SignedString(Secret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return
}
