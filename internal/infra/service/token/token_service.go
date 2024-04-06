package token

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/silasstoffel/account-service/internal/domain/auth"
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
