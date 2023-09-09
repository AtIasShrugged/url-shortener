package shorten_test

import (
	"context"
	"testing"
	"url-shortener/internal/model"
	"url-shortener/internal/shorten"
	"url-shortener/internal/storage/shortening"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Shorten(t *testing.T) {
	t.Run("generates shortening for a given URL", func(t *testing.T) {
		const url = "https://google.com"

		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = model.ShortenInput{RawURL: url}
		)

		shortening, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		require.NotEmpty(t, shortening.Identifier)
		assert.Equal(t, input.RawURL, shortening.OriginalURL)
		assert.NotZero(t, shortening.CreatedAt)
	})

	t.Run("uses custom identifier if provided", func(t *testing.T) {
		const identifier = "google"
		const url = "https://google.com"

		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = model.ShortenInput{
				RawURL:     url,
				Identifier: mo.Some(identifier),
			}
		)

		shortening, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		require.Equal(t, identifier, shortening.Identifier)
		assert.Equal(t, url, shortening.OriginalURL)
		assert.NotZero(t, shortening.CreatedAt)
	})

	t.Run("returns error if identifier is already taken", func(t *testing.T) {
		const identifier = "google"
		const url = "https://google.com"

		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = model.ShortenInput{
				RawURL:     url,
				Identifier: mo.Some(identifier),
			}
		)

		_, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		_, err = svc.Shorten(context.Background(), input)
		require.ErrorIs(t, err, model.ErrIdentifierExists)
	})
}
