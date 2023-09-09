package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/mo"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrIdentifierExists = errors.New("identifier already exists")
)

type Shortening struct {
	Identifier  string    `json:"identifier"`
	OriginalUrl string    `json:"original_url"`
	Visits      int64     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShortenInput struct {
	RawURL     string
	Identifier mo.Option[string]
	CreatedBy  string
}

type User struct {
	IsActive    bool   `json:"is_verified,omitempty"`
	GithubLogin string `json:"gh_login"`

	// TODO: maybe store it in something like vault
	GithubAccessKey string    `json:"gh_access_key,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	User `json:"user_data"`
}
