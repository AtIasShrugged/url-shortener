package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/model"

	"github.com/labstack/echo/v4"
)

type callbackProvider interface {
	GithubAuthCallback(ctx context.Context, sessionCode string) (*model.User, string, error)
}

func HandleGithubAuthCallback(cbProvider callbackProvider) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionCode := c.QueryParam("code")
		if sessionCode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "missing session code")
		}

		_, jwt, err := cbProvider.GithubAuthCallback(c.Request().Context(), sessionCode)
		if err != nil {
			log.Printf("error handling github auth callback: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		redirectUrl := fmt.Sprintf("%s/auth/token.html?token=%s", config.Get().BaseURL, jwt)
		return c.Redirect(http.StatusMovedPermanently, redirectUrl)
	}
}
