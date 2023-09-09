package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	BaseURL string
	Host    string
	Port    int
	DB      DBConfig
	Github  GithubConfig
	Auth    AuthConfig
}

func (c Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type DBConfig struct {
	DSN      string
	Database string
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
}

type AuthConfig struct {
	JWTSecretKey string
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		cfg = func() Config {
			return Config{
				BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
				Host:    getEnv("HOST", "0.0.0.0"),
				Port:    getEnvAsInt("PORT", 8080),
				DB: DBConfig{
					DSN:      getEnv("MONGODB_DSN", ""),
					Database: getEnv("MONGODB_DATABASE", ""),
				},
				Github: GithubConfig{
					ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
					ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
				},
				Auth: AuthConfig{
					JWTSecretKey: getEnv("JWT_SECRET", "secret"),
				},
			}
		}()
	})

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
