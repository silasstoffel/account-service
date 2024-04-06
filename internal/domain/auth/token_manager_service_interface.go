package auth

import "time"

type CreateTokenOutput struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   time.Time
	CreatedAt   time.Time
}

type VerifyTokenOutput struct {
	Sub       string
	ExpiresIn time.Time
}

const (
	ErrorParseToken   = "ERROR_PARSE_TOKEN"
	InvalidToken      = "INVALID_TOKEN"
	ErrorConvertToken = "ERROR_CONVERT_TOKEN"
)

type TokenManagerService interface {
	CreateToken(data string) (*CreateTokenOutput, error)
	VerifyToken(token string) (*VerifyTokenOutput, error)
}
