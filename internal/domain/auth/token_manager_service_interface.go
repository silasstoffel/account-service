package auth

import "time"

type CreateTokenOutput struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   time.Time
	CreatedAt   time.Time
}

type TokenManagerService interface {
	CreateToken(data string) (*CreateTokenOutput, error)
}
