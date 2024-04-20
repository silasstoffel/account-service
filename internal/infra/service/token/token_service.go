package token

import (
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
		return nil, exception.New(exception.ErrorParseToken, &err)
	}

	if !t.Valid {
		return nil, exception.New(exception.ErrorParseToken, &err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, exception.New(exception.ErrorConvertToken, &err)
	}

	return &auth.VerifyTokenOutput{
		Sub:       claims["sub"].(string),
		ExpiresIn: time.Unix(int64(claims["exp"].(float64)), 0),
	}, nil
}
