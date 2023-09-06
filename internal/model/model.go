package model

import (
	"errors"
	"time"

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
