package auth

import (
	"context"
	"fmt"
	"log"
	"url-shortener/internal/model"

	"github.com/google/go-github/v55/github"
)

type Storage interface {
	CreateOrUpdate(context.Context, model.User) (*model.User, error)
	Update(context.Context, model.User) error
	GetByGithubLogin(context.Context, string) (*model.User, error)
	Deactivate(context.Context, string) error
}

//go:generate moq --out=mock_github_client.gen.go --pkg=auth . GithubClient
type GithubClient interface {
	ExchangeCodeToAccessKey(ctx context.Context, clientID, clientSeret, code string) (string, error)
	GetUser(ctx context.Context, accessKey, user string) (*github.User, error)
}

type Service struct {
	github  GithubClient
	storage Storage

	ghClientID     string
	ghClientSecret string
}

func NewService(githubClient GithubClient, storage Storage, ghClientID, ghClientSecret string) *Service {
	return &Service{
		github:         githubClient,
		storage:        storage,
		ghClientID:     ghClientID,
		ghClientSecret: ghClientSecret,
	}
}

func (s *Service) GithubAuthLink() string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?scopes=user&client_id=%s", s.ghClientID)
}

func (s *Service) GithubAuthCallback(ctx context.Context, sessionCode string) (*model.User, string, error) {
	accessKey, err := s.github.ExchangeCodeToAccessKey(ctx, s.ghClientID, s.ghClientSecret, sessionCode)
	if err != nil {
		return nil, "", err
	}

	ghUser, err := s.github.GetUser(ctx, accessKey, "")
	if err != nil {
		return nil, "", err
	}

	user, err := s.RegisterUser(ctx, ghUser, accessKey)
	if err != nil {
		return nil, "", err
	}

	jwt, err := MakeJWT(*user)
	if err != nil {
		log.Printf("failed to make jwt: %v", err)
		return nil, "", err
	}

	return user, jwt, nil
}

func (s *Service) RegisterUser(ctx context.Context, ghUser *github.User, accessKey string) (*model.User, error) {
	user := model.User{
		GithubLogin:     ghUser.GetLogin(),
		IsActive:        true,
		GithubAccessKey: accessKey,
	}

	return s.storage.CreateOrUpdate(ctx, user)
}
