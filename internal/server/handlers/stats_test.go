package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/model"
	"url-shortener/internal/server/handlers"
	"url-shortener/internal/shorten"
	"url-shortener/internal/storage/shortening"

	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleStats(t *testing.T) {
	t.Run("returns shortening with given identifier", func(t *testing.T) {
		var (
			provider = shorten.NewService(shortening.NewInMemory())
			handler  = handlers.HandleStats(provider)
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest("GET", "/stats/abc", nil)
			e        = echo.New()
			c        = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		c.SetPath("/stats/:identifier")
		c.SetParamNames("identifier")
		c.SetParamValues("abc")

		_, err := provider.Shorten(
			context.Background(),
			model.ShortenInput{
				Identifier: mo.Some("abc"),
				RawURL:     "https://google.com",
				CreatedBy:  "user",
			},
		)
		require.NoError(t, err)

		require.NoError(t, handler(c))

		var s model.Shortening
		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&s))

		assert.Equal(t, "abc", s.Identifier)
		assert.Equal(t, "https://google.com", s.OriginalURL)
		assert.Equal(t, "user", s.CreatedBy)
		assert.Equal(t, int64(0), s.Visits)
	})

	t.Run("returns 404 if shortenung with given identifier does not exist", func(t *testing.T) {
		var (
			provider = shorten.NewService(shortening.NewInMemory())
			handler  = handlers.HandleStats(provider)
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest("GET", "/stats/abc", nil)
			e        = echo.New()
			c        = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		c.SetPath("/stats/:identifier")
		c.SetParamNames("identifier")
		c.SetParamValues("abc")

		var httpErr *echo.HTTPError
		require.ErrorAs(t, handler(c), &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})
}
