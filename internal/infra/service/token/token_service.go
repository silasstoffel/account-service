package token

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/silasstoffel/account-service/internal/domain/auth"
	"github.com/silasstoffel/account-service/internal/exception"
)

type TokenService struct {
	Secret           string
	EmittedBy        string
	ExpiresInMinutes int
}

func (ref *TokenService) CreateToken(data string) (*auth.CreateTokenOutput, error) {
	now := time.Now()
	iss := ref.EmittedBy
	if iss == "" {
		iss = "account-service"
	}
	expires := now.Add(time.Minute * time.Duration(ref.ExpiresInMinutes))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": data,
			"iss": iss,
			"exp": expires.Unix(),
			"iat": now.Unix(),
		})

	signed, err := token.SignedString([]byte(ref.Secret))
	if err != nil {
		log.Fatalln("Error when sign token", err)
		return nil, err
	}

	return &auth.CreateTokenOutput{
		AccessToken: signed,
		ExpiresIn:   expires,
		CreatedAt:   now,
	}, nil
}

func (ref *TokenService) VerifyToken(token string) (*auth.VerifyTokenOutput, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(ref.Secret), nil
	})

	if err != nil {
		message := "Error when parse token"
		log.Println(message, err)
		return nil, exception.New(auth.ErrorParseToken, message, err, exception.HttpUnauthorized)
	}

	if !t.Valid {
		message := "Invalid token"
		return nil, exception.New(auth.ErrorParseToken, message, err, exception.HttpUnauthorized)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		message := "Error when convert token"
		return nil, exception.New(auth.ErrorConvertToken, message, err, exception.HttpUnauthorized)
	}

	return &auth.VerifyTokenOutput{
		Sub:       claims["sub"].(string),
		ExpiresIn: time.Unix(int64(claims["exp"].(float64)), 0),
	}, nil
}
