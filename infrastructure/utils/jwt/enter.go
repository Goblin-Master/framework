package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

const (
	TOKEN_USER_ID         = "user_id"
	AUTH_ENUMS_ATOKEN     = "atoken"
	AUTH_ENUMS_RTOKEN     = "rtoken"
	ATOKEN_EFFECTIVE_TIME = time.Hour * 2
	RTOKEN_EFFECTIVE_TIME = time.Hour * 24 * 7
)

type MyClaims struct {
	UserID int64  `json:"user_id"`
	Class  string `json:"class"`
	jwt.RegisteredClaims
}

var mySecret = []byte("framework")

type TokenData struct {
	UserID int64
	Class  string
	Time   time.Duration
}

func GenToken(data TokenData) (string, error) {
	// 创建一个我们自己的声明
	claims := MyClaims{
		data.UserID,
		data.Class,
		jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(data.Time)), // 过期时间
			Issuer:    "framework",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	if tokenString == "" {
		return nil, errors.New("token为空")
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token已过期")
		}
		if strings.Contains(err.Error(), "signature is invalid") {
			return nil, errors.New("token无效")
		}
		if strings.Contains(err.Error(), "token contains an invalid") {
			return nil, errors.New("token非法")
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token无效")
}

// 用于验证令牌是否有效
func IdentifyToken(token string) (TokenData, error) {
	//解析token
	claim, err := ParseToken(token)
	var data TokenData
	if err != nil {
		return TokenData{}, err
	}
	data.UserID = claim.UserID
	data.Class = claim.Class
	if claim.Class == AUTH_ENUMS_RTOKEN {
		data.Time = RTOKEN_EFFECTIVE_TIME - time.Duration(time.Now().Unix()-claim.RegisteredClaims.NotBefore.Unix())
	} else {
		data.Time = ATOKEN_EFFECTIVE_TIME - time.Duration(time.Now().Unix()-claim.RegisteredClaims.NotBefore.Unix())
	}
	return data, nil
}

func FullToken(class string, user_id int64) (data TokenData) {
	data.UserID = user_id
	if class == AUTH_ENUMS_ATOKEN {
		data.Time = ATOKEN_EFFECTIVE_TIME
		data.Class = AUTH_ENUMS_ATOKEN
	} else {
		data.Time = RTOKEN_EFFECTIVE_TIME
		data.Class = AUTH_ENUMS_RTOKEN
	}
	return
}

// ForTest
//
//	@Description: 用于测试，可以自己生成一个自定义时间的atoken
//	@param user_id
//	@param duration
func ForTest(user_id int64, duration time.Duration) (string, error) {
	var data TokenData
	data.UserID = user_id
	data.Time = duration
	data.Class = AUTH_ENUMS_ATOKEN
	return GenToken(data)
}
