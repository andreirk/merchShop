package jwt

import "github.com/golang-jwt/jwt/v5"

type JwtData struct {
	Username string
}

type Jwt struct {
	Secret string
}

func NewJwt(secret string) *Jwt {
	return &Jwt{secret}
}

func (j *Jwt) Sign(data JwtData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Username,
	})
	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *Jwt) Parse(tokenString string) (*JwtData, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, false
	}
	email := token.Claims.(jwt.MapClaims)["email"].(string)
	return &JwtData{email}, token.Valid
}
