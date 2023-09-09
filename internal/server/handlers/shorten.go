package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"url-shortener/internal/config"
	"url-shortener/internal/model"
	"url-shortener/internal/shorten"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
)

type shortener interface {
	Shorten(context.Context, model.ShortenInput) (*model.Shortening, error)
}

type shortenRequest struct {
	URL        string `json:"url" validate:"required,url"`
	Identifier string `json:"identifier,omitempty" validate:"omitempty,alphanum"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url,omitempty"`
	Message  string `json:"message,omitempty"`
}

func HandleShorten(shortener shortener) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req shortenRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := c.Validate(req); err != nil {
			return err
		}

		userToken, ok := c.Get("user").(*jwt.Token)
		fmt.Println(userToken, ok)
		if !ok {
			log.Println("error: user is not presented in context")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		userClaims, ok := userToken.Claims.(*model.UserClaims)
		if !ok {
			log.Println("error: failed to get user claims from token")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		identifier := mo.None[string]()
		if strings.TrimSpace(req.Identifier) != "" {
			identifier = mo.Some(req.Identifier)
		}

		input := model.ShortenInput{
			RawURL:     req.URL,
			Identifier: identifier,
			CreatedBy:  userClaims.User.GithubLogin,
		}

		shortening, err := shortener.Shorten(c.Request().Context(), input)
		if err != nil {
			if errors.Is(err, model.ErrIdentifierExists) {
				return echo.NewHTTPError(http.StatusConflict, model.ErrIdentifierExists.Error())
			}

			log.Printf("error shortenin url %q: %v", req.URL, err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		shortURL, err := shorten.PrependBaseUrl(config.Get().BaseURL, shortening.Identifier)
		if err != nil {
			log.Printf("shorten.PrependBaseUrl: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(
			http.StatusOK,
			shortenResponse{ShortURL: shortURL},
		)
	}
}
