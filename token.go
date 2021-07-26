package portgate

import (
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
	"time"
)

func CreateToken(config *Config, givenKey string) (string, error) {
	// Our token will last 7 days
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: GetExpirationDateFrom(time.Now()).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return token.SignedString([]byte(config.jwtSecret))
}

func VerifyToken(config *Config, tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.jwtSecret), nil
	})

	if err == nil && token.Valid {
		return true, nil
	} else {
		return false, err
	}
}

func VerifyTokenFromCookie(config *Config, ctx *fasthttp.RequestCtx) bool {
	cookie := ctx.Request.Header.Cookie("_portgate_token")
	if cookie != nil {
		ok, _ := VerifyToken(config, string(cookie))
		return ok
	}
	return false
}

func GetExpirationDateFrom(date time.Time) time.Time {
	return date.AddDate(0, 0, 7)
}
