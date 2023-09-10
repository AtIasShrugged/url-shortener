package auth_test

import (
	"context"
	"errors"
	"testing"
	"url-shortener/internal/auth"
	"url-shortener/internal/model"
	"url-shortener/internal/storage/user"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/go-github/v55/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGithubAuthLink(t *testing.T) {
	t.Run("returns github auth link with client id", func(t *testing.T) {
		svc := auth.NewService(nil, nil, "client-id", "")
		assert.Equal(
			t,
			"https://github.com/login/oauth/authorize?scopes=user&client_id=client-id",
			svc.GithubAuthLink(),
		)
	})
}

func TestService_GithubAuthCallback(t *testing.T) {
	t.Run("returns user model and JWT", func(t *testing.T) {
		var (
			ghClient = &auth.GithubClientMock{
				ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSecret, code string) (string, error) {
					return "access-key", nil
				},
				GetUserFunc: func(ctx context.Context, accessToken, user string) (*github.User, error) {
					return &github.User{
						Login: github.String(gofakeit.Username()),
					}, nil
				},
			}
			userStorage = user.NewInMemory()
			svc         = auth.NewService(ghClient, userStorage, "", "")
		)

		user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
		require.NoError(t, err)
		assert.True(t, user.IsActive)
		assert.NotEmpty(t, token)
	})

	t.Run("returns error", func(t *testing.T) {
		t.Run("when exchanging code to access key fails", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSeret, code string) (string, error) {
						return "", errors.New("exchange code to access key error")
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
			)

			user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Empty(t, token)
		})

		t.Run("when getting user from github fails", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSeret, code string) (string, error) {
						return "access-key", nil
					},
					GetUserFunc: func(ctx context.Context, accessKey, user string) (*github.User, error) {
						return nil, errors.New("error getting the user")
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
			)

			user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Empty(t, token)
		})

		t.Run("when registering user fails", func(t *testing.T) {
			var (
				ghClient = &auth.GithubClientMock{
					ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID, clientSeret, code string) (string, error) {
						return "access-key", nil
					},
					GetUserFunc: func(ctx context.Context, accessKey, user string) (*github.User, error) {
						return nil, errors.New("get user error")
					},
				}
				userStorage = user.NewInMemory()
				svc         = auth.NewService(ghClient, userStorage, "", "")
			)

			user, token, err := svc.GithubAuthCallback(context.Background(), gofakeit.Numerify("code-###"))
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Empty(t, token)
		})
	})
}

func TestService_RegisterUser(t *testing.T) {
	t.Run("returns user model", func(t *testing.T) {
		t.Run("even is user already exists", func(t *testing.T) {
			var (
				userStorage  = user.NewInMemory()
				svc          = auth.NewService(nil, userStorage, "", "")
				ghUser       = &github.User{Login: github.String(gofakeit.Username())}
				existingUser = model.User{
					IsActive:    true,
					GithubLogin: ghUser.GetLogin(),
				}
			)

			_, err := userStorage.CreateOrUpdate(context.Background(), existingUser)
			require.NoError(t, err)

			user, err := svc.RegisterUser(context.Background(), ghUser, "")
			require.NoError(t, err)
			assert.Equal(t, existingUser.GithubLogin, user.GithubLogin)
			assert.Equal(t, existingUser.IsActive, user.IsActive)
		})
	})
}
